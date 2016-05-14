package nori

import "golang.org/x/net/context"

type Transport interface {
	String() string
	Init(context.Context) error
	Setup() error
	Close() error
	Consume(queueName string) (<-chan *Request, error)
	Ack(req *Request) error
	Nack(req *Request, requeue bool) error
	Reject(req *Request, requeue bool) error
}
