package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	// create a sample request
	r := httptest.NewRequest("POST", "/hometest", nil)

	form := New(r.PostForm)

	// because we have not passed any params form can't be invalid yet so
	if !form.Valid() {
		t.Error("Form is supposed to be valid no params/error params has been passed yet")
	}
}

func TestForm_Required(t *testing.T) {

	r := httptest.NewRequest("POST", "/hometest", nil)

	form := New(r.PostForm)

	form.Required("a", "b", "c")
	// because we do not have values passed for the params form can't be valid yet so
	if form.Valid() {
		t.Error("Form is supposed to be invalid the params doesnt exist")
	}

	//post some data so we test val
	postedData := url.Values{}
	postedData.Add("first_name", "John")
	postedData.Add("last_name", "smith")
	postedData.Add("sex", "male")

	r = httptest.NewRequest("POST", "/hometest", nil)
	r.PostForm = postedData
	form = New(r.PostForm)

	form.Required("first_name", "last_name", "sex")
	// because we do not have values passed for the params form can't be valid yet so
	if !form.Valid() {
		t.Error("Form is supposed to be valid the params exist")
	}
}

func TestForm_Has(t *testing.T) {

	r := httptest.NewRequest("POST", "/has", nil)

	form := New(r.PostForm)

	if form.Has("first_name", r) {
		t.Error("Expected false but got true, firstname doesn't exist.")
	}

	// create the data to post
	postedData := url.Values{}
	postedData.Add("first_name", "John")

	// add form data
	r = httptest.NewRequest("POST", "/has", nil)
	// posting a data to request
	r.PostForm = postedData
	form = New(r.PostForm)

	if !form.Has("first_name", r) {
		//t.Error(r.PostForm.Get("first_name"))
		t.Error("Expected true but got false, firstname exist.")
	}
}

func TestForm_MinLength(t *testing.T) {

	r := httptest.NewRequest("POST", "/has", nil)

	form := New(r.PostForm)

	if form.MinLength("first_name", 3, r) {
		t.Error("Expected false but got true, firstname doesn't exist. so length is 0")
	}

	postedData := url.Values{}
	postedData.Add("first_name", "John")

	// add form data
	r = httptest.NewRequest("POST", "/has", nil)

	r.PostForm = postedData
	form = New(r.PostForm)

	if !form.MinLength("first_name", 3, r) {
		t.Error("Expected true but got false, firstname exist and length is > 3.")
	}

}

func TestForm_IsEmail(t *testing.T) {

	postedData := url.Values{}
	postedData.Add("email", "mathias@hy")

	form := New(postedData)

	form.IsEmail("email")

	if form.Valid() {
		t.Error("Form is meant to be invalid email isn't an email format.")
	}

	postedData = url.Values{}
	postedData.Add("email", "mathias@gmail.com")

	form = New(postedData)

	form.IsEmail("email")

	if !form.Valid() {
		t.Error("Form is meant to be valid email is an email format.")
	}

}
