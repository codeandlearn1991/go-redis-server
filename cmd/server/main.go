package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
)

func main() {
	port := flag.String("port", "6379", "Port to start the Redis server on")

	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	if err := listen(logger, *port); err != nil {
		logger.Error("listen error", slog.Any("err", err))
		os.Exit(1)
	}
}

func listen(logger *slog.Logger, port string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	listenerConfig := net.ListenConfig{}
	listener, err := listenerConfig.Listen(ctx, "tcp", "0.0.0.0:"+port)
	if err != nil {
		return fmt.Errorf("creating listener: %w", err)
	}

	defer func() {
		if err := listener.Close(); err != nil {
			logger.Error("failed to close listener", slog.Any("err", err))
		}
	}()

	// Infinite for loop to listen for incoming connections.
	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("accept connection: %w", err)
		}

		// Handle connections in a goroutine, allowing us to accept
		// multiple connections to work with.
		go handleConn(ctx, logger, conn)
	}
}

func handleConn(ctx context.Context, logger *slog.Logger, conn net.Conn) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			buffer := make([]byte, 1024)
			n, err := conn.Read(buffer)
			if errors.Is(err, io.EOF) {
				return
			}

			if err != nil {
				logger.Error(
					"connection from",
					slog.String("addr", conn.RemoteAddr().String()),
					slog.Any("err", err))
				return
			}

			message := string(buffer[:n])
			logger.Info("message from client", slog.String("message", message))

			response := []byte("Hello from the server")
			_, err = conn.Write(response)
			if err != nil {
				logger.Error("responding to client", slog.Any("err", err))
				return
			}
		}
	}
}
