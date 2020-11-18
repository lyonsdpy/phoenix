// @Author : DAIPENGYUAN
// @File : worker_cron
// @Time : 2020/11/9 17:08 
// @Description : Crontab启动项

package handlers

import (
	"phoenix/conf"
	"phoenix/task"
	"phoenix/utils/log"
	"time"
)

func RunCronTask() {
	log.Logger.Info("Cron Task Started")
	task.CronTask.Start()
}

func ExitCronTask() {
	log.Logger.Warn("cron task shutdown now ...")
	ctx := task.CronTask.Stop()
	select {
	case <-ctx.Done():
		log.Logger.Warnf("cron task shutdown success")
	case <-time.After(conf.ShutdownTimeout):
		log.Logger.Warnf("cron task shutdown timeout of %s", conf.ShutdownTimeout.String())
	}
	log.Logger.Warn("cron task exiting")
}
