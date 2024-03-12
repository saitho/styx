package service

import (
	"flamingo.me/flamingo/v3/framework/flamingo"
	"fmt"
	"slices"
)

type ServiceEventManager struct {
	subscriptions map[string][]string
	logger        flamingo.Logger
}

func (m *ServiceEventManager) Inject(
	logger flamingo.Logger) {
	m.logger = logger
	m.subscriptions = map[string][]string{}
}

func (m *ServiceEventManager) Emit(eventName string, event interface{}) error {
	for _, subscriberHost := range m.subscriptions[eventName] {
		m.logger.Info(fmt.Sprintf("Emitting event %s to %s", eventName, subscriberHost))
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
