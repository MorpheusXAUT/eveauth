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

func SetupRoutes(controller *Controller) []Route {
	r := []Route{
		Route{
			Name:        "IndexGet",
			Methods:     []string{"GET"},
			Pattern:     "/",
			HandlerFunc: controller.IndexGetHandler,
		},
	}

	return r
}
