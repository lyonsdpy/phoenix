// @Author : DAIPENGYUAN
// @File : base
// @Time : 2020/9/27 10:56 
// @Description : 

package locker

import (
	"errors"
	"phoenix/conf"
	"time"
)

func getTimout(ch chan int) error {
	select {
	case ch <- 1:
		return nil
	case <-time.After(conf.LockerTimeout):
		return errors.New("get locker timeout")
	}
}

func releaseTimout(ch chan int) error {
	select {
	case <-ch:
		return nil
	case <-time.After(conf.LockerTimeout):
		return errors.New("release locker timeout")
	}
}