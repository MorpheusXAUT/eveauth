package web

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/morpheusxaut/eveauth/misc"

	"github.com/gorilla/mux"
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
		controller.SendRedirect(w, r, controller.Session.GetLoginRedirect(w, r), http.StatusSeeOther)
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
		controller.SendRedirect(w, r, controller.Session.GetLoginRedirect(w, r), http.StatusSeeOther)
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

	controller.SendRedirect(w, r, controller.Session.GetLoginRedirect(w, r), http.StatusSeeOther)
}

// LoginRegisterGetHandler displays the registration page of the web app
func (controller *Controller) LoginRegisterGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 2
	response["pageTitle"] = "Register"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	response["loggedIn"] = loggedIn

	if loggedIn {
		controller.SendRedirect(w, r, "/settings", http.StatusSeeOther)
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
		controller.SendRedirect(w, r, "/settings", http.StatusSeeOther)
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

	err = controller.Session.SendEmailVerification(w, r, username, email)
	if err != nil {
		misc.Logger.Warnf("Failed to send email verification: [%v]", err)

		response["status"] = 1
		response["result"] = fmt.Errorf("Failed to send email verification, please try again!")

		controller.SendResponse(w, r, "loginregister", response)

		return
	}

	response["status"] = 2
	response["result"] = "Verification email sent! Please use the provided link to verify your account!"

	controller.SendResponse(w, r, "login", response)
}

// LoginVerifyGetHandler handles the verification of email addresses as produced by the registration system
func (controller *Controller) LoginVerifyGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 2
	response["pageTitle"] = "Register"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	if loggedIn {
		controller.SendRedirect(w, r, "/settings", http.StatusSeeOther)
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

	controller.SendRedirect(w, r, "/settings/accounts", http.StatusSeeOther)
}

// LoginVerifyResendGetHandler allows the user to request re-sending the verification code to his email address
func (controller *Controller) LoginVerifyResendGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 2
	response["pageTitle"] = "Resend email verification"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	if loggedIn {
		controller.SendRedirect(w, r, "/settings", http.StatusSeeOther)
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
		controller.SendRedirect(w, r, "/settings", http.StatusSeeOther)
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

	username := r.FormValue("username")
	email := r.FormValue("email")

	if len(username) == 0 && len(email) == 0 {
		misc.Logger.Warnf("Received empty username or email")

		response["status"] = 1
		response["result"] = fmt.Errorf("Empty username or email, please try again!")

		controller.SendResponse(w, r, "loginverifyresend", response)

		return
	}

	err = controller.Session.ResendEmailVerification(w, r, username, email)
	if err != nil {
		misc.Logger.Warnf("Failed to resend email verification: [%v]", err)

		response["status"] = 1
		response["result"] = fmt.Errorf("Failed to resend email verification, please try again!")

		controller.SendResponse(w, r, "loginverifyresend", response)

		return
	}

	response["status"] = 2
	response["result"] = "Verification email resent! Please use the provided link to verify your account!"

	controller.SendResponse(w, r, "login", response)
}

// LoginResetGetHandler allows the user to reset their password
func (controller *Controller) LoginResetGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 2
	response["pageTitle"] = "Reset Password"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	if loggedIn {
		controller.SendRedirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	response["loggedIn"] = loggedIn
	response["status"] = 0
	response["result"] = nil

	controller.SendResponse(w, r, "loginreset", response)
}

