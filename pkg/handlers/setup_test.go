package handlers

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"github.com/mattyjr007/reservationbookings/pkg/config"
	"github.com/mattyjr007/reservationbookings/pkg/models"
	"github.com/mattyjr007/reservationbookings/pkg/render"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

var app config.AppConfig
var session *scs.SessionManager

// map of functions to do some processing and pass to the template
var functions = template.FuncMap{}
var filedirtemp = "./../../templates"

func getRoutes() http.Handler {
	//registers model.reservation so we can use it in session
	gob.Register(models.Reservation{})

	//change to true if in production mode
	app.Inproduction = false

	// Initialize a new session manager and configure the session lifetime.
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.Inproduction

	//store the session to an instance of app config
	app.Session = session

	// create the cache for template files
	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Println("Error getting template cache", err)
		log.Fatal(err) //this exits the application
		//return

	}
	// store template cache to struct
	app.TemplateCache = tc
	app.UseCache = true // so it uses CreateTestTemplateCache()

	// pass template cache to variable in render so cache can be used there
	render.NewTemplates(&app)

	// create a repository
	repo := NewRepo(&app)
	NewHandlers(repo)

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	//mux.Use(writeToConsole)
	//mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)
	mux.Get("/search-availability", Repo.Reservation)
	mux.Get("/make-reservation", Repo.MakeReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	//a sample json get request to return json output
	//mux.Get("/search-availability-json", Repo.AvailabilityJson)
	// Post Requests
	mux.Post("/search-availability", Repo.PostSearchAvailability)
	// a post request that send in data from javascript
	mux.Post("/search-availability-json", Repo.AvailabilityJson)
	// make reservation post request to explore validation
	mux.Post("/make-reservation", Repo.PostReservation)

	// configure the path to read in static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux

}

// NoSurf adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.Inproduction,
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

// CreateTemplateCache creates a map of templates
func CreateTestTemplateCache() (map[string]*template.Template, error) {

	// this stores the template in a cache
	myCache := map[string]*template.Template{}

	//get all template with .page.gohtml
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.gohtml", filedirtemp))
	if err != nil {
		return myCache, err
	}
	//loop through all the pages
	for _, page := range pages {
		// get the pages base/name
		name := filepath.Base(page)

		// pass the template file
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		// check if the template matches any layout
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.gohtml", filedirtemp))
		if err != nil {
			return myCache, err
		}

		// check if any layout is found and pass the 'pages' into the 'layout'
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.gohtml", filedirtemp))
			if err != nil {
				return myCache, err
			}
		}
		// store the cache templates in a map
		myCache[name] = ts

	}

	return myCache, nil
}
