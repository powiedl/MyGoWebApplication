package forms

type errors map[string][]string // ein Map mit einem string als Key und als value ein slice von strings, d. h. zu einem key, kann es mehrere Fehlermeldungen geben

// Add will add an error message to a specific form field
func (e errors) Add(field,message string) {
	e[field] = append(e[field],message)
}

// Get returns the first error message found for a specific form field
func (e errors)Get(field string) string {
	errorString := e[field]
	if len(errorString) == 0 { return "" }
	return errorString[0]
}