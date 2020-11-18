// @Author : DAIPENGYUAN
// @File : scp
// @Time : 2020/9/30 11:30
// @Description : scp操作相关接口

package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"phoenix/mod"
	"phoenix/service"
	"phoenix/task"
	"phoenix/utils/errno"
)

var ScpRoute = []mod.Route{
	{Method: "POST", Path: "/scp-get/sync", Handler: ScpGet, Comment: "scp-get同步请求接口"},
	{Method: "POST", Path: "/scp-get/async", Handler: ScpGetAsync, Comment: "scp-get异步"},
	{Method: "POST", Path: "/scp-get/job", Handler: ScpGetJob, Comment: "scp-get任务"},
	{Method: "POST", Path: "/scp-send/sync", Handler: ScpSend, Comment: "scp-send同步请求接口"},
	{Method: "POST", Path: "/scp-send/async", Handler: ScpSendAsync, Comment: "scp-send异步"},
	{Method: "POST", Path: "/scp-send/job", Handler: ScpSendJob, Comment: "scp-send任务"},
}

// @Tags SCP
// @Summery SCP-GET同步请求
// @Description scp从远端获取文件，返回文件的二进制
// @Accept  multipart/form-data
// @Param target header string true "SCP请求目标" default(127.0.0.1:22)
// @Param username header string true "用户名"
// @Param password header string true "密码"
// @Param filepath formData string true "远端的文件路径(含文件名)"
// @Success 200 {string} string "远端文件二进制"
// @Router /scp-get/sync [post]
func ScpGet(c *gin.Context) {
	var req mod.Request
	var err error
	req.Base = bindHeaderBase(c)
	req.Base.SrvType = mod.TSCPGET
	scpreq, err := bindFileForm(c)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.ScpGet.WithErr(err))
		return
	}
	req.ReqData.ScpReq = &scpreq
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.ScpGet.WithErr(err))
		return
	}
	srv := service.NewService(req.Base)
	fname, fctt, err := srv.ScpGetFile(req.ReqData.ScpReq.Filepath)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.ScpGet.WithErr(err))
		return
	}
	respFile(fname, fctt, c)
}

// @Tags SCP
// @Summery SCP-GET异步请求
// @Description scp从远端获取文件
// @Accept  multipart/form-data
// @Param target header string true "SCP请求目标" default(127.0.0.1:22)
// @Param username header string true "用户名"
// @Param password header string true "密码"
// @Param filepath formData string true "远端的文件路径(含文件名)"
// @Param cb_target header string true "http回调目标" default(http://127.0.0.1:8080/v1/common)
// @Param cb_method header string true "http回调方法" Enums(POST)
// @Success 200 {object} errno.Err "{"code":200,"message"："请求成功","data":{"taskid":"93363008-23fc-11eb-a41a-f875a41832dc"}}"
// @Router /scp-get/async [post]
func ScpGetAsync(c *gin.Context) {
	var req mod.Request
	var err error
	req.Base = bindHeaderBase(c)
	req.Job = bindHeaderJob(c)
	req.Base.SrvType = mod.TSCPGET
	req.Job.RetType = mod.RCALLBACK
	scpreq, err := bindFileForm(c)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.ScpGet.WithErr(err))
		return
	}
	req.ReqData.ScpReq = &scpreq
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.ScpGet.WithErr(err))
		return
	}
	tk, err := task.NewTask(req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.ScpGet.WithErr(err))
		return
	}
	c.JSON(http.StatusOK, errno.OK.ID(tk.ID))
	go tk.Run()
}

// @Tags SCP
// @Summery SCP-GET任务请求
// @Description scp从远端获取文件，通常用于定时备份配置文件等
// @Accept  multipart/form-data
// @Param target header string true "SCP请求目标" default(127.0.0.1:22)
// @Param username header string true "用户名"
// @Param password header string true "密码"
// @Param filepath formData string true "远端的文件路径(含文件名)"
// @Param spec header string true "任务时间间隔" Enums(@every 5s,@every 10s,@hourly,/5 * * * *)
// @Param cb_target header string true "http回调目标" default(http://127.0.0.1:8080/v1/common)
// @Param cb_method header string true "http回调方法" Enums(POST)
// @Success 200 {object} errno.Err "{"code":200,"message"："请求成功","data":{"taskid":"93363008-23fc-11eb-a41a-f875a41832dc"}}"
// @Router /scp-get/job [post]
func ScpGetJob(c *gin.Context) {
	var req mod.Request
	var err error
	req.Base = bindHeaderBase(c)
	req.Job = bindHeaderJob(c)
	req.Base.SrvType = mod.TSCPGET
	req.Job.RetType = mod.RCALLBACK
	scpreq, err := bindFileForm(c)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.ScpGet.WithErr(err))
		return
	}
	req.ReqData.ScpReq = &scpreq
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.ScpGet.WithErr(err))
		return
	}
	tk, err := task.NewTask(req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.ScpGet.WithErr(err))
		return
	}
	err = task.Tasks.AddJob(tk)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.ScpGet.WithErr(err))
		return
	}
	c.JSON(http.StatusOK, errno.OK.ID(tk.ID))
}

