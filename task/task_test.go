package task

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProtocolV1_JSONUnmarshal(t *testing.T) {
	jsonPayload := []byte(`{
		"expires":null,
		"utc":true,
		"args":[],
		"chord":null,
		"callbacks":null,
		"errbacks":null,
		"taskset":null,
		"id":"00000000-0000-0000-0000-000000000000",
		"retries":0,
		"task":"tasks.hello_world",
		"timelimit":[null,null],
		"eta":null,
		"kwargs":{}
	}`)

	var actual ProtocolV1
	assert.NoError(t, json.Unmarshal(jsonPayload, &actual))

	expected := ProtocolV1{
		Name:       "tasks.hello_world",
		ID:         "00000000-0000-0000-0000-000000000000",
		Args:       []string{},
		KWArgs:     map[string]interface{}{},
		Retries:    0,
		ETA:        nil,
		Expires:    nil,
		UTC:        true,
		Callbacks:  nil,
		Errbacks:   nil,
		TimeLimits: [2]*float64{nil, nil},
		TaskSet:    nil,
		Chord:      nil,
	}
	assert.Equal(t, expected, actual)
}
