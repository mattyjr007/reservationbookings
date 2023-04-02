package main

import (
	"net/http"
	"os"
	"testing"
)

// TestMain runs before the other test
func TestMain(m *testing.M) {

	// Exit the test, but run the test first with m.run()
	os.Exit(m.Run())
}

// create an object to satisfy implement the http.handler interface
type myHandler struct {
}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
