package lib

import (
	"flamingo.me/flamingo/v3/framework/config"

	"strings"
)

type Config struct {
	Modules map[string]StyxModule
}

func (c *Config) Inject(cfg *struct {
	CompleteConfig config.Map `inject:"config:styx"`
}) *Config {
	c.Modules = map[string]StyxModule{}
	if moduleString, ok := cfg.CompleteConfig.Get("modules"); moduleString != "" && ok {
		for _, moduleName := range strings.Split(moduleString.(string), ",") {
			c.Modules[moduleName] = StyxModule{
				ServiceName: moduleName,
			}
		}
	}
	return c
}
