package main

import (
	"log"

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
	worker.Transport = nori.NewAMQPTransport()
	worker.Tasks = []nori.Task{
		{
			Name: "add",
			Handler: func(ctx *nori.Context, req AddRequest) (*AddResponse, error) {
				return &AddResponse{Sum: req.A + req.B}, nil
			},
		},
	}

	if err := worker.Run(); err != nil {
		log.Fatal(err)
	}
}
