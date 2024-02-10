package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"flamingo.me/flamingo/v3/framework/flamingo"
)

type ModuleStatus int

const ( // iota is reset to 0
	MODULESTATUS_UNKNOWN         ModuleStatus = iota // c0 == 0
	MODULESTATUS_HTTP_ERROR                   = iota // c1 == 1
	MODULESTATUS_READY                        = iota // c2 == 2
	MODULESTATUS_INVALIDRESPONSE              = iota // c2 == 2
)

type StyxModule struct {
	ServiceName string

	LastStatus         ModuleStatus
	LastStatusReadable string
	LastStatusDate     time.Time

	eventManager *ModuleEventManager
	logger       flamingo.Logger
}

func (m *StyxModule) Inject(logger flamingo.Logger, eventManager *ModuleEventManager) {
	m.eventManager = eventManager
	m.logger = logger
}

type StatusResponse struct {
	Status string
}

type InitResponse struct {
	SubscribedEvents []string
}

func StatusText(status ModuleStatus) string {
	switch status {
	case MODULESTATUS_UNKNOWN:
		return "unknown"
	case MODULESTATUS_HTTP_ERROR:
		return "HTTP Error"
	case MODULESTATUS_INVALIDRESPONSE:
		return "Invalid Response"
	case MODULESTATUS_READY:
		return "ready"
	}
	return "???"
}

func (m *StyxModule) Init() error {
	m.logger.Info("Initializing service \"" + m.ServiceName + "\"")
	resp, err := http.Get("http://" + m.ServiceName + ":8844/_styx/init")
	if err != nil {
		m.updateStatus(MODULESTATUS_HTTP_ERROR)
		return fmt.Errorf("unable to connect to service init endpoint")
	}
	defer resp.Body.Close()

	target := &InitResponse{}
	if err = json.NewDecoder(resp.Body).Decode(target); err != nil {
		m.updateStatus(MODULESTATUS_INVALIDRESPONSE)
		return fmt.Errorf("unable to decode init response from init endpoint")
	}

	for _, eventName := range target.SubscribedEvents {
		// TODO: validate error
		_ = m.eventManager.Subscribe(m.ServiceName, eventName)
	}

	// Trigger status update
	m.Status()

	return nil
}

func (m *StyxModule) updateStatus(status ModuleStatus) {
	m.LastStatusDate = time.Now()
	m.LastStatus = status
	m.LastStatusReadable = StatusText(status)
}

func (m *StyxModule) Status() ModuleStatus {
	if m.LastStatusDate.Compare(time.Now().Add(time.Minute*-15)) == -1 {
		m.logger.Debugf("Getting status for service " + m.ServiceName)
		resp, err := http.Get("http://" + m.ServiceName + ":8844/_styx/status")
		if err != nil {
			m.updateStatus(MODULESTATUS_HTTP_ERROR)
		} else {
			defer resp.Body.Close()

			target := &StatusResponse{}
			if err = json.NewDecoder(resp.Body).Decode(target); err != nil {
				m.updateStatus(MODULESTATUS_INVALIDRESPONSE)
			} else if target.Status == "ready" {
				m.updateStatus(MODULESTATUS_READY)
			} else {
				m.updateStatus(MODULESTATUS_UNKNOWN)
			}
		}
	}
	return m.LastStatus
}
