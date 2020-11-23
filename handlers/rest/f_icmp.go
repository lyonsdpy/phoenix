// @Author : DAIPENGYUAN
// @File : icmp
// @Time : 2020/9/30 11:31
// @Description : ping相关接口

package rest

import (
	"net/http"
	"phoenix/mod"
	"phoenix/service"
	"phoenix/task"
	"phoenix/utils/errno"

	"github.com/gin-gonic/gin"
)

var IcmpRoute = []mod.Route{
	{Method: "POST", Path: "/icmp/sync", Handler: Icmp, Comment: "icmp同步请求接口"},
	{Method: "POST", Path: "/icmp/async", Handler: IcmpAsync, Comment: "icmp异步"},
	{Method: "POST", Path: "/icmp/job", Handler: IcmpJob, Comment: "icmp任务"},
}

// @Tags ICMP
// @Summery Icmp-Ping同步请求
// @Description 请求时的返回时间等于count*max_rtt
// @Accept  json
// @Produce json
// @Param icmp_req body mod.ICMPRequest true "ICMP的请求body"
// @Param target header string true "ICMP请求目标" default(www.baidu.com)
// @Success 200 {object} mod.ICMPResponse "ICMP返回信息"
// @Router /icmp/sync [post]
func Icmp(c *gin.Context) {
	var req mod.Request
	req.Base = bindHeaderBase(c)
	req.Base.SrvType = mod.TICMP
	err := c.BindJSON(&req.ReqData.IcmpReq)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.Icmp.WithErr(err))
		return
	}
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.Icmp.WithErr(err))
		return
	}
	srv := service.NewService(req.Base)
	rst, err := srv.Ping(req.ReqData.IcmpReq.Size, req.ReqData.IcmpReq.Count,
		req.ReqData.IcmpReq.MaxRtt, req.ReqData.IcmpReq.UseUdp)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.Icmp.WithErr(err))
		return
	}
	c.JSON(http.StatusOK, rst)
	return
}

// @Tags ICMP
// @Summery Icmp-Ping异步请求
// @Description 请求时的返回时间等于count*max_rtt
// @Accept  json
// @Produce json
// @Param icmp_req body mod.ICMPRequest true "ICMP的请求body"
// @Param target header string true "ICMP请求目标" default(www.baidu.com)
// @Param cb_target header string true "http回调目标" default(http://127.0.0.1:8080/v1/common)
// @Param cb_method header string true "http回调方法" Enums(POST)
// @Success 200 {object} errno.Err "{"code":200,"message"："请求成功","data":{"taskid":"93363008-23fc-11eb-a41a-f875a41832dc"}}"
// @Router /icmp/async [post]
func IcmpAsync(c *gin.Context) {
	var req mod.Request
	req.Base = bindHeaderBase(c)
	req.Job = bindHeaderJob(c)
	req.Base.SrvType = mod.TICMP
	req.Job.RetType = mod.RCALLBACK
	err := c.BindJSON(&req.ReqData.IcmpReq)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.Icmp.WithErr(err))
		return
	}
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.Icmp.WithErr(err))
		return
	}
	// 新增任务
	tk, err := task.NewTask(req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.Icmp.WithErr(err))
		return
	}
	c.JSON(http.StatusOK, errno.OK.ID(tk.ID))
	go tk.Run()
}

// @Tags ICMP
// @Summery Icmp-Ping任务请求
// @Description 任务的执行间隔不要小于count*max_rtt的值
// @Accept  json
// @Produce json
// @Param icmp_req body mod.ICMPRequest true "ICMP的请求body"
// @Param target header string true "ICMP请求目标" default(www.baidu.com)
// @Param spec header string true "任务时间间隔" Enums(@every 5s,@every 10s,@hourly,/5 * * * *)
// @Param cb_target header string true "http回调目标" default(http://127.0.0.1:8080/v1/common)
// @Param cb_method header string true "http回调方法" Enums(POST)
// @Success 200 {object} errno.Err "{"code":200,"message"："请求成功","data":{"taskid":"93363008-23fc-11eb-a41a-f875a41832dc"}}"
// @Router /icmp/job [post]
func IcmpJob(c *gin.Context) {
	var req mod.Request
	req.Base = bindHeaderBase(c)
	req.Job = bindHeaderJob(c)
	req.Base.SrvType = mod.TICMP
	req.Job.RetType = mod.RCALLBACK
	err := c.BindJSON(&req.ReqData.IcmpReq)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.Icmp.WithErr(err))
		return
	}
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.Icmp.WithErr(err))
		return
	}
	// 新增任务
	tk, err := task.NewTask(req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.Icmp.WithErr(err))
		return
	}
	err = task.Tasks.AddJob(tk)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.Icmp.WithErr(err))
		return
	}
	c.JSON(http.StatusOK, errno.OK.ID(tk.ID))
}
