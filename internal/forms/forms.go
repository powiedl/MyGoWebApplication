package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form is a type holding a general form struct including an url.Values object
type Form struct {
	url.Values
	Errors errors
}

// New is a function to initialize a form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Valid returns false in case of errors, otherwise true
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// Required checks for the existence of form fields in the post and ensures they are not empty
func (f* Form) Required(fields ...string) { // field is multiple strings, inside of the function field is of type []string - slice of strings
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field,"This field cannot be empty.")
		}
	}
}

// Has checks for the existence of a form field in the post and ensures that the field is not empty
func (f *Form)Has(field string) bool {
	formField := f.Get(field)
	if formField == "" {
		return false
	}
	return true
}

// MinLength returns false if the field value is shorter than a given length, otherwise true
func (f *Form)MinLength(field string, length int) bool {
	formField := f.Get(field)
	if len(strings.TrimSpace(formField)) < length {
		f.Errors.Add(field,fmt.Sprintf("This field must have at least %d characters)",length))
		return false
	}
	return true
}

// IsEmail checks if the value of a field is a valid email address
func (f *Form)IsEmail(field string){
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field,"This field must be a valid email address")
	}
}