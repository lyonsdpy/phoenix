// @Author : DAIPENGYUAN
// @File : callback
// @Time : 2020/11/4 14:58
// @Description :

package task

import (
	"errors"
	"fmt"
	"phoenix/conf"
	"phoenix/mod"
	"phoenix/utils/common"
	"phoenix/utils/errno"
	"phoenix/utils/httputil"
)

func (s Task) runCallback() error {
	switch s.Base.SrvType {
	case mod.TICMP:
		return s.icmpCallback()
	case mod.TSNMPBULKGET:
		return s.snmpGetCallback()
	case mod.TSNMPWALK:
		return s.snmpWalkCallback()
	case mod.TNC:
		return s.ncExecCallback()
	case mod.TNCCAP:
		return s.ncCapCallback()
	case mod.TSFTPGET:
		return s.sftpGetCallback()
	case mod.TSFTPSEND:
		return s.sftpSendCallback()
	case mod.TSCPGET:
		return s.scpGetCallback()
	case mod.TSCPSEND:
		return s.scpSendCallback()
	case mod.TSSHRUN:
		return s.sshRunCallback()
	case mod.TSSHEXEC:
		return s.sshExecCallback()
	case mod.TTELNET:
		// TODO:telnet尚未实现
		return s.notImplementCallback()
	default:
		return s.notImplementCallback()
	}
}

func (s Task) icmpCallback() error {
	var err error
	s.CallBack.CbContentType = httputil.TJSON
	resp, err := s.newService().Ping(s.ReqData.IcmpReq.Size, s.ReqData.IcmpReq.Count,
		s.ReqData.IcmpReq.MaxRtt, s.ReqData.IcmpReq.UseUdp)
	if err != nil {
		s.rawCallback([]byte(common.IndentJson(err.Error())))
		return err
	}
	_, err = s.rawCallback([]byte(common.IndentJson(resp)))
	return err
}

// snmpGetCallback:json返回
func (s Task) snmpGetCallback() error {
	s.CallBack.CbContentType = httputil.TJSON
	resp, err := s.newService().SnmpBulkGet(s.ReqData.SnmpReq)
	if err != nil {
		s.rawCallback([]byte(errno.SnmpGet.WithErr(err).Json()))
		return err
	}
	_, err = s.rawCallback([]byte(common.IndentJson(resp)))
	return err
}

func (s Task) snmpWalkCallback() error {
	s.CallBack.CbContentType = httputil.TJSON
	resp, err := s.newService().SnmpWalk(s.ReqData.SnmpReq[0])
	if err != nil {
		s.rawCallback([]byte(errno.SnmpWalk.WithErr(err).Json()))
		return err
	}
	_, err = s.rawCallback([]byte(common.IndentJson(resp)))
	return err
}

func (s Task) ncExecCallback() error {
	s.CallBack.CbContentType = httputil.TXML
	resp, err := s.newService().Nc(s.ReqData.NetconfReq)
	if err != nil {
		s.rawCallback([]byte(errno.Netconf.WithErr(err).Json()))
		return err
	}
	_, err = s.rawCallback([]byte(common.IndentXml(resp)))
	return err
}

func (s Task) ncCapCallback() error {
	s.CallBack.CbContentType = httputil.TXML
	resp, err := s.newService().NcCap()
	if err != nil {
		s.rawCallback([]byte(errno.NetconfCap.WithErr(err).Xml()))
		return err
	}
	_, err = s.rawCallback([]byte(common.IndentXml(resp)))
	return err
}

func (s Task) sshRunCallback() error {
	s.CallBack.CbContentType = httputil.TTEXT
	resp, err := s.newService().SSHRun(s.ReqData.SSHReq)
	if err != nil {
		s.CallBack.CbContentType = httputil.TJSON
		s.rawCallback([]byte(errno.SshRun.WithErr(err).Json()))
		return err
	}
	_, err = s.rawCallback(resp)
	return err
}

