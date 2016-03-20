package nori

import (
	"time"

	"golang.org/x/net/context"
)

type TaskFunc func(Request) (Response, error)

type Request interface {
	Context() context.Context

	GetArg(int) (interface{}, error)
	MustArg(int) interface{}
	Args() []interface{}

	GetKWArg(string) (interface{}, error)
	MustKWArg(string) interface{}
	KWArgs() map[string]interface{}

	ID() string
	ETA() *time.Time
	Expires() *time.Time
	IsUTC() bool

	NewResponse() Response
}

type Response interface {
	SetID(string)
	SetETA(time.Time)
	SetExpires(time.Time)
	SetIsUTC() bool

	SetBody(interface{})
}

type Task struct {
	Name     string
	Handler  TaskFunc
	Request  interface{}
	Response interface{}
}
