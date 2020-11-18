// @Author : DAIPENGYUAN
// @File : ssh_shell
// @Time : 2020/9/28 17:01 
// @Description : ssh shell连接池

package pool

import (
	"phoenix/conf"
	"phoenix/engines"
)

func GetShell(target, username, password string) (*engines.SSHShell, error) {
	key := getKey(target + username + password)
	ShellCache.DeleteExpired()
	if cac, ok := ShellCache.Get(key); ok {
		return cac.(*engines.SSHShell), nil
	}
	ng, err := engines.NewShell(target, username, password)
	if err != nil {
		return nil, err
	}
	ShellCache.Set(key, ng, conf.PoolTimeout)
	return ng, nil
}

func DelShell(target, username, password string) {
	key := getKey(target + username + password)
	ShellCache.Delete(key)
}
