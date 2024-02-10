package rpc

import (
	"context"
	"fmt"
	"math"
	"net"

	"flamingo.me/dingo"
	"flamingo.me/flamingo/v3/framework/config"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"

	"saitho.me/styx-app/src/lib"
)

type Module struct {
	UnimplementedEventServiceServer

	subscriptions map[string][]string
	logger        flamingo.Logger
	config        *lib.Config
	grpcConfig    config.Map
	eventManager  *lib.ServiceEventManager
}

func (m *Module) Inject(
	logger flamingo.Logger,
	eventManager *lib.ServiceEventManager,
	config *lib.Config,
	cfg *struct {
		GrpcConfig config.Map `inject:"config:styx.grpc"`
	}) {
	m.logger = logger
	m.config = config
	m.grpcConfig = cfg.GrpcConfig
	m.eventManager = eventManager
}

// Configure is the default Method a Module needs to implement
func (m *Module) Configure(injector *dingo.Injector) {
	port, ok := m.grpcConfig.Get("port")
	if !ok {
		m.logger.Fatal("failed to get RPC port")
		return
	}
	m.subscriptions = map[string][]string{}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", int(math.Round(port.(float64)))))
	if err != nil {
		m.logger.Fatal(fmt.Sprintf("failed to listen: %v", err))
		return
	}
	s := grpc.NewServer()
	RegisterEventServiceServer(s, m)
	m.logger.Info(fmt.Sprintf("GRPC server listening at %v", lis.Addr()))
	go func() {
		if err := s.Serve(lis); err != nil {
			m.logger.Fatal(fmt.Sprintf("failed to serve: %v", err))
		}
	}()
}

func (m *Module) Emit(ctx context.Context, in *Event) (*Empty, error) {
	return &Empty{}, m.eventManager.Emit(in.Name, in)
}

func (m *Module) Subscribe(ctx context.Context, in *SubscribeEvent) (*Empty, error) {
	p, _ := peer.FromContext(ctx)
	senderAddress := p.Addr.String()
	return &Empty{}, m.eventManager.Subscribe(senderAddress, in.EventName)
}

func (m *Module) Unsubscribe(ctx context.Context, in *SubscribeEvent) (*Empty, error) {
	// Not implemented
	return &Empty{}, nil
}
