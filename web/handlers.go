package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/morpheusxaut/eveauth/misc"
)

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

	state := misc.GenerateRandomString(32)

	response["ssoState"] = state
	response["ssoClientID"] = controller.Config.EVESSOClientID
	response["ssoCallbackURL"] = controller.Config.EVESSOCallbackURL

	err := controller.Session.SetSSOState(w, r, state)
	if err != nil {
		response["success"] = false
		response["error"] = err

		misc.Logger.Warnf("Failed to set SSO state: [%v]", err)
		controller.SendResponse(w, r, "login", response)
		return
	}

	err = r.ParseForm()
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
		case "ssostate":
			response["error"] = fmt.Errorf("Failed to verify EVE SSO login, please try again!")
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

func (controller *Controller) LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 2
	response["pageTitle"] = "LoginForm"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	response["loggedIn"] = loggedIn

	if loggedIn {
		http.Redirect(w, r, controller.Session.GetLoginRedirect(r), http.StatusSeeOther)
		return
	}

	response["success"] = true
	response["error"] = nil

	// TODO Send response to client
}

func (controller *Controller) LoginSSOGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 2
	response["pageTitle"] = "LoginSSO"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	response["loggedIn"] = loggedIn

	if loggedIn {
		http.Redirect(w, r, controller.Session.GetLoginRedirect(r), http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		misc.Logger.Warnf("Failed to parse form: [%v]", err)
		controller.SendRawError(w, http.StatusBadRequest, err)
		return
	}

	authorizationCode := r.FormValue("code")
	state := r.FormValue("state")

	if len(authorizationCode) == 0 || len(state) == 0 {
		misc.Logger.Warnf("Received empty authorization code or SSO state")
		controller.SendRawError(w, http.StatusBadRequest, fmt.Errorf("Received invalid response from EVE SSO login"))
		return
	}

	savedState := controller.Session.GetSSOState(r)
	if !strings.EqualFold(savedState, state) {
		misc.Logger.Warnf("Failed to verify SSO state")
		http.Redirect(w, r, "/login?error=ssoState", http.StatusSeeOther)
		return
	}

	response["success"] = true
	response["error"] = nil

	// TODO Send response to client
}

func (controller *Controller) LogoutGetHandler(w http.ResponseWriter, r *http.Request) {
	controller.Session.DestroySession(w, r)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

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
