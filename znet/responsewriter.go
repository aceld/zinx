package znet

import (
	"bufio"
	"net"
	"net/http"
)

type wsFakeWriter struct {
	conn *net.TCPConn
	rw   *bufio.ReadWriter
}

func newResponseWriter(conn *net.TCPConn) *wsFakeWriter {
	w := new(wsFakeWriter)
	w.conn = conn
	w.rw = bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	return w
}

func (w *wsFakeWriter) Header() http.Header {
	return nil
}

func (w *wsFakeWriter) Write(bytes []byte) (int, error) {
	return 0, nil
}

func (w *wsFakeWriter) WriteHeader(statusCode int) {
	//TODO implement me
	panic("implement me")
}
func (w *wsFakeWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.conn, w.rw, nil
}
