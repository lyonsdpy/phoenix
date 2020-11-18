// @Author : DAIPENGYUAN
// @File : sftp
// @Time : 2020/9/30 11:30
// @Description :

package rest

import (
	"net/http"
	"phoenix/mod"
	"phoenix/service"
	"phoenix/task"
	"phoenix/utils/errno"

	"github.com/gin-gonic/gin"
)

var SftpRoute = []mod.Route{
	{Method: "POST", Path: "/sftp-get/sync", Handler: SftpGet, Comment: "sftp-get同步请求接口"},
	{Method: "POST", Path: "/sftp-get/async", Handler: SftpGetAsync, Comment: "sftp-get异步"},
	{Method: "POST", Path: "/sftp-get/job", Handler: SftpGetJob, Comment: "sftp-get任务"},
	{Method: "POST", Path: "/sftp-send/sync", Handler: SftpSend, Comment: "sftp-send同步请求接口"},
	{Method: "POST", Path: "/sftp-send/async", Handler: SftpSendAsync, Comment: "sftp-send异步"},
	{Method: "POST", Path: "/sftp-send/job", Handler: SftpSendJob, Comment: "sftp-send任务"},
}

// @Tags SFTP
// @Summery SFTP-GET同步请求
// @Description sftp从远端获取文件，返回文件的二进制
// @Accept  multipart/form-data
// @Produce application/octet-stream
// @Param target header string true "sftp请求目标" default(127.0.0.1:22)
// @Param username header string true "用户名"
// @Param password header string true "密码"
// @Param filepath formData string true "远端的文件路径(含文件名)"
// @Success 200 {string} string "远端文件二进制"
// @Router /sftp-get/sync [post]
func SftpGet(c *gin.Context) {
	var req mod.Request
	var err error
	req.Base = bindHeaderBase(c)
	req.Base.SrvType = mod.TSFTPGET
	sftpreq, err := bindFileForm(c)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.ScpGet.WithErr(err))
		return
	}
	req.ReqData.SftpReq = &sftpreq
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.ScpGet.WithErr(err))
		return
	}
	srv := service.NewService(req.Base)
	fname, fctt, err := srv.SftpGetFile(req.ReqData.SftpReq.Filepath)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.ScpGet.WithErr(err))
		return
	}
	respFile(fname, fctt, c)
	return
}

// @Tags SFTP
// @Summery SFTP-GET异步请求
// @Description sftp从远端获取文件
// @Accept  multipart/form-data
// @Produce application/json
// @Param target header string true "sftp请求目标" default(127.0.0.1:22)
// @Param username header string true "用户名"
// @Param password header string true "密码"
// @Param filepath formData string true "远端的文件路径(含文件名)"
// @Param cb_target header string true "http回调目标" default(http://127.0.0.1:8080/v1/common)
// @Param cb_method header string true "http回调方法" Enums(POST)
// @Success 200 {object} errno.Err "{"code":200,"message"："请求成功","data":{"taskid":"93363008-23fc-11eb-a41a-f875a41832dc"}}"
// @Router /sftp-get/async [post]
func SftpGetAsync(c *gin.Context) {
	var req mod.Request
	var err error
	req.Base = bindHeaderBase(c)
	req.Job = bindHeaderJob(c)
	req.Base.SrvType = mod.TSFTPGET
	req.Job.RetType = mod.RCALLBACK
	sftpreq, err := bindFileForm(c)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SftpGet.WithErr(err))
		return
	}
	req.ReqData.SftpReq = &sftpreq
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SftpGet.WithErr(err))
		return
	}
	tk, err := task.NewTask(req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.SftpGet.WithErr(err))
		return
	}
	c.JSON(http.StatusOK, errno.OK.ID(tk.ID))
	go tk.Run()
}

// @Tags SFTP
// @Summery SFTP-GET任务请求
// @Description sftp从远端获取文件，通常用于定时备份配置文件等
// @Accept  multipart/form-data
// @Produce application/json
// @Param target header string true "SFTP请求目标" default(127.0.0.1:22)
// @Param username header string true "用户名"
// @Param password header string true "密码"
// @Param filepath formData string true "远端的文件路径(含文件名)"
// @Param spec header string true "任务时间间隔" Enums(@every 5s,@every 10s,@hourly,/5 * * * *)
// @Param cb_target header string true "http回调目标" default(http://127.0.0.1:8080/v1/common)
// @Param cb_method header string true "http回调方法" Enums(POST)
// @Success 200 {object} errno.Err "{"code":200,"message"："请求成功","data":{"taskid":"93363008-23fc-11eb-a41a-f875a41832dc"}}"
// @Router /sftp-get/job [post]
func SftpGetJob(c *gin.Context) {
	var req mod.Request
	var err error
	req.Base = bindHeaderBase(c)
	req.Job = bindHeaderJob(c)
	req.Base.SrvType = mod.TSFTPGET
	req.Job.RetType = mod.RCALLBACK
	sftpreq, err := bindFileForm(c)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SftpGet.WithErr(err))
		return
	}
	req.ReqData.SftpReq = &sftpreq
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SftpGet.WithErr(err))
		return
	}
	tk, err := task.NewTask(req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.SftpGet.WithErr(err))
		return
	}
	err = task.Tasks.AddJob(tk)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.SftpGet.WithErr(err))
		return
	}
	c.JSON(http.StatusOK, errno.OK.ID(tk.ID))
}

