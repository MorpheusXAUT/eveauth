package web

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/morpheusxaut/eveauth/session"
)

// Templates stores the parsed HTTP templates used by the web app
type Templates struct {
	template *template.Template
	session  *session.Controller
}

// SetupTemplates parses the HTTP templates from disk and stores them for later usage
func SetupTemplates(sess *session.Controller) *Templates {
	templates := &Templates{
		session: sess,
	}

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
		"IsResultNil":          func(res interface{}) bool { return templates.IsResultNil(res) },
		"HasUserRole":          func(role string) bool { return templates.HasUserRole(r, role) },
		"QueryCorporationName": func(i int64) string { return templates.QueryCorporationName(i) },
	}
}

// IsResultNil checks whether the given result is nil
func (templates *Templates) IsResultNil(result interface{}) bool {
	return (result == nil)
}

// HasUserRole checks whether the current user has a role with the given name granted
func (templates *Templates) HasUserRole(r *http.Request, role string) bool {
	return templates.session.HasUserRole(r, role)
}

// QueryCorporationName queries the database for the name of the corporation with the given ID
func (templates *Templates) QueryCorporationName(corporationID int64) string {
	corporationName, err := templates.session.QueryCorporationName(corporationID)
	if err != nil {
		return fmt.Sprintf("#%d", corporationID)
	}

	return corporationName
}
