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

	templates.template = template.Must(template.New("").Funcs(templates.TemplateFunctions(nil)).ParseGlob("web/template/*"))

	return templates
}

func (templates *WebTemplates) TemplateFunctions(r *http.Request) template.FuncMap {
	return template.FuncMap{}
}

func (templates *WebTemplates) ExecuteTemplates(w http.ResponseWriter, r *http.Request, template string, response map[string]interface{}) error {
	return templates.template.Funcs(templates.TemplateFunctions(r)).ExecuteTemplate(w, template, response)
}
