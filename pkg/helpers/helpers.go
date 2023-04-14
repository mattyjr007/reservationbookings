package helpers

import (
	"fmt"
	"github.com/mattyjr007/reservationbookings/pkg/config"
	"net/http"
	"runtime/debug"
)

var app *config.AppConfig

// NewHelpers connect it to the current app config, so we can use its properties
func NewHelpers(ap *config.AppConfig) {
	app = ap
}

func ClientError(w http.ResponseWriter, status int) {
	// just print the status
	app.InfoLog.Println("Client error with status of", status)
	// replies the client with the gotten error
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	// for server error we can trace where the error came from and get more information
	trace := fmt.Sprintf("%s\n %s", err, debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)

}
