package nori

import "golang.org/x/net/context"

type (
	// KWArgMap represents the keyword arguments of a task
	KWArgMap map[string]interface{}

	// kwargMapContextKey holds the context key of a KWArgMap
	kwargMapContextKey struct{}
)

// KWArgsFromContext extracts the KWArgMap from context
func KWArgsFromContext(ctx context.Context) (KWArgMap, bool) {
	val, ok := ctx.Value(kwargMapContextKey{}).(KWArgMap)
	return val, ok
}

// WithKWArgs returns a new context with the KWArgMap
func WithKWArgs(parent context.Context, val KWArgMap) context.Context {
	return context.WithValue(parent, kwargMapContextKey{}, val)
}
