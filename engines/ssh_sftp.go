// @Author : DAIPENGYUAN
// @File : sftp
// @Time : 2020/9/8 11:29 
// @Description : sftp引擎

package engines

import (
	"bufio"
	"bytes"
	"github.com/pkg/sftp"
	"path"
	"path/filepath"
)

func NewSFTP(target, user, password string) (*SSHSftp, error) {
	client, err := newSSHClient(target, user, password)
	if err != nil {
		return nil, err
	}
	s, err := sftp.NewClient(client)
	if err != nil {
		closeCLient(client)
		return nil, err
	}
	return &SSHSftp{sftp: s}, nil
}

type SSHSftp struct {
	sftp *sftp.Client
}

func (s *SSHSftp) RcvFile(remPath string) (string, []byte, error) {
	var buf = bytes.NewBuffer([]byte{})
	f, err := s.sftp.Open(remPath)
	if err != nil {
		return "", nil, err
	}
	defer f.Close()
	r := bufio.NewReader(f)
	_, err = r.WriteTo(buf)
	if err != nil {
		return "", nil, err
	}
	return filepath.Base(f.Name()), buf.Bytes(), nil
}

func (s *SSHSftp) SendFile(filename string, filectt []byte, remPath ...string) error {
	var p = filename
	if len(remPath) != 0 {
		p = path.Join(remPath[0], p)
	}
	f, err := s.sftp.Create(p)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)
	_, err = w.Write(filectt)
	if err != nil {
		return err
	}
	err = w.Flush()
	return err
}

func (s *SSHSftp) Close() error {
	if s.sftp != nil {
		return s.sftp.Close()
	}
	return nil
}
