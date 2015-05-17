package main

import (
	"github.com/jianyuan/nori/server"
	"github.com/jianyuan/nori/task"
	"golang.org/x/net/context"
)

func HelloWorld(ctx context.Context) error {
	return nil
}

func main() {
	s := server.New("tasks")
	s.RegisterTask(&task.Task{
		Func: HelloWorld,
		Name: "hello_world",
	})

	go s.RunManagementServer(":8080")
	s.Run()

	<-make(chan bool)
}
