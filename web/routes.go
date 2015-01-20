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
		Route{
			Name:        "LoginGet",
			Methods:     []string{"GET"},
			Pattern:     "/login",
			HandlerFunc: controller.LoginGetHandler,
		},
		Route{
			Name:        "LoginPost",
			Methods:     []string{"Post"},
			Pattern:     "/login",
			HandlerFunc: controller.LoginPostHandler,
		},
		Route{
			Name:        "LoginSSOGet",
			Methods:     []string{"GET"},
			Pattern:     "/login/sso",
			HandlerFunc: controller.LoginSSOGetHandler,
		},
		Route{
			Name:        "AuthorizeGet",
			Methods:     []string{"GET"},
			Pattern:     "/authorize",
			HandlerFunc: controller.AuthorizeGetHandler,
		},
	}

	return r
}
