package transport

import (
	"github.com/jianyuan/nori/message"
	"gopkg.in/tomb.v2"
)

type Driver interface {
	Name() string
	Tomb() *tomb.Tomb
	Setup() error
	Close() error
	Consume(string) (<-chan *message.Request, error)
	Reply(*message.Request, message.Response) error
}
