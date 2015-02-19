package web

import (
	"html/template"
	"net/http"
)

// Templates stores the parsed HTTP templates used by the web app
type Templates struct {
	template *template.Template
}

// SetupTemplates parses the HTTP templates from disk and stores them for later usage
func SetupTemplates() *Templates {
	templates := &Templates{}

	templates.template = template.Must(template.New("").Funcs(templates.TemplateFunctions(nil)).ParseGlob("app/templates/*"))

	return templates
}

// ReloadTemplates re-reads the HTTP templates from disk and refreshes the output
func (templates *Templates) ReloadTemplates() {
	templates.template = template.Must(template.New("").Funcs(templates.TemplateFunctions(nil)).ParseGlob("app/templates/*"))
}

// ExecuteTemplates performs all replacement in the HTTP templates and sends the response to the client
func (templates *Templates) ExecuteTemplates(w http.ResponseWriter, r *http.Request, template string, response map[string]interface{}) error {
	return templates.template.Funcs(templates.TemplateFunctions(r)).ExecuteTemplate(w, template, response)
}

// TemplateFunctions prepares a map of functions to be used within templates
func (templates *Templates) TemplateFunctions(r *http.Request) template.FuncMap {
	return template.FuncMap{
		"IsResultNil": func(r interface{}) bool { return templates.IsResultNil(r) },
	}
}

// IsResultNil checks whether the given result/interface is nil
func (templates *Templates) IsResultNil(r interface{}) bool {
	return (r == nil)
}
