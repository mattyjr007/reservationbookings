package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/mattyjr007/reservationbookings/pkg/config"
	"github.com/mattyjr007/reservationbookings/pkg/handlers"
	"github.com/mattyjr007/reservationbookings/pkg/models"
	"github.com/mattyjr007/reservationbookings/pkg/render"
	"log"
	"net/http"
	"time"
)

const portNumber = "localhost:8080"

// import the apps configs
var app config.AppConfig

// define the session variable outside main func incaase we need it
var session *scs.SessionManager

func main() {
	// what I am putting in session
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
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Println("Error getting template cache", err)
		log.Fatal(err) //this exits the application
		//return
	}
	// store template cache to struct
	app.TemplateCache = tc
	app.UseCache = false

	// pass template cache to variable in render so cache can be used there
	render.NewTemplates(&app)

	// create a repository
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	/*http.HandleFunc("/", handlers.Repo.Home)
	http.HandleFunc("/about", handlers.Repo.About)

	fmt.Printf("Server started on port: %s\n", portNumber)
	log.Fatal(http.ListenAndServe(portNumber, nil))*/
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	fmt.Printf("Server started on port: %s\n", portNumber)
	err = srv.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}

}
