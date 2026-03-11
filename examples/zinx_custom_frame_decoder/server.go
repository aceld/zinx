/**
 * Custom frame decoder example
 *
 * This example demonstrates how to use custom frame decoder for handling
 * packet splitting without length field (e.g., using '\r' as delimiter)
 *
 * Usage:
 * Run the server, then use telnet to test:
 *   telnet 127.0.0.1 7777
 *   Type some data ending with \r (press Enter in telnet sends \r\n)
 */
package main

import (
	"fmt"

	"github.com/aceld/zinx/ziface"
	"github.com/aceld/zinx/zlog"
	"github.com/aceld/zinx/znet"
)

// CustomFrameDecoder implements IFrameDecoder interface
// Splits packets by '\r' delimiter
type CustomFrameDecoder struct{}

// Decode splits the data by '\r' delimiter
func (d *CustomFrameDecoder) Decode(buff []byte) [][]byte {
	var result [][]byte
	start := 0

	for i := 0; i < len(buff); i++ {
		if buff[i] == '\r' {
			// Found delimiter, create a new frame
			if i > start {
				frame := make([]byte, i-start)
				copy(frame, buff[start:i])
				result = append(result, frame)
			}
			start = i + 1
		}
	}

	// Handle remaining data (without delimiter at the end)
	if start < len(buff) {
		frame := make([]byte, len(buff)-start)
		copy(frame, buff[start:])
		result = append(result, frame)
	}

	return result
}

func main() {
	// Create a custom frame decoder
	customDecoder := &CustomFrameDecoder{}

	// Create server with custom frame decoder
	server := znet.NewServer("CustomFrameDecoderServer")
	server.SetFrameDecoder(customDecoder)

	// Add route
	server.AddRouter(1, &Router{})

	zlog.Ins().InfoF("Starting server with custom frame decoder on :7777...")
	server.Serve()
}

// Router handle the request
type Router struct {
	znet.BaseRouter
}

func (r *Router) Handle(request ziface.IRequest) {
	zlog.Ins().Infof("Received data: %s", string(request.GetData()))
	fmt.Printf("Received data: %s\n", string(request.GetData()))
}
