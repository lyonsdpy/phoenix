// @Author : DAIPENGYUAN
// @File : vty
// @Time : 2020/9/26 17:48 
// @Description : 

package locker

import (
	"errors"
	"phoenix/conf"
)

var vTYLockChan = make(map[string]chan int)

func VTYGet(target string) error {
	var err error
	if v, ok := vTYLockChan[target]; ok {
		err = getTimout(v)
	} else {
		vTYLockChan[target] = make(chan int, conf.VTYLocker)
		err = getTimout(vTYLockChan[target])
	}
	return err
}

func VTYRelease(target string) error {
	var err error
	if v, ok := vTYLockChan[target]; ok {
		err = releaseTimout(v)
		if len(v) == 0 {
			close(vTYLockChan[target])
			delete(vTYLockChan, target)
		}
	} else {
		err = errors.New("release lock failed:locker does not exist")
	}
	return err
}