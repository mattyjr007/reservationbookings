package forms

// errors stores errors and the message
type errors map[string][]string

// Add appends error to map errors the name of the form field and error message
func (e errors) Add(field, message string) {
	// appends the error message
	e[field] = append(e[field], message)
}

// Get the first error message
func (e errors) Get(field string) string {
	//get the specific
	es := e[field]    // get the array of errors from field name
	if len(es) == 0 { // check if error array is empty
		return ""
	}
	// return the first element of the array, most recent error
	return es[0]
}
