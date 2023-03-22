package forms

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"net/http"
	"net/url"
	"strings"
)

// Form creates a custom form struct, embeds an url.Values object
type Form struct {
	url.Values
	Errors errors
}

// New creates a new form
func New(data url.Values) *Form {
	// return a pointer to a form
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

func (f *Form) Has(field string, r *http.Request) bool {
	// checks if the field of the form exist
	x := r.Form.Get(field)
	if x == "" {
		//	f.Errors.Add(field, "This field cannot be blank")
		return false
	}
	return true
}

// Required loops through a range of strings and get the values of the form and check before adding it's error
func (f *Form) Required(fields ...string) {

	for _, field := range fields {
		value := f.Get(field) // gets the value of the form through url.values (f.Values.Get())
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank") //adds the error message
		}
		// now let's check for min length
	}
}

// MinLength checks for a minimum length in a string
func (f *Form) MinLength(field string, length int, r *http.Request) bool {
	value := f.Get(field) // get the value of the field in form
	// checks for the length of the string
	if len(value) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long.", length))
		return false
	}
	return true
}

// IsEmail checks if an email is valid
func (f *Form) IsEmail(field string) {
	// check valid email address
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid Email address")
	}

}

// Valid return if there are no errors
func (f *Form) Valid() bool {
	return len(f.Errors) == 0

}
