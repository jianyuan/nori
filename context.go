package nori

import "golang.org/x/net/context"

type workerContextKey struct{}

func NewWorkerContext(ctx context.Context, w *Worker) context.Context {
	return context.WithValue(ctx, workerContextKey{}, w)
}

func WorkerFromContext(ctx context.Context) (*Worker, bool) {
	w, ok := ctx.Value(workerContextKey{}).(*Worker)
	return w, ok
}
