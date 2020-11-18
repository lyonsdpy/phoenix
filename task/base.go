// @Author : DAIPENGYUAN
// @File : types
// @Time : 2020/10/22 15:02
// @Description : 固定类型定义

package task

import (
	"errors"
	"phoenix/mod"
	"phoenix/service"
	"phoenix/utils/errno"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

// 创建新任务并检查参数是否合法，并设置部分默认参数
func NewTask(req mod.Request) (*Task, error) {
	// 检查请求参数的合法性
	if req.Job.RetType == 0 {
		req.Job.RetType = mod.RCALLBACK
	}
	if req.Job.RetType == mod.RCALLBACK {
		err := checkCallback(req.Job.CallBack)
		if err != nil {
			return nil, errors.New("create task failed," + err.Error())
		}
	}
	// 创建任务
	task := &Task{
		Job:     req.Job,
		Base:    req.Base,
		ID:      uuid.NewV1().String(),
		ReqData: req.ReqData,
		cronId:  0,
	}
	return task, nil
}

type Task struct {
	mod.Job
	Base mod.Base
	// 任务UUID,任务产生时自动生成
	ID      string      `json:"id" uri:"id" binding:"required,uuid"`
	ReqData mod.ReqData `json:"req_data,omitempty"`
	State   *TState     `json:"state,omitempty"`
	cronId  int         `json:"-"` // cron自动执行时的ID参数,用于cron的查看和删除
}

type TState struct {
	LastRun string `json:"last_run,omitempty"`
	Error   error  `json:"error,omitempty"`
}

func (s Task) newService() service.Service {
	return service.Service{Base: s.Base}
}

func (s Task) Run() {
	var err error

	if s.RetType == mod.RCALLBACK {
		err = s.runCallback()
	} else {
		err = s.notImplementCallback()
	}

	s.flushState(err)
}

// 从任务列表中更新任务的状态
func (s Task) flushState(err error) {
	if v, ok := Tasks[s.ID]; ok {
		if v.State == nil {
			v.State = new(TState)
		}
		v.State.LastRun = time.Now().Format("2006-01-02 15:04:05")
		if err != nil {
			v.State.Error = err
		}
	}
}

// 实现cron的run接口,任务周期执行的主入口
// 请求了service之后再通过Callback接口发送给远端目标

// 检查callback参数是否合法
func checkCallback(cb mod.CallBack) error {
	if cb.CbTarget == "" {
		return errors.New("callback check failed,target nil")
	}
	cb.CbMethod = strings.ToUpper(cb.CbMethod)
	if cb.CbMethod != "POST" {
		return errors.New("callback check failed,nil or wrong method")
	}
	return nil
}

func HttpResp(taskID string) *errno.Err {
	return errno.OK.ID(taskID)
}
