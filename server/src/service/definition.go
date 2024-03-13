package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"flamingo.me/flamingo/v3/framework/flamingo"
)

type Status int

const (
	StatusUnknown         Status = iota
	StatusHttpError              = iota
	StatusReady                  = iota
	StatusInvalidResponse        = iota
)

type StyxService struct {
	ServiceName string
	IpAddress   string
	Port        int
	Version     string
	Config      *StyxServiceConfig

	LastStatus         Status
	LastStatusReadable string
	LastStatusDate     time.Time

	logger flamingo.Logger
}

type StyxServiceConfig struct {
	SubscribedEvents []string
}

func (m *StyxService) CallEndpoint(endpoint string, ignoreStatusCode bool) (*http.Response, error) {
	if endpoint != "status" {
		// perform status ping if needed
		if m.pingStatus() != StatusReady {
			return nil, fmt.Errorf("service \"%s\" status is not ready", m.ServiceName)
		}
	}
	resp, err := http.Get(fmt.Sprintf("http://%s:%d/_styx/%s", m.IpAddress, m.Port, endpoint))
	if err != nil {
		return nil, fmt.Errorf("unable to connect to service \"%s\" endpoint \"%s\"", m.ServiceName, endpoint)
	}
	if !ignoreStatusCode && resp.StatusCode != 200 {
		return nil, fmt.Errorf("service \"%s\" endpoint \"%s\" returned HTTP code %d", m.ServiceName, endpoint, resp.StatusCode)
	}
	return resp, nil
}

type EventRequest struct {
	eventName string
	data      interface{}
}

func (m *StyxService) EmitEvent(eventName string, data interface{}) error {
	if m.pingStatus() != StatusReady {
		return fmt.Errorf("service \"%s\" status is not ready")
	}

	bytesData, err := json.Marshal(EventRequest{eventName: eventName, data: data})
	if err != nil {
		return fmt.Errorf("unable to marshal event data")
	}
	resp, err := http.Post(fmt.Sprintf("http://%s:8844/_styx/event", m.IpAddress), "application/json", bytes.NewBuffer(bytesData))
	if err != nil {
		return fmt.Errorf("unable to connect to events endpoint on service \"%s\"", m.ServiceName)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("events endpoint on service \"%s\" returned HTTP code %d", m.ServiceName, resp.StatusCode)
	}
	return nil
}

func (m *StyxService) Inject(logger flamingo.Logger) {
	m.logger = logger
}

type StatusResponse struct {
	Status string `yaml:"status,omitempty"`
}

type VersionResponse struct {
	Version string
}

func StatusText(status Status) string {
	switch status {
	case StatusUnknown:
		return "unknown"
	case StatusHttpError:
		return "HTTP Error"
	case StatusInvalidResponse:
		return "Invalid Response"
	case StatusReady:
		return "ready"
	}
	return "???"
}

func (m *StyxService) Init() error {
	m.logger.Info("Initializing service \"" + m.ServiceName + "\"")
	if m.Port == 0 {
		m.Port = 8844 // set default port
	}
	resp, err := m.CallEndpoint("init", false)
	if err != nil {
		m.updateStatus(StatusHttpError)
		return err
	}
	defer resp.Body.Close()
	serviceConfig := &StyxServiceConfig{}
	if err = json.NewDecoder(resp.Body).Decode(serviceConfig); err != nil {
		return err
	}
	m.Config = serviceConfig

	version, err := m.GetVersion()
	if err != nil {
		return err
	}
	if version != "" {
		// Cache settings by version
		if err = StoreServiceConfig(m.ServiceName, version, *serviceConfig); err != nil {
			return err
		}
	}

	return nil
}

func (m *StyxService) updateStatus(status Status) {
	m.LastStatusDate = time.Now()
	m.LastStatus = status
	m.LastStatusReadable = StatusText(status)
}

func (m *StyxService) GetVersion() (string, error) {
	resp, err := m.CallEndpoint("version", true)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 404 {
		return "", nil
	}
	target := &VersionResponse{}
	if err = json.NewDecoder(resp.Body).Decode(target); err != nil {
		return "", err
	}
	m.Version = target.Version
	return target.Version, nil
}

func (m *StyxService) pingStatus() Status {
	if m.LastStatusDate.Compare(time.Now().Add(time.Minute*-15)) == -1 {
		m.logger.Debugf("Getting status for service \"" + m.ServiceName + "\"")
		resp, err := m.CallEndpoint("status", false)
		if err != nil {
			m.updateStatus(StatusHttpError)
		} else {
			defer resp.Body.Close()

			target := &StatusResponse{}
			if err = json.NewDecoder(resp.Body).Decode(target); err != nil {
				m.updateStatus(StatusInvalidResponse)
			} else if target.Status == "ready" {
				m.updateStatus(StatusReady)
			} else {
				m.updateStatus(StatusUnknown)
			}
		}
	}
	return m.LastStatus
}
