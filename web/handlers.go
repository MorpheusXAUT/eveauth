package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/morpheusxaut/eveauth/misc"
)

// IndexGetHandler displays the index page of the web app
func (controller *Controller) IndexGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 1
	response["pageTitle"] = "Index"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	response["loggedIn"] = loggedIn
	response["success"] = true
	response["error"] = nil

	controller.SendResponse(w, r, "index", response)
}

// LoginGetHandler displays the login page of the web app
func (controller *Controller) LoginGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 2
	response["pageTitle"] = "Login"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	response["loggedIn"] = loggedIn

	if loggedIn {
		http.Redirect(w, r, controller.Session.GetLoginRedirect(r), http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		response["success"] = false
		response["error"] = err

		misc.Logger.Warnf("Failed to parse form: [%v]", err)
		controller.SendResponse(w, r, "login", response)
		return
	}

	errorType := r.FormValue("error")

	if len(errorType) > 0 {
		response["success"] = false

		switch strings.ToLower(errorType) {
		case "credentials":
			response["error"] = fmt.Errorf("Invalid username or password, please try again!")
			break
		case "parse":
			response["error"] = fmt.Errorf("Failed to parse request, please try again!")
			break
		default:
			response["error"] = errorType
		}

		controller.SendResponse(w, r, "login", response)
		return
	}

	response["success"] = true
	response["error"] = nil

	controller.SendResponse(w, r, "login", response)
}

