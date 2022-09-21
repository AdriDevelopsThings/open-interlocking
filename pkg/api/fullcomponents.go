package api

import (
	"encoding/json"
	"net/http"

	"github.com/adridevelopsthings/open-interlocking/pkg/components"
)

func FullComponents(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(components.GetFullComponents())
}
