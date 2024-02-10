package service

import (
	"context"
	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"net"

	"saitho.me/styx-app/src/proto"
)

type GrpcConfig struct {
	Port int
}

type RpcModule struct {
	proto.UnimplementedEventServiceServer

	subscriptions  map[string][]string
	logger         flamingo.Logger
	serviceManager *ServiceManager
	grpcConfig     *GrpcConfig
	eventManager   *ServiceEventManager
}

func (m *RpcModule) Inject(
	logger flamingo.Logger,
	eventManager *ServiceEventManager,
	serviceManager *ServiceManager,
	cfg *struct {
		GrpcConfig config.Map `inject:"config:styx.grpc"`
	}) {
	m.logger = logger
	m.serviceManager = serviceManager
	m.grpcConfig = &GrpcConfig{}
	if err := cfg.GrpcConfig.MapInto(m.grpcConfig); err != nil {
		m.logger.Fatal(err)
	}
	m.eventManager = eventManager
}

// Configure is the default Method a Module needs to implement
func (m *RpcModule) Configure(injector *dingo.Injector) {
	m.subscriptions = map[string][]string{}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", m.grpcConfig.Port))
	if err != nil {
		m.logger.Fatal(fmt.Sprintf("failed to listen: %v", err))
		return
	}
	s := grpc.NewServer()
	proto.RegisterEventServiceServer(s, m)
	m.logger.Info(fmt.Sprintf("GRPC server listening at %v", lis.Addr()))
	go func() {
		if err := s.Serve(lis); err != nil {
			m.logger.Fatal(fmt.Sprintf("failed to serve: %v", err))
		}
	}()
}

func (m *RpcModule) Emit(ctx context.Context, in *proto.Event) (*proto.Empty, error) {
	return &proto.Empty{}, m.eventManager.Emit(in.Name, in)
}

func (m *RpcModule) Subscribe(ctx context.Context, in *proto.SubscribeEvent) (*proto.Empty, error) {
	p, _ := peer.FromContext(ctx)
	senderAddress := p.Addr.String()
	return &proto.Empty{}, m.eventManager.Subscribe(senderAddress, in.EventName)
}

func (m *RpcModule) Unsubscribe(ctx context.Context, in *proto.SubscribeEvent) (*proto.Empty, error) {
	// Not implemented
	return &proto.Empty{}, nil
}
