package nori

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"golang.org/x/net/context"
)

func TestKWArgsContext(t *testing.T) {
	ctx := context.Background()

	kwargMap := make(KWArgMap)
	ctx = WithKWArgs(ctx, kwargMap)

	{
		read, ok := KWArgsFromContext(ctx)
		assert.Len(t, read, 0)
		assert.True(t, ok)
	}

	kwargKey := "argKey"
	kwargVal := "argVal"
	kwargMap[kwargKey] = kwargVal

	{
		read, ok := KWArgsFromContext(ctx)
		assert.Len(t, read, 1)
		assert.True(t, ok)
		actualVal, present := read[kwargKey]
		assert.Equal(t, kwargVal, actualVal)
		assert.True(t, present)
	}
}
