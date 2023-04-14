package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/mattyjr007/reservationbookings/pkg/config"
	"github.com/mattyjr007/reservationbookings/pkg/forms"
	"github.com/mattyjr007/reservationbookings/pkg/helpers"
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
	// create an empty reservation so we don't get errors
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	render.RenderTemplateN(w, r, "make-reservation.page.gohtml", &models.TemplateData{
		Form: forms.New(nil), // passes the form data
		Data: data,           // pass the empty reservation data
	})

}

// PostReservation route handler
func (s *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm() // this parses the form data r.Form
	if err != nil {
		//log.Fatal(err)
		helpers.ServerError(w, err)
		return
	}

	// save the body/ form data to an instance of the reservation type
	reservation := models.Reservation{
		FirstName: r.Form.Get("first-name"),
		LastName:  r.Form.Get("last-name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm) //

	// this adds the error of field and message also returns true if no error
	//form.Has("first-name", r)
	form.Required("first-name", "last-name", "email", "phone")
	form.MinLength("first-name", 3, r)
	form.IsEmail("email")

	// check if there are no errors entirely ,and it is valid form
	if !form.Valid() {
		// if there are errors it stores the reserv data in our template data
		// doing the above let the user know what we have entered
		data := make(map[string]interface{})
		data["reservation"] = reservation

		// render the form with the error message and prefilled correct values
		render.RenderTemplateN(w, r, "make-reservation.page.gohtml", &models.TemplateData{
			Form: form, // passes the form data
			Data: data, // pass the data
		})
		return

	}
	// after checking if the form is valid we will store the reservation detail in session and take them to reservation summary
	// store the details in session
	s.App.Session.Put(r.Context(), "reservation", reservation) // but tell the program in main what type of data to store
	//redirect to summary page
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther) //303
}

// ReservationSummary route handler displays reservation details after making reservation
func (s *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})
	// get reservation from the session
	reservation, ok := s.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {
		// this holds if reservation data is not present
		//log.Println("Error please make a reservation first")
		s.App.ErrorLog.Println("Error please make a reservation first")
		// now store an error message in session
		s.App.Session.Put(r.Context(), "error", "Please make a reservation first !!")
		// redirect them to homepage
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

		return // prevents the rest of the code running

	} else {
		//remove the reservation data from session
		s.App.Session.Remove(r.Context(), "reservation")
		data["reservation"] = reservation

		render.RenderTemplateN(w, r, "reservation-summary.page.gohtml", &models.TemplateData{
			Data: data,
		})
	}

}

// PostSearchAvailability route handler
func (s *Repository) PostSearchAvailability(w http.ResponseWriter, r *http.Request) {

	start := r.Form.Get("start")
	end := r.Form.Get("end")
	_, err := w.Write([]byte(fmt.Sprintf("Welcome the start date is %s and end date is %s", start, end)))
	if err != nil {
		//log.Fatal(err)
		helpers.ServerError(w, err)
		return
	}
	//http.Redirect(w, r, "/", 302)
}

type jsonRespone struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJson is a GET/POST operation to return some JSON data
func (s *Repository) AvailabilityJson(w http.ResponseWriter, r *http.Request) {
	// create a sample data response
	resp := jsonRespone{
		OK:      true,
		Message: "Available!",
	}

	// convert the struct to json data
	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		//log.Fatal(err)
		helpers.ServerError(w, err)
		return
	}
	// set the header to let the application know what we are sending
	w.Header().Set("Content-Type", "application/json")
	// return the json response with the writer
	w.Write(out)

	//http.Redirect(w, r, "/", 302)
}
