package lib

import (
	"strings"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

type Config struct {
	Services map[string]*StyxService
}

type ConfigModule struct {
	injector *dingo.Injector
}

func GetConfig(
	logger flamingo.Logger,
	eventManager *ServiceEventManager,
	cfg *struct {
		CompleteConfig config.Map `inject:"config:styx"`
	},
) *Config {
	c := &Config{}
	c.Services = map[string]*StyxService{}
	if serviceString, ok := cfg.CompleteConfig.Get("services"); serviceString != "" && ok {
		for _, serviceName := range strings.Split(serviceString.(string), ",") {
			service := &StyxService{ServiceName: serviceName}
			service.Inject(logger, eventManager)
			if err := service.Init(); err != nil {
				logger.Error(err)
			}
			c.Services[serviceName] = service
		}
	}
	return c
}

func (c *ConfigModule) Configure(injector *dingo.Injector) {
	injector.Bind(new(Config)).ToProvider(GetConfig).In(dingo.Singleton)
}
