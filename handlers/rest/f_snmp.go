// @Author : DAIPENGYUAN
// @File : snmp
// @Time : 2020/9/30 11:29
// @Description : snmp相关接口

package rest

import (
	"net/http"
	"phoenix/mod"
	"phoenix/service"
	"phoenix/task"
	"phoenix/utils/errno"
	"phoenix/utils/log"

	"github.com/gin-gonic/gin"
	_ "github.com/lyonsdpy/gosnmp"
)

var SnmpRoute = []mod.Route{
	{Method: "POST", Path: "/snmp-typelist", Handler: SnmpValueList, Comment: "获取snmp返回值类型对应表"},
	{Method: "POST", Path: "/snmp-bulkget/sync", Handler: SnmpGet, Comment: "snmp-get同步请求接口"},
	{Method: "POST", Path: "/snmp-bulkget/async", Handler: SnmpGetAsync, Comment: "snmp-get异步"},
	{Method: "POST", Path: "/snmp-bulkget/job", Handler: SnmpGetJob, Comment: "snmp-get任务"},
	{Method: "POST", Path: "/snmp-walk/sync", Handler: SnmpWalk, Comment: "snmp-walk同步请求接口"},
	{Method: "POST", Path: "/snmp-walk/async", Handler: SnmpWalkAsync, Comment: "snmp-walk异步"},
	{Method: "POST", Path: "/snmp-walk/job", Handler: SnmpWalkJob, Comment: "snmp-walk任务"},
}

// @Tags SNMP
// @Summery 获取SNMP数据类型对应关系
// @Description 获取SNMP数据类型对应关系
// @Accept  application/json
// @Produce application/json
// @Param target header string true "snmp请求目标" default(127.0.0.1:161)
// @Param community header string true "团体字"
// @Success 200 {string} string "{1:"integer",2:"bitstring"...}"
// @Router /snmp-typelist [post]
func SnmpValueList(c *gin.Context) {
	srv := service.NewService(mod.Request{}.Base)
	v := srv.SnmpTypeList()
	c.JSON(http.StatusOK, v)
}

// @Tags SNMP
// @Summery SNMP-BULKGET获取数据方法
// @Description 一次性获取多个单独的OID对应的值
// @Accept  application/json
// @Produce application/json
// @Param target header string true "snmp请求目标" default(127.0.0.1:161)
// @Param community header string true "团体字"
// @Param oid_list body array true "snmp请求oid列表,至少为1"
// @Success 200 {array} gosnmp.SnmpPDU "[{"name":"1.3.6.4.1.2.0","type":1,"value":"gi1/0/1"},{...}]"
// @Router /snmp-bulkget/sync [post]
func SnmpGet(c *gin.Context) {
	var req mod.Request
	req.Base = bindHeaderBase(c)
	req.Base.SrvType = mod.TSNMPBULKGET
	err := c.BindJSON(&req.ReqData.SnmpReq)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SnmpGet.WithErr(err))
		return
	}
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SnmpGet.WithErr(err))
		return
	}
	srv := service.NewService(req.Base)
	rst, err := srv.SnmpBulkGet(req.ReqData.SnmpReq)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.SnmpGet.WithErr(err))
		return
	}
	c.JSON(http.StatusOK, rst)
	return
}

// @Tags SNMP
// @Summery SNMP-BULKGET异步请求
// @Description 一次性获取多个单独的OID对应的值
// @Accept  application/json
// @Produce application/json
// @Param target header string true "SNMP请求目标" default(127.0.0.1:161)
// @Param cb_target header string true "http回调目标" default(http://127.0.0.1:8080/v1/common)
// @Param cb_method header string true "http回调方法" Enums(POST)
// @Param community header string true "团体字"
// @Param req_xml body string true "snmp请求oid列表,至少为1"
// @Success 200 {object} errno.Err
// @Router /snmp-bulkget/async [post]
func SnmpGetAsync(c *gin.Context) {
	var req mod.Request
	req.Base = bindHeaderBase(c)
	req.Job = bindHeaderJob(c)
	req.Base.SrvType = mod.TSNMPBULKGET
	req.Job.RetType = mod.RCALLBACK
	err := c.BindJSON(&req.ReqData.SnmpReq)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SnmpGet.WithErr(err))
		return
	}
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SnmpGet.WithErr(err))
		return
	}
	tk, err := task.NewTask(req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.SnmpGet.WithErr(err))
		return
	}
	go tk.Run()
	c.JSON(http.StatusOK, errno.OK.ID(tk.ID))
}

