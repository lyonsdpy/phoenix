// @Author : DAIPENGYUAN
// @File : basic
// @Time : 2020/9/30 14:23
// @Description :

package mod

import (
	"errors"
	"phoenix/utils/common"
	"strings"
)

// API请求结构体
type Request struct {
	Base    Base    `json:"base,omitempty"`
	Job     Job     `json:"job,omitempty"`
	ReqData ReqData `json:"req_data,omitempty"`
}

func (s Request) Check() error {
	switch s.Base.SrvType {
	case TUNKNOWN:
		return errors.New("service check failed,unknown service type")
	case TICMP:
		if s.ReqData.IcmpReq == nil {
			return errors.New("icmp service check failed,empty request data")
		}
		if s.ReqData.IcmpReq.MaxRtt == 0 {
			s.ReqData.IcmpReq.MaxRtt = 1000
		}
		if s.ReqData.IcmpReq.Count == 0 {
			s.ReqData.IcmpReq.Count = 1
		}
		if s.ReqData.IcmpReq.Size == 0 {
			s.ReqData.IcmpReq.Size = 64
		}
	case TSNMPTYPELIST:
		break
	case TSNMPBULKGET:
		if len(s.ReqData.SnmpReq) == 0 {
			return errors.New("snmp-get service check failed,empty oid list")
		}
		if s.Base.Community == "" {
			return errors.New("snmp-get service check failed,empty community")
		}
	case TSNMPWALK:
		if len(s.ReqData.SnmpReq) != 1 {
			return errors.New("snmp-walk service check failed,oid length !=1")
		}
		if s.Base.Community == "" {
			return errors.New("snmp-walk service check failed,empty community")
		}
		if !strings.HasPrefix(s.ReqData.SnmpReq[0], ".") {
			s.ReqData.SnmpReq[0] = "." + s.ReqData.SnmpReq[0]
		}
	case TNC:
		if len(s.ReqData.NetconfReq) == 0 {
			return errors.New("netconf service check failed,empty xml string")
		}
		if !common.IsValidXML([]byte(s.ReqData.NetconfReq)) {
			return errors.New("netconf service check failed,xmlstr is not a valid xml")
		}
		if s.Base.Username == "" || s.Base.Password == "" {
			return errors.New("netconf service check failed,empty username or password")
		}
	case TNCCAP:
		if s.Base.Username == "" || s.Base.Password == "" {
			return errors.New("netconf-capaliblity service check failed,empty username or password")
		}
	case TSSHRUN:
	case TSSHEXEC:
		if len(s.ReqData.SSHReq) == 0 {
			return errors.New("ssh service check failed,cmdlist lengh==0")
		}
		if s.Base.Username == "" || s.Base.Password == "" {
			return errors.New("ssh service check failed,empty username or password")
		}
	case TSCPGET:
		if s.ReqData.ScpReq.Filepath == "" {
			return errors.New("scp-get service check failed,empty filepath")
		}
		if s.Base.Username == "" || s.Base.Password == "" {
			return errors.New("scp-get service check failed,empty username or password")
		}
	case TSFTPGET:
		if s.ReqData.SftpReq.Filepath == "" {
			return errors.New("sftp-get service check failed,empty filepath")
		}
		if s.Base.Username == "" || s.Base.Password == "" {
			return errors.New("sftp-get service check failed,empty username or password")
		}
	case TSCPSEND:
		if s.ReqData.ScpReq.Filepath == "" {
			return errors.New("scp-send service check failed,empty filepath")
		}
		if len(s.ReqData.ScpReq.Filecontent) == 0 {
			return errors.New("scp-send service check failed,empty filecontent")
		}
		if s.Base.Username == "" || s.Base.Password == "" {
			return errors.New("scp-send service check failed,empty username or password")
		}
	case TSFTPSEND:
		if s.ReqData.SftpReq.Filepath == "" {
			return errors.New("sftp-send service check failed,empty filepath")
		}
		if len(s.ReqData.SftpReq.Filecontent) == 0 {
			return errors.New("sftp-send service check failed,empty filecontent")
		}
		if s.Base.Username == "" || s.Base.Password == "" {
			return errors.New("sftp-send service check failed,empty username or password")
		}
	case TTELNET:
		return errors.New("telnet not implement")
	case THTTP:
		return errors.New("httputil not implement")
	default:
		return errors.New("create task failed,service type unimplement")
	}
	return nil
}

// 认证相关参数
type Base struct {
	Target string `json:"target"`
	// 任务类型 :1=icmp-ping,2=snmp-get,3=snmp-walk,4=netconf-exec,5=netconf-capability
	// 6=sftp-getFile,7=sftp-sendFile,8=ssh-run,9=ssh-exec,10=telnet
	SrvType   TTYPE    `json:"srv_type"`
	Username  string   `json:"username,omitempty"`
	Password  string   `json:"password,omitempty"`
	Community string   `json:"community,omitempty"`
	NcCaps    []string `json:"nc_caps,omitempty"`
}

// 调用端请求的通用参数，如果有备注则会在handler处理中加入Headers
type Job struct {
	// 返回值类型 :1=callback 2=cache
	RetType  RType    `json:"ret_type,omitempty"`
	Spec     string   `json:"spec,omitempty"`
	Comment  string   `json:"comment,omitempty"`
	CallBack CallBack `json:"callback,omitempty"`
}

// HTTP回调信息
type CallBack struct {
	CbTarget      string            `json:"cb_target,omitempty"`
	CbParams      string            `json:"cb_params,omitempty"`
	CbMethod      string            `json:"cb_method,omitempty"`
	CbContentType string            `json:"cb_content_type,omitempty"`
	CbHeaders     map[string]string `json:"cb_headers,omitempty"`
}

// 请求数据内容
type ReqData struct {
	IcmpReq    *ICMPRequest  `json:"icmp_req,omitempty"`
	NetconfReq string        `json:"netconf_req,omitempty"`
	ScpReq     *SFileRequest `json:"scp_req,omitempty"`
	SftpReq    *SFileRequest `json:"sftp_req,omitempty"`
	SnmpReq    []string      `json:"snmp_req,omitempty"`
	SSHReq     []string      `json:"ssh_req,omitempty"`
}

// 常规返回结果
type AsyncNormalResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	JobId   string `json:"job_id,omitempty"` // 任务ID,在处理任务时会有值
}

// ICMPRequest:ICMP请求结构体
// MaxRtt:max round trip time,最大往返时间,毫秒
type ICMPRequest struct {
	Size   int   `json:"size,omitempty"`
	Count  int   `json:"count,omitempty"`
	MaxRtt int64 `json:"max_rtt,omitempty"`
	UseUdp bool  `json:"use_udp,omitempty"`
}

// ICMPResponse:ICMP返回结构体
type ICMPResponse struct {
	PacketsRecv  int      `json:"packets_recv"`
	PacketsSent  int      `json:"packets_sent"`
	PacketLoss   float64  `json:"packet_loss"`
	RemoteIpAddr string   `json:"remote_ip_addr"`
	RemoteAddr   string   `json:"remote_addr"`
	Rtts         []string `json:"rtts"`
	MinRtt       string   `json:"min_rtt"`
	MaxRtt       string   `json:"max_rtt"`
	AvgRtt       string   `json:"avg_rtt"`
	StdDevRtt    string   `json:"std_dev_rtt"`
}

type SFileRequest struct {
	Filepath    string `json:"filepath,omitempty"`
	Filecontent []byte `json:"filecontent,omitempty"`
}
