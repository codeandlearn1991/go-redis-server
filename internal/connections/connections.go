package connections

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"strings"

	"github.com/codeandlearn1991/go-redis-server/internal/commands"
	"github.com/codeandlearn1991/go-redis-server/internal/resp"
)

type Handler struct {
	logger *slog.Logger
}

func NewHandler(logger *slog.Logger) *Handler {
	return &Handler{logger: logger}
}

func (h *Handler) Handle(ctx context.Context, conn net.Conn) {
	defer func(c net.Conn) {
		if err := c.Close(); err != nil {
			h.logger.Error(
				"Failed to close connection",
				slog.String("addr", c.RemoteAddr().String()),
				slog.String("err", err.Error()),
			)
		}
	}(conn)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			msg, err := resp.Deserilize(conn)
			if err != nil {
				if errors.Is(err, io.EOF) {
					return
				}
				h.logger.Error("Deserilization failed", slog.String("err", err.Error()))
				return
			}

			if msg.Type != resp.Array || len(msg.Array) == 0 {
				errMsg, err := resp.Serialize(resp.NewError("ERR expected array value"))
				if err != nil {
					h.logger.Error("Serialize err message", slog.String("err", err.Error()))
					return
				}

				if _, err := conn.Write([]byte(errMsg)); err != nil {
					h.logger.Error("Writing connection failed", slog.String("err", err.Error()))
					return
				}

				continue
			}

			reply := handleCommand(strings.ToUpper(msg.Array[0].String), msg.Array[0:]...)

			replyMsg, err := resp.Serialize(reply)
			if err != nil {
				h.logger.Error("Serialize err message", slog.String("err", err.Error()))
				return
			}

			if _, err := conn.Write([]byte(replyMsg)); err != nil {
				h.logger.Error("Writing connection failed", slog.String("err", err.Error()))
				return
			}
		}
	}
}

func handleCommand(cmd string, args ...*resp.Value) *resp.Value {
	switch cmd {
	case "PING":
		return commands.Ping(args...)
	case "ECHO":
		return commands.Echo(args...)
	default:
		return resp.NewError("ERR unknown command '" + args[0].String + "'")
	}
}
