package zpack

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/aceld/zinx/ziface"
)

// run in terminal:
// go test -v ./znet -run=TestDataPack

// This function is responsible for testing the functionality of data packet splitting and packaging.
func TestDataPack(t *testing.T) {
	// Create a TCP server socket.
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err:", err)
		return
	}

	// Create a server goroutine, responsible for reading and parsing the data from the client goroutine that may contain sticky packets.
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept err:", err)
			}

			// Handle client requests
			go func(conn net.Conn) {
				// Create a packet splitting and packaging object dp.
				dp := Factory().NewPack(ziface.ZinxDataPack)
				for {
					// 1. Read the head part of the stream first.
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData) // ReadFull will fill msg until it's full
					if err != nil {
						fmt.Println("read head error")
					}
					// Unpack the headData byte stream into msg.
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack err:", err)
						return
					}

					if msgHead.GetDataLen() > 0 {
						// msg has data, read data again.
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetDataLen())

						// Read the byte stream from io based on dataLen.
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err:", err)
							return
						}

						fmt.Println("==> Recv Msg: ID=", msg.ID, ", len=", msg.DataLen, ", data=", string(msg.Data))
					}
				}
			}(conn)
		}
	}()

	// Client goroutine, responsible for simulating data containing sticky packets and sending it to the server.
	go func() {
		conn, err := net.Dial("tcp", "127.0.0.1:7777")
		if err != nil {
			fmt.Println("client dial err:", err)
			return
		}

		// Create a packet splitting and packaging object dp.
		dp := Factory().NewPack(ziface.ZinxDataPack)

		// Package msg1.
		msg1 := &Message{
			ID:      0,
			DataLen: 5,
			Data:    []byte{'h', 'e', 'l', 'l', 'o'},
		}

		sendData1, err := dp.Pack(msg1)
		if err != nil {
			fmt.Println("client pack msg1 err:", err)
			return
		}

		// Package msg2.
		msg2 := &Message{
			ID:      1,
			DataLen: 7,
			Data:    []byte{'w', 'o', 'r', 'l', 'd', '!', '!'},
		}
		sendData2, err := dp.Pack(msg2)
		if err != nil {
			fmt.Println("client temp msg2 err:", err)
			return
		}

		// Concatenate sendData1 and sendData2 to create a sticky packet.
		sendData1 = append(sendData1, sendData2...)

		// Write data to the server.
		conn.Write(sendData1)
	}()

	// Block the client.
	select {
	case <-time.After(time.Second):
		return
	}
}
