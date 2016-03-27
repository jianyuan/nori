package amqp

import (
	"errors"
	"fmt"

	"github.com/streadway/amqp"
)

type amqpAdmin struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func toAMQPTable(src map[string]interface{}) amqp.Table {
	if src == nil {
		return nil
	}
	dest := make(amqp.Table, len(src))
	for key, val := range src {
		if table, ok := val.(map[string]interface{}); ok {
			val = toAMQPTable(table)
		}
		dest[key] = val
	}
	return dest
}

func (a *amqpAdmin) maybeOpen() error {
	if a.ch == nil {
		ch, err := a.conn.Channel()
		if err != nil {
			return err
		}
		a.ch = ch
	}
	return nil
}

func (a *amqpAdmin) DeclareExchange(e *Exchange) error {
	if e == nil {
		return errors.New("amqp: Exchange is nil")
	}
	if err := a.maybeOpen(); err != nil {
		return err
	}
	return a.ch.ExchangeDeclare(
		e.Name,              // name string
		e.Kind,              // kind string
		e.Durable,           // durable bool
		e.AutoDelete,        // autoDelete bool
		false,               // internal bool
		false,               //	noWait bool
		toAMQPTable(e.Args), //	args amqp.Table
	)
}

func (a *amqpAdmin) DeleteExchange(name string) error {
	if err := a.maybeOpen(); err != nil {
		return err
	}
	return a.ch.ExchangeDelete(
		name,  // name string
		false, // ifUnused bool
		false, // noWait bool
	)
}

func (a *amqpAdmin) DeclareQueue(q *Queue) error {
	if q == nil {
		return errors.New("amqp: Queue is nil")
	}
	if err := a.maybeOpen(); err != nil {
		return err
	}
	_, err := a.ch.QueueDeclare(
		q.Name,              // name string
		q.Durable,           // durable bool
		q.AutoDelete,        // autoDelete bool
		q.Exclusive,         // exclusive bool
		false,               // noWait bool
		toAMQPTable(q.Args), // args amqp.Table
	)
	return err
}

func (a *amqpAdmin) DeclareAnonymousQueue() (*Queue, error) {
	if err := a.maybeOpen(); err != nil {
		return nil, err
	}
	q, err := a.ch.QueueDeclare(
		"",    // name string
		false, // durable bool
		true,  // autoDelete bool
		true,  // exclusive bool
		false, // noWait bool
		nil,   // args amqp.Table
	)
	if err != nil {
		return nil, err
	}
	return NewQueue(
		q.Name, // name string
		false,  // durable bool
		true,   // exclusive bool
		true,   // autoDelete bool
		nil,    // args map[string]interface{}
	)
}

func (a *amqpAdmin) DeleteQueue(name string, ifUnused, ifEmpty bool) error {
	if err := a.maybeOpen(); err != nil {
		return err
	}
	_, err := a.ch.QueueDelete(
		name,     // name string
		ifUnused, // ifUnused bool
		ifEmpty,  // ifEmpty bool
		false,    // noWait bool
	)
	return err
}

func (a *amqpAdmin) PurgeQueue(name string, noWait bool) error {
	if err := a.maybeOpen(); err != nil {
		return err
	}
	_, err := a.ch.QueuePurge(
		name,   // name string
		noWait, // noWait bool
	)
	return err
}

func (a *amqpAdmin) DeclareBinding(b *Binding) error {
	if b == nil {
		return errors.New("amqp: Binding is nil")
	}
	if err := a.maybeOpen(); err != nil {
		return err
	}

	switch dest := b.Destination.(type) {
	case nil:
		return errors.New("amqp: Binding destination is nil")
	case *Exchange:
		return a.ch.ExchangeBind(
			dest.Name,           // destination string
			b.RoutingKey,        // key string
			b.Exchange.Name,     // source string
			false,               // noWait bool
			toAMQPTable(b.Args), // args amqp.Table
		)
	case *Queue:
		return a.ch.QueueBind(
			dest.Name,           // name string
			b.RoutingKey,        // key string
			b.Exchange.Name,     // exchange string
			false,               // noWait bool
			toAMQPTable(b.Args), // args amqp.Table
		)
	default:
		return fmt.Errorf("amqp: Unsupported binding destination %T", dest)
	}
}

func (a *amqpAdmin) RemoveBinding(b *Binding) error {
	if b == nil {
		return errors.New("amqp: Binding is nil")
	}
	if err := a.maybeOpen(); err != nil {
		return err
	}

	switch dest := b.Destination.(type) {
	case nil:
		return errors.New("amqp: Binding destination is nil")
	case *Exchange:
		return a.ch.ExchangeUnbind(
			dest.Name,           // destination string
			b.RoutingKey,        // key string
			b.Exchange.Name,     // source string
			false,               // noWait bool
			toAMQPTable(b.Args), // args amqp.Table
		)
	case *Queue:
		return a.ch.QueueUnbind(
			dest.Name,           // name string
			b.RoutingKey,        // key string
			b.Exchange.Name,     // exchange string
			toAMQPTable(b.Args), // args amqp.Table
		)
	default:
		return fmt.Errorf("amqp: Unsupported binding destination %T", dest)
	}
}

func (a *amqpAdmin) Close() error {
	if err := a.ch.Close(); err != nil {
		return err
	}
	a.ch = nil
	return nil
}

func NewAMQPAdmin(conn *amqp.Connection) (Admin, error) {
	return &amqpAdmin{
		conn: conn,
	}, nil
}
