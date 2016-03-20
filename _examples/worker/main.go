package main

import "github.com/jianyuan/nori"

func HelloWorld(req nori.Request) (nori.Response, error) {
	return nil, nil
}

func main() {
	s := nori.NewServer("tasks")

	s.RegisterTask(&nori.Task{
		Name:    "hello_world",
		Handler: HelloWorld,
	})

	go s.RunManagementServer(":8080")
	s.Run()

	<-make(chan bool)
}
