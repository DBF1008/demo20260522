package config

import (
	"flag"
	"log"

	"gitee.com/cristiane/micro-mall-api/config/setting"
	"gitee.com/cristiane/micro-mall-api/vars"
	"github.com/go-ini/ini"
)

const (
	ConfFileName            = "./etc/app.ini"
	SectionServer           = "web-server"
	SectionLogger           = "web-logger"
	SectionJwt              = "web-jwt"
	SectionRateLimit        = "web-rate-limit"
	SectionTransactionSeata = "transaction-seata"
)

var (
	cfg      *ini.File
	flagConf = flag.String("web_conf_file", "", "Set app config.")
)

var defaultSectionLoaders = map[string]func(){
	SectionServer: func() {
		vars.ServerSetting = new(setting.ServerSettingS)
		MapConfig(SectionServer, vars.ServerSetting)
	},
	SectionRateLimit: func() {
		vars.RateLimitSetting = new(setting.RateLimitSettingS)
		MapConfig(SectionRateLimit, vars.RateLimitSetting)
	},
	SectionLogger: func() {
		vars.LoggerSetting = new(setting.LoggerSettingS)
		MapConfig(SectionLogger, vars.LoggerSetting)
	},
	SectionJwt: func() {
		vars.JwtSetting = new(setting.JwtSettingS)
		MapConfig(SectionJwt, vars.JwtSetting)
	},
	SectionTransactionSeata: func() {
		vars.TransactionSeataSetting = new(setting.TransactionSeataSettingS)
		MapConfig(SectionTransactionSeata, vars.TransactionSeataSetting)
	},
}

func LoadDefaultConfig(application *vars.Application) error {
	flag.Parse()

	loadedCfg, err := ini.Load(configFilePath())
	if err != nil {
		return err
	}
	cfg = loadedCfg

	for _, sectionName := range cfg.SectionStrings() {
		if loader, ok := defaultSectionLoaders[sectionName]; ok {
			loader()
		}
	}
	return nil
}

func configFilePath() string {
	if *flagConf != "" {
		return *flagConf
	}
	return ConfFileName
}

func MapConfig(section string, v interface{}) {
	log.Printf("[info] Load default config %s", section)
	sec, err := cfg.GetSection(section)
	if err != nil {
		log.Fatalf("[err] Fail to parse '%s': %v", section, err)
	}
	if err = sec.MapTo(v); err != nil {
		log.Fatalf("[err] %s section map to setting err: %v", section, err)
	}
}
