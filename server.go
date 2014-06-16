package main

import (
	"fmt"
	"encoding/json"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/sdboyer/mock_druplatform_api/acquia"
	"net/http"
)

var fml = fmt.Println

type ServerIdentifier interface {
	Port() int
	Name() string
}

var servers = make(map[int]ServerIdentifier, 0)

func main() {
	router := mux.NewRouter()
	router.Headers("Content-Type", "application/json")

	router.HandleFunc("/", ServerListHandler)
	router.HandleFunc("/server/create", hhCreateServer).Methods("POST")

	n := negroni.New()
	n.UseHandler(router)
	n.Run(":10233")
}

type createServerRequest struct {
	port int `json:"port"`
	server_type string  `json:"server_type"`
	version string `json:"version"`
}

func hhCreateServer(w http.ResponseWriter, r *http.Request) {
	j := &createServerRequest{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(j); err != nil {
		panic(err)
	}

	an := negroni.New()
	an.UseHandler(acquia.NewRouter(acquia.NewServerInstance("default")))

	// TODO need to figure out how to get the port ahead of time
	go an.Run(":10234")
}

func ServerListHandler(w http.ResponseWriter, r *http.Request) {

}
