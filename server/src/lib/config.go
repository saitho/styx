package lib

import (
	"strings"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

type Config struct {
	Modules map[string]*StyxModule
}

type ConfigModule struct {
	injector *dingo.Injector
}

func GetConfig(
	logger flamingo.Logger,
	eventManager *ModuleEventManager,
	cfg *struct {
		CompleteConfig config.Map `inject:"config:styx"`
	},
) *Config {
	c := &Config{}
	c.Modules = map[string]*StyxModule{}
	if moduleString, ok := cfg.CompleteConfig.Get("modules"); moduleString != "" && ok {
		for _, moduleName := range strings.Split(moduleString.(string), ",") {
			module := &StyxModule{ServiceName: moduleName}
			module.Inject(logger, eventManager)
			if err := module.Init(); err != nil {
				logger.Error(err)
			}
			c.Modules[moduleName] = module
		}
	}
	return c
}

func (c *ConfigModule) Configure(injector *dingo.Injector) {
	injector.Bind(new(Config)).ToProvider(GetConfig).In(dingo.Singleton)
}
