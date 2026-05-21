package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"gitee.com/cristiane/micro-mall-api/internal/config"
	"gitee.com/cristiane/micro-mall-api/internal/logging"
	"gitee.com/cristiane/micro-mall-api/pkg/util/kprocess"
	"gitee.com/cristiane/micro-mall-api/vars"
	"github.com/robfig/cron/v3"
	"github.com/seata/seata-go/pkg/client"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const localAddr = "0.0.0.0:"

func RunApplication(application *vars.WEBApplication) {
	if application == nil || application.Application == nil {
		panic("webApplication is nil or application is nil")
	}
	if application.Name == "" {
		logging.Fatal("Application name can't not be empty")
	}

	application.Type = vars.AppTypeWeb
	vars.App = application
	if err := runApp(application); err != nil {
		logging.Infof("App exit over: %v\n", err)
	}
	logging.Info("App exit over")
}

func runApp(webApp *vars.WEBApplication) error {
	if err := loadWebConfig(webApp); err != nil {
		return err
	}
	if err := setupWebApp(webApp); err != nil {
		return err
	}

	next, err := startUpControl(vars.ServerSetting.PIDFile)
	if err != nil {
		return err
	}
	if !next {
		return nil
	}

	startCronTasks(webApp.RegisterTasks)

	return serveHTTP(webApp)
}

func loadWebConfig(webApp *vars.WEBApplication) error {
	if err := config.LoadDefaultConfig(webApp.Application); err != nil {
		return err
	}
	if webApp.LoadConfig != nil {
		if err := webApp.LoadConfig(); err != nil {
			return err
		}
	}
	return nil
}

func setupWebApp(webApp *vars.WEBApplication) error {
	if err := initApplication(webApp.Application); err != nil {
		return err
	}
	if err := setupWEBVars(webApp); err != nil {
		return err
	}
	if webApp.SetupVars != nil {
		if err := webApp.SetupVars(); err != nil {
			return fmt.Errorf("App.SetupVars err: %v", err)
		}
	}
	return nil
}

func startCronTasks(register func() []vars.CronTask) {
	if register == nil {
		return
	}

	cronTasks := register()
	if len(cronTasks) == 0 {
		return
	}

	cn := cron.New(cron.WithSeconds())
	for _, task := range cronTasks {
		if task.TaskFunc == nil {
			continue
		}
		if _, err := cn.AddFunc(task.Cron, task.TaskFunc); err != nil {
			logging.Fatalf("App run cron task err: %v", err)
		}
	}
	cn.Start()
	logging.Info("App run cron task")
}

func serveHTTP(webApp *vars.WEBApplication) error {
	if webApp.RegisterHttpRoute == nil {
		logging.Fatalf("App RegisterHttpRoute nil ??")
	}

	kp := new(kprocess.KProcess)
	addr := resolveListenAddr(webApp)
	network := resolveNetwork()
	ln, err := kp.Listen(network, addr, vars.ServerSetting.PIDFile)
	if err != nil {
		logging.Fatalf("App kprocess listen err: %v", err)
	}
	logging.Infof("server process listen network: %v, addr: %v\n", network, addr)

	server := http.Server{
		Handler:      buildHTTPHandler(webApp),
		ReadTimeout:  time.Duration(vars.ServerSetting.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(vars.ServerSetting.WriteTimeout) * time.Second,
		IdleTimeout:  time.Duration(vars.ServerSetting.IdleTimeout) * time.Second,
	}

	return serveWithProcessControl(kp, &server, ln, webApp.Application)
}

func resolveListenAddr(webApp *vars.WEBApplication) string {
	if webApp.EndPort != 0 {
		return localAddr + strconv.Itoa(webApp.EndPort)
	}
	if vars.ServerSetting != nil && vars.ServerSetting.EndPort != 0 {
		return localAddr + strconv.Itoa(vars.ServerSetting.EndPort)
	}
	return ""
}

func resolveNetwork() string {
	if vars.ServerSetting != nil && vars.ServerSetting.Network != "" {
		return vars.ServerSetting.Network
	}
	return "tcp"
}

func buildHTTPHandler(webApp *vars.WEBApplication) http.Handler {
	ginEngine := webApp.RegisterHttpRoute()
	if vars.ServerSetting != nil && vars.ServerSetting.SupportH2 {
		logging.Info("server http handler support h2")
		return h2c.NewHandler(ginEngine, &http2.Server{IdleTimeout: time.Duration(vars.ServerSetting.IdleTimeout) * time.Second})
	}
	return ginEngine
}

func serveWithProcessControl(kp *kprocess.KProcess, server *http.Server, ln net.Listener, application *vars.Application) error {
	serverCloseCh := make(chan struct{})
	go func() {
		defer close(serverCloseCh)
		if err := server.Serve(ln); err != nil {
			logging.Infof("App run Serve: %v\n", err)
		}
	}()

	select {
	case <-kp.Exit():
	case <-serverCloseCh:
	}

	appPrepareForceExit()
	if err := server.Shutdown(context.Background()); err != nil {
		logging.Infof("App server Shutdown: %v\n", err)
	}
	logging.Info("App server Shutdown ok")

	return appShutdown(application)
}

func setupWEBVars(webApp *vars.WEBApplication) error {
	if err := setupCommonVars(webApp); err != nil {
		return err
	}
	if vars.TransactionSeataSetting != nil && vars.TransactionSeataSetting.Enable {
		if vars.TransactionSeataSetting.ConfFile == "" {
			return fmt.Errorf("transaction seata loaded config file null")
		}
		client.InitPath(vars.TransactionSeataSetting.ConfFile)
	}
	return nil
}
