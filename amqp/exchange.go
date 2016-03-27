package amqp

import "sync"

type Exchange struct {
	Kind       string
	Name       string
	Durable    bool
	AutoDelete bool

	mu   sync.Mutex
	Args map[string]interface{}
}

func NewExchange(
	kind string,
	name string,
	durable bool,
	autoDelete bool,
	args map[string]interface{},
) (*Exchange, error) {
	return &Exchange{
		Kind:       kind,
		Name:       name,
		Durable:    durable,
		AutoDelete: autoDelete,
		Args:       args,
	}, nil
}

func (*Exchange) Bindable() {}

func (e *Exchange) SetArg(key string, val interface{}) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.Args == nil {
		e.Args = make(map[string]interface{})
	}
	e.Args[key] = val
}
