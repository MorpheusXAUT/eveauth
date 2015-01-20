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
			Methods:     []string{"POST"},
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
			Name:        "LogoutGet",
			Methods:     []string{"GET"},
			Pattern:     "/logout",
			HandlerFunc: controller.LogoutGetHandler,
		},
		Route{
			Name:        "AuthorizeGet",
			Methods:     []string{"GET"},
			Pattern:     "/authorize",
			HandlerFunc: controller.AuthorizeGetHandler,
		},
		Route{
			Name:        "SettingsGet",
			Methods:     []string{"GET"},
			Pattern:     "/settings",
			HandlerFunc: controller.SettingsGetHandler,
		},
		Route{
			Name:        "SettingsAccountsGet",
			Methods:     []string{"GET"},
			Pattern:     "/settings/accounts",
			HandlerFunc: controller.SettingsAccountsGetHandler,
		},
		Route{
			Name:        "SettingsAPIKeysGet",
			Methods:     []string{"GET"},
			Pattern:     "/settings/apikeys",
			HandlerFunc: controller.SettingsAPIKeysGetHandler,
		},
		Route{
			Name:        "SettingsCharactersGet",
			Methods:     []string{"GET"},
			Pattern:     "/settings/characters",
			HandlerFunc: controller.SettingsCharactersGetHandler,
		},
		Route{
			Name:        "LegalGet",
			Methods:     []string{"GET"},
			Pattern:     "/legal",
			HandlerFunc: controller.LegalGetHandler,
		},
	}

	return r
}
