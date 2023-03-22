package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mattyjr007/reservationbookings/pkg/config"
	"github.com/mattyjr007/reservationbookings/pkg/handlers"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(writeToConsole)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Get("/majors-suite", handlers.Repo.Majors)
	mux.Get("/search-availability", handlers.Repo.Reservation)
	mux.Get("/make-reservation", handlers.Repo.MakeReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)

	//a sample json get request to return json output
	//mux.Get("/search-availability-json", handlers.Repo.AvailabilityJson)
	// Post Requests
	mux.Post("/search-availability", handlers.Repo.PostSearchAvailability)
	// a post request that send in data from javascript
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJson)
	// make reservation post request to explore validation
	mux.Post("/make-reservation", handlers.Repo.PostReservation)

	// configure the path to read in static files
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))
	return mux
}
