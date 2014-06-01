package main

import (
  "github.com/codegangsta/negroni"
  "github.com/gorilla/mux"
  "github.com/sdboyer/mock_druplatform_api/acquia"
  "net/http"
  "fmt"
)

var fml = fmt.Println


func main() {
  // Spawn an Acquia API server on the default port
  go negroni.New().UseHandler(acquia.NewRouter()).Run(":10234")

  router := mux.NewRouter()

  router.HandleFunc("/", ServerListHandler)

  n := negroni.New()
  n.UseHandler(router)
  n.Run(":10233")
}

func ServerListHandler(w http.ResponseWriter, r *http.Request) {

}
