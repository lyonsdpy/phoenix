// @Author : DAIPENGYUAN
// @File : snmp
// @Time : 2020/9/26 17:48
// @Description :

package locker

import (
	"errors"
	"phoenix/conf"
)

var sNMPLockChan = make(map[string]chan int)

func SnmpGet(target string) error {
	var err error
	if v, ok := sNMPLockChan[target]; ok {
		err = getTimout(v)
	} else {
		sNMPLockChan[target] = make(chan int, conf.SNMPLocker)
		err = getTimout(sNMPLockChan[target])
	}
	return err
}

func SnmpRelease(target string) error {
	var err error
	if v, ok := sNMPLockChan[target]; ok {
		err = releaseTimout(v)
		if len(v) == 0 {
			close(sNMPLockChan[target])
			delete(sNMPLockChan, target)
		}
	} else {
		err = errors.New("release lock failed:locker don't exist")
	}
	return err
}
