// @Author : DAIPENGYUAN
// @File : snetconf
// @Time : 2020/9/30 11:31
// @Description :

package rest

import (
	"io/ioutil"
	"net/http"
	"phoenix/mod"
	"phoenix/service"
	"phoenix/task"
	"phoenix/utils/errno"

	"github.com/gin-gonic/gin"
)

var SnetconfRoute = []mod.Route{
	{Method: "POST", Path: "/snetconf/sync", Handler: Snetconf, Comment: "snetconf同步请求接口"},
	{Method: "POST", Path: "/snetconf/async", Handler: SnetconfAsync, Comment: "snetconf异步"},
	{Method: "POST", Path: "/snetconf/job", Handler: SnetconfJob, Comment: "snetconf任务"},
	{Method: "POST", Path: "/snetconf-cap/sync", Handler: SnetconfCap, Comment: "snetconf-cap同步请求接口"},
	{Method: "POST", Path: "/snetconf-cap/async", Handler: SnetconfCapAsync, Comment: "snetconf-cap异步"},
	{Method: "POST", Path: "/snetconf-cap/job", Handler: SnetconfCapJob, Comment: "snetconf-cap任务"},
}

// @Tags SNetconf
// @Summery SNetconf同步请求
// @Description netconf via ssh的交互
// @Accept  text/xml
// @Produce text/xml
// @Param target header string true "snetconf请求目标" default(127.0.0.1:830)
// @Param username header string true "用户名"
// @Param password header string true "密码"
// @Param req_xml body string true "netconf请求xml"
// @Success 200 {string} string "远端xml"
// @Router /snetconf/sync [post]
func Snetconf(c *gin.Context) {
	var req mod.Request
	var err error
	req.Base = bindHeaderBase(c)
	req.Base.SrvType = mod.TNC
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Error(err)
		c.XML(http.StatusBadRequest, errno.Netconf.WithErr(err))
		return
	}
	req.ReqData.NetconfReq = string(b)
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.XML(http.StatusBadRequest, errno.Netconf.WithErr(err))
		return
	}
	srv := service.NewService(req.Base)
	rst, err := srv.Nc(req.ReqData.NetconfReq)
	if err != nil {
		c.Error(err)
		c.XML(http.StatusInternalServerError, errno.Netconf.WithErr(err))
		return
	}
	c.XML(200, rst)
}

// @Tags SNetconf
// @Summery SNetconfT异步请求
// @Description SNetconf via ssh的异步交互
// @Accept  text/xml
// @Produce text/xml
// @Param target header string true "snetconf请求目标" default(127.0.0.1:830)
// @Param username header string true "用户名"
// @Param password header string true "密码"
// @Param cb_target header string true "http回调目标" default(http://127.0.0.1:8080/v1/common)
// @Param cb_method header string true "http回调方法" Enums(POST)
// @Param req_xml body string true "netconf请求xml"
// @Success 200 {object} errno.Err
// @Router /snetconf/async [post]
func SnetconfAsync(c *gin.Context) {
	var req mod.Request
	var err error
	req.Base = bindHeaderBase(c)
	req.Job = bindHeaderJob(c)
	req.Base.SrvType = mod.TNC
	req.Job.RetType = mod.RCALLBACK
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Error(err)
		c.XML(http.StatusBadRequest, errno.Netconf.WithErr(err))
		return
	}
	req.ReqData.NetconfReq = string(b)
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.XML(http.StatusBadRequest, errno.Netconf.WithErr(err))
		return
	}
	tk, err := task.NewTask(req)
	if err != nil {
		c.Error(err)
		c.XML(http.StatusInternalServerError, errno.Netconf.WithErr(err))
		return
	}
	c.XML(http.StatusOK, errno.OK.ID(tk.ID))
	go tk.Run()
}

