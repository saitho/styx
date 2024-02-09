package lib

import (
	"encoding/json"
	"net/http"
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
}

type StatusResponse struct {
	Status string
}

func (m *StyxModule) Status() ModuleStatus {
	// http://localhost:8844/_styx/status
	resp, err := http.Get("http://" + m.ServiceName + ":8844/_styx/status")
	if err != nil {
		return MODULESTATUS_HTTP_ERROR
	}
	defer resp.Body.Close()

	target := &StatusResponse{}
	if err = json.NewDecoder(resp.Body).Decode(target); err != nil {
		return MODULESTATUS_INVALIDRESPONSE
	}
	if target.Status == "ready" {
		return MODULESTATUS_READY
	}
	return MODULESTATUS_UNKNOWN
}

func (m *StyxModule) StatusText() string {
	switch m.Status() {
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
