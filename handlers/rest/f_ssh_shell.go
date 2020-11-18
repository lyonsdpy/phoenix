// @Author : DAIPENGYUAN
// @File : ssh
// @Time : 2020/9/30 11:30
// @Description :

package rest

import (
	"io"
	"net/http"

	"phoenix/mod"
	"phoenix/service"
	"phoenix/task"
	"phoenix/utils/errno"

	"github.com/gin-gonic/gin"
)

var SSHShellRoute = []mod.Route{
	{Method: "POST", Path: "/ssh-shell/sync", Handler: SSHShell, Comment: "ssh-shell同步请求接口"},
	{Method: "POST", Path: "/ssh-shell/async", Handler: SSHShellAsync, Comment: "ssh-shell异步"},
	{Method: "POST", Path: "/ssh-shell/job", Handler: SSHShellJob, Comment: "ssh-shell任务"},
}

// @Tags SSH
// @Summery SSH-EXEC执行命令
// @Description EXEC可以多次执行，系统提供缓存
// @Accept  application/json
// @Param target header string true "ssh请求目标" default(127.0.0.1:22)
// @Param username header string true "用户名"
// @Param password header string true "密码"
// @Param cmd_list body array true "执行的命令列表,至少为1"
// @Success 200 {string} string "返回设备命令执行的回显结果"
// @Router /ssh-shell/sync [post]
func SSHShell(c *gin.Context) {
	var req mod.Request
	req.Base = bindHeaderBase(c)
	req.Base.SrvType = mod.TSSHEXEC
	err := c.BindJSON(&req.ReqData.SSHReq)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SshShell.WithErr(err))
		return
	}
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SshShell.WithErr(err))
		return
	}
	srv := service.NewService(req.Base)
	rst, err := srv.SSHExec(req.ReqData.SSHReq)
	if err != nil && err != io.EOF {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.SshShell.WithErr(err))
		return
	}
	c.String(200, "%s", string(rst))
}

// @Tags SSH
// @Summery SSH-EXEC执行命令(异步)
// @Description EXEC可以多次执行，系统提供缓存
// @Accept  application/json
// @Produce application/json
// @Param target header string true "ssh请求目标" default(127.0.0.1:22)
// @Param cb_target header string true "http回调目标" default(http://127.0.0.1:8080/v1/common)
// @Param cb_method header string true "http回调方法" Enums(POST)
// @Param username header string true "用户名"
// @Param password header string true "密码"
// @Param cmd_list body string true "执行的命令列表,至少为1"
// @Success 200 {object} errno.Err
// @Router /ssh-shell/async [post]
func SSHShellAsync(c *gin.Context) {
	var req mod.Request
	req.Base = bindHeaderBase(c)
	req.Job = bindHeaderJob(c)
	req.Base.SrvType = mod.TSSHEXEC
	req.Job.RetType = mod.RCALLBACK
	err := c.BindJSON(&req.ReqData.SSHReq)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SshShell.WithErr(err))
		return
	}
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SshShell.WithErr(err))
		return
	}
	tk, err := task.NewTask(req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.SshShell.WithErr(err))
		return
	}
	c.JSON(http.StatusOK, errno.OK.ID(tk.ID))
	go tk.Run()
}

// @Tags SSH
// @Summery SSH-EXEC执行命令(任务)
// @Description EXEC可以多次执行，系统提供缓存
// @Accept  application/json
// @Produce application/json
// @Param target header string true "ssh请求目标" default(127.0.0.1:22)
// @Param cb_target header string true "http回调目标" default(http://127.0.0.1:8080/v1/common)
// @Param cb_method header string true "http回调方法" Enums(POST)
// @Param username header string true "用户名"
// @Param password header string true "密码"
// @Param cmd_list body string true "执行的命令列表,至少为1"
// @Param spec header string true "任务时间间隔" Enums(@every 5s,@every 10s,@hourly,/5 * * * *)
// @Success 200 {object} errno.Err
// @Router /ssh-shell/job [post]
func SSHShellJob(c *gin.Context) {
	var req mod.Request
	req.Base = bindHeaderBase(c)
	req.Job = bindHeaderJob(c)
	req.Base.SrvType = mod.TSSHEXEC
	req.Job.RetType = mod.RCALLBACK
	err := c.BindJSON(&req.ReqData.SSHReq)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SshShell.WithErr(err))
		return
	}
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SshShell.WithErr(err))
		return
	}
	tk, err := task.NewTask(req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.SshShell.WithErr(err))
		return
	}
	err = task.Tasks.AddJob(tk)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.SshShell.WithErr(err))
		return
	}
	c.JSON(http.StatusOK, errno.OK.ID(tk.ID))
}
