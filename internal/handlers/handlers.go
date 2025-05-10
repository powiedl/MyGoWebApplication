package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/powiedl/myGoWebApplication/internal/config"
	"github.com/powiedl/myGoWebApplication/internal/driver"
	"github.com/powiedl/myGoWebApplication/internal/forms"
	"github.com/powiedl/myGoWebApplication/internal/helpers"
	"github.com/powiedl/myGoWebApplication/internal/models"
	"github.com/powiedl/myGoWebApplication/internal/render"
	"github.com/powiedl/myGoWebApplication/internal/repository"
	"github.com/powiedl/myGoWebApplication/internal/repository/dbrepo"
)

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB repository.DatabaseRepo
}

// Repo - the repository used by the handlers
var Repo *Repository

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig,db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB: dbrepo.NewPostgresRepo(db.SQL, a),
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
	//m.DB.AllUsers()
	// remoteIp := r.RemoteAddr
	//m.App.Session.Put(r.Context(),"remote_ip",remoteIp)
	render.Template(w,r,"home-page.template.html",&models.TemplateData{}) // Pointer to an TemplateData struct where the property StringMap is set to the sidekickMap
}

// About is the handler for the about page
func (m *Repository) About(w http.ResponseWriter,r *http.Request) {
	log.Println("Handling About page")
	// send the result or any prepared data to the template
	render.Template(w,r,"about-page.template.html",&models.TemplateData{}) // Pointer to an TemplateData struct where the property StringMap is set to the sidekickMap
}

// Contact is the handler for the contact page
func (m *Repository) Contact(w http.ResponseWriter,r *http.Request) {
	log.Println("Handling Contact page")
	// send the result or any prepared data to the template
	render.Template(w,r,"contact-page.template.html",&models.TemplateData{}) // Pointer to an TemplateData struct where the property StringMap is set to the sidekickMap
}

// Couple is the handler for the couple page
func (m *Repository) Couple(w http.ResponseWriter,r *http.Request) {
	log.Println("Handling Couple page")
	// send the result or any prepared data to the template
	render.Template(w,r,"couple-page.template.html",&models.TemplateData{}) // Pointer to an TemplateData struct where the property StringMap is set to the sidekickMap
}

// Eremite is the handler for the eremite page
func (m *Repository) Eremite(w http.ResponseWriter,r *http.Request) {
	log.Println("Handling Eremite page")
	// send the result or any prepared data to the template
	render.Template(w,r,"eremite-page.template.html",&models.TemplateData{}) // Pointer to an TemplateData struct where the property StringMap is set to the sidekickMap
}

// Family is the handler for the family page
func (m *Repository) Family(w http.ResponseWriter,r *http.Request) {
	log.Println("Handling Family page")
	// send the result or any prepared data to the template
	render.Template(w,r,"family-page.template.html",&models.TemplateData{}) // Pointer to an TemplateData struct where the property StringMap is set to the sidekickMap
}

// Reservation is the handler for the check-availability page
func (m *Repository) Reservation(w http.ResponseWriter,r *http.Request) {
	log.Println("Handling check-availability page (Reservation route)")
	// send the result or any prepared data to the template
	render.Template(w,r,"check-availability-page.template.html",&models.TemplateData{}) // Pointer to an TemplateData struct where the property StringMap is set to the sidekickMap
}

// Reservation is the POST handler for the check-availability page
func (m *Repository) PostReservation(w http.ResponseWriter,r *http.Request) {
	log.Println("Handling POST check-availability page (POST Reservation route)")
	start := r.Form.Get("startingDate")
	end := r.Form.Get("endingDate")
	layout := "2006-01-02"
	startDate, err := time.Parse(layout,start)
	if err != nil {
		helpers.ServerError(w,err)
		return 
	}
	
	endDate, err := time.Parse(layout,end)
	if err != nil {
		helpers.ServerError(w,err)
		return 
	}

	bungalows,err := m.DB.SearchAvailabilityByDatesForAllBungalows(startDate,endDate)
	if err != nil {
		helpers.ServerError(w,err)
		return
	}
	if len(bungalows) == 0 {
		m.App.Session.Put(r.Context(),"error","No bungalow is available at that time")
		http.Redirect(w,r,"/reservation",http.StatusSeeOther)
		return
	}
	data := make(map[string]interface{})
	data["bungalows"]=bungalows
	
	res := models.Reservation{
		StartDate: startDate,
		EndDate:endDate,
	}
	m.App.Session.Put(r.Context(),"reservation",res)
	render.Template(w,r,"choose-bungalow-page.template.html",&models.TemplateData{
		Data:data,
	})

	// send the result or any prepared data to the template
	w.Write([]byte(fmt.Sprintf("SUCCESS - you sent a post to the reservation page: start=%s, end=%s",start,end)))
}

