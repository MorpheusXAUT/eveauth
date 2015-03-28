package web

import (
	"net/http"
)

// Route stores information about a web route being handled
type Route struct {
	// Name represents a name for the web route
	Name string
	// Methods contains all HTTP methods available to access this route
	Methods []string
	// Pattern defines the URL pattern used to match this route
	Pattern string
	// HandlerFunc represents the web handler function to call for this route
	HandlerFunc http.HandlerFunc
}

// SetupRoutes initialises all used web routes and returns them for the router
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
			Name:        "LoginRegisterGet",
			Methods:     []string{"GET"},
			Pattern:     "/login/register",
			HandlerFunc: controller.LoginRegisterGetHandler,
		},
		Route{
			Name:        "LoginRegisterPost",
			Methods:     []string{"POST"},
			Pattern:     "/login/register",
			HandlerFunc: controller.LoginRegisterPostHandler,
		},
		Route{
			Name:        "LoginVerifyGet",
			Methods:     []string{"GET"},
			Pattern:     "/login/verify",
			HandlerFunc: controller.LoginVerifyGetHandler,
		},
		Route{
			Name:        "LoginVerifyResendGet",
			Methods:     []string{"GET"},
			Pattern:     "/login/verify/resend",
			HandlerFunc: controller.LoginVerifyResendGetHandler,
		},
		Route{
			Name:        "LoginVerifyResendPost",
			Methods:     []string{"POST"},
			Pattern:     "/login/verify/resend",
			HandlerFunc: controller.LoginVerifyResendPostHandler,
		},
		Route{
			Name:        "LoginResetGet",
			Methods:     []string{"GET"},
			Pattern:     "/login/reset",
			HandlerFunc: controller.LoginResetGetHandler,
		},
		Route{
			Name:        "LoginResetPost",
			Methods:     []string{"POST"},
			Pattern:     "/login/reset",
			HandlerFunc: controller.LoginResetPostHandler,
		},
		Route{
			Name:        "LoginResetVerifyGet",
			Methods:     []string{"GET"},
			Pattern:     "/login/reset/verify",
			HandlerFunc: controller.LoginResetVerifyGetHandler,
		},
		Route{
			Name:        "LoginResetVerifyPost",
			Methods:     []string{"POST"},
			Pattern:     "/login/reset/verify",
			HandlerFunc: controller.LoginResetVerifyPostHandler,
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
			Name:        "SettingsPut",
			Methods:     []string{"PUT"},
			Pattern:     "/settings",
			HandlerFunc: controller.SettingsPutHandler,
		},
		Route{
			Name:        "SettingsAccountsGet",
			Methods:     []string{"GET"},
			Pattern:     "/settings/accounts",
			HandlerFunc: controller.SettingsAccountsGetHandler,
		},
		Route{
			Name:        "SettingsAccountsPut",
			Methods:     []string{"PUT"},
			Pattern:     "/settings/accounts",
			HandlerFunc: controller.SettingsAccountsPutHandler,
		},
		Route{
			Name:        "SettingsCharactersGet",
			Methods:     []string{"GET"},
			Pattern:     "/settings/characters",
			HandlerFunc: controller.SettingsCharactersGetHandler,
		},
		Route{
			Name:        "SettingsCharactersPut",
			Methods:     []string{"PUT"},
			Pattern:     "/settings/characters",
			HandlerFunc: controller.SettingsCharactersPutHandler,
		},
		Route{
			Name:        "AdminUsersGet",
			Methods:     []string{"GET"},
			Pattern:     "/admin/users",
			HandlerFunc: controller.AdminUsersGetHandler,
		},
		Route{
			Name:        "AdminGroupsGet",
			Methods:     []string{"GET"},
			Pattern:     "/admin/groups",
			HandlerFunc: controller.AdminGroupsGetHandler,
		},
		Route{
			Name:        "AdminRolesGet",
			Methods:     []string{"GET"},
			Pattern:     "/admin/roles",
			HandlerFunc: controller.AdminRolesGetHandler,
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
