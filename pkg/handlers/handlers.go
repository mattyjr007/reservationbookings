package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/mattyjr007/reservationbookings/pkg/config"
	"github.com/mattyjr007/reservationbookings/pkg/models"
	"github.com/mattyjr007/reservationbookings/pkg/render"
	"log"
	"net/http"
)

// Create a Repository to pass components

type Repository struct {
	App *config.AppConfig
}

// Repo the Repository Used by The Handlers
var Repo *Repository

// NewRepo takes in the configuration and create a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home Handler with receiver func to connect to Repo struct
func (s *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// get remote ip address and store it in the session
	remoteIP := r.RemoteAddr
	s.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplateN(w, r, "home.page.gohtml", &models.TemplateData{})

}

// About page handler function
func (s *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello Again..."
	stringMap["name"] = "Samuel mathias"

	//get the remote IP and store it in stringMap to display in the about page
	remoteIP := s.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplateN(w, r, "about.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
	})

}

// Generals route handler
func (s *Repository) Generals(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplateN(w, r, "generals.page.gohtml", &models.TemplateData{})

}

// Majors route handler
func (s *Repository) Majors(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplateN(w, r, "majors.page.gohtml", &models.TemplateData{})

}

// Reservation route handler
func (s *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplateN(w, r, "search-availability.page.gohtml", &models.TemplateData{})

}

// MakeReservation route handler
func (s *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplateN(w, r, "make-reservation.page.gohtml", &models.TemplateData{})

}

// PostReservation route handler
func (s *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {

	start := r.Form.Get("start")
	end := r.Form.Get("end")
	_, err := w.Write([]byte(fmt.Sprintf("Welcome the start date is %s and end date is %s", start, end)))
	if err != nil {
		log.Fatal(err)
	}
	//http.Redirect(w, r, "/", 302)
}

type jsonRespone struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJson is a get operation to return some JSON data
func (s *Repository) AvailabilityJson(w http.ResponseWriter, r *http.Request) {
	// create a sample data response
	resp := jsonRespone{
		OK:      true,
		Message: "Available!",
	}

	// convert the struct to json data
	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Fatal(err)
	}
	// set the header to let the application know what we are sending
	w.Header().Set("Content-Type", "application/json")
	// return the json response with the writer
	w.Write(out)

	//http.Redirect(w, r, "/", 302)
}