// MakeReservation is the handler for the make-reservation page
func (m *Repository) MakeReservation(w http.ResponseWriter,r *http.Request) {
	log.Println("Handling MakeReservation page")
	res,ok := m.App.Session.Get(r.Context(),"reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(),"error","No reservation data in this session available.")
	}

	bungalow, err := m.DB.GetBungalowById(res.BungalowID)
	if err != nil {
		m.App.Session.Put(r.Context(),"error","cannot find bungalow")
		helpers.ServerError(w,err)
		return
	}
	res.Bungalow = bungalow	

	// write the reservation (now with the bungalow details back to the session, so we can use it again in other routes
	resJSON,err:=json.Marshal(res)
	if err == nil {
		log.Println("  reservation (make):",string(resJSON))
	} else {
		log.Println("Unable to convert resevation to json, so it can't logged to the console")
	}
  m.App.Session.Put(r.Context(),"reservation",res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	data := make(map[string]any)
	data["reservation"] = res
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed
	

	// send the result or any prepared data to the template
	render.Template(w,r,"make-reservation-page.template.html",&models.TemplateData{
		Form: forms.New(nil),
		Data: data,
		StringMap: stringMap,
	}) // Pointer to an TemplateData struct where the property StringMap is set to the sidekickMap
}

// Reservation is the POST handler for the check-availability page
func (m *Repository) PostMakeReservation(w http.ResponseWriter,r *http.Request) {
	log.Println("Handling POST make-reservation page")
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w,err)
		//log.Println("!!! Error parsing form data at make-reservation, Error:",err)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout,sd)
	if err != nil {
		helpers.ServerError(w,err)
		return 
	}
	
	endDate, err := time.Parse(layout,ed)
	if err != nil {
		helpers.ServerError(w,err)
		return 
	}

	bungalowID,err := strconv.Atoi(r.Form.Get("bungalow_id"))
	if err != nil {
		helpers.ServerError(w,err)
		return 
	}

	log.Println("Successfully parsed form data at make-reservation")
	reservation := models.Reservation{
		FullName: r.Form.Get("full_name"),
		Email: r.Form.Get("email"),
		Phone: r.Form.Get("phone"),
		StartDate: startDate,
		EndDate: endDate,
		BungalowID: bungalowID,
	}

	form := forms.New(r.PostForm)

	//form.Has("full_name",r)
	form.Required("full_name","email")
	form.MinLength("full_name",2)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]any)
		data["reservation"] = reservation

		render.Template(w,r,"make-reservation-page.template.html",&models.TemplateData{
			Form:form,
			Data:data,
		})
		return
	}
	newReservationId,err := m.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w,err)
		return
	}
	
	restriction := models.BungalowRestriction{
		ID:0,
		StartDate:startDate,
		EndDate:endDate,
		BungalowID:bungalowID,
		ReservationID: newReservationId,
		RestrictionID:1, 
	}

	err = m.DB.InsertBungalowRestriction(restriction)
	if err != nil {
		helpers.ServerError(w,err)
		return
	}

	/*
	resJSON,err:=json.Marshal(reservation)
	if err == nil {
		log.Println("  reservation:",string(resJSON))
	} else {
		log.Println("Unable to convert resevation to json, so it can't logged to the console")
	}
*/

	m.App.Session.Put(r.Context(),"reservation",reservation)
	http.Redirect(w,r,"/reservation-overview",http.StatusSeeOther)
}

