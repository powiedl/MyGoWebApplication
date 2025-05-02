package main

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/powiedl/myGoWebApplication/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig
	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// works - expected return type
	case http.Handler:
		// works too
	default:
		t.Error(fmt.Sprintf("Type mismatch: Expected chi.Mux, got %T",v))
	}
}