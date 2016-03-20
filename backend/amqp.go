package backend

type AMQPBackend struct {
}

func (AMQPBackend) Name() string { return "AMQPBackend" }

var _ Driver = (*AMQPBackend)(nil)
