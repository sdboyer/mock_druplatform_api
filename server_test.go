package main

import (
	"testing"
	"io/ioutil"
	"bytes"
	"net"
	"net/http"
	"encoding/json"
)

func TestListen(t *testing.T) {
	l, err := net.Listen("tcp", "0.0.0.0:5551")
	if err != nil {
		t.Error("Should be able to open port")
	}

	_, err = net.Listen("tcp", "0.0.0.0:5551")
	t.Log(err)
	if err == nil {
		t.Error("Port should already be grabbed")
	}

	l.Close()

	l2, err := net.Listen("tcp", "0.0.0.0:5551")
	if err != nil {
		t.Error("Should be able to open port again")
	}
	l2.Close()

	_, err = l2.Accept()
	t.Log(err)
}

func TestCreateServer(t *testing.T) {
	go setUpMainHttpd()

	c := new(http.Client)
	body, err := json.Marshal(createServerRequest{ServerType: "acquia", Version: "1.0"})
	if err != nil {
		t.Error("Failed to marshal request JSON: ", err)
	}

	resp, err := c.Post("http://localhost:10233/acquia", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Error(err)
	}

	t.Log(resp)
	respbody, err := ioutil.ReadAll(resp.Body)
	t.Log(string(respbody))
}
