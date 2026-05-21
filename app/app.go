package app

import (
	"flag"
	"fmt"
	"os"
	"sync"
	"time"

	"gitee.com/cristiane/micro-mall-api/config/setting"
	"gitee.com/cristiane/micro-mall-api/internal/config"
	"gitee.com/cristiane/micro-mall-api/internal/logging"
	"gitee.com/cristiane/micro-mall-api/internal/util/startup"
	varsInternal "gitee.com/cristiane/micro-mall-api/internal/vars"
	"gitee.com/cristiane/micro-mall-api/vars"
	"gitee.com/kelvins-io/common/log"
)

var (
	flagLoggerLevel = flag.String("logger_level", "", "set logger level eg: debug,warn,error,info")
	flagLoggerPath  = flag.String("logger_path", "", "set logger root path eg: /tmp/kelvins-app")
	flagEnv         = flag.String("env", "", "set exec environment eg: dev,test,prod")
)

func initApplication(application *vars.Application) error {
	flag.Parse()

	application.LoggerPath = resolveStringSetting(
		config.DefaultLoggerRootPath,
		application.LoggerPath,
		loggerRootPath(),
		*flagLoggerPath,
	)
	application.LoggerLevel = resolveStringSetting(
		config.DefaultLoggerLevel,
		application.LoggerLevel,
		loggerLevel(),
		*flagLoggerLevel,
	)
	application.Environment = resolveStringSetting(
		config.DefaultEnvironmentRelease,
		application.Environment,
		serverEnvironment(),
		*flagEnv,
	)

	vars.LoggerLevel = application.LoggerLevel
	vars.Environment = application.Environment

	if vars.ServerSetting == nil {
		vars.ServerSetting = new(setting.ServerSettingS)
	}

	if err := log.InitGlobalConfig(application.LoggerPath, application.LoggerLevel, application.Name); err != nil {
		return fmt.Errorf("log.InitGlobalConfig: %v", err)
	}

	return nil
}

func resolveStringSetting(defaultValue string, values ...string) string {
	result := defaultValue
	for _, value := range values {
		if value != "" {
			result = value
		}
	}
	return result
}

func loggerRootPath() string {
	if vars.LoggerSetting == nil {
		return ""
	}
	return vars.LoggerSetting.RootPath
}

func loggerLevel() string {
	if vars.LoggerSetting == nil {
		return ""
	}
	return vars.LoggerSetting.Level
}

func serverEnvironment() string {
	if vars.ServerSetting == nil {
		return ""
	}
	return vars.ServerSetting.Environment
}

var appCloseChOnce sync.Once

func appShutdown(application *vars.Application) error {
	if appCloseCh != nil {
		appCloseChOnce.Do(func() {
			close(appCloseCh)
		})
	}
	if application.StopFunc != nil {
		return application.StopFunc()
	}
	return nil
}

func appPrepareForceExit() {
	if !execStopFunc {
		return
	}
	time.AfterFunc(10*time.Second, func() {
		logging.Info("App server Shutdown timeout, force exit")
		os.Exit(1)
	})
}

func setupCommonVars(application *vars.WEBApplication) error {
	var err error
	vars.ErrorLogger, err = log.GetErrLogger("err")
	if err != nil {
		return err
	}
	varsInternal.ErrorLogger = vars.ErrorLogger

	vars.BusinessLogger, err = log.GetBusinessLogger("business")
	if err != nil {
		return err
	}

	vars.AccessLogger, err = log.GetAccessLogger("access")
	if err != nil {
		return err
	}
	if vars.ServerSetting.PIDFile == "" {
		wd, _ := os.Getwd()
		vars.ServerSetting.PIDFile = fmt.Sprintf("%s/%s.pid", wd, application.Name)
	}
	return nil
}

var execStopFunc bool

var appCloseCh = make(chan struct{})

func startUpControl(pidFile string) (next bool, err error) {
	vars.AppCloseCh = appCloseCh
	next, err = startup.ParseCliCommand(pidFile)
	if next {
		execStopFunc = true
	}
	return
}
