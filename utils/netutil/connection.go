// @Author : DAIPENGYUAN
// @File : net
// @Time : 2020/9/9 10:04 
// @Description : 网络通用工具

package netutil

import (
	"fmt"
	"net"
	"time"
)

func NewConn(conn net.Conn, readTimeout, writeTimeout time.Duration) Conn {
	return Conn{
		conn,
		readTimeout,
		writeTimeout,
	}
}

type Conn struct {
	net.Conn
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func (c Conn) Read(b []byte) (n int, err error) {
	c.Conn.SetReadDeadline(time.Now().Add(c.ReadTimeout))
	n, err = c.Conn.Read(b)
	if err != nil {
		fmt.Println(err)
	}
	return n, err
}

func (c Conn) Write(b []byte) (n int, err error) {
	c.Conn.SetWriteDeadline(time.Now().Add(c.WriteTimeout))
	return c.Conn.Write(b)
}

// 普通拨号网络连接
// target = "127.0.0.1:22"
// protocol = "tcp"
func DialConn(target, protocol string, timeout time.Duration) (net.Conn, error) {
	conn, err := net.DialTimeout(protocol, target, timeout)
	if err != nil {
		return nil, err
	} else {
		return conn, nil
	}
}

// 带读写超时的网络连接
func DialTimeoutConn(target, protocol string,
	timeout, readTimeout, writeTimeout time.Duration) (net.Conn, error) {
	conn, err := net.DialTimeout(protocol, target, timeout)
	if err != nil {
		return nil, err
	}
	connTimeout := &Conn{conn, readTimeout, writeTimeout}
	return connTimeout, nil
}
