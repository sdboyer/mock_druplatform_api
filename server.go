package main

import (
	//"log"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/codegangsta/negroni"
	"github.com/sdboyer/mock_druplatform_api/acquia"
	"net/http"
	"net"
	"time"
	"strconv"
)

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
//
// (Copied from net/http, b/c they did not deign to export)
type tcpKeepAliveListener struct {
    *net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
    tc, err := ln.AcceptTCP()
    if err != nil {
        return
    }
    tc.SetKeepAlive(true)
    tc.SetKeepAlivePeriod(3 * time.Minute)
    return tc, nil
}

type ServerInstance struct {
	Listener tcpKeepAliveListener
	MockApp
}

type MockApp interface {
	Router() *mux.Router
	Version() string
	Serve(net.Listener)
}

var servers = make([]ServerInstance, 0)

type createServerRequest struct {
	ServerType string `json:"server_type"`
	Version string `json:"version"`
}

type createServerResponse struct {
	Port int `json:"port"`
	ServerType string `json:"server_type"`
	Version string `json:"version"`
}

func hhCreateAcquiaServer(w http.ResponseWriter, r *http.Request) {
	j := &createServerRequest{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(j); err != nil {
		panic(err)
	}

	an := negroni.New()
	app := acquia.NewServerInstance("default")
	an.UseHandler(app.Router())

	laddr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:0") // listen to all the things
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		panic(err)
	}

	kal := tcpKeepAliveListener{l}
	si := ServerInstance{Listener: kal, MockApp: app}

	resp, err := json.Marshal(createServerResponse{
		Port: kal.Addr().(*net.TCPAddr).Port,
		ServerType: "acquia",
		Version: si.Version(),
	})
	if err != nil {
		panic(err)
	}

	go app.Serve(kal)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Location", "acquia/" + strconv.Itoa(kal.Addr().(*net.TCPAddr).Port))
	w.WriteHeader(201)
	w.Write(resp)

	servers = append(servers, si)
}

func hhListAcquiaServers(w http.ResponseWriter, r *http.Request) {

}

func hhListServerTypes(w http.ResponseWriter, r *http.Request) {

}
