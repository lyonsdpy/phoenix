// @Author : DAIPENGYUAN
// @File : ssh_netconf
// @Time : 2020/9/27 13:19
// @Description : netconf相关操作服务,有连接池相关操作

package service

import (
	"github.com/lyonsdpy/go-netconf/netconf"
	"phoenix/locker"
	"phoenix/pool"
)

func (s Service) Nc(xmlbyte string) (*netconf.RPCReply, error) {
	err := locker.VTYGet(s.Target)
	if err != nil {
		return nil, err
	}
	defer locker.VTYRelease(s.Target)
	ng, err := pool.GetNetconf(s.Target, s.Username, s.Password, s.NcCaps...)
	if err != nil {
		return nil, err
	}
	reply, err := ng.Exec(xmlbyte)
	if err != nil {
		ng.Close()
		pool.DelNetconf(s.Target, s.Username, s.Password)
		return nil, err
	}
	return reply, nil
}

func (s Service) NcCap() ([]string, error) {
	err := locker.VTYGet(s.Target)
	if err != nil {
		return nil, err
	}
	defer locker.VTYRelease(s.Target)
	ng, err := pool.GetNetconf(s.Target, s.Username, s.Password)
	if err != nil {
		return nil, err
	}
	return ng.Capabilities(), nil
}
