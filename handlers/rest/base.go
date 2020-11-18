// @Author : DAIPENGYUAN
// @File : base
// @Time : 2020/10/29 9:28
// @Description :

package rest

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"phoenix/mod"
	"phoenix/utils/encrypt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func bindHeaderBase(c *gin.Context) mod.Base {
	var r mod.Base
	srvtype, _ := strconv.Atoi(c.Request.Header.Get("srv_type"))
	r.SrvType = mod.TTYPE(srvtype)
	r.Target = c.Request.Header.Get("target")
	r.Username = c.Request.Header.Get("username")
	password := c.Request.Header.Get("password")
	r.Password = encrypt.FastDecrypt(password)
	r.Community = c.Request.Header.Get("community")
	r.NcCaps = c.Request.Header.Values("nc_caps")
	return r
}

func bindHeaderJob(c *gin.Context) mod.Job {
	var r mod.Job
	rtype, _ := strconv.Atoi(c.Request.Header.Get("ret_type"))
	r.RetType = mod.RType(rtype)
	r.Spec = c.Request.Header.Get("spec")
	r.Comment = c.Request.Header.Get("comment")
	r.CallBack.CbTarget = c.Request.Header.Get("cb_target")
	r.CallBack.CbParams = c.Request.Header.Get("cb_params")
	r.CallBack.CbMethod = c.Request.Header.Get("cb_method")
	hds := c.Request.Header.Values("cb_headers")
	for _, v := range hds {
		if t, k, v := parseKeyVal(v); t {
			r.CallBack.CbHeaders[k] = v
		}
	}
	return r
}

func bindFileForm(c *gin.Context) (mod.SFileRequest, error) {
	var req mod.SFileRequest
	err := c.Request.ParseMultipartForm(32 << 20)
	if err != nil {
		return mod.SFileRequest{}, err
	}
	form := c.Request.MultipartForm
	for field, fvalue := range form.File {
		if field == "filecontent" {
			fd, _ := fvalue[0].Open()
			fb, _ := ioutil.ReadAll(fd)
			req.Filecontent = fb
		}
	}
	for field, valList := range form.Value {
		if field == "filepath" {
			req.Filepath = valList[0]
		}
	}
	if req.Filepath == "" {
		return mod.SFileRequest{}, errors.New("bind file form failed,empty filepath")
	}
	return req, nil
}

// 返回文件时的快捷方法
func respFile(fname string, fctt []byte, c *gin.Context) {
	hd := c.Writer.Header()
	//hd.Set("Content-Type", "binary/octet-stream")
	hd.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fname))
	hd.Set("Pragma", "No-cache")
	hd.Set("Cache-Control", "no-cache")
	hd.Set("Expires", "0")
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Write(fctt)
	c.Writer.Flush()
}

// 从http的form请求中解析出request请求参数
func parseFormData(r *http.Request) (mod.Request, error) {
	var req mod.Request
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return mod.Request{}, err
	}
	form := r.MultipartForm
	for field, fvalue := range form.File {
		for _, file := range fvalue {
			switch strings.ToLower(field) {
			case "sftp_req":
				f, _ := file.Open()
				fb, _ := ioutil.ReadAll(f)
				req.ReqData.SftpReq.Filecontent = fb
			case "scp_req":
				f, _ := file.Open()
				fb, _ := ioutil.ReadAll(f)
				req.ReqData.ScpReq.Filecontent = fb
			}
		}
	}
	for field, valList := range form.Value {
	for2:
		for _, val := range valList {
			switch strings.ToLower(field) {
			// base相关内容
			case "target":
				req.Base.Target = val
				break for2
			case "username":
				req.Base.Username = val
				break for2
			case "password":
				req.Base.Password = val
				break for2
			case "community":
				req.Base.Community = val
				break for2
			case "srv_type":
				vali, _ := strconv.Atoi(val)
				req.Base.SrvType = mod.TTYPE(vali)
			case "nc_caps":
				req.Base.NcCaps = append(req.Base.NcCaps, val)
			// job相关内容
			case "spec":
				req.Job.Spec = val
				break for2
			case "comment":
				req.Job.Comment = val
				break for2
			case "cb_target":
				req.Job.CallBack.CbTarget = val
				break for2
			case "cb_params":
				req.Job.CallBack.CbParams = val
				break for2
			case "cb_method":
				req.Job.CallBack.CbMethod = val
				break for2
			case "cb_content_type":
				req.Job.CallBack.CbContentType = val
				break for2
			case "cb_headers":
				if req.Job.CallBack.CbHeaders == nil {
					req.Job.CallBack.CbHeaders = make(map[string]string)
				}
				if t, k, v := parseKeyVal(val); t {
					req.Job.CallBack.CbHeaders[k] = v
				}
			// Req-Data相关内容
			case "ssh_req":
				for _, v := range strings.Split(val, "\n") {
					req.ReqData.SSHReq = append(req.ReqData.SSHReq, v)
				}
				break for2
			case "sftp_req":
				if t, k, v := parseKeyVal(val); t && k == "filepath" {
					req.ReqData.SftpReq.Filepath = v
				}
				break for2
			case "scp_req":
				if t, k, v := parseKeyVal(val); t && k == "filepath" {
					req.ReqData.ScpReq.Filepath = v
				}
				break for2
			case "snmp_req":
				for _, v := range strings.Split(val, "\n") {
					req.ReqData.SnmpReq = append(req.ReqData.SSHReq, v)
				}
				break for2
			case "netconf_req":
				req.ReqData.NetconfReq = val
				break for2
			case "icmp_req":
				if t, k, v := parseKeyVal(val); t {
					switch k {
					case "size":
						vn, _ := strconv.Atoi(v)
						req.ReqData.IcmpReq.Size = vn
					case "count":
						vn, _ := strconv.Atoi(v)
						req.ReqData.IcmpReq.Count = vn
					case "max_rtt":
						vn, _ := strconv.Atoi(v)
						req.ReqData.IcmpReq.MaxRtt = int64(vn)
					case "use_udp":
						if strings.ToLower(v) == "true" {
							req.ReqData.IcmpReq.UseUdp = true
						}
					}
				}
			}
		}
	}
	return req, nil
}

func parseKeyVal(val string) (bool, string, string) {
	vList := strings.Split(val, "=")
	if len(vList) == 2 {
		return true, vList[0], vList[1]
	}
	return false, "", ""
}
