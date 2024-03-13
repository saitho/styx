package service

import "fmt"

// key = string (service:version)
var serviceConfigCache = map[string]StyxServiceConfig{}

func FetchServiceConfig(service string, version string) (StyxServiceConfig, error) {
	val, ok := serviceConfigCache[service+":"+version]
	if !ok {
		return val, fmt.Errorf("no service config found for service \"%s\" and version \"%s\"", service, version)
	}
	return val, nil
}

func StoreServiceConfig(service string, version string, config StyxServiceConfig) error {
	serviceConfigCache[service+":"+version] = config
	return nil
}
