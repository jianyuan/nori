package nori

import (
	"errors"
	"fmt"
)

type Worker struct {
	Name      string
	Transport Transport
	Tasks     []Task
}

// NewWorker creates a new Worker with sane defaults
func NewWorker() *Worker {
	return &Worker{}
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

	ctx := newContext(w)
	if err := w.Transport.Init(ctx); err != nil {
		return fmt.Errorf("nori.Worker: transport initialization failed: %s", err)
	}

	return nil
}
