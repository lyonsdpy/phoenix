// @Author : DAIPENGYUAN
// @File : scp
// @Time : 2020/9/8 11:29 
// @Description : scp引擎

package engines

import (
	"bytes"
	"github.com/hnakamur/go-scp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func NewScp(target, user, password string) (*SSHScp, error) {
	client, err := newSSHClient(target, user, password)
	if err != nil {
		return nil, err
	}
	scpClient := scp.NewSCP(client)
	return &SSHScp{scpClient, client}, nil
}

type SSHScp struct {
	scp       *scp.SCP
	sshClient *ssh.Client
}

// 从远端接收文件返回文件[]byte
func (s *SSHScp) RcvFile(fpath string) (string, []byte, error) {
	var buf = bytes.NewBuffer([]byte{})
	fInfo, err := s.scp.Receive(fpath, buf)
	if err != nil {
		return "", nil, err
	}
	return fInfo.Name(), buf.Bytes(), nil
}

// 将fileInfo中的文件[]byte发送到远端remPath中
func (s *SSHScp) SendFile(fpath string, filectt []byte) error {
	var (
		w     = ioutil.NopCloser(bytes.NewReader(filectt))
		fname = filepath.Base(fpath)
		fp    = fpath
	)
	remInfo := scp.NewFileInfo(fname, int64(len(filectt)), os.ModePerm, time.Now(), time.Now())
	return s.scp.Send(remInfo, w, fp)
}

func (s *SSHScp) Close() error {
	if s.sshClient != nil {
		return s.sshClient.Close()
	}
	return nil
}