// @Tags SFTP
// @Summery SFTP-SEND同步请求
// @Description sftp将文件保存到远端
// @Accept  multipart/form-data
// @Produce application/json
// @Param target header string true "SFTP请求目标" default(127.0.0.1:22)
// @Param username header string true "用户名"
// @Param password header string true "密码"
// @Param filepath formData string true "远端的文件路径(含文件名)"
// @Param filecontent formData file true "远端文件内容"
// @Success 200 {object} errno.Err "{"code":200,"message":"请求成功"}"
// @Router /sftp-send/sync [post]
func SftpSend(c *gin.Context) {
	var req mod.Request
	var err error
	req.Base = bindHeaderBase(c)
	req.Job = bindHeaderJob(c)
	req.Base.SrvType = mod.TSFTPSEND
	req.Job.RetType = mod.RCALLBACK
	sftpreq, err := bindFileForm(c)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SftpSend.WithErr(err))
		return
	}
	req.ReqData.SftpReq = &sftpreq
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SftpSend.WithErr(err))
		return
	}
	srv := service.NewService(req.Base)
	err = srv.SftpSendFile(req.ReqData.SftpReq.Filepath, req.ReqData.SftpReq.Filecontent)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.SftpSend.WithErr(err))
		return
	}
	c.JSON(http.StatusOK, errno.OK.ID())
}

// @Tags SFTP
// @Summery SFTP-SEND异步请求
// @Description sftp将文件保存到远端
// @Accept  multipart/form-data
// @Produce application/json
// @Param target header string true "SFTP请求目标" default(127.0.0.1:22)
// @Param username header string true "用户名"
// @Param password header string true "密码"
// @Param filepath formData string true "远端的文件路径(含文件名)"
// @Param filecontent formData file true "远端文件内容"
// @Param cb_target header string true "http回调目标" default(http://127.0.0.1:8080/v1/common)
// @Param cb_method header string true "http回调方法" Enums(POST)
// @Success 200 {object} errno.Err "{"code":200,"message"："请求成功","data":{"taskid":"93363008-23fc-11eb-a41a-f875a41832dc"}}"
// @Router /sftp-send/async [post]
func SftpSendAsync(c *gin.Context) {
	var req mod.Request
	var err error
	req.Base = bindHeaderBase(c)
	req.Job = bindHeaderJob(c)
	req.Base.SrvType = mod.TSFTPSEND
	req.Job.RetType = mod.RCALLBACK
	sftpreq, err := bindFileForm(c)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SftpSend.WithErr(err))
		return
	}
	req.ReqData.SftpReq = &sftpreq
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SftpSend.WithErr(err))
		return
	}
	tk, err := task.NewTask(req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.SftpSend.WithErr(err))
		return
	}
	c.JSON(http.StatusOK, errno.OK.ID(tk.ID))
	go tk.Run()
}

// @Tags SFTP
// @Summery SFTP-SEND任务请求
// @Description sftp将文件保存到远端
// @Accept multipart/form-data
// @Produce application/json
// @Param target header string true "SCP请求目标" default(127.0.0.1:22)
// @Param username header string true "用户名"
// @Param password header string true "密码"
// @Param filepath formData string true "远端的文件路径(含文件名)"
// @Param filecontent formData file true "远端文件内容"
// @Param spec header string true "任务时间间隔" Enums(@every 5s,@every 10s,@hourly,/5 * * * *)
// @Param cb_target header string true "http回调目标" default(http://127.0.0.1:8080/v1/common)
// @Param cb_method header string true "http回调方法" Enums(POST)
// @Success 200 {object} errno.Err "{"code":200,"message"："请求成功","data":{"taskid":"93363008-23fc-11eb-a41a-f875a41832dc"}}"
// @Router /sftp-send/job [post]
func SftpSendJob(c *gin.Context) {
	var req mod.Request
	var err error
	req.Base = bindHeaderBase(c)
	req.Job = bindHeaderJob(c)
	req.Base.SrvType = mod.TSFTPSEND
	req.Job.RetType = mod.RCALLBACK
	sftpreq, err := bindFileForm(c)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SftpSend.WithErr(err))
		return
	}
	req.ReqData.SftpReq = &sftpreq
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SftpSend.WithErr(err))
		return
	}
	tk, err := task.NewTask(req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.SftpSend.WithErr(err))
		return
	}
	err = task.Tasks.AddJob(tk)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.SftpSend.WithErr(err))
		return
	}
	c.JSON(http.StatusOK, errno.OK.ID(tk.ID))
}
