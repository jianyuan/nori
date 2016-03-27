package amqp

import (
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/require"
)

func TestToAMQPTable(t *testing.T) {
	input := map[string]interface{}{
		"string": "a",
		"int":    1,
		"map": map[string]interface{}{
			"string": "b",
			"another_map": map[string]interface{}{
				"string": "c",
			},
		},
	}
	output := amqp.Table{
		"string": "a",
		"int":    1,
		"map": amqp.Table{
			"string": "b",
			"another_map": amqp.Table{
				"string": "c",
			},
		},
	}
	require.Equal(t, output, toAMQPTable(input))
}
