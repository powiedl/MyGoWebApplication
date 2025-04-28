package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/powiedl/myGoWebApplication/internal/config"
	"github.com/powiedl/myGoWebApplication/internal/models"
	"github.com/powiedl/myGoWebApplication/internal/render"
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
	render.RenderTemplate(w,r,"home-page.template.html",&models.TemplateData{}) // &TemplateData{} - Pointer to an empty TemplateData struct
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter,r *http.Request) {
	log.Println("Handling About page")
	// send the result or any prepared data to the template
	render.RenderTemplate(w,r,"about-page.template.html",&models.TemplateData{}) // Pointer to an TemplateData struct where the property StringMap is set to the sidekickMap
}

// Contact is the handler for the contact page
func (m *Repository) Contact(w http.ResponseWriter,r *http.Request) {
	log.Println("Handling Contact page")
	// send the result or any prepared data to the template
	render.RenderTemplate(w,r,"contact-page.template.html",&models.TemplateData{}) // Pointer to an TemplateData struct where the property StringMap is set to the sidekickMap
}

// Couple is the handler for the couple page
func (m *Repository) Couple(w http.ResponseWriter,r *http.Request) {
	log.Println("Handling Couple page")
	// send the result or any prepared data to the template
	render.RenderTemplate(w,r,"couple-page.template.html",&models.TemplateData{}) // Pointer to an TemplateData struct where the property StringMap is set to the sidekickMap
}

// Eremite is the handler for the eremite page
func (m *Repository) Eremite(w http.ResponseWriter,r *http.Request) {
	log.Println("Handling Eremite page")
	// send the result or any prepared data to the template
	render.RenderTemplate(w,r,"eremite-page.template.html",&models.TemplateData{}) // Pointer to an TemplateData struct where the property StringMap is set to the sidekickMap
}

// Family is the handler for the family page
func (m *Repository) Family(w http.ResponseWriter,r *http.Request) {
	log.Println("Handling Family page")
	// send the result or any prepared data to the template
	render.RenderTemplate(w,r,"family-page.template.html",&models.TemplateData{}) // Pointer to an TemplateData struct where the property StringMap is set to the sidekickMap
}

// Reservation is the handler for the check-availability page
func (m *Repository) Reservation(w http.ResponseWriter,r *http.Request) {
	log.Println("Handling check-availability page (Reservation route)")
	// send the result or any prepared data to the template
	render.RenderTemplate(w,r,"check-availability-page.template.html",&models.TemplateData{}) // Pointer to an TemplateData struct where the property StringMap is set to the sidekickMap
}

// Reservation is the handler for the check-availability page
func (m *Repository) PostReservation(w http.ResponseWriter,r *http.Request) {
	log.Println("Handling POST check-availability page (POST Reservation route)")
	start := r.Form.Get("startingDate")
	end := r.Form.Get("endingDate")
	// send the result or any prepared data to the template
	w.Write([]byte(fmt.Sprintf("SUCCESS - you sent a post to the reservation page: start=%s, end=%s",start,end)))
}

// MakeReservation is the handler for the make-reservation page
func (m *Repository) MakeReservation(w http.ResponseWriter,r *http.Request) {
	log.Println("Handling MakeReservation page")
	// send the result or any prepared data to the template
	render.RenderTemplate(w,r,"make-reservation-page.template.html",&models.TemplateData{}) // Pointer to an TemplateData struct where the property StringMap is set to the sidekickMap
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// ReservationJSON is the handler for the reservation-json page
func (m *Repository) ReservationJSON(w http.ResponseWriter,r *http.Request) {
	log.Println("Handling ReservationJSON page")
	resp := jsonResponse{
		OK: true,
		Message: "It's available!",
	}
	output, err := json.MarshalIndent(resp,"","  ")
	if err != nil {
		log.Println("Error converting data to JSON:",err)
	}
	// send the result
	w.Header().Set("Content-Type","application/json") 
	w.Write(output)
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
