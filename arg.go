package nori

import "golang.org/x/net/context"

type (
	ArgMap map[int]interface{}

	argMapContextKey struct{}
)

func ArgsFromContext(ctx context.Context) (ArgMap, bool) {
	val, ok := ctx.Value(argMapContextKey{}).(ArgMap)
	return val, ok
}

func WithArgs(parent context.Context, val ArgMap) context.Context {
	return context.WithValue(parent, argMapContextKey{}, val)
}