// @Tags SNMP
// @Summery SNMP-BULKGET任务请求
// @Description 一次性获取多个单独的OID对应的值
// @Accept  application/json
// @Produce application/json
// @Param target header string true "SNMP请求目标" default(127.0.0.1:161)
// @Param cb_target header string true "http回调目标" default(http://127.0.0.1:8080/v1/common)
// @Param cb_method header string true "http回调方法" Enums(POST)
// @Param community header string true "团体字"
// @Param spec header string true "任务时间间隔" Enums(@every 5s,@every 10s,@hourly,/5 * * * *)
// @Param req_xml body string true "snmp请求oid列表,至少为1"
// @Success 200 {object} errno.Err
// @Router /snmp-bulkget/job [post]
func SnmpGetJob(c *gin.Context) {
	var req mod.Request
	var oids []string
	err := c.ShouldBindJSON(&oids)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SnmpGet.WithErr(err))
		return
	}
	req.ReqData.SnmpReq = oids
	req.Base = bindHeaderBase(c)
	req.Job = bindHeaderJob(c)
	req.Base.SrvType = mod.TSNMPBULKGET
	req.Job.RetType = mod.RCALLBACK
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SnmpGet.WithErr(err))
		return
	}
	tk, err := task.NewTask(req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.SnmpGet.WithErr(err))
		return
	}
	err = task.Tasks.AddJob(tk)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.SnmpGet.WithErr(err))
		return
	}
	c.JSON(http.StatusOK, errno.OK.ID(tk.ID))
}

// @Tags SNMP
// @Summery SNMP-WALK获取数据方法
// @Description walk一个oid下的所有值
// @Accept  application/json
// @Produce application/json
// @Param target header string true "snmp请求目标" default(127.0.0.1:161)
// @Param community header string true "团体字"
// @Param oid_list body array true "snmp请求oid列表,只能为1"
// @Success 200 {array} gosnmp.SnmpPDU "[{"name":"1.3.6.4.1.2.0","type":1,"value":"gi1/0/1"},{...}]"
// @Router /snmp-walk/sync [post]
func SnmpWalk(c *gin.Context) {
	var req mod.Request
	req.Base = bindHeaderBase(c)
	req.Base.SrvType = mod.TSNMPWALK
	err := c.BindJSON(&req.ReqData.SnmpReq)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SnmpWalk.WithErr(err))
		return
	}
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SnmpWalk.WithErr(err))
		return
	}
	log.Logger.Debug(req.ReqData.SnmpReq)
	srv := service.NewService(req.Base)
	rst, err := srv.SnmpWalk(req.ReqData.SnmpReq[0])
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.SnmpWalk.WithErr(err))
		return
	}
	c.JSON(http.StatusOK, rst)
	return
}

// @Tags SNMP
// @Summery SNMP-WALK异步请求
// @Description walk一个oid下的所有值
// @Accept  application/json
// @Produce application/json
// @Param target header string true "SNMP请求目标" default(127.0.0.1:161)
// @Param cb_target header string true "http回调目标" default(http://127.0.0.1:8080/v1/common)
// @Param cb_method header string true "http回调方法" Enums(POST)
// @Param community header string true "团体字"
// @Param req_xml body string true "snmp请求oid列表,只能为1"
// @Success 200 {object} errno.Err
// @Router /snmp-walk/async [post]
func SnmpWalkAsync(c *gin.Context) {
	var req mod.Request
	req.Base = bindHeaderBase(c)
	req.Job = bindHeaderJob(c)
	req.Base.SrvType = mod.TSNMPWALK
	req.Job.RetType = mod.RCALLBACK
	err := c.BindJSON(&req.ReqData.SnmpReq)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SnmpWalk.WithErr(err))
		return
	}
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SnmpWalk.WithErr(err))
		return
	}
	tk, err := task.NewTask(req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.SnmpWalk.WithErr(err))
		return
	}
	c.JSON(http.StatusOK, errno.OK.ID(tk.ID))
	go tk.Run()
}

// @Tags SNMP
// @Summery SNMP-WALK任务请求
// @Description walk一个oid下的所有值
// @Accept  application/json
// @Produce application/json
// @Param target header string true "SNMP请求目标" default(127.0.0.1:161)
// @Param cb_target header string true "http回调目标" default(http://127.0.0.1:8080/v1/common)
// @Param cb_method header string true "http回调方法" Enums(POST)
// @Param community header string true "团体字"
// @Param spec header string true "任务时间间隔" Enums(@every 5s,@every 10s,@hourly,/5 * * * *)
// @Param req_xml body string true "snmp请求oid列表,只能为1"
// @Success 200 {object} errno.Err
// @Router /snmp-walk/job [post]
func SnmpWalkJob(c *gin.Context) {
	var req mod.Request
	req.Base = bindHeaderBase(c)
	req.Job = bindHeaderJob(c)
	req.Base.SrvType = mod.TSNMPWALK
	req.Job.RetType = mod.RCALLBACK
	err := c.BindJSON(&req.ReqData.SnmpReq)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SnmpWalk.WithErr(err))
		return
	}
	err = req.Check()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.SnmpWalk.WithErr(err))
		return
	}
	tk, err := task.NewTask(req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.SnmpWalk.WithErr(err))
		return
	}
	err = task.Tasks.AddJob(tk)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.SnmpWalk.WithErr(err))
		return
	}
	c.JSON(http.StatusOK, errno.OK.ID(tk.ID))
}
