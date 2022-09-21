package api

import (
	"net/http"
	"strings"

	"github.com/adridevelopsthings/open-interlocking/pkg/authorization"
	"github.com/adridevelopsthings/open-interlocking/pkg/components"
	"github.com/gorilla/mux"
)

func UpdateBlockOccupied(w http.ResponseWriter, req *http.Request) {
	ret := authorization.CheckAuthorization(req.Header.Get("authorization"), "occupy", true)
	if ret != 0 {
		http.Error(w, http.StatusText(ret), ret)
		return
	}
	params := mux.Vars(req)

	var from_block *components.Block
	var from_switch *components.RailroadSwitch
	var to_block *components.Block
	var to_switch *components.RailroadSwitch

	from := params["from"]
	to := params["to"]

	if strings.HasPrefix(from, "B") {
		from_block = components.GetBlockByName(from)
	} else if strings.HasPrefix(from, "W") {
		from_switch = components.GetSwitchByName(from)
	}

	if strings.HasPrefix(to, "B") {
		to_block = components.GetBlockByName(to)
	} else if strings.HasPrefix(to, "W") {
		to_switch = components.GetSwitchByName(to)
	}

	if (from_block == nil && from_switch == nil) || (to_block == nil && to_switch == nil) {
		http.Error(w, ObjectNotFoundError.Name, ObjectNotFoundError.Http_error)
		return
	}
	components.OccupyBlock(from_block, from_switch, to_block, to_switch)
	w.Write([]byte("success"))
}
