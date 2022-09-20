package pkg

import (
	"fmt"
	"net/http"
	"os"

	"github.com/adridevelopsthings/open-interlocking/pkg/api"
	"github.com/gorilla/mux"
)

func RunServer() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/connection/{signal1:[a-zA-Z\\d]+}/{signal2:[a-zA-Z\\d]+}", api.Connection).Methods("GET", "POST", "DELETE")
	rtr.HandleFunc("/connection/{signal1:[a-zA-Z\\d]+}/{signal2:[a-zA-Z\\d]+}/delete", api.ConnectionDelete).Methods("POST")
	rtr.HandleFunc("/{kind:[a-z_]+}/{name:[a-zA-Z\\d]+}", api.GetState).Methods("GET", "POST")
	http.Handle("/", rtr)
	host := os.Getenv("OPEN_INTERLOCKING_HOST")
	if host == "" {
		host = ":8000"
	}
	fmt.Printf("Started http server: %q\n", host)
	err := http.ListenAndServe(host, nil)
	if err != nil {
		fmt.Printf("Error while starting http server: %v\n", err)
	}
}
