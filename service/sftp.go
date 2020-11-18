// @Author : DAIPENGYUAN
// @File : ssh_sftp
// @Time : 2020/9/28 13:41
// @Description : sftp服务

package service

import (
	"phoenix/engines"
	"phoenix/locker"
)

func (s Service) SftpGetFile(fpath string) (string, []byte, error) {
	err := locker.VTYGet(s.Target)
	if err != nil {
		return "", nil, err
	}
	defer locker.VTYRelease(s.Target)
	ng, err := engines.NewSFTP(s.Target, s.Username, s.Password)
	if err != nil {
		return "", nil, err
	}
	defer ng.Close()
	return ng.RcvFile(fpath)
}

func (s Service) SftpSendFile(fpath string, filectt []byte) error {
	err := locker.VTYGet(s.Target)
	if err != nil {
		return err
	}
	defer locker.VTYRelease(s.Target)
	ng, err := engines.NewSFTP(s.Target, s.Username, s.Password)
	if err != nil {
		return err
	}
	err = ng.SendFile(fpath, filectt)
	return err
}
