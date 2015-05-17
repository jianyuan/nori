package nori

import (
	"encoding/json"

	"golang.org/x/net/context"
)

type (
	// ArgMap represents the positional arguments of a task
	ArgMap map[int]interface{}

	// argMapContextKey holds the context key of an ArgMap
	argMapContextKey struct{}
)

// ArgsFromContext extracts the ArgMap from context
func ArgsFromContext(ctx context.Context) (ArgMap, bool) {
	val, ok := ctx.Value(argMapContextKey{}).(ArgMap)
	return val, ok
}

// WithArgs returns a new context with the ArgMap
func WithArgs(parent context.Context, val ArgMap) context.Context {
	return context.WithValue(parent, argMapContextKey{}, val)
}

// MarshalJSON marshals ArgMap into a JSON array
func (m ArgMap) MarshalJSON() ([]byte, error) {
	maxKey := 0
	for i := range m {
		if i > maxKey {
			maxKey = i
		}
	}

	argList := make([]interface{}, maxKey+1)
	for i, val := range m {
		argList[i] = val
	}

	return json.Marshal(argList)
}

// UnmarshalJSON unmarshals JSON array into ArgMap
func (m *ArgMap) UnmarshalJSON(b []byte) error {
	var argList []interface{}
	if err := json.Unmarshal(b, &argList); err != nil {
		return err
	}

	// Check if map is initialized
	if *m == nil {
		*m = make(ArgMap, 0)
	}

	for i, val := range argList {
		(*m)[i] = val
	}

	return nil
}
