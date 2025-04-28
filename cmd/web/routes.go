package main

import (
	"net/http"

	// go get github.com/bmizerany/pat
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/powiedl/myGoWebApplication/internal/config"
	"github.com/powiedl/myGoWebApplication/internal/handlers"
)

func routes(app *config.AppConfig) http.Handler {
  mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)  // include chi's Recoverer Middleware - gracefully absorb panics and prints the stack trace
	//mux.Use(hitLogger) // include our own middleware in the call handling
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	// mux.Get("/",http.HandlerFunc(handlers.Repo.Home))
	// mux.Get("/about",http.HandlerFunc(handlers.Repo.About))
	mux.Get("/",handlers.Repo.Home)
	mux.Get("/about",handlers.Repo.About)
	mux.Get("/contact",handlers.Repo.Contact)
	mux.Get("/eremite",handlers.Repo.Eremite)
	mux.Get("/couple",handlers.Repo.Couple)
	mux.Get("/family",handlers.Repo.Family)
	mux.Get("/reservation",handlers.Repo.Reservation)
	mux.Post("/reservation",handlers.Repo.PostReservation)
	mux.Post("/reservation-json",handlers.Repo.ReservationJSON)
	mux.Get("/make-reservation",handlers.Repo.MakeReservation)
	

	fileServer := http.FileServer(http.Dir(app.Basedir + "static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	
	// #region 5-38
  /*
	mux := pat.New()

  mux.Get("/",http.HandlerFunc(handlers.Repo.Home))
  mux.Get("/about",http.HandlerFunc(handlers.Repo.About))
  */
	// #endregion

	return mux
}