type jsonResponse struct {
	OK           bool   `json:"ok"`
	Message      string `json:"message"`
	Bungalow_id  string `json:"bungalow_id"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
}

// ReservationJSON is the handler for the reservation-json page
func (m *Repository) ReservationJSON(w http.ResponseWriter,r *http.Request) {
	log.Println("Handling ReservationJSON page")

	bungalowID, err := strconv.Atoi(r.Form.Get("bungalow_id"))
	if err != nil {
		helpers.ServerError(w,err)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	layout :="2006-01-02"

	startDate, err := time.Parse(layout,sd)
	if err != nil {
		helpers.ServerError(w,err)
		return
	}
	endDate,err := time.Parse(layout,ed)
	if err != nil {
		helpers.ServerError(w,err)
		return
	}

	available,err := m.DB.SearchAvailabilityByDatesByBungalowID(startDate,endDate,bungalowID)
	var resp jsonResponse
	if err != nil {
		helpers.ServerError(w,err)
		resp = jsonResponse{
			OK: false,
			Message: "Error querying database",
		}
	} else {
  	resp = jsonResponse{
	  	OK: available,
		  Message: "",
			Bungalow_id:strconv.Itoa(bungalowID),
			StartDate: sd,
			EndDate: ed,
	  }
	}

	output, err := json.MarshalIndent(resp,"","  ")
	if err != nil {
		helpers.ServerError(w,err)
		//return
	}
	// send the result
	w.Header().Set("Content-Type","application/json") 
	w.Write(output)
}

// ReservationOverview displays the reservation summary page
func (m *Repository) ReservationOverview(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling ReservationOverview page")

  // copied from MakeReservation - reservation.Bungalow does not get "transported" in the session, don't understand why, so we have to do it this way
	res,ok := m.App.Session.Get(r.Context(),"reservation").(models.Reservation)
	if !ok {
		m.App.ErrorLog.Println("!!! Unable to get the reservation information out of the session")
		m.App.Session.Put(r.Context(),"error","No reservation data in this session available.")
		http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
	}

	/*
	resJSON,err:=json.Marshal(res)
	if err == nil {
		log.Println("  reservation (overview):",string(resJSON))
	} else {
		log.Println("Unable to convert reservation to json, so it can't logged to the console")
	}
*/
	bungalow, err := m.DB.GetBungalowById(res.BungalowID)
	if err != nil {
		m.App.Session.Put(r.Context(),"error","cannot find bungalow")
		helpers.ServerError(w,err)
		return
	}
	res.Bungalow = bungalow	

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	data := make(map[string]any)
	data["reservation"] = res
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed
	
	m.App.Session.Remove(r.Context(),"reservation") // alternativ könnte man oben statt Get Pop verwenden,
	                                                // dann müsste man hier nicht extra removen
	//log.Println("Reservation information (out of the session)",reservation)

	render.Template(w,r,"reservation-overview-page.template.html",&models.TemplateData{
		Data:data,
		StringMap: stringMap,
	}) 
}

// ChooseBungalow displays a list of available bungalows and lets the user choose a bungalow
func (m *Repository) ChooseBungalow(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling ChooseBungalow page")
	exploded := strings.Split(r.RequestURI,"/")
	bungalowId,err := strconv.Atoi(exploded[2]) // die URL lautet /choose-bungalow/{id} - darum ist die id das 3. Element (also 2, weil das erste Element ja den Index 0 hat)
	if err != nil {
		m.App.Session.Put(r.Context(),"error","Missing parameter from URL")
		http.Redirect(w,r,"/",http.StatusSeeOther)
		return
	}
	res,ok := m.App.Session.Get(r.Context(),"reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(),"error","Cannot get reservation back from session")
		http.Redirect(w,r,"/",http.StatusSeeOther)
		return
	}
	res.BungalowID = bungalowId
	m.App.Session.Put(r.Context(),"reservation",res)
	// hier ist in der Reservierung in der Session die Bungalow-ID, das Start- und das Enddatum vorhanden
	
	http.Redirect(w,r,"/make-reservation",http.StatusSeeOther)
}

// BookBungalow takes URL parameters from the get request, creates a reservation, stores it in the session and redirects to make-reservation
func (m *Repository) BookBungalow(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling BookBungalow request")
  bungalowId,err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		helpers.ServerError(w,err)
		m.App.Session.Put(r.Context(),"error","Something went wrong while trying to create the booking ...")
		http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
	}

	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")
	layout := "2006-01-02"
	startDate,err := time.Parse(layout,sd)
	if err != nil {
		helpers.ServerError(w,err)
		m.App.Session.Put(r.Context(),"error","Something went wrong while trying to create the booking ...")
		http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
	}

	endDate,err := time.Parse(layout,ed)
	if err != nil {
		helpers.ServerError(w,err)
		m.App.Session.Put(r.Context(),"error","Something went wrong while trying to create the booking ...")
		http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
	}

	var res models.Reservation

	bungalow, err := m.DB.GetBungalowById(bungalowId)
	if err != nil {
		m.App.Session.Put(r.Context(),"error","Cannot find bungalow!")
		http.Redirect(w,r,"/",http.StatusSeeOther)
	}
	res.BungalowID = bungalowId
	res.Bungalow = bungalow
	res.StartDate = startDate
	res.EndDate = endDate

	m.App.Session.Put(r.Context(),"reservation",res)
	http.Redirect(w,r,"/make-reservation",http.StatusSeeOther)
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
