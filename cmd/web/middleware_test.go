package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {

	// create a handler to pass to nosurf
	var myH myHandler
	h := NoSurf(&myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing this is the output we expect
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is %T", v))
	}

}

func TestSessionLoad(t *testing.T) {

	// create a handler to pass to nosurf
	var myH myHandler
	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing this is the output we expect
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is %T", v))
	}

}
