package resp_test

import (
	"strings"
	"testing"

	"github.com/codeandlearn1991/go-redis-server/internal/resp"
	"github.com/stretchr/testify/assert"
)

func TestDeserialize(t *testing.T) {
	decodingTestCases := []struct {
		name        string
		input       string
		expectedErr string
		expectedVal *resp.Value
	}{
		{
			name:  "deserialize array",
			input: "*1\r\n:123\r\n",
			expectedVal: &resp.Value{
				Type:  resp.Array,
				Array: []*resp.Value{{Type: resp.Integer, Integer: 123}},
			},
		},
		{
			name:  "deserialize null array",
			input: "*-1\r\n",
			expectedVal: &resp.Value{
				Type:   resp.Array,
				IsNull: true,
			},
		},
		{
			name:  "deserialize simple string",
			input: "+OK\r\n",
			expectedVal: &resp.Value{
				Type:   resp.SimpleString,
				String: "OK",
			},
		},
		{
			name:  "deserialize error",
			input: "-ERR invalid command\r\n",
			expectedVal: &resp.Value{
				Type:   resp.Error,
				String: "ERR invalid command",
			},
		},
		{
			name:  "deserialize integer",
			input: ":123\r\n",
			expectedVal: &resp.Value{
				Type:    resp.Integer,
				Integer: 123,
			},
		},
		{
			name:  "deserialize bulk string",
			input: "$5\r\nhello\r\n",
			expectedVal: &resp.Value{
				Type:   resp.BulkString,
				String: "hello",
			},
		},
		{
			name:  "deserialize null bulk string",
			input: "$-1\r\n",
			expectedVal: &resp.Value{
				Type:   resp.BulkString,
				IsNull: true,
			},
		},
		{
			name:        "no first byte",
			input:       "",
			expectedErr: "read resp first byte",
		},
		{
			name:        "error read line bytes",
			input:       "*invalid",
			expectedErr: "read line bytes",
		},
		{
			name:        "error line not terminated",
			input:       "*invalid\n",
			expectedErr: "line not terminated with",
		},
		{
			name:        "error parse array num",
			input:       "*x\r\n",
			expectedErr: "parse the num elements",
		},

		{
			name:        "error decoding arr el",
			input:       "*1\r\n:x",
			expectedErr: "deserialize array element",
		},
		{
			name:        "error bulk string len",
			input:       "$x\r\nhello\r\n",
			expectedErr: "parse string len",
		},
		{
			name:        "error read bulk string",
			input:       "$5\r\n\r\n",
			expectedErr: "read bulk string",
		},
		{
			name:        "error bulk string not terminated",
			input:       "$5\r\nhell\r\n",
			expectedErr: "bulk string not terminated correctly",
		},
	}

	for _, tc := range decodingTestCases {
		val, err := resp.Deserilize(strings.NewReader(tc.input))

		if tc.expectedErr != "" {
			assert.Error(t, err)
			assert.ErrorContains(t, err, tc.expectedErr)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedVal, val)
		}
	}
}
