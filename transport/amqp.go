package transport

import (
	"encoding/json"
	"fmt"

	"github.com/jianyuan/nori/message"
	"github.com/jianyuan/nori/protocol"
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
	rawMsgs, err := srv.consume(name)
	if err != nil {
		return nil, err
	}

	msgChan := make(chan *message.Request)
	srv.tomb.Go(func() error {
		for {
			select {
			case <-srv.tomb.Dying():
				return nil

			case rawMsg, ok := <-rawMsgs:
				if !ok {
					log.Warnln("Channel closed")
					return nil
				}

				msg, err := srv.parseDelivery(rawMsg)
				if err != nil {
					log.Warnln("Error parsing message:", err)
					rawMsg.Nack(false, false)
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
		celeryTask := protocol.CeleryTask{}
		if err := json.Unmarshal(d.Body, &celeryTask); err != nil {
			return nil, err
		}

		// TODO other celery fields
		return &message.Request{
			HandlerName: celeryTask.Name,
			ID:          celeryTask.ID,
			Args:        celeryTask.Args,
			KWArgs:      celeryTask.KWArgs,
			ETA:         celeryTask.ETA,
			ExpiresAt:   celeryTask.ExpiresAt,
			IsUTC:       celeryTask.IsUTC,
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

func NewAMQPTransport(url string) Driver {
	return &AMQPTransport{
		URL:          url,
		ExchangeName: "celery",
		ExchangeKind: "direct",
		tomb:         new(tomb.Tomb),
	}
}
