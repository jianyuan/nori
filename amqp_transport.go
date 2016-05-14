package nori

type AMQPTransport struct{}

func NewAMQPTransport() Transport {
	return &AMQPTransport{}
}

func (at *AMQPTransport) Init(ctx *Context) error {
	return nil
}
