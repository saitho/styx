package lib

import (
	"fmt"
	"slices"

	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
)

type ServiceEventManager struct {
	subscriptions map[string][]string
	logger        flamingo.Logger
	config        config.Map
}

func (m *ServiceEventManager) Inject(
	logger flamingo.Logger,
	cfg *struct {
		GrpcConfig config.Map `inject:"config:styx.grpc"`
	}) {
	m.logger = logger
	m.config = cfg.GrpcConfig
	m.subscriptions = map[string][]string{}
}

func (m *ServiceEventManager) Emit(eventName string, event interface{}) error {
	for _, subscription := range m.subscriptions[eventName] {
		for _, subscriberAddress := range subscription {
			m.logger.Info(fmt.Sprintf("Emitting event %s to %s", eventName, subscriberAddress))
		}
	}
	return nil
}

func (m *ServiceEventManager) Subscribe(senderAddress string, eventName string) error {
	if !slices.Contains(m.subscriptions[eventName], senderAddress) {
		m.logger.Info(fmt.Sprintf("Subscribing %s to %s", senderAddress, eventName))
		m.subscriptions[eventName] = append(m.subscriptions[eventName], senderAddress)
	}
	return nil
}
