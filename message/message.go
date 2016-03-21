package message

import (
	"time"

	"golang.org/x/net/context"
)

type Request struct {
	Ctx       context.Context
	TaskName  string
	ID        string
	Args      []interface{}
	KWArgs    map[string]interface{}
	ETA       *time.Time
	ExpiresAt *time.Time
	IsUTC     bool

	ReplyTo *string
	// TODO other celery fields
}

func (req *Request) MustArg(pos int) interface{} {
	if pos >= len(req.Args) {
		panic("Arg missing")
	}
	return req.Args[pos]
}

func (req *Request) MustKWArg(key string) interface{} {
	if val, ok := req.KWArgs[key]; ok {
		return val
	}
	panic("KWArg missing")
}

func (req *Request) NewResponse() Response {
	return &response{
		ID:      req.ID,
		Status:  Success,
		ReplyTo: req.ReplyTo,
	}
}

func NewRequest() *Request {
	return &Request{
		Ctx:    context.Background(),
		KWArgs: make(map[string]interface{}),
	}
}

type Response interface {
	SetID(string)
	SetStatus(State)
	SetBody(interface{}) error
	GetID() string
	GetStatus() State
	GetBody() interface{}
	GetReplyTo() *string
}

type response struct {
	ID     string
	Status State
	Body   interface{}

	ReplyTo *string
}

func (resp *response) SetID(id string) {
	resp.ID = id
}

func (resp *response) SetStatus(status State) {
	resp.Status = status
}

func (resp *response) SetBody(body interface{}) error {
	resp.Body = body
	return nil
}

func (resp *response) GetID() string {
	return resp.ID
}

func (resp *response) GetStatus() State {
	return resp.Status
}

func (resp *response) GetBody() interface{} {
	return resp.Body
}

func (resp *response) GetReplyTo() *string {
	return resp.ReplyTo
}
