package main

import (
	"net/http"

	// go get github.com/bmizerany/pat
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/powiedl/myGoWebApplication/pkg/config"
	"github.com/powiedl/myGoWebApplication/pkg/handlers"
)

func routes(app *config.AppConfig) http.Handler {
  mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)  // include chi's Recoverer Middleware - gracefully absorb panics and prints the stack trace
	//mux.Use(hitLogger) // include our own middleware in the call handling
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/",http.HandlerFunc(handlers.Repo.Home))
	mux.Get("/about",http.HandlerFunc(handlers.Repo.About))
	// #region 5-38
  /*
	mux := pat.New()

  mux.Get("/",http.HandlerFunc(handlers.Repo.Home))
  mux.Get("/about",http.HandlerFunc(handlers.Repo.About))
  */
	// #endregion

	return mux
}