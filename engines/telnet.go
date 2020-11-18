// @Author : DAIPENGYUAN
// @File : telnet
// @Time : 2020/9/8 11:30 
// @Description : telnet引擎

package engines

import (
	"errors"
	"fmt"
	"net"
	"phoenix/conf"
	"phoenix/utils/netutil"
	"strings"
	"time"
)

// telnet建立网络连接
func NewTelnet(target, user, password string) (*Telnet, error) {
	//var buf = make([]byte, 4096)
	conn, err := netutil.DialTimeoutConn(target, "tcp",
		conf.NetTimeout, conf.ReadTimeout, conf.WriteTimeout)
	if err != nil {
		return nil, err
	}
	if !telnetHandshake(conn, user, password) {
		conn.Close()
		return nil, errors.New("telnet handshake failed")
	}
	time.Sleep(conf.ReadSleep)
	//n, err := conn.Read(buf)
	//if err != nil {
	//	conn.Close()
	//	return nil, err
	//}
	//fmt.Println(n)
	//fmt.Println(string(buf))
	return &Telnet{conn: conn}, nil
}

type Telnet struct {
	conn net.Conn
}

func (s *Telnet) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

func (s *Telnet) Exec(cmds ...string) (string, error) {
	var (
		err      error
		cmdBytes = []byte("\n" + strings.Join(cmds, "\n") + "\n" + msgSeperator + "\n")
	)
	defer func() {
		if err != nil {
			s.Close()
		}
	}()
	_, err = s.conn.Write(cmdBytes)
	if err != nil {
		return "", nil
	}
	readOut, err := cliRead(s.conn)
	return string(readOut), err
}

func telnetHandshake(conn net.Conn, user, pass string) bool {
	var buf [4096]byte
	n, err := conn.Read(buf[0:])
	if nil != err {
		fmt.Println(err)
		return false
	}
	buf[1] = 252
	buf[4] = 252
	buf[7] = 252
	buf[10] = 252
	n, err = conn.Write(buf[0:n])
	if nil != err {
		fmt.Println(err)
		return false
	}

	n, err = conn.Read(buf[0:])
	if nil != err {
		fmt.Println(err)
		return false
	}
	buf[1] = 252
	buf[4] = 251
	buf[7] = 252
	buf[10] = 254
	buf[13] = 252
	n, err = conn.Write(buf[0:n])
	if nil != err {
		fmt.Println(err)
		return false
	}

	n, err = conn.Read(buf[0:])
	if nil != err {
		fmt.Println(err)
		return false
	}
	buf[1] = 252
	buf[4] = 252
	n, err = conn.Write(buf[0:n])
	if nil != err {
		fmt.Println(err)
		return false
	}

	n, err = conn.Read(buf[0:])
	if nil != err {
		fmt.Println(err)
		return false
	}
	//if false == s.IsAuthentication {
	//	return true
	//}

	n, err = conn.Write([]byte(user + "\n"))
	if nil != err {
		fmt.Println(err)
		return false
	}
	time.Sleep(conf.ReadSleep)

	n, err = conn.Read(buf[0:])
	if nil != err {
		fmt.Println(err)
		return false
	}
	fmt.Println(string(buf[0:n]))

	n, err = conn.Write([]byte(pass + "\n"))
	if nil != err {
		fmt.Println(err)
		return false
	}
	time.Sleep(conf.ReadSleep)

	n, err = conn.Read(buf[0:])
	if nil != err {
		fmt.Println(err)
		return false
	}
	return true
}
