package main

import (
	"log"

	"golang.org/x/net/context"

	"github.com/jianyuan/nori"
)

type AddRequest struct {
	A int
	B int
}

type AddResponse struct {
	Sum int
}

func main() {
	worker := nori.NewWorker()
	worker.Transport = nori.NewAMQPTransport("amqp://guest:guest@localhost:5672/")
	worker.Tasks = []nori.Task{
		{
			Name: "add",
			Handler: func(ctx context.Context, req AddRequest) (*AddResponse, error) {
				return &AddResponse{Sum: req.A + req.B}, nil
			},
		},
	}

	if err := worker.Run(); err != nil {
		log.Fatal(err)
	}
	if err := worker.Wait(); err != nil {
		log.Fatal(err)
	}
}
