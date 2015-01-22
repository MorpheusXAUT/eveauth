package web

import (
	"html/template"
	"net/http"
)

type Templates struct {
	template *template.Template
}

func SetupTemplates() *Templates {
	templates := &Templates{}

	templates.template = template.Must(template.New("").Funcs(templates.TemplateFunctions(nil)).ParseGlob("web/templates/*"))

	return templates
}

func (templates *Templates) ReloadTemplates() {
	templates.template = template.Must(template.New("").Funcs(templates.TemplateFunctions(nil)).ParseGlob("web/templates/*"))
}

func (templates *Templates) ExecuteTemplates(w http.ResponseWriter, r *http.Request, template string, response map[string]interface{}) error {
	return templates.template.Funcs(templates.TemplateFunctions(r)).ExecuteTemplate(w, template, response)
}

func (templates *Templates) TemplateFunctions(r *http.Request) template.FuncMap {
	return template.FuncMap{
		"IsErrorNil": func(e interface{}) bool { return templates.IsErrorNil(e) },
	}
}

func (templates *Templates) IsErrorNil(e interface{}) bool {
	return (e == nil)
}
