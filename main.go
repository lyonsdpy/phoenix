// @Author : DAIPENGYUAN
// @File : main
// @Time : 2020/9/30 15:42
// @Description :
// TODO:规范化JSON和XML返回值

package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"phoenix/conf"
	"phoenix/handlers"
	"phoenix/utils/log"
	"runtime"
	"syscall"
	"time"

	"github.com/kwanhur/go-svc/svc"
)

var (
	flagSet = flag.NewFlagSet("phoenix", flag.ExitOnError)
	// 运行参数
	port     = flagSet.Int("port", 8080, "服务端口,默认8080")
	loglevel = flagSet.String("level", "info", "服务级别,默认info")
	logpath  = flagSet.String("log", "./phoenix.log", "日志文件目录结构,默认值为./phoenix.log")
	pidpath  = flagSet.String("pid", "./phoenix.pid", "pid文件保存路径")
	version  = flagSet.Bool("version", false, "查看版本信息")
	// 编译参数,给定的为默认参数,具体在makefile中赋值
	Version   = "0.0.0"
	CommitID  = "0000000"
	BuildTime = time.Time{}.Format("2006-01-02 15:04")
)

type program struct {
	env    svc.Environment
	worker *handlers.Workerd
}

func (p *program) Init(env svc.Environment) error {
	p.env = env
	if env.IsWindowsService() {
		dir := filepath.Dir(os.Args[0])
		return os.Chdir(dir)
	}

	if !flagSet.Parsed() {
		flagSet.Parse(os.Args[1:])
	}
	conf.Port = *port
	conf.Pidpath = *pidpath
	conf.Logpath = *logpath
	conf.LogLevel = *loglevel
	if *version {
		fmt.Printf("software version:%s\n", Version)
		fmt.Printf("commit id:%s\n", CommitID)
		fmt.Printf("built by %s %s/%s at %s\n", runtime.Version(), runtime.GOOS, runtime.GOARCH, BuildTime)
		os.Exit(2)
	}

	daemon := &handlers.Workerd{}
	if err := daemon.Init(); err != nil {
		return err
	}
	p.worker = daemon
	return nil
}

func (p *program) Start() error {
	log.Logger.Infof("starting")
	if p.worker != nil {
		p.worker.Main()
	} else {
		log.Logger.Warn("worker is nil")
	}
	return nil
}

func (p *program) Stop() error {
	log.Logger.Warn("stopping")
	if p.worker != nil {
		if err := p.worker.Exit(); err != nil {
			return err
		}
	} else {
		log.Logger.Warn("worker is nil")
	}
	return nil
}

func (p *program) Reload(signal os.Signal) {
	log.Logger.Infof("got signal:%s", signal.String())
	switch signal {
	case syscall.SIGHUP:
		// reload config
		p.worker.Reload()
	}
}

// @title phoenix采集agent
// @version 1.0.0
// @description 可用于监控采集,网管信息获取的统一入口.
// @description 后续计划支持prometheus，grpc，etcd.

// @contact.name 代澎源
// @contact.url https://github.com/lyonsdpy/phoenix
// @contact.email lyonsdpy@163.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 127.0.0.1:8080
// @BasePath /v1
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	pg := &program{}
	svc.Notify(syscall.SIGHUP, pg.Reload)
	if err := svc.Run(pg, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT); err != nil {
		fmt.Printf("run exit with err:%s", err)
		os.Exit(2)
	} else {
		log.Logger.Info("bye :-)")
	}
}
