// @Author : DAIPENGYUAN
// @File : errno
// @Time : 2020/11/4 22:55 
// @Description : 错误码
// TODO:实现注册错误码方法

package errno

import (
	"encoding/json"
	"encoding/xml"
)

const (
	OK     ErrType = 10200
	ScpGet ErrType = iota + 10500
	ScpSend
	SftpGet
	SftpSend
	SnmpGet
	SnmpWalk
	Netconf
	NetconfCap
	Icmp
	SshRun
	SshShell
	Telnet
	Job
	NotImplement
)

var Emap = map[ErrType]string{
	OK:           "请求成功",
	ScpGet:       "SCP获取文件失败",
	ScpSend:      "SCP发送文件失败",
	SftpGet:      "SFTP获取文件失败",
	SftpSend:     "SFTP发送文件失败",
	SnmpGet:      "SNMP GET执行失败",
	SnmpWalk:     "SNMP WALK执行失败",
	Netconf:      "NETCONF执行失败",
	NetconfCap:   "NETCONF能力集获取失败",
	Icmp:         "ICMP请求失败",
	SshRun:       "SSH执行命令失败",
	SshShell:     "SSH交互式命令失败",
	Telnet:       "TELNET执行命令失败",
	Job:          "任务操作失败",
	NotImplement: "尚未实现",
}

type ErrType int

func (e ErrType) String() string {
	return Emap[e]
}

func (e ErrType) Int() int {
	return int(e)
}

func (e ErrType) ID(id ...string) *Err {
	var r = Err{
		Code:    e.Int(),
		Message: e.String(),
	}
	if len(id) != 0 {
		r.ID = id[0]
	}
	return &r
}

func (e ErrType) WithErr(err error) *Err {
	var eString string
	if err != nil {
		eString = err.Error()
	}
	var r = Err{
		Code:    e.Int(),
		Message: e.String(),
		Error:   eString,
	}
	return &r
}

type Err struct {
	Code    int    `json:"code" xml:"code"`
	Message string `json:"message" xml:"message"`
	Error   string `json:"error,omitempty" xml:"error,omitempty"`
	ID      string `json:"id,omitempty" xml:"data,omitempty"`
}

func (e *Err) Json() string {
	ret, err := json.MarshalIndent(e, "", "\t")
	if err != nil {
		eb, _ := json.MarshalIndent(err.Error(), "", "\t")
		return string(eb)
	}
	return string(ret)
}

func (e *Err) Xml() string {
	r, err := xml.MarshalIndent(e, "", "\t")
	if err != nil {
		eb, _ := xml.MarshalIndent(err.Error(), "", "\t")
		return string(eb)
	}
	return string(r)
}
