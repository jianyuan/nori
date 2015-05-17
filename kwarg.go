package nori

import "golang.org/x/net/context"

type (
	KWArgMap map[string]interface{}

	kwargMapContextKey struct{}
)

func KWArgsFromContext(ctx context.Context) (KWArgMap, bool) {
	val, ok := ctx.Value(kwargMapContextKey{}).(KWArgMap)
	return val, ok
}

func WithKWArgs(parent context.Context, val KWArgMap) context.Context {
	return context.WithValue(parent, kwargMapContextKey{}, val)
}