func (s Task) sshExecCallback() error {
	s.CallBack.CbContentType = httputil.TTEXT
	resp, err := s.newService().SSHExec(s.ReqData.SSHReq)
	if err != nil {
		s.CallBack.CbContentType = httputil.TJSON
		s.rawCallback([]byte(errno.SshShell.WithErr(err).Json()))
		return err
	}
	_, err = s.rawCallback(resp)
	return err
}

func (s Task) sftpGetCallback() error {
	fname, fctt, err := s.newService().SftpGetFile(s.ReqData.SftpReq.Filepath)
	if err != nil {
		s.CallBack.CbContentType = httputil.TJSON
		s.rawCallback([]byte(errno.SftpGet.WithErr(err).Json()))
		return err
	}
	s.CallBack.CbContentType = httputil.TSTREAM
	s.CallBack.CbHeaders = make(map[string]string)
	s.CallBack.CbHeaders["Content-Disposition"] = fmt.Sprintf(`attachment;filename="%s"`, fname)
	s.CallBack.CbHeaders["Pragma"] = "No-cache"
	s.CallBack.CbHeaders["Cache-Control"] = "no-cache"
	s.CallBack.CbHeaders["Expires"] = "0"
	_, err = s.rawCallback(fctt)
	return err
}

func (s Task) sftpSendCallback() error {
	s.CallBack.CbContentType = httputil.TJSON
	err := s.newService().SftpSendFile(s.ReqData.SftpReq.Filepath, s.ReqData.SftpReq.Filecontent)
	if err != nil {
		s.rawCallback([]byte(errno.SftpSend.WithErr(err).Json()))
		return err
	}
	_, err = s.rawCallback([]byte(errno.OK.ID().Json()))
	return err
}

func (s Task) scpGetCallback() error {
	fname, fctt, err := s.newService().ScpGetFile(s.ReqData.ScpReq.Filepath)
	if err != nil {
		s.CallBack.CbContentType = httputil.TJSON
		s.rawCallback([]byte(errno.ScpGet.WithErr(err).Json()))
		return err
	}
	s.CallBack.CbContentType = httputil.TSTREAM
	s.CallBack.CbHeaders = make(map[string]string)
	s.CallBack.CbHeaders["Content-Disposition"] = fmt.Sprintf(`attachment;filename="%s"`, fname)
	s.CallBack.CbHeaders["Pragma"] = "No-cache"
	s.CallBack.CbHeaders["Cache-Control"] = "no-cache"
	s.CallBack.CbHeaders["Expires"] = "0"
	_, err = s.rawCallback(fctt)
	return err
}

func (s Task) scpSendCallback() error {
	s.CallBack.CbContentType = httputil.TJSON
	err := s.newService().ScpSendFile(s.ReqData.ScpReq.Filepath, s.ReqData.ScpReq.Filecontent)
	if err != nil {
		s.rawCallback([]byte(errno.ScpSend.WithErr(err).Json()))
		return err
	}
	_, err = s.rawCallback([]byte(errno.OK.ID().Json()))
	return err
}

func (s Task) notImplementCallback() error {
	s.CallBack.CbContentType = httputil.TJSON
	s.rawCallback([]byte(errno.NotImplement.ID().Json()))
	return errors.New("function not implement")
}

func (s Task) rawCallback(body []byte) ([]byte, error) {
	if s.CallBack.CbHeaders == nil {
		s.CallBack.CbHeaders = make(map[string]string)
	}
	s.CallBack.CbHeaders["ID"] = s.ID
	url := s.CallBack.CbTarget
	if s.CallBack.CbParams != "" {
		url += "?" + s.CallBack.CbParams
	}
	httpClient := &httputil.Client{
		URI:         url,
		Method:      s.CallBack.CbMethod,
		Headers:     s.CallBack.CbHeaders,
		ContentType: s.CallBack.CbContentType,
		Timeout:     conf.NetTimeout,
	}
	resp, err := httpClient.RAW(body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
