// @Author : DAIPENGYUAN
// @File : ssh_run
// @Time : 2020/9/28 13:40
// @Description : ssh shell的服务,其中有带pool缓存的方法

package service

import (
	"phoenix/engines"
	"phoenix/locker"
	"phoenix/pool"
)

// 直接运行命令，仅单次调用
func (s Service) SSHRun(cmdlist []string) ([]byte, error) {
	err := locker.VTYGet(s.Target)
	if err != nil {
		return nil, err
	}
	defer locker.VTYRelease(s.Target)
	return engines.SSHRun(s.Target, s.Username, s.Password, cmdlist...)
}

// 模拟交互式命令行，创建连接池，再连接超时前可以复用连接
func (s Service) SSHExec(cmdlist []string) ([]byte, error) {
	err := locker.VTYGet(s.Target)
	if err != nil {
		return nil, err
	}
	defer locker.VTYRelease(s.Target)
	ng, err := pool.GetShell(s.Target, s.Username, s.Password)
	if err != nil {
		return nil, err
	}
	reply, err := ng.Exec(cmdlist...)
	if err != nil {
		ng.Close()
		pool.DelShell(s.Target, s.Username, s.Password)
	}
	return reply, err
}
