package amqp

type ConnectionFactory interface {
	URL() string
	Create() (Connection, error)
	Close() error

	NotifyCreate(chan Connection) chan Connection
	NotifyClose(chan Connection) chan Connection
}

type ConnectionListener interface {
	OnCreate(Connection)
	OnClose(Connection)
}

type Connection interface {
	CreateChannel() (Channel, error)
	Close() error
	IsOpen() bool

	NotifyCreateChannel(chan Channel) chan Channel
	NotifyCloseChannel(chan Channel) chan Channel
}

type Channel interface {
}

type ChannelListener interface {
	OnCreate(Channel)
	OnClose(Channel)
}
