// @Author : DAIPENGYUAN
// @File : base
// @Time : 2020/10/27 14:55 
// @Description : 

package service

import "phoenix/mod"

func NewService(Base mod.Base) Service {
	return Service{Base}
}

type Service struct {
	mod.Base
}
