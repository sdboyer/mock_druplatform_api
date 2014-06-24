package main

import (
	"fmt"

	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

var fml = fmt.Println

func main() {
	setUpMainHttpd()
}

func setUpMainHttpd() {
	router := mux.NewRouter()
	router.HandleFunc("/", hhListServerTypes).Methods("GET")
	router.HandleFunc("/acquia", hhListAcquiaServers).Methods("GET")
	router.HandleFunc("/acquia", hhCreateAcquiaServer).Methods("POST")

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.Use(negroni.NewLogger())
	n.UseHandler(router)
	n.Run(":10233")
}
