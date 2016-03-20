package nori

import "github.com/jianyuan/nori/message"

type TaskFunc func(message.Request) (message.Response, error)

type Task struct {
	Name     string
	Handler  TaskFunc
	Request  interface{}
	Response interface{}
}
