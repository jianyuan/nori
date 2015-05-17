package nori

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"golang.org/x/net/context"
)

func TestArgsContext(t *testing.T) {
	ctx := context.Background()

	argMap := make(ArgMap)
	ctx = WithArgs(ctx, argMap)

	{
		read, ok := ArgsFromContext(ctx)
		assert.Len(t, read, 0)
		assert.True(t, ok)
	}

	argPos := 0
	argVal := "arg0"
	argMap[argPos] = argVal

	{
		read, ok := ArgsFromContext(ctx)
		assert.Len(t, read, 1)
		assert.True(t, ok)
		actualVal, present := read[argPos]
		assert.Equal(t, argVal, actualVal)
		assert.True(t, present)
	}

}
