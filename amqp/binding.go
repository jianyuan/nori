package amqp

import "sync"

type Bindable interface {
	Bindable()
}

type Binding struct {
	Destination Bindable
	Exchange    *Exchange
	RoutingKey  string

	mu   sync.RWMutex
	Args map[string]interface{}
}

func NewBinding(
	dest Bindable,
	exchange *Exchange,
	routingKey string,
	args map[string]interface{},
) (*Binding, error) {
	return &Binding{
		Destination: dest,
		Exchange:    exchange,
		RoutingKey:  routingKey,
		Args:        args,
	}, nil
}

func (b *Binding) SetArg(key string, val interface{}) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.Args == nil {
		b.Args = make(map[string]interface{})
	}
	b.Args[key] = val
}