// LoginResetPostHandler allows the user to reset their password
func (controller *Controller) LoginResetPostHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 2
	response["pageTitle"] = "Reset Password"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	if loggedIn {
		controller.SendRedirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	response["loggedIn"] = loggedIn

	err := r.ParseForm()
	if err != nil {
		misc.Logger.Warnf("Failed to parse form: [%v]", err)

		response["status"] = 1
		response["result"] = fmt.Errorf("Failed to parse form, please try again!")

		controller.SendResponse(w, r, "loginreset", response)

		return
	}

	username := r.FormValue("username")
	email := r.FormValue("email")

	if len(username) == 0 && len(email) == 0 {
		misc.Logger.Warnf("Received empty username or email")

		response["status"] = 1
		response["result"] = fmt.Errorf("Empty username or email, please try again!")

		controller.SendResponse(w, r, "loginreset", response)

		return
	}

	err = controller.Session.SendPasswordReset(w, r, username, email)
	if err != nil {
		misc.Logger.Warnf("Failed to send password reset: [%v]", err)

		response["status"] = 1
		response["result"] = fmt.Errorf("Failed to send password reset, please try again!")

		controller.SendResponse(w, r, "loginreset", response)

		return
	}

	response["status"] = 2
	response["result"] = "Password reset mail sent! Please use the provided link to change your password!"

	controller.SendResponse(w, r, "loginreset", response)
}

// LoginResetVerifyGetHandler provides the user with a form to reset their password
func (controller *Controller) LoginResetVerifyGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 2
	response["pageTitle"] = "Reset Password"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	if loggedIn {
		controller.SendRedirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	response["loggedIn"] = loggedIn

	err := r.ParseForm()
	if err != nil {
		misc.Logger.Warnf("Failed to parse form: [%v]", err)

		response["status"] = 1
		response["result"] = fmt.Errorf("Failed to parse form, please try again!")

		controller.SendResponse(w, r, "loginreset", response)

		return
	}

	email := r.FormValue("email")
	username := r.FormValue("username")
	verification := r.FormValue("verification")

	if len(email) == 0 || len(username) == 0 || len(verification) == 0 {
		misc.Logger.Warnf("Received empty email, username or verification code")

		response["status"] = 1
		response["result"] = fmt.Errorf("Empty email, username or verification code, please try again!")

		controller.SendResponse(w, r, "loginreset", response)

		return
	}

	response["status"] = 0
	response["result"] = nil
	response["email"] = email
	response["username"] = username
	response["verification"] = verification

	controller.SendResponse(w, r, "loginresetverify", response)
}

// LoginResetVerifyPostHandler updates the user's password as per choice
func (controller *Controller) LoginResetVerifyPostHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 2
	response["pageTitle"] = "Reset Password"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	if loggedIn {
		controller.SendRedirect(w, r, "/settings", http.StatusSeeOther)
		return
	}

	response["loggedIn"] = loggedIn

	err := r.ParseForm()
	if err != nil {
		misc.Logger.Warnf("Failed to parse form: [%v]", err)

		response["status"] = 1
		response["result"] = fmt.Errorf("Failed to parse form, please try again!")

		controller.SendResponse(w, r, "loginreset", response)

		return
	}

	email := r.FormValue("email")
	username := r.FormValue("username")
	verification := r.FormValue("verification")
	password := r.FormValue("password")

	if len(email) == 0 || len(username) == 0 || len(verification) == 0 || len(password) == 0 {
		misc.Logger.Warnf("Received empty email, username, verification, old or password")

		response["status"] = 1
		response["result"] = fmt.Errorf("Empty email, username, verification, old or password, please try again!")

		controller.SendResponse(w, r, "loginreset", response)

		return
	}

	err = controller.Session.VerifyPasswordReset(w, r, email, username, verification, password)
	if err != nil {
		misc.Logger.Warnf("Failed to verify password reset: [%v]", err)

		response["status"] = 1
		response["result"] = fmt.Errorf("Failed to reset password, please try again!")

		controller.SendResponse(w, r, "loginreset", response)

		return
	}

	response["status"] = 2
	response["result"] = "Successfully changed password!"

	controller.SendResponse(w, r, "login", response)
}

