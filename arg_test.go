package nori

import (
	"encoding/json"
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

func TestArgMap_ImplementsJSONMarshalerAndUnmarshaler(t *testing.T) {
	assert.Implements(t, (*json.Marshaler)(nil), &ArgMap{})
	assert.Implements(t, (*json.Unmarshaler)(nil), &ArgMap{})
}

func TestArgMap_JSONMarshaler(t *testing.T) {
	m := ArgMap{
		0: "arg0",
		2: 2,
		3: 3.01,
	}

	jsonBytes, err := json.Marshal(m)
	assert.NoError(t, err)
	assert.Equal(t, `["arg0",null,2,3.01]`, string(jsonBytes))
}

func TestArgMap_JSONUnmarshaler(t *testing.T) {
	var actualM ArgMap
	jsonBytes := []byte(`["arg0",null,2,3.01]`)
	assert.NoError(t, json.Unmarshal(jsonBytes, &actualM))
	assert.Len(t, actualM, 4)
	assert.EqualValues(t, "arg0", actualM[0])
	assert.EqualValues(t, nil, actualM[1])
	assert.EqualValues(t, 2, actualM[2])
	assert.EqualValues(t, 3.01, actualM[3])
}
