package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
  // hier kann man die notwendigen Variablen für die Tests erzeugen und initialisieren ...

	os.Exit(m.Run()) // end yourself, but before this let all tests run
}

type myHandler struct{}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter,r *http.Request) {
	
}
