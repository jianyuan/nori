package main

import (
	"log"

	"golang.org/x/net/context"

	"github.com/jianyuan/nori"
	"github.com/jianyuan/nori/message"
	"github.com/jianyuan/nori/transport"
)

func Ping(req *message.Request) (message.Response, error) {
	resp := req.NewResponse()
	resp.SetBody("Pong!")
	return resp, nil
}

func Add(req *message.Request) (message.Response, error) {
	resp := req.NewResponse()
	a := int(req.MustArg(0).(float64))
	b := int(req.MustArg(1).(float64))
	resp.SetBody(a + b)
	return resp, nil
}

func main() {
	ctx := context.Background()

	s, err := nori.NewServer(ctx, &nori.Configuration{
		Name:      "tasks",
		Transport: transport.NewAMQPTransport("amqp://guest:guest@localhost:5672//"),
	})
	if err != nil {
		log.Panicln("Server configuration error:", err)
	}

	s.RegisterTask(&nori.Task{
		Name:    "ping",
		Handler: Ping,
	})
	s.RegisterTask(&nori.Task{
		Name:    "add",
		Handler: Add,
	})

	if err := s.Run(); err != nil {
		log.Panicln("Can't start worker:", err)
	}

	// sigs := make(chan os.Signal, 1)
	// signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// go func() {
	// 	<-sigs
	// 	s.Stop()
	// }()

	if err := s.Wait(); err != nil {
		log.Panicln("Worker terminated prematurely:", err)
	}
}
