/*
	zd为zinx distributed
	ZDConn是zinx node节点收发消息的基础连接结构
*/
package zdnet

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/aceld/zinx/utils"
)

type ZDConn struct {
	Ip   string
	Port int
	Conn net.Conn
}

/*
	创建新连接
*/
func NewZDConn(ip string, port int) *ZDConn {
	zdConn := &ZDConn{
		Ip:   ip,
		Port: port,
	}

	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("Dial ", address, "err = ", err)
		return nil
	}

	zdConn.Conn = conn

	return zdConn
}

/*
	读取数据
*/
func (zdconn *ZDConn) Recv() []byte {
	readData := make([]byte, 0)

	for {
		buf := make([]byte, utils.ZD_CONN_BUFSIZE)
		len, err := zdconn.Conn.Read(buf)
		if len < 1 && err != nil {
			if err == io.EOF {
				break
			}
		} else {
			if len == utils.ZD_CONN_BUFSIZE {
				readData = append(readData, buf...)
				//设置读取超时
				zdconn.Conn.SetDeadline(time.Now().Add(utils.ZD_CONN_READTIMEOUT * time.Millisecond))
				//简单睡眠，等待缓冲区，防止大数据灌溉
				time.Sleep(time.Millisecond)
			} else {
				readData = append(readData, buf[:len]...)
				break
			}

			if err == io.EOF {
				break
			}
		}
	}

	return readData
}

/*
	发送数据
*/
func (zdconn *ZDConn) Send(writeData []byte) error {
	for {
		_, err := zdconn.Conn.Write(writeData)
		if err != nil {
			return err
		} else {
			return nil
		}
	}
}

func (zdconn *ZDConn) Close() {
	if zdconn.Conn != nil {
		zdconn.Conn.Close()
	}
}
