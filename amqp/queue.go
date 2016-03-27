package amqp

import "sync"

type Queue struct {
	Name       string
	Durable    bool
	Exclusive  bool
	AutoDelete bool

	mu   sync.RWMutex
	Args map[string]interface{}
}

func NewQueue(
	name string,
	durable bool,
	exclusive bool,
	autoDelete bool,
	args map[string]interface{},
) (*Queue, error) {
	return &Queue{
		Name:       name,
		Durable:    durable,
		Exclusive:  exclusive,
		AutoDelete: autoDelete,
		Args:       args,
	}, nil
}

func (*Queue) Bindable() {}

func (q *Queue) SetArg(key string, val interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if q.Args == nil {
		q.Args = make(map[string]interface{})
	}
	q.Args[key] = val
}
