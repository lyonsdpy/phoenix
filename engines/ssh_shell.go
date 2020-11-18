// @Author : DAIPENGYUAN
// @File : ssh
// @Time : 2020/9/8 11:28 
// @Description : ssh的shell交互式引擎
// XXX:如果对端的回显内容中没有命令本身,则不应该使用shell交互式命令行

package engines

import (
	"golang.org/x/crypto/ssh"
	"io"
	"phoenix/conf"
	"strings"
	"time"
)

func NewShell(target, user, password string) (*SSHShell, error) {
	client, err := newSSHTimeoutClient(target, user, password)
	if err != nil {
		return nil, err
	}
	session, err := client.NewSession()
	if err != nil {
		closeCLient(client)
		return nil, err
	}
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 28800,
		ssh.TTY_OP_OSPEED: 28800,
	}
	session.RequestPty("xterm", 400, 256, modes)
	readWriteCloser, err := setShellReadWriteCloser(session)
	if err != nil {
		session.Close()
		return nil, err
	}
	time.Sleep(conf.ReadSleep)
	return &SSHShell{session: session, readWriterCloser: readWriteCloser}, nil
}

type SSHShell struct {
	session          *ssh.Session
	readWriterCloser io.ReadWriteCloser
}

func (s *SSHShell) Exec(cmds ...string) ([]byte, error) {
	var (
		err  error
		data = []byte("\n" + strings.Join(cmds, "\n") + "\n")
	)
	defer func() {
		if err != nil {
			s.Close()
		}
	}()
	cliWrite(s.readWriterCloser, data)
	readOut, err := cliRead(s.readWriterCloser)
	return readOut, err
}

func (s *SSHShell) Close() error {
	if s.session != nil {
		return s.session.Close()
	}
	return nil
}

func setShellReadWriteCloser(session *ssh.Session) (io.ReadWriteCloser, error) {
	w, err := session.StdinPipe()
	if err != nil {
		return nil, err
	}
	r, err := session.StdoutPipe()
	if err != nil {
		return nil, err
	}
	err = session.Shell()
	if err != nil {
		return nil, err
	}
	return NewReadWriteCloser(r, w), nil
}

type ReadWriteCloser struct {
	io.Reader
	io.WriteCloser
}

func NewReadWriteCloser(r io.Reader, w io.WriteCloser) *ReadWriteCloser {
	return &ReadWriteCloser{r, w}
}
