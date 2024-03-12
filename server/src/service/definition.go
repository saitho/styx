package service

import (
	"encoding/json"
	"flamingo.me/flamingo/v3/framework/flamingo"
	"fmt"
	"net/http"
	"time"
)

type ServiceStatus int

const (
	ServiceStatusUnknown         ServiceStatus = iota
	ServiceStatusHttpError                     = iota
	ServiceStatusReady                         = iota
	ServiceStatusInvalidResponse               = iota
)

type StyxService struct {
	ServiceName string

	LastStatus         ServiceStatus
	LastStatusReadable string
	LastStatusDate     time.Time

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

	// Trigger status update
	m.Status()

	return nil
}

func (m *StyxService) updateStatus(status ServiceStatus) {
	m.LastStatusDate = time.Now()
	m.LastStatus = status
	m.LastStatusReadable = StatusText(status)
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
	}
	return m.LastStatus
}
