package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/sdboyer/mock_druplatform_api/acquia"
	"net/http"
)

var fml = fmt.Println

func main() {
	// Spawn an Acquia API server on the default port
	aq := acquia.NewServerInstance()

	an := negroni.New()
	an.UseHandler(acquia.NewRouter(aq))

	go an.Run(":10234")

	router := mux.NewRouter()

	router.HandleFunc("/", ServerListHandler)

	n := negroni.New()
	n.UseHandler(router)
	n.Run(":10233")
}

func ServerListHandler(w http.ResponseWriter, r *http.Request) {

}