// LogoutGetHandler destroys the user's current session and thus logs him out
func (controller *Controller) LogoutGetHandler(w http.ResponseWriter, r *http.Request) {
	controller.Session.DestroySession(w, r)

	controller.SendRedirect(w, r, "/", http.StatusSeeOther)
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

		controller.SendRedirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	err := r.ParseForm()
	if err != nil {
		misc.Logger.Warnf("Failed to parse form: [%v]", err)

		response["status"] = 1
		response["result"] = fmt.Errorf("Failed to parse form, please try again!")

		controller.SendResponse(w, r, "index", response)

		return
	}

	app := r.FormValue("app")
	callback := r.FormValue("callback")
	auth := r.FormValue("auth")

	if len(app) == 0 || len(callback) == 0 || len(auth) == 0 {
		misc.Logger.Warnf("Received empty app, callback or auth")

		response["status"] = 1
		response["result"] = fmt.Errorf("Empty app, callback or auth, please try again!")

		controller.SendResponse(w, r, "index", response)

		return
	}

	application, err := controller.Session.VerifyApplication(app, callback, auth)
	if err != nil {
		misc.Logger.Warnf("Failed to verify app authentication: [%v]", err)

		response["status"] = 1
		response["result"] = fmt.Errorf("Failed to authenticate app, please try again!")

		controller.SendResponse(w, r, "index", response)

		return
	}

	encryptedPayload, err := controller.Session.EncodeUserPermissions(r, application)
	if err != nil {
		misc.Logger.Warnf("Failed to encode user permissions: [%v]", err)

		response["status"] = 1
		response["result"] = fmt.Errorf("Failed to encode user permissions, please try again!")

		controller.SendResponse(w, r, "index", response)

		return
	}

	callbackURL, err := url.Parse(callback)
	if err != nil {
		misc.Logger.Warnf("Failed to parse callback URL: [%v]", err)

		response["status"] = 1
		response["result"] = fmt.Errorf("Failed to parse callback URL, please try again!")

		controller.SendResponse(w, r, "index", response)

		return
	}

	callbackPayload := url.Values{}
	callbackPayload.Add("permissions", encryptedPayload)

	callbackURL.RawQuery = callbackPayload.Encode()

	controller.SendRedirect(w, r, callbackURL.String(), http.StatusSeeOther)
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

		controller.SendRedirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	user, err := controller.Session.GetUser(r)
	if err != nil {
		misc.Logger.Warnf("Failed to load user: [%v]", err)

		response["status"] = 1
		response["result"] = fmt.Errorf("Failed to load user, please try again!")

		controller.SendResponse(w, r, "settings", response)

		return
	}

	response["user"] = user
	response["status"] = 0
	response["result"] = nil

	controller.SendResponse(w, r, "settings", response)
}

