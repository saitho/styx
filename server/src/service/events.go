package service

import (
	"context"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"saitho.me/styx-app/src/proto"
	"slices"
	"time"
)

type ServiceEventManager struct {
	subscriptions map[string][]string
	logger        flamingo.Logger
	Config        *GrpcConfig
}

func (m *ServiceEventManager) Inject(
	logger flamingo.Logger,
	cfg *struct {
		GrpcConfig config.Map `inject:"config:styx.grpc"`
	}) {
	m.logger = logger
	m.Config = &GrpcConfig{}
	if err := cfg.GrpcConfig.MapInto(m.Config); err != nil {
		m.logger.Fatal(err)
	}
	m.subscriptions = map[string][]string{}
}

func (m *ServiceEventManager) Emit(eventName string, event interface{}) error {
	for _, subscriberHost := range m.subscriptions[eventName] {
		m.logger.Info(fmt.Sprintf("Emitting event %s to %s", eventName, subscriberHost))

		// Set up a connection to the server.
		conn, err := grpc.Dial(subscriberHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			m.logger.Error(err)
			return err
		}
		defer conn.Close()
		client := proto.NewEventServiceClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, err = client.Emit(ctx, event.(*proto.Event))
		if err != nil {
			m.logger.Error(err)
			return err
		}
	}
	return nil
}

func (m *ServiceEventManager) Subscribe(senderAddress string, eventName string) error {
	if !slices.Contains(m.subscriptions[eventName], senderAddress) {
		m.logger.Info(fmt.Sprintf("Subscribing %s to %s", senderAddress, eventName))
		m.subscriptions[eventName] = append(m.subscriptions[eventName], senderAddress)
	}
	m.Emit(eventName, &proto.Event{Name: eventName, Data: "foo=bar"})
	return nil
}
