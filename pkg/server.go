package pkg

import (
	"fmt"
	"net/http"

	"github.com/adridevelopsthings/open-interlocking/pkg/api"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func RunServer(host string) {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/connection/{signal1:S\\d+}/{signal2:S\\d+|-}", api.Connection).Methods("GET", "POST", "DELETE")
	rtr.HandleFunc("/{kind:[a-z_]+}/{name:[a-zA-Z\\d]+}", api.GetState).Methods("GET", "POST")
	rtr.HandleFunc("/block/occupy/{from:[WB]\\d+|-}/{to:[WB]\\d+|-}/{action:join|leave}", api.UpdateBlockOccupied).Methods("POST")
	rtr.HandleFunc("/fullcomponents", api.FullComponents).Methods("GET")
	http.Handle("/", cors.AllowAll().Handler(rtr))
	fmt.Printf("Started http server: %q\n", host)
	err := http.ListenAndServe(host, nil)
	if err != nil {
		fmt.Printf("Error while starting http server: %v\n", err)
	}
}
