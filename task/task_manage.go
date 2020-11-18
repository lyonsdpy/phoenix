// @Author: Perry
// @Date  : 2020/10/2
// @Desc  : 任务管理
// 1-记录api调用时需要用到的信息
// 2-记录后台任务运行的ID
// 原则是可以通过前台信息能够还原出后台任务

package task

import (
	"errors"
	"github.com/robfig/cron/v3"
)

var (
	Tasks    = make(tasks)
	CronTask = cron.New()
)

type tasks map[string]*Task // key为job-id

// 列出当前所有的任务
func (s tasks) GetAllJob() []*Task {
	var tList []*Task
	for _, v := range s {
		tList = append(tList, v)
	}
	return tList
}

func (s tasks) GetJob(id string) (*Task, error) {
	if v, ok := s[id]; ok {
		return v, nil
	}
	return nil, errors.New("task not found")
}

// 增加计划任务
func (s tasks) AddJob(j *Task) error {
	if err := checkSpec(j.Spec); err != nil {
		return errors.New("add job failed," + err.Error())
	}
	entryID, err := CronTask.AddJob(j.Spec, j)
	if err != nil {
		return err
	}
	j.cronId = int(entryID)
	s[j.ID] = j
	return nil
}

// 删除计划任务
func (s tasks) DelJob(jId string) error {
	if job, ok := s[jId]; ok {
		CronTask.Remove(cron.EntryID(job.cronId))
		delete(s, jId)
		return nil
	}
	return errors.New("delete job failed,given job id does not exist")
}
