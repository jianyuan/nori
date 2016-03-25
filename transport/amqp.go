package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/jianyuan/nori/message"
	"github.com/jianyuan/nori/protocol"
	"github.com/kr/pretty"
	"github.com/streadway/amqp"
	"gopkg.in/tomb.v2"
)

type AMQPTransport struct {
	URL          string
	ExchangeName string
	ExchangeKind string
	tomb         *tomb.Tomb
	conn         *amqp.Connection
	channel      *amqp.Channel
}

func (AMQPTransport) Name() string { return "AMQPTransport" }

func (srv *AMQPTransport) Setup() error {
	// TODO advanced config
	conn, err := amqp.Dial(srv.URL)
	if err != nil {
		return err
	}
	srv.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	srv.channel = ch

	if err := srv.channel.ExchangeDeclare(
		srv.ExchangeName, // name
		srv.ExchangeKind, // kind
		true,             // durable
		false,            // autoDelete
		false,            // internal
		false,            // noWait
		nil,              // args
	); err != nil {
		return err
	}

	return nil
}

func (srv *AMQPTransport) Consume(name string) (<-chan *message.Request, error) {
	deliveryChan, err := srv.consume(name)
	if err != nil {
		return nil, err
	}

	msgChan := make(chan *message.Request)
	srv.tomb.Go(func() error {
		for {
			select {
			case <-srv.tomb.Dying():
				return nil

			case delivery, ok := <-deliveryChan:
				if !ok {
					log.Warnln("Channel closed")
					return nil
				}

				msg, err := srv.parseDelivery(delivery)
				if err != nil {
					log.Warnln("Error parsing message:", err)
					delivery.Nack(false, false)
					continue
				}
				msgChan <- msg
			}
		}
		return nil
	})
	return msgChan, nil
}

func (srv *AMQPTransport) consume(name string) (<-chan amqp.Delivery, error) {
	q, err := srv.channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // autoDelete,
		false, // exclusive
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return nil, err
	}

	if err := srv.channel.QueueBind(
		q.Name,           // name
		q.Name,           // key
		srv.ExchangeName, // exchange
		false,            // noWait
		nil,              // args
	); err != nil {
		return nil, err
	}

	msgs, err := srv.channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // autoAck
		false,  // exclusive
		false,  // noLocal
		false,  // noWait
		nil,    // args
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

func (AMQPTransport) parseDelivery(d amqp.Delivery) (*message.Request, error) {
	switch d.ContentType {
	case "application/json":
		// TODO refactor
		var celeryTask protocol.CeleryTask
		if err := json.Unmarshal(d.Body, &celeryTask); err != nil {
			return nil, err
		}

		pretty.Println("CeleryTask:", celeryTask)

		// TODO other celery fields
		return &message.Request{
			TaskName:  celeryTask.Name,
			ID:        celeryTask.ID,
			Args:      celeryTask.Args,
			KWArgs:    celeryTask.KWArgs,
			ETA:       celeryTask.ETA,
			ExpiresAt: celeryTask.ExpiresAt,
			IsUTC:     celeryTask.IsUTC,
			ReplyTo:   &d.ReplyTo,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported content type %q", d.ContentType)
	}
}

func (srv *AMQPTransport) Tomb() *tomb.Tomb {
	return srv.tomb
}

func (srv *AMQPTransport) Close() error {
	if srv.conn != nil {
		if err := srv.conn.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (srv *AMQPTransport) Reply(req *message.Request, resp message.Response) error {
	replyTo := resp.GetReplyTo()
	if replyTo == nil {
		return errors.New("AMQPTransport: no reply queue specified")
	}

	body, err := messageResponseBytes(resp)
	if err != nil {
		return err
	}

	return srv.channel.Publish(
		"",       // exchange
		*replyTo, // key
		true,     // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			DeliveryMode:  amqp.Persistent,
			CorrelationId: resp.GetID(),
			Timestamp:     time.Now().UTC(),
			Body:          body,
		})
}

func NewAMQPTransport(url string) Driver {
	return &AMQPTransport{
		URL:          url,
		ExchangeName: "celery",
		ExchangeKind: "direct",
		tomb:         new(tomb.Tomb),
	}
}

func messageResponseBytes(resp message.Response) ([]byte, error) {
	p, err := protocol.NewCeleryResult(resp)
	if err != nil {
		return nil, err
	}
	return json.Marshal(p)
}
