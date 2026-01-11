package commands_test

import (
	"testing"

	"github.com/codeandlearn1991/go-redis-server/internal/commands"
	"github.com/codeandlearn1991/go-redis-server/internal/resp"
	"github.com/stretchr/testify/assert"
)

func Test_Ping(t *testing.T) {
	testCases := []struct {
		name     string
		args     []*resp.Value
		expected *resp.Value
	}{
		{
			name: "args len greater than 2",
			args: []*resp.Value{
				{
					Type:   resp.SimpleString,
					String: "Value 1",
				},
				{
					Type:   resp.SimpleString,
					String: "Value 2",
				},
				{
					Type:   resp.SimpleString,
					String: "Value 3",
				},
			},
			expected: resp.NewError("ERR wrong number of arguments for 'ping' command"),
		},
		{
			name: "args len equal to 2",
			args: []*resp.Value{
				{
					Type:   resp.SimpleString,
					String: "echo",
				},
				{
					Type:   resp.SimpleString,
					String: "hello",
				},
			},
			expected: resp.NewBulkString("hello"),
		},
		{
			name: "pong response",
			args: []*resp.Value{
				{
					Type:   resp.SimpleString,
					String: "echo",
				},
			},
			expected: resp.NewSimpleString("PONG"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := commands.Ping(tc.args...)

			assert.Equal(t, tc.expected, got)
		})
	}
}

func Test_Echo(t *testing.T) {
	testCases := []struct {
		name     string
		args     []*resp.Value
		expected *resp.Value
	}{
		{
			name:     "args len not 2",
			args:     []*resp.Value{},
			expected: resp.NewError("ERR wrong number of arguments for 'echo' command"),
		},
		{
			name: "success",
			args: []*resp.Value{
				{
					Type:   resp.SimpleString,
					String: "echo",
				},
				{
					Type:   resp.SimpleString,
					String: "hello",
				},
			},
			expected: resp.NewBulkString("hello"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := commands.Echo(tc.args...)

			assert.Equal(t, tc.expected, got)
		})
	}
}
