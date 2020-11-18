// @Author : DAIPENGYUAN
// @File : worker
// @Time : 2020/11/9 16:23
// @Description : gin任务启动

package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"phoenix/conf"
	"phoenix/handlers/rest"
	"phoenix/utils/log"
	"strconv"
)

var server *http.Server

func RunGin() {
	gin.SetMode(gin.ReleaseMode)

	engine := gin.New()
	engine.Use(gin.Recovery())
	rest.SetupRoutes(engine)

	server = &http.Server{
		Addr:    ":" + strconv.Itoa(conf.Port),
		Handler: engine,
	}
	log.ResetLogger(conf.LogLevel, conf.Logpath)
	log.Logger.Infof("Rest Server Started at Port:%d, loglevel=%s", conf.Port, conf.LogLevel)

	if err := server.ListenAndServe(); err != nil {
		log.Logger.Info(err.Error())
	}
}

func ExitGin() {
	if server != nil {
		log.Logger.Warn("rest server shutdown now ...")
		ctx, cancel := context.WithTimeout(context.Background(), conf.ShutdownTimeout)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Logger.Error(err)
			return
		}
		select {
		case <-ctx.Done():
			log.Logger.Warnf("rest server shutdown timeout of %s", conf.ShutdownTimeout.String())
		}
		log.Logger.Warn("rest server exiting")
	}
	return
}
