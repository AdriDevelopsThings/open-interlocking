package api

import (
	"net/http"

	"github.com/adridevelopsthings/open-interlocking/pkg/components"
	"github.com/gorilla/mux"
)

func getState(kind string, name string, ack bool) (bool, error) {
	switch kind {
	case "signal":
		signal := components.GetSignalByName(name)
		if signal == nil {
			return false, &ObjectNotFoundError
		}
		if ack {
			signal.Acknowledged = true
		}
		return signal.State, nil
	case "distant_signal":
		distant_signal := components.GetDistantSignalByName(name)
		if distant_signal == nil {
			return false, &ObjectNotFoundError
		}
		if ack {
			distant_signal.Acknowledged = true
		}
		return distant_signal.State, nil
	case "switch":
		rswitch := components.GetSwitchByName(name)
		if rswitch == nil {
			return false, &ObjectNotFoundError
		}
		if ack {
			rswitch.Acknowledged = true
		}
		return rswitch.State, nil
	default:
		return false, &ObjectNotFoundError
	}
}

func GetState(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	kind := params["kind"]
	name := params["name"]
	ack := false
	if req.Method == "POST" {
		ack = true
	}
	state, err := getState(kind, name, ack)

	if err != nil {
		serr := err.(*OpenInterlockingError)
		http.Error(w, serr.Name, serr.Http_error)
		return
	}

	var sstate string
	if state {
		sstate = "true"
	} else {
		sstate = "false"
	}
	w.Write([]byte(sstate))
}
