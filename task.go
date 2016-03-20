package nori

import "github.com/jianyuan/nori/message"

type TaskHandlerFunc func(*message.Request) (message.Response, error)

type Task struct {
	Name     string
	Handler  TaskHandlerFunc
	Request  interface{}
	Response interface{}
}
