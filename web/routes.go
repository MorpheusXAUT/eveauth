package web

import (
	"net/http"
)

type Route struct {
	Name        string
	Methods     []string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var routes = []Route{
	Route{
		Name:        "IndexGet",
		Methods:     []string{"GET"},
		Pattern:     "/",
		HandlerFunc: Router.IndexGetHandler,
	},
}
