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
	mux.Post("/make-reservation",handlers.Repo.PostMakeReservation)
	mux.Get("/reservation-overview",handlers.Repo.ReservationOverview)
	mux.Get("/choose-bungalow/{id}",handlers.Repo.ChooseBungalow)
	mux.Get("/book-bungalow",handlers.Repo.BookBungalow)
	mux.Get("/user/login",handlers.Repo.ShowLogin)
	mux.Post("/user/login",handlers.Repo.PostShowLogin)
	mux.Get("/user/logout",handlers.Repo.Logout)
	mux.Route("/admin",func(mux chi.Router){
		//mux.Use(Auth)
		mux.Get("/dashboard",handlers.Repo.AdminDashboard)
		mux.Get("/reservations-new",handlers.Repo.AdminNewReservations)
		mux.Get("/reservations-all",handlers.Repo.AdminAllReservations)
		mux.Get("/reservations-calendar",handlers.Repo.AdminReservationsCalendar)
		mux.Get("/reservations/{src}/{id}",handlers.Repo.AdminShowReservation)
		mux.Post("/reservations/{src}/{id}",handlers.Repo.AdminPostShowReservation)
	})

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
