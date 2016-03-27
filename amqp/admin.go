package amqp

type Admin interface {
	DeclareExchange(*Exchange) error
	DeleteExchange(string) error

	DeclareQueue(*Queue) error
	DeclareAnonymousQueue() (*Queue, error)
	DeleteQueue(name string, ifUnused, ifEmpty bool) error
	PurgeQueue(name string, noWait bool) error

	DeclareBinding(*Binding) error
	RemoveBinding(*Binding) error

	Close() error
}
