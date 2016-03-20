package main

import "github.com/jianyuan/nori"

func Ping(req nori.Request) (nori.Response, error) {
	resp := req.NewResponse()
	resp.SetBody("Pong!")
	return resp, nil
}

func Add(req nori.Request) (nori.Response, error) {
	resp := req.NewResponse()
	resp.SetBody(req.MustArg(0).(int) + req.MustArg(1).(int))
	return resp, nil
}

func main() {
	s := nori.NewServer("tasks")

	s.RegisterTask(&nori.Task{
		Name:    "ping",
		Handler: Ping,
	})
	s.RegisterTask(&nori.Task{
		Name:    "add",
		Handler: Add,
	})

	go s.RunManagementServer(":8080")
	s.Run()

	<-make(chan bool)
}
