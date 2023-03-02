package handlers

import (
	"github.com/mattyjr007/reservationbookings/pkg/config"
	"github.com/mattyjr007/reservationbookings/pkg/models"
	"github.com/mattyjr007/reservationbookings/pkg/render"
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

	render.RenderTemplateN(w, "home.page.gohtml", &models.TemplateData{})

}

// About page handler function
func (s *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello Again..."
	stringMap["name"] = "Samuel mathias"

	//get the remote IP and store it in stringMap to display in the about page
	remoteIP := s.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplateN(w, "about.page.gohtml", &models.TemplateData{
		StringMap: stringMap,
	})

}

// generals route handler
func (s *Repository) Generals(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplateN(w, "generals.page.gohtml", &models.TemplateData{})

}

// Majors route handler
func (s *Repository) Majors(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplateN(w, "majors.page.gohtml", &models.TemplateData{})

}

// Reservation route handler
func (s *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplateN(w, "search-availability.page.gohtml", &models.TemplateData{})

}

// MakeReservation route handler
func (s *Repository) MakeReservation(w http.ResponseWriter, r *http.Request) {

	render.RenderTemplateN(w, "make-reservation.page.gohtml", &models.TemplateData{})

}
