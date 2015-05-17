package nori

import (
	"encoding/json"

	"golang.org/x/net/context"
)

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

func (m ArgMap) MarshalJSON() ([]byte, error) {
	maxKey := 0
	for i, _ := range m {
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
