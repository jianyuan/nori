package amqp

import (
	"errors"
	"fmt"
	"sync"

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

// Connection

type amqpConnection struct {
	listener ConnectionListener
	conn     *amqp.Connection

	muNotify         sync.Mutex
	createChannelChs []chan<- Channel
	closeChannelChs  []chan<- Channel
}

func (c *amqpConnection) CreateChannel() (Channel, error) {
	ch, err := c.conn.Channel()
	if err != nil {
		return nil, err
	}
	wrappedCh := newAMQPChannel(ch, c)

	c.OnCreate(wrappedCh)

	return wrappedCh, nil
}

func (c *amqpConnection) Close() error {
	if err := c.conn.Close(); err != nil {
		return err
	}

	if c.listener != nil {
		c.listener.OnClose(c)
	}

	return nil
}

func (c *amqpConnection) IsOpen() bool {
	return c.conn != nil
}

func (c *amqpConnection) NotifyCreateChannel(ch chan Channel) chan Channel {
	c.muNotify.Lock()
	defer c.muNotify.Unlock()
	c.createChannelChs = append(c.createChannelChs, ch)
	return ch
}

func (c *amqpConnection) NotifyCloseChannel(ch chan Channel) chan Channel {
	c.muNotify.Lock()
	defer c.muNotify.Unlock()
	c.closeChannelChs = append(c.closeChannelChs, ch)
	return ch
}

func (c *amqpConnection) OnCreate(channel Channel) {
	c.muNotify.Lock()
	defer c.muNotify.Unlock()
	for _, ch := range c.createChannelChs {
		ch <- channel
	}
}

func (c *amqpConnection) OnClose(channel Channel) {
	c.muNotify.Lock()
	defer c.muNotify.Unlock()
	for _, ch := range c.closeChannelChs {
		ch <- channel
	}
}

func newAMQPConnection(conn *amqp.Connection, listener ConnectionListener) Connection {
	return &amqpConnection{
		conn:     conn,
		listener: listener,
	}
}

var _ Connection = (*amqpConnection)(nil)
var _ ChannelListener = (*amqpConnection)(nil)

// Channel

type amqpChannel struct {
	listener ChannelListener
	ch       *amqp.Channel
}

func newAMQPChannel(ch *amqp.Channel, listener ChannelListener) Channel {
	return &amqpChannel{
		ch:       ch,
		listener: listener,
	}
}

var _ Channel = (*amqpChannel)(nil)

// Connection factory

type singleAMQPConnectionFactory struct {
	url string

	mu   sync.Mutex
	conn Connection

	muNotify  sync.Mutex
	createChs []chan<- Connection
	closeChs  []chan<- Connection
}

func (s *singleAMQPConnectionFactory) URL() string {
	return s.url
}

func (s *singleAMQPConnectionFactory) Create() (Connection, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.conn == nil {
		conn, err := amqp.Dial(s.url)
		if err != nil {
			return nil, err
		}
		s.conn = newAMQPConnection(conn, s)

		s.OnCreate(s.conn)
	}

	return s.conn, nil
}

func (s *singleAMQPConnectionFactory) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			return err
		}
		s.conn = nil
	}

	return nil
}

func (s *singleAMQPConnectionFactory) NotifyCreate(ch chan Connection) chan Connection {
	s.muNotify.Lock()
	defer s.muNotify.Unlock()
	s.createChs = append(s.createChs, ch)
	return ch
}

func (s *singleAMQPConnectionFactory) NotifyClose(ch chan Connection) chan Connection {
	s.muNotify.Lock()
	defer s.muNotify.Unlock()
	s.closeChs = append(s.closeChs, ch)
	return ch
}

func (s *singleAMQPConnectionFactory) OnCreate(conn Connection) {
	s.muNotify.Lock()
	defer s.muNotify.Unlock()
	for _, ch := range s.createChs {
		ch <- s.conn
	}
}
func (s *singleAMQPConnectionFactory) OnClose(conn Connection) {
	s.muNotify.Lock()
	defer s.muNotify.Unlock()
	for _, ch := range s.closeChs {
		ch <- s.conn
	}
}

func NewSingleAMQPConnectionFactory(url string) (ConnectionFactory, error) {
	return &singleAMQPConnectionFactory{
		url: url,
	}, nil
}

var _ ConnectionFactory = (*singleAMQPConnectionFactory)(nil)
var _ ConnectionListener = (*singleAMQPConnectionFactory)(nil)
