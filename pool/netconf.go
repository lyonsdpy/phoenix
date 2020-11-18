// @Author : DAIPENGYUAN
// @File : netconf
// @Time : 2020/9/28 15:38 
// @Description : netconf连接池

package pool

import (
	"phoenix/conf"
	"phoenix/engines"
)

func GetNetconf(target, username, password string, capabilities ...string) (*engines.SSHNetconf, error) {
	key := getKey(target + username + password)
	NetconfCache.DeleteExpired()
	if cac, ok := NetconfCache.Get(key); ok {
		return cac.(*engines.SSHNetconf), nil
	}
	ng, err := engines.NewNetconf(target, username, password, capabilities...)
	if err != nil {
		return nil, err
	}
	NetconfCache.Set(key, ng, conf.PoolTimeout)
	return ng, nil
}

func DelNetconf(target, username, password string) {
	key := getKey(target + username + password)
	NetconfCache.Delete(key)
}
