package transport

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/jianyuan/nori/log"
	"github.com/jianyuan/nori/message"
	"github.com/jianyuan/nori/protocol"
	"github.com/kr/pretty"
	"github.com/streadway/amqp"
	"gopkg.in/tomb.v2"
)

type AMQPTransport struct {
	context.Context
	URL          string
	ExchangeName string
	ExchangeKind string
	tomb         *tomb.Tomb
	conn         *amqp.Connection
	channel      *amqp.Channel
}

func (AMQPTransport) Name() string { return "AMQPTransport" }

func (t *AMQPTransport) Init(ctx context.Context) error {
	t.Context = ctx
	return nil
}

func (t *AMQPTransport) Setup() error {
	// TODO advanced config
	conn, err := amqp.Dial(t.URL)
	if err != nil {
		return err
	}
	t.conn = conn

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	t.channel = ch

	if err := t.channel.ExchangeDeclare(
		t.ExchangeName, // name
		t.ExchangeKind, // kind
		true,           // durable
		false,          // autoDelete
		false,          // internal
		false,          // noWait
		nil,            // args
	); err != nil {
		return err
	}

	return nil
}

func (t *AMQPTransport) Consume(name string) (<-chan *message.Request, error) {
	deliveryChan, err := t.consume(name)
	if err != nil {
		return nil, err
	}

	msgChan := make(chan *message.Request)
	t.tomb.Go(func() error {
		for {
			select {
			case <-t.tomb.Dying():
				return nil

			case delivery, ok := <-deliveryChan:
				if !ok {
					log.FromContext(t).Warnln("Channel closed")
					return nil
				}

				msg, err := t.parseDelivery(delivery)
				if err != nil {
					log.FromContext(t).Warnln("Error parsing message:", err)
					delivery.Nack(false, false)
					continue
				}
				msgChan <- msg
			}
		}
	})
	return msgChan, nil
}

func (t *AMQPTransport) consume(name string) (<-chan amqp.Delivery, error) {
	q, err := t.channel.QueueDeclare(
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

	if err := t.channel.QueueBind(
		q.Name,         // name
		q.Name,         // key
		t.ExchangeName, // exchange
		false,          // noWait
		nil,            // args
	); err != nil {
		return nil, err
	}

	msgs, err := t.channel.Consume(
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

func (t *AMQPTransport) Tomb() *tomb.Tomb {
	return t.tomb
}

func (t *AMQPTransport) Close() error {
	if t.conn != nil {
		if err := t.conn.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (t *AMQPTransport) Reply(req *message.Request, resp message.Response) error {
	replyTo := resp.GetReplyTo()
	if replyTo == nil {
		return errors.New("AMQPTransport: no reply queue specified")
	}

	body, err := messageResponseBytes(resp)
	if err != nil {
		return err
	}

	return t.channel.Publish(
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
