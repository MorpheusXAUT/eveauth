package web

import (
	"fmt"
	"net/http"
	"net/url"
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
	response["status"] = 0
	response["result"] = nil

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

	response["status"] = 0
	response["result"] = nil

	controller.SendResponse(w, r, "login", response)
}

// LoginPostHandler handles submitted data from the login page and verifies the user's credentials
func (controller *Controller) LoginPostHandler(w http.ResponseWriter, r *http.Request) {
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
		misc.Logger.Warnf("Failed to parse form: [%v]", err)

		response["status"] = 1
		response["result"] = fmt.Errorf("Failed to parse form, please try again!")

		controller.SendResponse(w, r, "login", response)

		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if len(username) == 0 || len(password) == 0 {
		misc.Logger.Warnf("Received empty username or password")

		response["status"] = 1
		response["result"] = fmt.Errorf("Empty username or password, please try again!")

		controller.SendResponse(w, r, "login", response)

		return
	}

	err = controller.Session.Authenticate(w, r, username, password)
	if err != nil {
		misc.Logger.Warnf("Failed to authenticate user: [%v]", err)

		response["status"] = 1
		response["result"] = fmt.Errorf("Invalid username or password, please try again!")

		controller.SendResponse(w, r, "login", response)

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

	response["loggedIn"] = loggedIn

	if loggedIn {
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	response["status"] = 0
	response["result"] = nil

	controller.SendResponse(w, r, "loginregister", response)
}

// LoginRegisterPostHandler handles submitted data from the registration page and creates a new user
func (controller *Controller) LoginRegisterPostHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 2
	response["pageTitle"] = "Register"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	response["loggedIn"] = loggedIn

	if loggedIn {
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		misc.Logger.Warnf("Failed to parse form: [%v]", err)

		response["status"] = 1
		response["result"] = fmt.Errorf("Failed to parse form, please try again!")

		controller.SendResponse(w, r, "loginregister", response)

		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	if len(username) == 0 || len(email) == 0 || len(password) == 0 {
		misc.Logger.Warnf("Received empty username, email or password")

		response["status"] = 1
		response["result"] = fmt.Errorf("Empty username, email or password, please try again!")

		controller.SendResponse(w, r, "loginregister", response)

		return
	}

	err = controller.Session.CreateNewUser(w, r, username, email, password)
	if err != nil {
		misc.Logger.Warnf("Failed to create new user: [%v]", err)

		response["status"] = 1
		response["result"] = fmt.Errorf("Failed to create new user, please try again!")

		controller.SendResponse(w, r, "loginregister", response)

		return
	}

	err = controller.Session.SendEmailVerification(username, email)
	if err != nil {
		misc.Logger.Warnf("Failed to send email verification: [%v]", err)

		response["status"] = 1
		response["result"] = fmt.Errorf("Failed to send email verification, please try again!")

		controller.SendResponse(w, r, "loginregister", response)

		return
	}

	http.Redirect(w, r, "/settings/accounts", http.StatusSeeOther)
}

// LoginVerifyGetHandler handles the verification of email addresses as produced by the registration system
func (controller *Controller) LoginVerifyGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 2
	response["pageTitle"] = "Register"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	if loggedIn {
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	response["loggedIn"] = loggedIn

	err := r.ParseForm()
	if err != nil {
		misc.Logger.Warnf("Failed to parse form: [%v]", err)

		response["status"] = 1
		response["result"] = fmt.Errorf("Failed to parse form, please try again!")

		controller.SendResponse(w, r, "loginregister", response)

		return
	}

	email := r.FormValue("email")
	verification := r.FormValue("verification")

	if len(email) == 0 || len(verification) == 0 {
		misc.Logger.Warnf("Received empty email or verification code")

		response["status"] = 1
		response["result"] = fmt.Errorf("Empty email or verification code, please try again!")

		controller.SendResponse(w, r, "loginregister", response)

		return
	}

	err = controller.Session.VerifyEmail(w, r, email, verification)
	if err != nil {
		misc.Logger.Warnf("Failed to verify email: [%v]", err)

		response["status"] = 1
		response["result"] = fmt.Errorf("Failed to verify email, please try again!")

		controller.SendResponse(w, r, "loginregister", response)

		return
	}

	http.Redirect(w, r, "/settings/accounts", http.StatusSeeOther)
}

// LoginVerifyResendGetHandler allows the user to request re-sending the verification code to his email address
func (controller *Controller) LoginVerifyResendGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 2
	response["pageTitle"] = "Resend email verification"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	if loggedIn {
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	response["loggedIn"] = loggedIn
	response["status"] = 0
	response["result"] = nil

	controller.SendResponse(w, r, "loginverifyresend", response)
}

// LoginVerifyResendPostHandler resends the activiation email as requested by the user
func (controller *Controller) LoginVerifyResendPostHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 2
	response["pageTitle"] = "Resend email verification"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	if loggedIn {
		http.Redirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	response["loggedIn"] = loggedIn

	// TODO resend email verification
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
	response["status"] = 0
	response["result"] = nil

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
		err := controller.Session.SetLoginRedirect(w, r, r.URL.String())
		if err != nil {
			misc.Logger.Warnf("Failed to set login redirect: [%v]", err)
			controller.SendRawError(w, http.StatusInternalServerError, err)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		misc.Logger.Warnf("Failed to parse form: [%v]", err)

		response["status"] = 1
		response["result"] = fmt.Errorf("Failed to parse form, please try again!")

		controller.SendResponse(w, r, "authorize", response)

		return
	}

	app := r.FormValue("app")
	callback := r.FormValue("callback")
	auth := r.FormValue("auth")

	if len(app) == 0 || len(callback) == 0 || len(auth) == 0 {
		misc.Logger.Warnf("Received empty app, callback or auth")

		response["status"] = 1
		response["result"] = fmt.Errorf("Empty app, callback or auth, please try again!")

		controller.SendResponse(w, r, "authorize", response)

		return
	}

	application, err := controller.Session.VerifyApplication(app, callback, auth)
	if err != nil {
		misc.Logger.Warnf("Failed to verify app authentication: [%v]", err)

		response["status"] = 1
		response["result"] = fmt.Errorf("Failed to authenticate app, please try again!")

		controller.SendResponse(w, r, "authorize", response)

		return
	}

	encryptedPayload, err := controller.Session.EncodeUserPermissions(r, application)
	if err != nil {
		misc.Logger.Warnf("Failed to encode user permissions: [%v]", err)

		response["status"] = 1
		response["result"] = fmt.Errorf("Failed to encode user permissions, please try again!")

		controller.SendResponse(w, r, "authorize", response)

		return
	}

	callbackURL, err := url.Parse(callback)
	if err != nil {
		misc.Logger.Warnf("Failed to parse callback URL: [%v]", err)

		response["status"] = 1
		response["result"] = fmt.Errorf("Failed to parse callback URL, please try again!")

		controller.SendResponse(w, r, "authorize", response)

		return
	}

	callbackPayload := url.Values{}
	callbackPayload.Add("permissions", encryptedPayload)

	callbackURL.RawQuery = callbackPayload.Encode()

	http.Redirect(w, r, callbackURL.String(), http.StatusSeeOther)
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

	response["status"] = 0
	response["result"] = nil

	controller.SendResponse(w, r, "settings", response)
}

// SettingsAccountsGetHandler displays the currently associated accounts and lets the user modify them
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

	user, err := controller.Session.GetUser(r)
	if err != nil {
		response["status"] = 1
		response["result"] = fmt.Errorf("Failed to load user data, please try again!")

		controller.SendResponse(w, r, "settingsaccounts", response)
	}

	response["accounts"] = user.Accounts
	response["status"] = 0
	response["result"] = nil

	controller.SendResponse(w, r, "settingsaccounts", response)
}

// SettingsAccountsPutHandler handles AJAX requests used to update the user's accounts settings
func (controller *Controller) SettingsAccountsPutHandler(w http.ResponseWriter, r *http.Request) {
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

	err := r.ParseForm()
	if err != nil {
		misc.Logger.Warnf("Failed to parse form: [%v]", err)

		response["status"] = 1
		response["result"] = "Failed to parse form, please try again!"

		controller.SendJSONResponse(w, r, response)

		return
	}

	command := r.FormValue("command")
	apiKeyID := r.FormValue("apiKeyID")
	apivCode := r.FormValue("apivCode")

	if len(command) == 0 || len(apiKeyID) == 0 || len(apivCode) == 0 {
		misc.Logger.Warnf("Received empty command, apiKeyID or apivCode")

		response["status"] = 1
		response["result"] = "Empty API Key ID or vCode, please try again!"

		controller.SendJSONResponse(w, r, response)

		return
	}

	switch strings.ToLower(command) {
	case "apikeyadd":
		err = controller.Session.SaveAPIKey(w, r, apiKeyID, apivCode)
		if err != nil {
			misc.Logger.Warnf("Failed to save API key: [%v]", err)

			response["status"] = 1

			if strings.Contains(err.Error(), "Duplicate entry") {
				response["result"] = "An API key with this key ID already exists in database, please try again!"
			} else {
				response["result"] = "Failed to save API key, please try again!"
			}

			controller.SendJSONResponse(w, r, response)

			return
		}
	}

	response["status"] = 0
	response["result"] = nil

	controller.SendJSONResponse(w, r, response)
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

	response["status"] = 0
	response["result"] = nil

	controller.SendResponse(w, r, "settingscharacters", response)
}

// LegalGetHandler displays some legal information as well as copyright disclaimers and contact info
func (controller *Controller) LegalGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 5
	response["pageTitle"] = "Legal"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	response["loggedIn"] = loggedIn
	response["status"] = 0
	response["result"] = nil

	controller.SendResponse(w, r, "legal", response)
}
