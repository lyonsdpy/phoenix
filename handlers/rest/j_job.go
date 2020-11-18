// @Author : DAIPENGYUAN
// @File : j_job
// @Time : 2020/10/21 17:11
// @Description : 任务管理相关接口

package rest

import (
	"net/http"
	"phoenix/mod"
	"phoenix/task"
	"phoenix/utils/errno"

	"github.com/gin-gonic/gin"
)

var JobRoute = []mod.Route{
	{Method: "GET", Path: "/joblist", Handler: JobList, Comment: "列出所有的任务"},
	{Method: "GET", Path: "/job/:id", Handler: JobGet, Comment: "查询任务详情"},
	{Method: "DELETE", Path: "/job/:id", Handler: JobDelete, Comment: "删除任务"},
}

// @Tags JOB
// @Summery 任务查找,查询单个任务
// @Description 任务查找,查询单个任务
// @Accept  application/json
// @Produce application/json
// @Success 200 {string} string "返回设备命令执行的回显结果"
// @Router /ssh-shell/sync [post]
func JobGet(c *gin.Context) {
	var jbind task.Task
	err := c.ShouldBindUri(&jbind)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.Job.WithErr(err))
		return
	}
	job, err := task.Tasks.GetJob(jbind.ID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, errno.Job.WithErr(err))
		return
	}
	c.JSON(http.StatusOK, job)
}

func JobList(c *gin.Context) {
	jobs := task.Tasks.GetAllJob()
	c.JSON(http.StatusOK, jobs)
}

func JobDelete(c *gin.Context) {
	var jbind task.Task
	err := c.ShouldBindUri(&jbind)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.Job.WithErr(err))
		return
	}
	err = task.Tasks.DelJob(jbind.ID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, errno.Job.WithErr(err))
		return
	}
	c.JSON(http.StatusOK, errno.OK.ID())
}
