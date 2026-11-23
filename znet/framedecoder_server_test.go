package znet

import (
	"encoding/binary"
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/aceld/zinx/zconf"
	"github.com/aceld/zinx/ziface"
)

// lengthFieldDecoder configures the connection frame decoder with a 2-byte
// length field at offset 0 that counts the payload only, and strips it.
// (为连接配置断粘包解码器：头部2字节长度字段，仅表示负载长度，解码后跳过)
type lengthFieldDecoder struct {
	maxFrameLength uint64
}

func (d *lengthFieldDecoder) GetLengthField() *ziface.LengthField {
	return &ziface.LengthField{
		MaxFrameLength:      d.maxFrameLength,
		LengthFieldOffset:   0,
		LengthFieldLength:   2,
		LengthAdjustment:    0,
		InitialBytesToStrip: 2,
	}
}

// Intercept keeps the message untouched, only framing matters here
// (拦截方法不修改消息，这里只关心断粘包)
func (d *lengthFieldDecoder) Intercept(chain ziface.IChain) ziface.IcResp {
	return chain.Proceed(chain.Request())
}

type collectingRouter struct {
	BaseRouter
	recv chan []byte
}

func (r *collectingRouter) Handle(req ziface.IRequest) {
	r.recv <- req.GetData()
}

// A peer that sends one frame above MaxFrameLength must not silence the
// connection: the frames it sends afterwards still reach their router.
// (对端发送一个超过MaxFrameLength的帧后，之后发送的合法帧仍应被路由处理)
func TestServerKeepsDeliveringFramesAfterOneOversizedFrame(t *testing.T) {
	// Ask the kernel for a free port, then let the server bind it.
	// (先取一个空闲端口再交给server监听)
	probe, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("probe listen: %v", err)
	}
	port := probe.Addr().(*net.TCPAddr).Port
	_ = probe.Close()

	config := &zconf.Config{Host: "127.0.0.1", TCPPort: port}

	// NewUserConfServer copies config into zconf.GlobalObject, which the rest of
	// the package reads. Put the previous values back so sibling tests still see
	// the configuration they expect.
	// (NewUserConfServer 会把配置写进全局对象, 测试结束后还原, 避免影响同包其他测试)
	prevHost, prevPort := zconf.GlobalObject.Host, zconf.GlobalObject.TCPPort
	t.Cleanup(func() {
		zconf.GlobalObject.Host, zconf.GlobalObject.TCPPort = prevHost, prevPort
	})

	server := NewUserConfServer(config)
	server.SetDecoder(&lengthFieldDecoder{maxFrameLength: 16})

	router := &collectingRouter{recv: make(chan []byte, 4)}
	server.AddRouter(0, router)

	server.Start()
	defer server.Stop()
	time.Sleep(300 * time.Millisecond)

	conn, err := net.Dial("tcp", "127.0.0.1:"+strconv.Itoa(port))
	if err != nil {
		t.Fatalf("dial server: %v", err)
	}
	defer conn.Close()

	// One frame the protocol allows to declare but the peer configured too short:
	// declared payload 100 bytes, so the frame is 102 bytes on the wire.
	// (一个超过上限的帧：声明100字节负载，实际线上102字节)
	oversized := make([]byte, 102)
	binary.BigEndian.PutUint16(oversized[0:2], 100)
	_ = conn.SetWriteDeadline(time.Now().Add(2 * time.Second))
	if _, err := conn.Write(oversized); err != nil {
		t.Fatalf("write oversized frame: %v", err)
	}

	// A well-formed frame right behind it. (紧随其后的合法帧)
	valid := []byte{0x00, 0x03, 'a', 'b', 'c'}
	if _, err := conn.Write(valid); err != nil {
		t.Fatalf("write valid frame: %v", err)
	}

	select {
	case got := <-router.recv:
		if string(got) != "abc" {
			t.Fatalf("router received %q, want %q", got, "abc")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("no frame was routed after an oversized frame; the connection went silent")
	}
}
