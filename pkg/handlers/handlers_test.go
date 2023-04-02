package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"gq", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"ms", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"sa", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"mar", "/make-reservation", "GET", []postData{}, http.StatusOK},
	//{"rs", "/reservation-summary", "GET", []postData{}, http.StatusOK},
	{"sap", "/search-availability", "POST", []postData{
		{"start", "2020-01-01"},
		{"end", "2020-01-02"},
	}, http.StatusOK},

	{"sapJ", "/search-availability-json", "POST", []postData{
		{"start", "2020-01-01"},
		{"end", "2020-01-02"},
	}, http.StatusOK},

	{"mrp", "/make-reservation", "POST", []postData{
		{"first-name", "John"},
		{"last-name", "walter"},
		{"email", "walter@gmail.com"},
		{"phone", "23456670831"},
	}, http.StatusOK},
}

func TestNewHandlers(t *testing.T) {
	routes := getRoutes()

	// create a server in testing
	ts := httptest.NewTLSServer(routes)
	defer ts.Close() // closes the port after the function is finished

	for _, e := range theTests {

		if e.method == "GET" {
			resp, err := ts.Client().Get(ts.URL + e.url) //get the port the test server is listening on and append it

			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}

		} else {
			// build a temporary value that holds information for a post request
			values := url.Values{}
			for _, p := range e.params {
				values.Add(p.key, p.value)
			}
			resp, err := ts.Client().PostForm(ts.URL+e.url, values)

			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}

		}
	}
}