// LoginPostHandler handles submitted data from the login page and verifies the user's credentials
func (controller *Controller) LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := controller.Session.IsLoggedIn(w, r)

	if loggedIn {
		http.Redirect(w, r, controller.Session.GetLoginRedirect(r), http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		misc.Logger.Warnf("Failed to parse form: [%v]", err)
		http.Redirect(w, r, "/login?error=parse", http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if len(username) == 0 || len(password) == 0 {
		misc.Logger.Warnf("Received empty username or password")
		http.Redirect(w, r, "/login?error=parse", http.StatusSeeOther)
		return
	}

	err = controller.Session.Authenticate(w, r, username, password)
	if err != nil {
		misc.Logger.Warnf("Failed to authenticate user: [%v]", err)
		http.Redirect(w, r, "/login?error=credentials", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, controller.Session.GetLoginRedirect(r), http.StatusSeeOther)
}

// LoginRegisterGetHandler displays the registration page of the web app
func (controller *Controller) LoginRegisterGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 2
	response["pageTitle"] = "Register"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	if loggedIn {
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		response["success"] = false
		response["error"] = err

		misc.Logger.Warnf("Failed to parse form: [%v]", err)
		controller.SendResponse(w, r, "loginregister", response)
		return
	}

	errorType := r.FormValue("error")

	if len(errorType) > 0 {
		response["success"] = false

		switch strings.ToLower(errorType) {
		case "duplicate":
			response["error"] = fmt.Errorf("An account with this name or email already exists, please try again!")
			break
		case "parse":
			response["error"] = fmt.Errorf("Failed to parse request, please try again!")
			break
		default:
			response["error"] = errorType
		}

		controller.SendResponse(w, r, "loginregister", response)
		return
	}

	response["loggedIn"] = loggedIn
	response["success"] = true
	response["error"] = nil

	controller.SendResponse(w, r, "loginregister", response)
}

// LoginRegisterPostHandler handles submitted data from the registration page and creates a new user
func (controller *Controller) LoginRegisterPostHandler(w http.ResponseWriter, r *http.Request) {
	loggedIn := controller.Session.IsLoggedIn(w, r)

	if loggedIn {
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		misc.Logger.Warnf("Failed to parse form: [%v]", err)
		http.Redirect(w, r, "/login/register?error=parse", http.StatusSeeOther)
		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if len(username) == 0 || len(email) == 0 || len(password) == 0 {
		misc.Logger.Warnf("Received empty username, email or password")
		http.Redirect(w, r, "/login?error=parse", http.StatusSeeOther)
		return
	}

	err = controller.Session.CreateNewUser(w, r, username, email, password)
	if err != nil {
		misc.Logger.Warnf("Failed to create new user: [%v]", err)
		http.Redirect(w, r, "/login?error=duplicate", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, controller.Session.GetLoginRedirect(r), http.StatusSeeOther)
}

// LoginResetGetHandler allows the user to reset their password
func (controller *Controller) LoginResetGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 2
	response["pageTitle"] = "Reset"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	if loggedIn {
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	response["loggedIn"] = loggedIn
	response["success"] = true
	response["error"] = nil

	controller.SendResponse(w, r, "loginreset", response)
}

// LogoutGetHandler destroys the user's current session and thus logs him out
func (controller *Controller) LogoutGetHandler(w http.ResponseWriter, r *http.Request) {
	controller.Session.DestroySession(w, r)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// AuthorizeGetHandler provides an endpoint for applications to request authorization and query user permissions
func (controller *Controller) AuthorizeGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 3
	response["pageTitle"] = "Authorize"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	response["loggedIn"] = loggedIn

	if !loggedIn {
		err := controller.Session.SetLoginRedirect(w, r, "/authorize")
		if err != nil {
			misc.Logger.Warnf("Failed to set login redirect: [%v]", err)
			controller.SendRawError(w, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	response["success"] = true
	response["error"] = nil

	controller.SendResponse(w, r, "authorize", response)
}

// SettingsGetHandler provides the user with some basic settings for his account
func (controller *Controller) SettingsGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 4
	response["pageTitle"] = "Settings"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	response["loggedIn"] = loggedIn

	if !loggedIn {
		err := controller.Session.SetLoginRedirect(w, r, "/settings")
		if err != nil {
			misc.Logger.Warnf("Failed to set login redirect: [%v]", err)
			controller.SendRawError(w, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	response["success"] = true
	response["error"] = nil

	controller.SendResponse(w, r, "settings", response)
}

// SettingsAccountsGetHandler allows the user to manage the EVE accounts linked to his auth-user
func (controller *Controller) SettingsAccountsGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 4
	response["pageTitle"] = "Accounts"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	response["loggedIn"] = loggedIn

	if !loggedIn {
		err := controller.Session.SetLoginRedirect(w, r, "/settings/accounts")
		if err != nil {
			misc.Logger.Warnf("Failed to set login redirect: [%v]", err)
			controller.SendRawError(w, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	response["success"] = true
	response["error"] = nil

	controller.SendResponse(w, r, "settingsaccounts", response)
}

// SettingsAPIKeysGetHandler allows the user to manage the API keys associated with his account
func (controller *Controller) SettingsAPIKeysGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 4
	response["pageTitle"] = "API Keys"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	response["loggedIn"] = loggedIn

	if !loggedIn {
		err := controller.Session.SetLoginRedirect(w, r, "/settings/apikeys")
		if err != nil {
			misc.Logger.Warnf("Failed to set login redirect: [%v]", err)
			controller.SendRawError(w, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	response["success"] = true
	response["error"] = nil

	controller.SendResponse(w, r, "settingsapikeys", response)
}

// SettingsCharactersGetHandler allows the user to manage the EVE Online characters associated with his account
func (controller *Controller) SettingsCharactersGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 4
	response["pageTitle"] = "Characters"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	response["loggedIn"] = loggedIn

	if !loggedIn {
		err := controller.Session.SetLoginRedirect(w, r, "/settings/characters")
		if err != nil {
			misc.Logger.Warnf("Failed to set login redirect: [%v]", err)
			controller.SendRawError(w, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	response["success"] = true
	response["error"] = nil

	controller.SendResponse(w, r, "settingscharacters", response)
}

// LegalGetHandler displays some legal information as well as copyright disclaimers and contact info
func (controller *Controller) LegalGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 5
	response["pageTitle"] = "Legal"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	response["loggedIn"] = loggedIn
	response["success"] = true
	response["error"] = nil

	controller.SendResponse(w, r, "legal", response)
}
