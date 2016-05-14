package nori

import (
	"errors"
	"reflect"
)

// Task describes a unit of work
type Task struct {
	// Name of task
	Name string

	// Handler to invoke when the task is executed
	Handler interface{}

	reqType  reflect.Type
	respType reflect.Type

	didInit bool
}

func (t *Task) init() error {
	if t.didInit {
		return nil
	}

	if t.Name == "" {
		return errors.New("nori.Task: name is blank")
	}

	if t.Handler == nil {
		return errors.New("nori.Task: handler is nil")
	}

	hv := reflect.ValueOf(t.Handler)
	if hv.Kind() != reflect.Func {
		return errors.New("nori.Task: handler is not a function")
	}

	ht := hv.Type()

	// Determine request type
	// TODO: advanced validation
	ct := reflect.TypeOf((*Context)(nil))
	for i := 0; i < ht.NumIn(); i++ {
		it := ht.In(i)
		if it != ct {
			t.reqType = it
			break
		}
	}

	// Determine response type
	// TODO: advanced validation
	et := reflect.TypeOf((*error)(nil)).Elem()
	for i := 0; i < ht.NumOut(); i++ {
		ot := ht.Out(i)
		if ot != et {
			t.respType = ot
			break
		}
	}

	t.didInit = true
	return nil
}
