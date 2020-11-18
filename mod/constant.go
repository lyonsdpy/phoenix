// @Author : DAIPENGYUAN
// @File : constant
// @Time : 2020/10/28 9:22 
// @Description : 

package mod

const (
	TUNKNOWN TTYPE = iota
	TICMP
	TSNMPTYPELIST
	TSNMPBULKGET
	TSNMPWALK
	TNC
	TNCCAP
	TSCPGET
	TSCPSEND
	TSFTPGET
	TSFTPSEND
	TSSHRUN
	TSSHEXEC
	TTELNET
	THTTP
)
const (
	RUNKNOWN RType = iota
	RCALLBACK
	RCACHE
)

type TTYPE int

func (s TTYPE) String() string {
	switch s {
	case TUNKNOWN:
		return "UNKNOWN"
	case TICMP:
		return "ICMP"
	case TSNMPBULKGET:
		return "SNMPGET"
	case TSNMPTYPELIST:
		return "TSNMPTYPELIST"
	case TSNMPWALK:
		return "SNMPWALK"
	case TNC:
		return "NETCONF"
	case TNCCAP:
		return "NETCONFCAP"
	case TSFTPGET:
		return "SFTPGET"
	case TSFTPSEND:
		return "SFTPSEND"
	case TSSHRUN:
		return "SSHRUN"
	case TSSHEXEC:
		return "SSHEXEC"
	case TTELNET:
		return "TELNET"
	default:
		return ""
	}
}

type RType int

func (s RType) String() string {
	switch s {
	case RUNKNOWN:
		return "UNKNOWN"
	case RCALLBACK:
		return "CALLBACK"
	case RCACHE:
		return "CACHE"
	default:
		return ""
	}
}
