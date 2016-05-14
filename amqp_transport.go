package nori

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/streadway/amqp"
	"golang.org/x/net/context"
	"gopkg.in/tomb.v2"
)

var (
	DefaultAMQPExchangeName = "nori"
	DefaultAMQPExchangeKind = "direct"
)

type amqpDeliveryCtxKey struct{}

type AMQPTransport struct {
	URL          string
	ExchangeName string
	ExchangeKind string

	worker *Worker
	tomb   tomb.Tomb

	conn *amqp.Connection
	ch   *amqp.Channel

	didInit bool
}

func NewAMQPTransport(url string) Transport {
	return &AMQPTransport{
		URL:          url,
		ExchangeName: DefaultAMQPExchangeName,
		ExchangeKind: DefaultAMQPExchangeKind,
	}
}

func (t *AMQPTransport) String() string {
	return fmt.Sprintf("AMQPTransport: %s", t.URL)
}

func (t *AMQPTransport) Init(ctx context.Context) error {
	if t.didInit {
		return errors.New("nori.AMQPTransport.Init: already initialized")
	}

	worker, ok := WorkerFromContext(ctx)
	if !ok {
		return errors.New("nori.AMQPTransport.Init: nori.Worker not in context")
	}

	t.worker = worker
	t.didInit = true
	return nil
}

func (t *AMQPTransport) Setup() error {
	conn, err := amqp.Dial(t.URL)
	if err != nil {
		return fmt.Errorf("nori.AMQPTransport.Setup: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("nori.AMQPTransport.Setup: %s", err)
	}

	if t.ExchangeName != "" {
		if err := ch.ExchangeDeclare(
			t.ExchangeName, // name
			t.ExchangeKind, // kind
			true,           // durable
			false,          // autoDelete
			false,          // internal
			false,          // noWait
			nil,            // args
		); err != nil {
			return fmt.Errorf("nori.AMQPTransport.Setup: %s", err)
		}
	}

	t.conn = conn
	t.ch = ch

	return nil
}

func (t *AMQPTransport) Consume(queueName string) (<-chan *Request, error) {
	deliveryChan, err := t.consume(queueName)
	if err != nil {
		return nil, fmt.Errorf("nori.AMQPTransport.Consume: %s", err)
	}

	reqChan := make(chan *Request)
	t.tomb.Go(func() error {
		for {
			select {
			case <-t.tomb.Dying():
				return nil

			case delivery, ok := <-deliveryChan:
				if !ok {
					log.Println("nori.AMQPTransport.Consume: channel closed")
					return nil
				}

				msg, err := t.parseDelivery(delivery)
				if err != nil {
					log.Println("nori.AMQPTransport.Consume: error parsing message:", err)
					delivery.Nack(false, false)
					continue
				}
				reqChan <- msg
			}
		}
	})
	return reqChan, nil
}

func (t *AMQPTransport) consume(queueName string) (<-chan amqp.Delivery, error) {
	q, err := t.ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // autoDelete
		false,     // exclusive
		false,     // noWait
		nil,       // args
	)
	if err != nil {
		return nil, err
	}

	if err := t.ch.QueueBind(
		q.Name,         // name
		q.Name,         // key
		t.ExchangeName, // exchange
		false,          // noWait
		nil,            // args
	); err != nil {
		return nil, err
	}

	msgs, err := t.ch.Consume(
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

func (t *AMQPTransport) Close() error {
	if t.conn != nil {
		return t.conn.Close()
	}
	return nil
}

func (t *AMQPTransport) parseDelivery(d amqp.Delivery) (*Request, error) {
	switch d.ContentType {
	case "application/json":
		var req Request
		if err := json.Unmarshal(d.Body, &req); err != nil {
			return nil, err
		}
		req.Context = context.WithValue(context.Background(), amqpDeliveryCtxKey{}, d)
		return &req, nil

	case "":
		return nil, errors.New("content type is not specified")

	default:
		return nil, fmt.Errorf("unsupported content type: %s", d.ContentType)
	}
}

func amqpDeliveryFromContext(ctx context.Context) (amqp.Delivery, bool) {
	d, ok := ctx.Value(amqpDeliveryCtxKey{}).(amqp.Delivery)
	return d, ok
}

func (t *AMQPTransport) Ack(req *Request) error {
	if d, ok := amqpDeliveryFromContext(req.Context); ok {
		return d.Ack(false)
	}
	return errors.New("amqp delivery not in context")
}

func (t *AMQPTransport) Nack(req *Request, requeue bool) error {
	if d, ok := amqpDeliveryFromContext(req.Context); ok {
		return d.Nack(false, requeue)
	}
	return errors.New("amqp delivery not in context")
}

func (t *AMQPTransport) Reject(req *Request, requeue bool) error {
	if d, ok := amqpDeliveryFromContext(req.Context); ok {
		return d.Reject(requeue)
	}
	return errors.New("amqp delivery not in context")
}
