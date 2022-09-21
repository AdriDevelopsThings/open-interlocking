package api

import (
	"encoding/json"
	"net/http"

	"github.com/adridevelopsthings/open-interlocking/pkg/authorization"
	"github.com/adridevelopsthings/open-interlocking/pkg/components"
	"github.com/gorilla/mux"
)

func ConnectionDelete(w http.ResponseWriter, req *http.Request) {
	ret := authorization.CheckAuthorization(req.Header.Get("authorization"), "connection/desolve", true)
	if ret != 0 {
		http.Error(w, http.StatusText(ret), ret)
		return
	}
	params := mux.Vars(req)
	signal1 := components.GetSignalByName(params["signal1"])
	signal2 := components.GetSignalByName(params["signal2"])
	connection := components.GetConnectionBySignals(signal1, signal2)
	if connection == nil {
		http.Error(w, ObjectNotFoundError.Name, ObjectNotFoundError.Http_error)
		return
	}
	if connection.State == components.ConnectionNotSet || connection.State == components.ConnectionDesolvingSignals {
		http.Error(w, RailroadConnectionWrongStateError.Name, RailroadConnectionWrongStateError.Http_error)
		return
	}
	if !connection.Desolve() {
		http.Error(w, RailroadBlockOccupiedError.Name, RailroadBlockOccupiedError.Http_error)
		return
	}
	json.NewEncoder(w).Encode(connection)
}

func Connection(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	signal1 := components.GetSignalByName(params["signal1"])
	signal2 := components.GetSignalByName(params["signal2"])
	var connection *components.RailroadConnection

	if req.Method == "GET" {
		connection = components.GetConnectionBySignals(signal1, signal2)
		if connection == nil {
			http.Error(w, ObjectNotFoundError.Name, ObjectNotFoundError.Http_error)
			return
		}
	}

	if req.Method == "POST" {
		ret := authorization.CheckAuthorization(req.Header.Get("authorization"), "connection/set", true)
		if ret != 0 {
			http.Error(w, http.StatusText(ret), ret)
			return
		}
		signal1 := components.GetSignalByName(params["signal1"])
		signal2 := components.GetSignalByName(params["signal2"])
		if signal1 == nil || signal2 == nil {
			http.Error(w, ObjectNotFoundError.Name, ObjectNotFoundError.Http_error)
			return
		}
		connection = components.GenerateConnection(signal1, signal2)
		if connection == nil {
			http.Error(w, RailroadConnectionApplyingError.Name, RailroadConnectionApplyingError.Http_error)
			return
		}
	}

	if req.Method == "DELETE" {
		ConnectionDelete(w, req)
		return
	}

	json.NewEncoder(w).Encode(connection)
}
