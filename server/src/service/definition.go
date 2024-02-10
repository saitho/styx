package service

import (
	"context"
	"encoding/json"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
	"saitho.me/styx-app/src/proto"
	"time"
)

type ServiceStatus int

const (
	ServiceStatusUnknown         ServiceStatus = iota
	ServiceStatusHttpError                     = iota
	ServiceStatusReady                         = iota
	ServiceStatusInvalidResponse               = iota
	ServiceStatusNoRpcConnection               = iota
	ServiceStatusNoRpcResponse                 = iota
)

type StyxService struct {
	ServiceName string
	GrpcHost    string

	LastStatus            ServiceStatus
	LastStatusReadable    string
	LastStatusDate        time.Time
	LastRpcStatus         ServiceStatus
	LastRpcStatusReadable string
	LastRpcStatusDate     time.Time

	eventManager *ServiceEventManager
	logger       flamingo.Logger
}

func (m *StyxService) Inject(logger flamingo.Logger, eventManager *ServiceEventManager) {
	m.eventManager = eventManager
	m.logger = logger
}

type StatusResponse struct {
	Status string
}

type InitResponse struct {
	SubscribedEvents []string
}

func StatusText(status ServiceStatus) string {
	switch status {
	case ServiceStatusUnknown:
		return "unknown"
	case ServiceStatusHttpError:
		return "HTTP Error"
	case ServiceStatusInvalidResponse:
		return "Invalid Response"
	case ServiceStatusReady:
		return "ready"
	case ServiceStatusNoRpcConnection:
		return "No RPC Connection"
	case ServiceStatusNoRpcResponse:
		return "No RPC Response"
	}
	return "???"
}

func (m *StyxService) Init() error {
	m.logger.Info("Initializing service \"" + m.ServiceName + "\"")
	resp, err := http.Get("http://" + m.ServiceName + ":8844/_styx/init")
	if err != nil {
		m.updateStatus(ServiceStatusHttpError)
		return fmt.Errorf("unable to connect to service init endpoint")
	}
	defer resp.Body.Close()

	target := &InitResponse{}
	if err = json.NewDecoder(resp.Body).Decode(target); err != nil {
		m.updateStatus(ServiceStatusInvalidResponse)
		return fmt.Errorf("unable to decode init response from init endpoint")
	}

	for _, eventName := range target.SubscribedEvents {
		// TODO: validate error
		_ = m.eventManager.Subscribe(m.GrpcHost, eventName)
	}

	// Trigger status update
	m.Status()

	return nil
}

func (m *StyxService) updateStatus(status ServiceStatus) {
	m.LastStatusDate = time.Now()
	m.LastStatus = status
	m.LastStatusReadable = StatusText(status)
}

func (m *StyxService) updateRpcStatus(status ServiceStatus) {
	m.LastRpcStatusDate = time.Now()
	m.LastRpcStatus = status
	m.LastRpcStatusReadable = StatusText(status)
}

func (m *StyxService) Status() ServiceStatus {
	if m.LastStatusDate.Compare(time.Now().Add(time.Minute*-15)) == -1 {
		m.logger.Debugf("Getting status for service \"" + m.ServiceName + "\"")
		resp, err := http.Get("http://" + m.ServiceName + ":8844/_styx/status")
		if err != nil {
			m.updateStatus(ServiceStatusHttpError)
		} else {
			defer resp.Body.Close()

			target := &StatusResponse{}
			if err = json.NewDecoder(resp.Body).Decode(target); err != nil {
				m.updateStatus(ServiceStatusInvalidResponse)
			} else if target.Status == "ready" {
				m.updateStatus(ServiceStatusReady)
			} else {
				m.updateStatus(ServiceStatusUnknown)
			}
		}

		m.logger.Debugf("Getting RPC status for service \"" + m.ServiceName + "\"")
		conn, err := grpc.Dial(m.GrpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			m.updateRpcStatus(ServiceStatusNoRpcConnection)
		} else {
			defer conn.Close()
			// Contact the server and print out its response.
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			c := proto.NewServiceProviderClient(conn)
			if _, err := c.Ping(ctx, &proto.Empty{}); err != nil {
				m.updateRpcStatus(ServiceStatusNoRpcResponse)
			}
		}
	}
	return m.LastStatus
}
