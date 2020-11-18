// @Author : DAIPENGYUAN
// @File : ssh_scp
// @Time : 2020/9/28 13:40
// @Description : scp服务

package service

import (
	"phoenix/engines"
	"phoenix/locker"
)

func (s Service) ScpGetFile(fpath string) (string, []byte, error) {
	err := locker.VTYGet(s.Target)
	if err != nil {
		return "", nil, err
	}
	defer locker.VTYRelease(s.Target)
	ng, err := engines.NewScp(s.Target, s.Username, s.Password)
	if err != nil {
		return "", nil, err
	}
	defer ng.Close()
	return ng.RcvFile(fpath)
}

func (s Service) ScpSendFile(fpath string, filectt []byte) error {
	err := locker.VTYGet(s.Target)
	if err != nil {
		return err
	}
	defer locker.VTYRelease(s.Target)
	ng, err := engines.NewScp(s.Target, s.Username, s.Password)
	if err != nil {
		return err
	}
	err = ng.SendFile(fpath, filectt)
	return err
}
