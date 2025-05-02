package handlers

import (
	"encoding/gob"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"github.com/powiedl/myGoWebApplication/internal/config"
	"github.com/powiedl/myGoWebApplication/internal/models"
	"github.com/powiedl/myGoWebApplication/internal/render"
)

var app config.AppConfig
var session *scs.SessionManager

func getRoutes() http.Handler {
		// Data to be available in the session
		gob.Register(models.Reservation{})

		// don't forget to change to tru in Production
		app.InProduction = true // changed to true for the test
		app.Basedir = "../../"

		infoLog := log.New(os.Stdout,"[INFO]\t",log.Ldate|log.Ltime)
		app.InfoLog = infoLog
		errorLog := log.New(os.Stdout,"[ERROR]\t",log.Ldate|log.Ltime|log.Lshortfile)
		app.ErrorLog = errorLog
	
		session = scs.New()
		session.Lifetime = 24 * time.Hour
		session.Cookie.Persist = true
		session.Cookie.SameSite = http.SameSiteLaxMode
		session.Cookie.Secure = app.InProduction
		app.Session = session
	
		tc, err := CreateTestTemplateCache(&app) // needed to change to CreateTestTemplateCache - to use the version inside of this test
		if err != nil {
			log.Fatalln("cannot create template cache")
		}
	
		app.TemplateCache = tc
		app.UseCache = false // disable cache
		//log.Println("app.TemplateCache",app.TemplateCache)
		//log.Println("app.Basedir",app.Basedir)
	
		repo := NewRepo(&app) // create a new repository "based on" app
		NewHandlers(repo)
	
		render.NewTemplates(&app) // call render.NewTemplates with the address of the app variable (which means, that the parameter is a pointer)
		
		// everything went fine, copied the content of the func routes
		mux := chi.NewRouter()

		mux.Use(middleware.Recoverer)  // include chi's Recoverer Middleware - gracefully absorb panics and prints the stack trace
		//mux.Use(hitLogger) // include our own middleware in the call handling
		//mux.Use(NoSurf) // exclude NoSurf in the handler test, to simplify the post tests - otherwise we would need to provide a vaild csrf_token to the test ...
		mux.Use(SessionLoad)
	
		// mux.Get("/",http.HandlerFunc(handlers.Repo.Home))
		// mux.Get("/about",http.HandlerFunc(handlers.Repo.About))
		mux.Get("/",Repo.Home)
		mux.Get("/about",Repo.About)
		mux.Get("/contact",Repo.Contact)
		mux.Get("/eremite",Repo.Eremite)
		mux.Get("/couple",Repo.Couple)
		mux.Get("/family",Repo.Family)
		mux.Get("/reservation",Repo.Reservation)
		mux.Post("/reservation",Repo.PostReservation)
		mux.Post("/reservation-json",Repo.ReservationJSON)
		mux.Get("/make-reservation",Repo.MakeReservation)
		mux.Post("/make-reservation",Repo.PostMakeReservation)
		mux.Get("/reservation-overview",Repo.ReservationOverview)
		
	
		fileServer := http.FileServer(http.Dir(app.Basedir + "static/"))
		mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

		return mux
	}

	// #region copied from middleware.go
	func NoSurf(next http.Handler) http.Handler {
		csrfHandler :=nosurf.New(next)
		csrfHandler.SetBaseCookie(http.Cookie{
			HttpOnly:true,
			Path:"/",
			Secure:app.InProduction,
			SameSite:http.SameSiteLaxMode,
		})
		
		return csrfHandler
	}
	
	// SessionLoad loads saved session data for each request
	func SessionLoad(next http.Handler) http.Handler {
		return session.LoadAndSave(next)
	}
// #endregion

	// #region copied from render.go
	// CreateTestTemplateCache creates a map and stores the templates in for caching - copied and renamed from render.go
func CreateTestTemplateCache(app *config.AppConfig) (map[string]*template.Template,error) {
	//log.Println("createTemplateCache...")
	theCache := map[string]*template.Template{}

	// get all available files *-page.template.html from folder ../../templates
	pages,err := filepath.Glob(app.Basedir + "templates/*-page.template.html")
	log.Println("CreateTemplateCache, #pages",len(pages))
	if err != nil {
		return theCache,err
	}

	// range through the slice of *-page.template.html
	for _,page := range pages {
		name := filepath.Base(page)
		//log.Println("  Processing page:",name)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			log.Printf("  Processing page: '%s' -   ERROR: %s",name,err)
			return theCache,nil
		}
		matches, err := filepath.Glob(app.Basedir + "templates/*-layout.template.html") // get all layout files
		if err != nil {
			return theCache, err
		}

		if len(matches)>0 {
			// at least one base layout was found - append the layouts to the currently processed page
			ts, err = ts.ParseGlob(app.Basedir + "templates/*-layout.template.html")
			if err != nil {
				log.Printf("  Processing templates for page: '%s' -   ERROR: %s",name,err)
				return theCache,nil
			}	
		}
		
		theCache[name] = ts
	}
	return theCache,nil
}
// #endregion