// @Tags SNetconf
// @Summery SNetconf任务请求
// @Description SNetconf via ssh的异步交互
// @Accept  text/xml
// @Produce text/xml
// @Param target header string true "snetconf请求目标" default(127.0.0.1:830)
// @Param username header string true "用户名"
// @Param password header string true "密码"
// @Param spec header string true "任务时间间隔" Enums(@every 5s,@every 10s,@hourly,/5 * * * *)
// @Param cb_target header string true "http回调目标" default(http://127.0.0.1:8080/v1/common)
// @Param cb_method header string true "http回调方法" Enums(POST)
// @Param req_xml body string true "netconf请求xml"
// @Success 200 {object} errno.Err
// @Router /snetconf/job [post]
func SnetconfJob(c *gin.Context) {
	var req mod.Request
	var err error
	req.Base = bindHeaderBase(c)
	req.Job = bindHeaderJob(c)
	req.Base.SrvType = mod.TNC
	req.Job.RetType = mod.RCALLBACK
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.Error(err)
		c.XML(http.StatusBadRequest, errno.Netconf.WithErr(err))
		return
	}
	req.ReqData.NetconfReq = string(b)
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.XML(http.StatusBadRequest, errno.Netconf.WithErr(err))
		return
	}
	tk, err := task.NewTask(req)
	if err != nil {
		c.Error(err)
		c.XML(http.StatusInternalServerError, errno.Netconf.WithErr(err))
		return
	}
	err = task.Tasks.AddJob(tk)
	if err != nil {
		c.Error(err)
		c.XML(http.StatusInternalServerError, errno.Netconf.WithErr(err))
		return
	}
	c.XML(http.StatusOK, errno.OK.ID(tk.ID))
}

// @Tags SNetconf
// @Summery SNetconf同步请求能力集
// @Description netconf via ssh的交互,获取能力集
// @Produce text/xml
// @Param target header string true "snetconf请求目标" default(127.0.0.1:830)
// @Param username header string true "用户名"
// @Param password header string true "密码"
// @Success 200 {string} string "远端xml"
// @Router /snetconf-cap/sync [post]
func SnetconfCap(c *gin.Context) {
	var req mod.Request
	var err error
	req.Base = bindHeaderBase(c)
	req.Base.SrvType = mod.TNCCAP
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.XML(http.StatusBadRequest, errno.NetconfCap.WithErr(err))
		return
	}
	srv := service.NewService(req.Base)
	rst, err := srv.NcCap()
	if err != nil {
		c.Error(err)
		c.XML(http.StatusInternalServerError, errno.NetconfCap.WithErr(err))
		return
	}
	c.XML(200, rst)
}

// @Tags SNetconf
// @Summery SNetconfT异步请求获取能力集
// @Description SNetconf via ssh获取能力集的异步交互
// @Produce text/xml
// @Param target header string true "snetconf请求目标" default(127.0.0.1:830)
// @Param username header string true "用户名"
// @Param password header string true "密码"
// @Param cb_target header string true "http回调目标" default(http://127.0.0.1:8080/v1/common)
// @Param cb_method header string true "http回调方法" Enums(POST)
// @Success 200 {object} errno.Err
// @Router /snetconf-cap/async [post]
func SnetconfCapAsync(c *gin.Context) {
	var req mod.Request
	var err error
	req.Base = bindHeaderBase(c)
	req.Job = bindHeaderJob(c)
	req.Base.SrvType = mod.TNCCAP
	req.Job.RetType = mod.RCALLBACK
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.XML(http.StatusBadRequest, errno.NetconfCap.WithErr(err))
		return
	}
	tk, err := task.NewTask(req)
	if err != nil {
		c.Error(err)
		c.XML(http.StatusInternalServerError, errno.NetconfCap.WithErr(err))
		return
	}
	c.XML(http.StatusOK, errno.OK.ID(tk.ID))
	go tk.Run()
}

// @Tags SNetconf
// @Summery SNetconf任务请求获取能力集
// @Description SNetconf via ssh获取能力集的异步交互
// @Produce text/xml
// @Param target header string true "snetconf请求目标" default(127.0.0.1:830)
// @Param username header string true "用户名"
// @Param password header string true "密码"
// @Param spec header string true "任务时间间隔" Enums(@every 5s,@every 10s,@hourly,/5 * * * *)
// @Param cb_target header string true "http回调目标" default(http://127.0.0.1:8080/v1/common)
// @Param cb_method header string true "http回调方法" Enums(POST)
// @Success 200 {object} errno.Err
// @Router /snetconf-cap/job [post]
func SnetconfCapJob(c *gin.Context) {
	var req mod.Request
	var err error
	req.Base = bindHeaderBase(c)
	req.Job = bindHeaderJob(c)
	req.Base.SrvType = mod.TNCCAP
	req.Job.RetType = mod.RCALLBACK
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.XML(http.StatusBadRequest, errno.NetconfCap.WithErr(err))
		return
	}
	tk, err := task.NewTask(req)
	if err != nil {
		c.Error(err)
		c.XML(http.StatusInternalServerError, errno.NetconfCap.WithErr(err))
		return
	}
	err = task.Tasks.AddJob(tk)
	if err != nil {
		c.Error(err)
		c.XML(http.StatusInternalServerError, errno.NetconfCap.WithErr(err))
		return
	}
	c.XML(http.StatusOK, errno.OK.ID(tk.ID))
}
