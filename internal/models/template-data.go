package models

import "github.com/powiedl/myGoWebApplication/internal/forms"

// TemplateData holds any kind of data sent from handlers to templates
type TemplateData struct {
	StringMap  map[string]string       // string data - map with key type string and value type string
	IntMap     map[string]int          // int data - map with key type string and value type int
	FloatMap   map[string]float64
	Data       map[string]any          // generic data - map with key type stirng and value type any (so value can hold any data)
	                                   // any is the modern way to express the empty interface (interface{}) 
	CSRFToken  string                  // to prevent CSRF (Cross-Site Request Forgery) attacks
	Flash      string                  // some text that should be displayed in a dialog at the client
	Warning    string                  // some warning that should be displayed at the client (e.g. if some form fields need to be filled out)
	Error      string                  // some error that should be displayed at the client
	Form       *forms.Form
}
