package znet

import (
	"bufio"
	"net"
	"net/http"
	"github.com/aceld/zinx/zlog"
)

type customWriterRecorder struct {
	conn *net.TCPConn
	rw   *bufio.ReadWriter
}

func newResponseWriter(conn *net.TCPConn, reader *bufio.Reader) *customWriterRecorder {
	w := new(customWriterRecorder)
	w.conn = conn
	w.rw = bufio.NewReadWriter(reader, bufio.NewWriter(conn))

	return w
}

func (w *customWriterRecorder) Header() http.Header {
	zlog.Ins().ErrorF("Header has not been implemented")
	return http.Header{}
}

func (w *customWriterRecorder) Write(bytes []byte) (int, error) {
	return w.rw.Write(bytes)
}

func (w *customWriterRecorder) WriteHeader(statusCode int) {
	//TODO implement me
	zlog.Ins().ErrorF("WriteHeader has not been implemented")
}

func (w *customWriterRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.conn, w.rw, nil
}
