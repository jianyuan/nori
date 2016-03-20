package message

import (
	"encoding/json"
	"time"

	"golang.org/x/net/context"
)

type Request struct {
	Ctx         context.Context
	HandlerName string
	ID          string
	Args        []interface{}
	KWArgs      map[string]interface{}
	ETA         *time.Time
	ExpiresAt   *time.Time
	IsUTC       bool
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
		ID: req.ID,
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
	SetStatus(string)
	SetBody(interface{}) error
}

type response struct {
	ID     string
	Status string
	Body   []byte
}

func (resp *response) SetID(id string) {
	resp.ID = id
}

func (resp *response) SetStatus(status string) {
	resp.Status = status
}

func (resp *response) SetBody(body interface{}) error {
	// TODO
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	resp.Body = b
	return nil
}
