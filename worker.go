package nori

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/net/context"
	"gopkg.in/tomb.v2"
)

type Worker struct {
	Name      string
	Transport Transport
	Tasks     []Task

	tomb tomb.Tomb
}

// NewWorker creates a new Worker with sane defaults
func NewWorker() *Worker {
	return &Worker{
		Name: "nori",
	}
}

func (w *Worker) Run() error {
	if w.Transport == nil {
		return errors.New("nori.Worker: no transport provided")
	}

	for _, t := range w.Tasks {
		if err := t.init(); err != nil {
			return err
		}
	}

	ctx := NewWorkerContext(context.Background(), w)
	if err := w.Transport.Init(ctx); err != nil {
		return fmt.Errorf("nori.Worker.Run: transport initialization failed: %s", err)
	}

	w.printInfo()

	w.tomb.Go(w.run)

	return nil
}

func (w *Worker) run() error {
	// TODO better retry mechanism
	for {
		select {
		case err := <-w.setupTransport():
			if err != nil {
				log.Println("nori.Worker.Run: transport setup error:", err)
				time.Sleep(5 * time.Second)
				continue
			}

			log.Println("Transport connected")

			w.consume()

		case <-w.tomb.Dying():
			log.Println("Cancelled")
			break

		case <-time.After(5 * time.Second):
			log.Println("nori.Worker.Run: transport setup error: timed out")
			time.Sleep(5 * time.Second)
		}
	}

	return w.Transport.Close()
}

func (w *Worker) setupTransport() <-chan error {
	errChan := make(chan error)
	w.tomb.Go(func() error {
		errChan <- w.Transport.Setup()
		return nil
	})
	return errChan
}

func (w *Worker) consume() {
	reqChan, err := w.Transport.Consume(w.Name)
	if err != nil {
		log.Println("nori.Worker.Run: consume error:", err)
		return
	}

	log.Printf("Consuming messages from queue %q", w.Name)

	for {
		select {
		case req := <-reqChan:
			log.Println("Got request:", req)
			log.Println(w.Transport.Ack(req))
			// TODO actually call the task

		case <-w.tomb.Dying():
			return
		}

	}
}

func (w *Worker) Wait() error {
	return w.tomb.Wait()
}

func (w *Worker) Stop() {
	w.tomb.Kill(nil)
}

func (w *Worker) printInfo() {
	if hostname, err := os.Hostname(); err == nil {
		log.Println("Hostname:", hostname)
	}

	log.Println("Registered tasks:")
	for _, t := range w.Tasks {
		log.Println("  -", t.Name)
	}

	log.Println("Transport:", w.Transport)
}
