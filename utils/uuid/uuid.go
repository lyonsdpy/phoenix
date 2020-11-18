// @Author : DAIPENGYUAN
// @File : utiluuid
// @Time : 2020/9/7 15:09
// @Description : uuid工具

package uuid

import uuid "github.com/satori/go.uuid"

func NewUuidV1() string {
	return uuid.NewV1().String()
}

func NewUuidV4() string {
	return uuid.NewV4().String()
}