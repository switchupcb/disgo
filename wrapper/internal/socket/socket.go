package socket

import (
	"compress/zlib"
	"context"
	"fmt"

	"github.com/goccy/go-json"

	"github.com/switchupcb/websocket"
)

// Read reads a JSON payload from conn into dst.
//
// Read handles zlib-stream compressed payloads when necessary.
func Read(ctx context.Context, conn *websocket.Conn, dst any) error {
	messageType, reader, err := conn.Reader(ctx)
	if err != nil {
		return err
	}

	// determine the reader based on the message type.
	switch messageType {
	case websocket.MessageText:
		// reuse buffers in between calls to avoid allocations.
		b := get()
		defer put(b)

		// read the message.
		if _, err := b.ReadFrom(reader); err != nil {
			return err
		}

		// unmarshal the message into dst.
		if err = json.Unmarshal(b.Bytes(), &dst); err != nil {
			return err
		}

	case websocket.MessageBinary:
		zlibReader, err := zlib.NewReader(reader)
		if err != nil {
			return err
		}
		defer zlibReader.Close()

		// unmarshal the message into dst.
		decoder := json.NewDecoder(zlibReader)
		if err = decoder.Decode(&dst); err != nil {
			return err
		}

	default:
		return fmt.Errorf("received unknown message type from connection: %v", messageType)
	}

	return nil
}

// Write writes a JSON payload from dst to conn.
func Write(ctx context.Context, conn *websocket.Conn, m websocket.MessageType, dst any) error {
	writer, err := conn.Writer(ctx, m)
	if err != nil {
		return err
	}

	if err := json.NewEncoder(writer).Encode(dst); err != nil {
		return err
	}

	return writer.Close()
}