// @Tags SCP
// @Summery SCP-SEND同步请求
// @Description scp将文件保存到远端
// @Accept  multipart/form-data
// @Produce json
// @Param target header string true "SCP请求目标" default(127.0.0.1:22)
// @Param username header string true "用户名"
// @Param password header string true "密码"
// @Param filepath formData string true "远端的文件路径(含文件名)"
// @Param filecontent formData file true "远端文件内容"
// @Success 200 {object} errno.Err "{"code":200,"message":"请求成功"}"
// @Router /scp-send/sync [post]
func ScpSend(c *gin.Context) {
	var req mod.Request
	var err error
	req.Base = bindHeaderBase(c)
	req.Job = bindHeaderJob(c)
	req.Base.SrvType = mod.TSCPSEND
	req.Job.RetType = mod.RCALLBACK
	scpreq, err := bindFileForm(c)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.ScpSend.WithErr(err))
		return
	}
	req.ReqData.ScpReq = &scpreq
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.ScpSend.WithErr(err))
		return
	}
	srv := service.NewService(req.Base)
	err = srv.SftpSendFile(req.ReqData.ScpReq.Filepath, req.ReqData.ScpReq.Filecontent)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.ScpSend.WithErr(err))
		return
	}
	c.JSON(http.StatusOK, errno.OK.ID())
}

// @Tags SCP
// @Summery SCP-SEND异步请求
// @Description scp将文件保存到远端
// @Accept  multipart/form-data
// @Param target header string true "SCP请求目标" default(127.0.0.1:22)
// @Param username header string true "用户名"
// @Param password header string true "密码"
// @Param filepath formData string true "远端的文件路径(含文件名)"
// @Param filecontent formData file true "远端文件内容"
// @Param cb_target header string true "http回调目标" default(http://127.0.0.1:8080/v1/common)
// @Param cb_method header string true "http回调方法" Enums(POST)
// @Success 200 {object} errno.Err "{"code":200,"message"："请求成功","data":{"taskid":"93363008-23fc-11eb-a41a-f875a41832dc"}}"
// @Router /scp-send/async [post]
func ScpSendAsync(c *gin.Context) {
	var req mod.Request
	var err error
	req.Base = bindHeaderBase(c)
	req.Job = bindHeaderJob(c)
	req.Base.SrvType = mod.TSCPSEND
	req.Job.RetType = mod.RCALLBACK
	scpreq, err := bindFileForm(c)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.ScpSend.WithErr(err))
		return
	}
	req.ReqData.ScpReq = &scpreq
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.ScpSend.WithErr(err))
		return
	}
	tk, err := task.NewTask(req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.ScpSend.WithErr(err))
		return
	}
	c.JSON(http.StatusOK, errno.OK.ID(tk.ID))
	go tk.Run()
}

// @Tags SCP
// @Summery SCP-SEND任务请求
// @Description scp将文件保存到远端
// @Accept  multipart/form-data
// @Param target header string true "SCP请求目标" default(127.0.0.1:22)
// @Param username header string true "用户名"
// @Param password header string true "密码"
// @Param filepath formData string true "远端的文件路径(含文件名)"
// @Param filecontent formData file true "远端文件内容"
// @Param spec header string true "任务时间间隔" Enums(@every 5s,@every 10s,@hourly,/5 * * * *)
// @Param cb_target header string true "http回调目标" default(http://127.0.0.1:8080/v1/common)
// @Param cb_method header string true "http回调方法" Enums(POST)
// @Success 200 {object} errno.Err "{"code":200,"message"："请求成功","data":{"taskid":"93363008-23fc-11eb-a41a-f875a41832dc"}}"
// @Router /scp-send/job [post]
func ScpSendJob(c *gin.Context) {
	var req mod.Request
	var err error
	req.Base = bindHeaderBase(c)
	req.Job = bindHeaderJob(c)
	req.Base.SrvType = mod.TSCPSEND
	req.Job.RetType = mod.RCALLBACK
	scpreq, err := bindFileForm(c)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.ScpSend.WithErr(err))
		return
	}
	req.ReqData.ScpReq = &scpreq
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.ScpSend.WithErr(err))
		return
	}
	tk, err := task.NewTask(req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.ScpSend.WithErr(err))
		return
	}
	err = task.Tasks.AddJob(tk)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.ScpSend.WithErr(err))
		return
	}
	c.JSON(http.StatusOK, errno.OK.ID(tk.ID))
}
