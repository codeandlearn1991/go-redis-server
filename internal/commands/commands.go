package commands

import "github.com/codeandlearn1991/go-redis-server/internal/resp"

func Ping(args ...*resp.Value) *resp.Value {
	if len(args) > 2 {
		return resp.NewError("ERR wrong number of arguments for 'ping' command")
	}

	if len(args) == 2 {
		return resp.NewBulkString(args[1].String)
	}

	return resp.NewSimpleString("PONG")
}

func Echo(args ...*resp.Value) *resp.Value {
	if len(args) != 2 {
		return resp.NewError("ERR wrong number of arguments for 'echo' command")
	}
	return resp.NewBulkString(args[1].String)
}
