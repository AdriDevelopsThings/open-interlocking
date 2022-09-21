package api

import (
	"fmt"
	"net/http"
)

type OpenInterlockingError struct {
	Name       string
	Http_error int
}

func (e *OpenInterlockingError) Error() string {
	return fmt.Sprintf("Error while request: %s: %d", e.Name, e.Http_error)
}

var ObjectNotFoundError = OpenInterlockingError{"Object wasn't found.", http.StatusNotFound}
var RailroadConnectionApplyingError = OpenInterlockingError{"RailroadConnectionApplyingError", http.StatusConflict}
var RailroadConnectionWrongStateError = OpenInterlockingError{"RailroadConnectionWrongStateError", http.StatusConflict}
var RailroadBlockOccupiedError = OpenInterlockingError{"RailroadBlockOccupiedError", http.StatusConflict}
