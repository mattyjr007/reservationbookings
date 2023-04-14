package render

import (
	"github.com/mattyjr007/reservationbookings/pkg/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData
	// get session
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	//call our adddefaultdata function
	session.Put(r.Context(), "flash", "123")
	result := AddDefaultData(&td, r)

	if result.Flash != "123" {
		t.Error("Failed")
	}
}

func getSession() (*http.Request, error) {
	// create a request to a dummy url to create request
	r := httptest.NewRequest("GET", "/sample", nil)

	ctx := r.Context()
	// put session data
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}

func TestRenderTemplateN(t *testing.T) {
	filedirtemp = "./../../templates"

	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	// create thr request val
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var ww mywriter

	err = RenderTemplateN(&ww, r, "home.page.gohtml", &models.TemplateData{})
	if err != nil {
		t.Error(err)
	}

	err = RenderTemplateN(&ww, r, "doesnt-exist.page.gohtml", &models.TemplateData{})
	if err == nil {
		t.Error("A non existent template was rendered")
	}
}

func TestNewTemplates(t *testing.T) {
	NewTemplates(app)
}

func TestCreateTemplateCache(t *testing.T) {
	filedirtemp = "./../../templates"

	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}
