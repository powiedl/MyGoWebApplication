package handlers

import (
	"log"
	"net/http"

	"github.com/powiedl/myGoWebApplication/pkg/config"
	"github.com/powiedl/myGoWebApplication/pkg/models"
	"github.com/powiedl/myGoWebApplication/pkg/render"
)

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// Repo - the repository used by the handlers
var Repo *Repository

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// #region 4-34
// Home is the handler for the home page
func (m *Repository) Home(w http.ResponseWriter,r *http.Request) {
	log.Println("Handling Home page")
	remoteIp := r.RemoteAddr
	m.App.Session.Put(r.Context(),"remote_ip",remoteIp)
	render.RenderTemplate(w,"home-page.template.html",&models.TemplateData{}) // &TemplateData{} - Pointer to an empty TemplateData struct
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter,r *http.Request) {
	log.Println("Handling About page")
	// some data or calculation takes place
	sidekickMap := make(map[string]string)
	sidekickMap["morty"] = "Ooh, wee!"

	remoteIp := m.App.Session.GetString(r.Context(),"remote_ip")
	sidekickMap["remote_ip"] = remoteIp
	
	// send the result or any prepared data to the template
	render.RenderTemplate(w,"about-page.template.html",&models.TemplateData{StringMap:sidekickMap,}) // Pointer to an TemplateData struct where the property StringMap is set to the sidekickMap
}
// #endregion

// #region bis inkl 04-33
/*
// Home is the handler for the home page
func Home(w http.ResponseWriter,r *http.Request) {
	log.Println("Handling Home page")
	render.RenderTemplate(w,"home-page.template.html")
}

// About is the handler for the about page
func About(w http.ResponseWriter,r *http.Request) {
	log.Println("Handling About page")
	render.RenderTemplate(w,"about-page.template.html")
}
*/
// #endregion