// SettingsPutHandler handles AJAX requests used to update the user's settings
func (controller *Controller) SettingsPutHandler(w http.ResponseWriter, r *http.Request) {
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

		controller.SendRedirect(w, r, "/login", http.StatusSeeOther)
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
	settingsEditOldPassword := r.FormValue("settingsEditOldPassword")
	settingsEditEmail := r.FormValue("settingsEditEmail")
	settingsEditNewPassword := r.FormValue("settingsEditNewPassword")
	settingsEditNewPasswordConfirmation := r.FormValue("settingsEditNewPasswordConfirmation")

	if len(command) == 0 || len(settingsEditOldPassword) == 0 || len(settingsEditEmail) == 0 {
		misc.Logger.Warnf("Received empty command, old password or email address")

		response["status"] = 1
		response["result"] = "Empty old password or email address, please try again!"

		controller.SendJSONResponse(w, r, response)

		return
	}

	if len(settingsEditNewPassword) > 0 && !strings.EqualFold(settingsEditNewPassword, settingsEditNewPasswordConfirmation) {
		misc.Logger.Warnf("New passwords didn't match, update cancelled")

		response["status"] = 1
		response["result"] = "New passwords didn't match, please try again!"

		controller.SendJSONResponse(w, r, response)

		return
	}

	switch strings.ToLower(command) {
	case "settingsedit":
		_, err = controller.Session.UpdateUser(w, r, settingsEditEmail, settingsEditOldPassword, settingsEditNewPassword)
		if err != nil {
			misc.Logger.Warnf("Failed to edit settings: [%v]", err)

			response["status"] = 1
			response["result"] = "Failed to edit settings, please try again!"

			controller.SendJSONResponse(w, r, response)

			return
		}
	}

	response["status"] = 0
	response["result"] = nil

	controller.SendJSONResponse(w, r, response)
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

		controller.SendRedirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	accounts, err := controller.Session.GetUserAccounts(r)
	if err != nil {
		misc.Logger.Warnf("Failed to get user accounts: [%v]", err)

		response["status"] = 1
		response["result"] = fmt.Errorf("Failed to load user data, please try again!")

		controller.SendResponse(w, r, "settingsaccounts", response)
	}

	response["accounts"] = accounts
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

		controller.SendRedirect(w, r, "/login", http.StatusSeeOther)
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

	if len(command) == 0 || len(apiKeyID) == 0 {
		misc.Logger.Warnf("Received empty command or apiKeyID")

		response["status"] = 1
		response["result"] = "Empty API Key ID, please try again!"

		controller.SendJSONResponse(w, r, response)

		return
	}

	switch strings.ToLower(command) {
	case "apikeyadd":
		apivCode := r.FormValue("apivCode")
		if len(apivCode) == 0 {
			misc.Logger.Warnf("Received empty apivCode")

			response["status"] = 1
			response["result"] = "Empty API vCode, please try again!"

			controller.SendJSONResponse(w, r, response)

			return
		}

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
	case "apikeydelete":
		err = controller.Session.DeleteAPIKey(w, r, apiKeyID)
		if err != nil {
			misc.Logger.Warnf("Failed to delete API key: [%v]", err)

			response["status"] = 1
			response["result"] = "Failed to delete API key, please try again!"

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

		controller.SendRedirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	characters, err := controller.Session.GetUserCharacters(r)
	if err != nil {
		misc.Logger.Warnf("Failed to get user characters: [%v]", err)

		response["status"] = 1
		response["result"] = fmt.Errorf("Failed to load user data, please try again!")

		controller.SendResponse(w, r, "settingscharacters", response)
	}

	response["characters"] = characters
	response["status"] = 0
	response["result"] = nil

	controller.SendResponse(w, r, "settingscharacters", response)
}

// SettingsCharactersPutHandler handles AJAX requests used to update the user's character settings
func (controller *Controller) SettingsCharactersPutHandler(w http.ResponseWriter, r *http.Request) {
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

		controller.SendRedirect(w, r, "/login", http.StatusSeeOther)
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
	characterID := r.FormValue("characterID")

	if len(command) == 0 || len(characterID) == 0 {
		misc.Logger.Warnf("Received empty command or characterID")

		response["status"] = 1
		response["result"] = "Empty character ID, please try again!"

		controller.SendJSONResponse(w, r, response)

		return
	}

	switch strings.ToLower(command) {
	case "charactersetdefault":
		err = controller.Session.SetDefaultCharacter(w, r, characterID)
		if err != nil {
			misc.Logger.Warnf("Failed to set default character: [%v]", err)

			response["status"] = 1
			response["result"] = "Failed to set default character, please try again!"

			controller.SendJSONResponse(w, r, response)

			return
		}
	}

	response["status"] = 0
	response["result"] = nil

	controller.SendJSONResponse(w, r, response)
}

// AdminUsersGetHandler allows administrators to modify users and assign new groups and roles
func (controller *Controller) AdminUsersGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 6
	response["pageTitle"] = "User Administration"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	response["loggedIn"] = loggedIn

	if !loggedIn {
		err := controller.Session.SetLoginRedirect(w, r, "/admin/users")
		if err != nil {
			misc.Logger.Warnf("Failed to set login redirect: [%v]", err)
			controller.SendRawError(w, http.StatusInternalServerError, err)
			return
		}

		controller.SendRedirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !controller.Session.HasUserRole(r, "admin.users") {
		misc.Logger.Warnf("Unauthorized access to user administration")

		response["status"] = 1
		response["result"] = "You don't have access to this page!"

		controller.SendResponse(w, r, "index", response)
		return
	}

	users, err := controller.Session.LoadAllUsers()
	if err != nil {
		misc.Logger.Warnf("Failed to load all users: [%v]")

		response["status"] = 1
		response["result"] = "Failed to load users, please try again!"

		controller.SendResponse(w, r, "adminusers", response)
		return
	}

	response["users"] = users
	response["status"] = 0
	response["result"] = nil

	controller.SendResponse(w, r, "adminusers", response)
}

// AdminUsersPostHandler updates the user and adds new roles and groups
func (controller *Controller) AdminUsersPostHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 6
	response["pageTitle"] = "User Administration"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	response["loggedIn"] = loggedIn

	if !loggedIn {
		err := controller.Session.SetLoginRedirect(w, r, "/admin/users")
		if err != nil {
			misc.Logger.Warnf("Failed to set login redirect: [%v]", err)
			controller.SendRawError(w, http.StatusInternalServerError, err)
			return
		}

		controller.SendRedirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !controller.Session.HasUserRole(r, "admin.users") {
		misc.Logger.Warnf("Unauthorized access to user administration")

		response["status"] = 1
		response["result"] = "You don't have access to this page!"

		controller.SendResponse(w, r, "adminusers", response)
		return
	}

	err := r.ParseForm()
	if err != nil {
		misc.Logger.Warnf("Failed to parse form: [%v]", err)

		response["status"] = 1
		response["result"] = fmt.Errorf("Failed to parse form, please try again!")

		controller.SendResponse(w, r, "adminusers", response)
		return
	}

	command := r.FormValue("command")
	userID, err := strconv.ParseInt(r.FormValue("userID"), 10, 64)
	if err != nil {
		misc.Logger.Warnf("Failed to parse user ID: [%v]", err)

		response["status"] = 1
		response["result"] = "Failed to parse user ID, please try again!"

		controller.SendResponse(w, r, "adminusers", response)
		return
	}

	if len(command) == 0 {
		misc.Logger.Warnf("Received empty command")

		response["status"] = 1
		response["result"] = "Empty command, please try again!"

		controller.SendResponse(w, r, "adminusers", response)
		return
	}

	switch strings.ToLower(command) {
	case "adminuserdetailsaddgroup":
		group := r.FormValue("adminUserDetailsAddGroupGroup")

		if len(group) == 0 {
			misc.Logger.Warnf("Received empty group ID")

			response["status"] = 1
			response["result"] = "Empty group, please try again!"

			controller.SendResponse(w, r, "adminusers", response)
			return
		}

		groupID, err := strconv.ParseInt(group, 10, 64)
		if err != nil {
			misc.Logger.Warnf("Failed to parse group ID: [%v]", err)

			response["status"] = 1
			response["result"] = "Failed to parse group ID, please try again!"

			controller.SendResponse(w, r, "adminusers", response)
			return
		}

		err = controller.Session.AddGroupToUser(userID, groupID)
		if err != nil {
			misc.Logger.Warnf("Failed to add group to user: [%v]", err)

			response["status"] = 1
			response["result"] = "Failed to add group to user, please try again!"

			controller.SendResponse(w, r, "adminusers", response)
			return
		}
	case "adminuserdetailsadduserrole":
		role := r.FormValue("adminUserDetailsAddUserRoleRole")
		granted := r.FormValue("adminUserDetailsAddUserRoleGranted")

		if len(role) == 0 {
			misc.Logger.Warnf("Received empty role")

			response["status"] = 1
			response["result"] = "Empty role, please try again!"

			controller.SendResponse(w, r, "adminusers", response)
			return
		}

		roleID, err := strconv.ParseInt(role, 10, 64)
		if err != nil {
			misc.Logger.Warnf("Failed to parse role ID: [%v]", err)

			response["status"] = 1
			response["result"] = "Failed to parse role ID, please try again!"

			controller.SendResponse(w, r, "adminusers", response)
			return
		}

		roleGranted := false
		if len(granted) > 0 && strings.EqualFold(granted, "on") {
			roleGranted = true
		}

		err = controller.Session.AddUserRoleToUser(userID, roleID, roleGranted)
		if err != nil {
			misc.Logger.Warnf("Failed to add user role to user: [%v]", err)

			response["status"] = 1
			response["result"] = "Failed to add user role to user, please try again!"

			controller.SendResponse(w, r, "adminusers", response)
			return
		}
	}

	controller.SendRedirect(w, r, fmt.Sprintf("/admin/user/%d", userID), http.StatusSeeOther)
}

// AdminUsersPutHandler updates the user and removes roles and groups
func (controller *Controller) AdminUsersPutHandler(w http.ResponseWriter, r *http.Request) {

}

// AdminUserDetailsGetHandler allows administrators to view details of a user
func (controller *Controller) AdminUserDetailsGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 6
	response["pageTitle"] = "User Details"

	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["userid"], 10, 64)
	if err != nil {
		misc.Logger.Warnf("Failed to parse userID: [%v]", err)

		response["status"] = 1
		response["result"] = "Failed to parse user ID, please try again!"

		controller.SendResponse(w, r, "adminuserdetails", response)
		return
	}

	loggedIn := controller.Session.IsLoggedIn(w, r)

	response["loggedIn"] = loggedIn

	if !loggedIn {
		err = controller.Session.SetLoginRedirect(w, r, fmt.Sprintf("/admin/user/%d", userID))
		if err != nil {
			misc.Logger.Warnf("Failed to set login redirect: [%v]", err)
			controller.SendRawError(w, http.StatusInternalServerError, err)
			return
		}

		controller.SendRedirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !controller.Session.HasUserRole(r, "admin.users") {
		misc.Logger.Warnf("Unauthorized access to user administration")

		response["status"] = 1
		response["result"] = "You don't have access to this page!"

		controller.SendResponse(w, r, "index", response)
		return
	}

	user, err := controller.Session.LoadUserFromUserID(userID)
	if err != nil {
		misc.Logger.Warnf("Failed to load user: [%v]", err)

		response["status"] = 1
		response["result"] = "Failed to load user details, please try again!"

		controller.SendResponse(w, r, "adminuserdetails", response)
		return
	}

	response["user"] = user

	availableGroups, err := controller.Session.LoadAvailableGroupsForUser(user.ID)
	if err != nil {
		misc.Logger.Warnf("Failed to load available groups: [%v]", err)

		response["status"] = 1
		response["result"] = "Failed to load available groups, please try again!"

		controller.SendResponse(w, r, "adminuserdetails", response)
		return
	}

	response["availableGroups"] = availableGroups

	availableUserRoles, err := controller.Session.LoadAvailableUserRolesForUser(user.ID)
	if err != nil {
		misc.Logger.Warnf("Failed to load available user roles: [%v]", err)

		response["status"] = 1
		response["result"] = "Failed to load available user roles, please try again!"

		controller.SendResponse(w, r, "adminuserdetails", response)
		return
	}

	response["availableUserRoles"] = availableUserRoles
	response["status"] = 0
	response["result"] = nil

	controller.SendResponse(w, r, "adminuserdetails", response)
}

// AdminGroupsGetHandler allows administrators to modify groups and assign new roles
func (controller *Controller) AdminGroupsGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 6
	response["pageTitle"] = "Group Administration"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	response["loggedIn"] = loggedIn

	if !loggedIn {
		err := controller.Session.SetLoginRedirect(w, r, "/admin/groups")
		if err != nil {
			misc.Logger.Warnf("Failed to set login redirect: [%v]", err)
			controller.SendRawError(w, http.StatusInternalServerError, err)
			return
		}

		controller.SendRedirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !controller.Session.HasUserRole(r, "admin.groups") {
		misc.Logger.Warnf("Unauthorized access to group administration")

		response["status"] = 1
		response["result"] = "You don't have access to this page!"

		controller.SendResponse(w, r, "index", response)
		return
	}

	groups, err := controller.Session.LoadAllGroups()
	if err != nil {
		misc.Logger.Warnf("Failed to load all groups: [%v]")

		response["status"] = 1
		response["result"] = "Failed to load groups, please try again!"

		controller.SendResponse(w, r, "admingroups", response)
		return
	}

	response["groups"] = groups
	response["status"] = 0
	response["result"] = nil

	controller.SendResponse(w, r, "admingroups", response)
}

// AdminGroupsPostHandler allows creation of new groups and adds roles to existing ones
func (controller *Controller) AdminGroupsPostHandler(w http.ResponseWriter, r *http.Request) {

}

// AdminGroupsPutHandler updates a groups and removes existing roles
func (controller *Controller) AdminGroupsPutHandler(w http.ResponseWriter, r *http.Request) {

}

// AdminGroupDetailsGetHandler allows administrators to view details of a group
func (controller *Controller) AdminGroupDetailsGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 6
	response["pageTitle"] = "Group Details"

	vars := mux.Vars(r)
	groupID, err := strconv.ParseInt(vars["groupid"], 10, 64)
	if err != nil {
		misc.Logger.Warnf("Failed to parse groupID: [%v]", err)

		response["status"] = 1
		response["result"] = "Failed to parse group ID, please try again!"

		controller.SendResponse(w, r, "admingroupdetails", response)
		return
	}

	loggedIn := controller.Session.IsLoggedIn(w, r)

	response["loggedIn"] = loggedIn

	if !loggedIn {
		err = controller.Session.SetLoginRedirect(w, r, fmt.Sprintf("/admin/group/%d", groupID))
		if err != nil {
			misc.Logger.Warnf("Failed to set login redirect: [%v]", err)
			controller.SendRawError(w, http.StatusInternalServerError, err)
			return
		}

		controller.SendRedirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !controller.Session.HasUserRole(r, "admin.groups") {
		misc.Logger.Warnf("Unauthorized access to group administration")

		response["status"] = 1
		response["result"] = "You don't have access to this page!"

		controller.SendResponse(w, r, "index", response)
		return
	}

	group, err := controller.Session.LoadGroupFromGroupID(groupID)
	if err != nil {
		misc.Logger.Warnf("Failed to load group: [%v]", err)

		response["status"] = 1
		response["result"] = "Failed to load group details, please try again!"

		controller.SendResponse(w, r, "admingroupdetails", response)
		return
	}

	response["group"] = group

	availableGroupRoles, err := controller.Session.LoadAvailableGroupRolesForGroup(group.ID)
	if err != nil {
		misc.Logger.Warnf("Failed to load available group roles: [%v]", err)

		response["status"] = 1
		response["result"] = "Failed to load available group roles, please try again!"

		controller.SendResponse(w, r, "admingroupdetails", response)
		return
	}

	response["availableGroupRoles"] = availableGroupRoles
	response["status"] = 0
	response["result"] = nil

	controller.SendResponse(w, r, "admingroupdetails", response)
}

// AdminRolesGetHandler allows administrators to modify roles
func (controller *Controller) AdminRolesGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})
	response["pageType"] = 6
	response["pageTitle"] = "Role Administration"

	loggedIn := controller.Session.IsLoggedIn(w, r)

	response["loggedIn"] = loggedIn

	if !loggedIn {
		err := controller.Session.SetLoginRedirect(w, r, "/admin/roles")
		if err != nil {
			misc.Logger.Warnf("Failed to set login redirect: [%v]", err)
			controller.SendRawError(w, http.StatusInternalServerError, err)
			return
		}

		controller.SendRedirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if !controller.Session.HasUserRole(r, "admin.roles") {
		misc.Logger.Warnf("Unauthorized access to role administration")

		response["status"] = 1
		response["result"] = "You don't have access to this page!"

		controller.SendResponse(w, r, "index", response)
		return
	}

	roles, err := controller.Session.LoadAllRoles()
	if err != nil {
		misc.Logger.Warnf("Failed to load all roles: [%v]")

		response["status"] = 1
		response["result"] = "Failed to load roles, please try again!"

		controller.SendResponse(w, r, "adminroles", response)
		return
	}

	response["roles"] = roles
	response["status"] = 0
	response["result"] = nil

	controller.SendResponse(w, r, "adminroles", response)
}

// AdminRolesPostHandler allows creation of new roles
func (controller *Controller) AdminRolesPostHandler(w http.ResponseWriter, r *http.Request) {

}

// AdminRolesPutHandler allows removal of existing roles
func (controller *Controller) AdminRolesPutHandler(w http.ResponseWriter, r *http.Request) {

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
