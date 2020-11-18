// @Author : DAIPENGYUAN
// @File : common
// @Time : 2020/9/8 11:09 
// @Description : 连接池管理，目前管理netconf,ssh-shell,telnet的连接

package pool

import (
	"crypto/sha1"
	"fmt"
	"github.com/patrickmn/go-cache"
	"time"
)

const (
	DefaultExpire        time.Duration = 300
	DefaultCleanInterval time.Duration = 60
)

var (
	ShellCache   = cache.New(DefaultExpire, DefaultCleanInterval)
	NetconfCache = cache.New(DefaultExpire, DefaultCleanInterval) // Netconf连接池
	TelnetCache  = cache.New(DefaultExpire, DefaultCleanInterval)
)

func getKey(origin string) string {
	Sha1Inst := sha1.New()
	Sha1Inst.Write([]byte(origin))
	Result := Sha1Inst.Sum([]byte(""))
	return fmt.Sprintf("%x", Result)
}
