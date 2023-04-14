package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/mattyjr007/reservationbookings/pkg/config"
	"github.com/mattyjr007/reservationbookings/pkg/handlers"
	"github.com/mattyjr007/reservationbookings/pkg/helpers"
	"github.com/mattyjr007/reservationbookings/pkg/models"
	"github.com/mattyjr007/reservationbookings/pkg/render"
	"log"
	"net/http"
	"os"
	"time"
)

const portNumber = "localhost:8080"

// import the apps configs
var app config.AppConfig

// define the session variable outside main func incaase we need it
var session *scs.SessionManager

// define the infoLog for logging infos
var infoLog *log.Logger

// define the infoLog for logging infos
var errorLog *log.Logger

func main() {

	err := run()
	if err != nil {
		log.Fatal(err) // it'll stop the application if any error
	}
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

func run() error {

	// what I am putting in session
	gob.Register(models.Reservation{})

	//change to true if in production mode
	app.Inproduction = false

	// add a logger to log in Error
	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	//add an error log
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

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
		return err
	}
	// store template cache to struct
	app.TemplateCache = tc
	app.UseCache = false

	// pass template cache to variable in render so cache can be used there
	render.NewTemplates(&app)

	// create a repository
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	// pass the app pointer
	helpers.NewHelpers(&app)

	return nil
}
