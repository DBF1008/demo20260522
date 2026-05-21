package startup

import (
	"gitee.com/cristiane/micro-mall-api/config"
	"gitee.com/cristiane/micro-mall-api/config/setting"
	"gitee.com/cristiane/micro-mall-api/vars"
)

const (
	SectionEmailConfig = "email-config"
	SectionG2Cache     = "micro-mall-g2cache"
)

func LoadConfig() error {
	vars.G2CacheSetting = new(setting.G2CacheSettingS)
	vars.EmailConfigSetting = new(vars.EmailConfigSettingS)

	config.MapConfig(SectionG2Cache, vars.G2CacheSetting)
	config.MapConfig(SectionEmailConfig, vars.EmailConfigSetting)

	return nil
}
