// @Author : DAIPENGYUAN
// @File : worker
// @Time : 2020/11/9 16:23 
// @Description : svc任务启动器,加载各任务

package handlers

import (
	"io/ioutil"
	"os"
	"phoenix/conf"
	"phoenix/utils/log"
	"strconv"
)

type Workerd struct {
	hasPid bool
}

func (w *Workerd) Init() error {
	// 初始化PID文件
	if !w.hasPid {
		pid := strconv.Itoa(os.Getpid())
		if err := ioutil.WriteFile(conf.Pidpath, []byte(pid), 644); err != nil {
			return err
		}
		w.hasPid = true
	}
	return nil
}

// Reload:用于在线重加载，通常用于热加载配置文件等内容
// 本项目现阶段不涉及配置文件的热加载,后续涉及到需要热重载的内容在此处引入
func (w *Workerd) Reload() {
	log.Logger.Info("reload starting")
	if err := w.Init(); err != nil {
		log.Logger.Error(err)
	}
	log.Logger.Info("reload end")
}

func (w *Workerd) Main() {
	go log.ResetLogger(conf.LogLevel, conf.Logpath)
	go RunGin()
	go RunCronTask()
}

func (w *Workerd) Exit() error {
	var err error

	ExitGin()
	ExitCronTask()

	if err = os.Remove(conf.Pidpath); err != nil {
		return err
	}
	return nil
}
