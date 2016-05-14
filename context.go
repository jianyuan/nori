package nori

type Context struct {
	Worker *Worker
}

func newContext(w *Worker) *Context {
	return &Context{
		Worker: w,
	}
}
