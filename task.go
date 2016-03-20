package nori

import (
	"time"

	"golang.org/x/net/context"
)

type ProtocolV1 struct {
	Name       string      `json:"task"`
	ID         string      `json:"id"`
	Args       ArgMap      `json:"args,omitempty"`
	KWArgs     KWArgMap    `json:"kwargs,omitempty"`
	Retries    int         `json:"retries,omitempty"`
	ETA        *time.Time  `json:"eta,omitempty"`
	Expires    *time.Time  `json:"expires,omitempty"`
	UTC        bool        `json:"utc,omitempty"`
	Callbacks  []string    `json:"callbacks,omitempty"`
	Errbacks   []string    `json:"errbacks,omitempty"`
	TimeLimits [2]*float64 `json:"timelimit,omitempty"`
	TaskSet    *string     `json:"taskset,omitempty"`
	Chord      *string     `json:"chord,omitempty"`
}

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
