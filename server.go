package main

import (
	//"log"
	"encoding/json"
	"github.com/codegangsta/negroni"
	"github.com/sdboyer/mock_druplatform_api/acquia"
	"net/http"
	"net"
	"time"
)

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
//
// Copied from net/http, b/c they did not deign to export
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
	HttpServer *http.Server
}

type ApiServlet interface {
	Name() string
	Version() string
}

var servers = make([]ServerInstance, 0)

type createServerRequest struct {
	ServerType string `json:"server_type"`
	Version string `json:"version"`
}
type createServerResponse struct{
	Port int `json:"port"`
	ServerType string `json:"server_type"`
	Version string `json:"version"`
}

func hhCreateServer(w http.ResponseWriter, r *http.Request) {
	j := &createServerRequest{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(j); err != nil {
		panic(err)
	}

	an := negroni.New()
	an.UseHandler(acquia.NewRouter(acquia.NewServerInstance("default")))

	laddr, err := net.ResolveTCPAddr("tcp", "0.0.0.0:0") // listen to all the things
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		panic(err)
	}
	kal := tcpKeepAliveListener{l}

	// The Addr prop shouldn't actually be used, but set it to avoid triggering defaults
	srv := &http.Server{Addr: laddr.String(), Handler: an}

	si := ServerInstance{Listener: kal, HttpServer: srv}
	resp, err := json.Marshal(createServerResponse{
		Port: kal.Addr().(*net.TCPAddr).Port,
		ServerType: "acquia",
		Version: "1.0",
	})

	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(resp) // auto-sends 200 response
	servers = append(servers, si)
	go srv.Serve(kal)
}

func hhListServers(w http.ResponseWriter, r *http.Request) {

}
