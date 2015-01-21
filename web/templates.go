package web

import (
	"html/template"
	"net/http"
)

type WebTemplates struct {
	template *template.Template
}

func SetupTemplates() *WebTemplates {
	templates := &WebTemplates{}

	templates.template = template.Must(template.New("").Funcs(templates.TemplateFunctions(nil)).ParseGlob("web/templates/*"))

	return templates
}

func (templates *WebTemplates) ReloadTemplates() {
	templates.template = template.Must(template.New("").Funcs(templates.TemplateFunctions(nil)).ParseGlob("web/templates/*"))
}

func (templates *WebTemplates) ExecuteTemplates(w http.ResponseWriter, r *http.Request, template string, response map[string]interface{}) error {
	return templates.template.Funcs(templates.TemplateFunctions(r)).ExecuteTemplate(w, template, response)
}

func (templates *WebTemplates) TemplateFunctions(r *http.Request) template.FuncMap {
	return template.FuncMap{
		"IsErrorNil": func(e interface{}) bool { return templates.IsErrorNil(e) },
	}
}

func (templates *WebTemplates) IsErrorNil(e interface{}) bool {
	return (e == nil)
}
