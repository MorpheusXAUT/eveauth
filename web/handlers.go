package web

import (
	"net/http"
)

func (controller *Controller) IndexGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	loggedIn, err := controller.Session.IsLoggedIn(w, r)
	if err != nil {
		response["success"] = false
		response["error"] = err

		// TODO Send response to client
	}

	response["loggedIn"] = loggedIn
	response["pageType"] = 1
	response["pageTitle"] = "Index"
	response["success"] = true
	response["error"] = nil

	controller.SendResponse(w, r, "index", response)
}

func (controller *Controller) LoginGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	loggedIn, err := controller.Session.IsLoggedIn(w, r)
	if err != nil {
		response["success"] = false
		response["error"] = err

		// TODO Send response to client
	}

	response["loggedIn"] = loggedIn
	response["pageType"] = 2
	response["pageTitle"] = "Login"

	if loggedIn {
		redirect, err := controller.Session.GetLoginRedirect(r)
		if err != nil {
			response["success"] = false
			response["error"] = err

			// TODO Send response to client
		}

		http.Redirect(w, r, redirect, http.StatusSeeOther)
		return
	}

	response["success"] = true
	response["error"] = nil

	controller.SendResponse(w, r, "login", response)
}

func (controller *Controller) LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	loggedIn, err := controller.Session.IsLoggedIn(w, r)
	if err != nil {
		response["success"] = false
		response["error"] = err

		// TODO Send response to client
	}

	response["loggedIn"] = loggedIn
	response["pageType"] = 2

	if loggedIn {
		redirect, err := controller.Session.GetLoginRedirect(r)
		if err != nil {
			response["success"] = false
			response["error"] = err

			// TODO Send response to client
		}

		http.Redirect(w, r, redirect, http.StatusSeeOther)
		return
	}

	response["success"] = true
	response["error"] = nil

	// TODO Send response to client
}

func (controller *Controller) LoginSSOGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	loggedIn, err := controller.Session.IsLoggedIn(w, r)
	if err != nil {
		response["success"] = false
		response["error"] = err

		// TODO Send response to client
	}

	response["loggedIn"] = loggedIn
	response["pageType"] = 2

	if loggedIn {
		redirect, err := controller.Session.GetLoginRedirect(r)
		if err != nil {
			response["success"] = false
			response["error"] = err

			// TODO Send response to client
		}

		http.Redirect(w, r, redirect, http.StatusSeeOther)
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

	loggedIn, err := controller.Session.IsLoggedIn(w, r)
	if err != nil {
		response["success"] = false
		response["error"] = err

		// TODO Send response to client
	}

	response["loggedIn"] = loggedIn
	response["pageType"] = 3
	response["pageTitle"] = "Authorize"

	if !loggedIn {
		err = controller.Session.SetLoginRedirect(w, r, "/authorize")
		if err != nil {
			response["success"] = false
			response["error"] = err

			// TODO Send response to client
		}
	}

	response["success"] = true
	response["error"] = nil

	controller.SendResponse(w, r, "authorize", response)
}

func (controller *Controller) SettingsGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	loggedIn, err := controller.Session.IsLoggedIn(w, r)
	if err != nil {
		response["success"] = false
		response["error"] = err

		// TODO Send response to client
	}

	response["loggedIn"] = loggedIn
	response["pageType"] = 4
	response["pageTitle"] = "Settings"

	if !loggedIn {
		err = controller.Session.SetLoginRedirect(w, r, "/settings")
		if err != nil {
			response["success"] = false
			response["error"] = err

			// TODO Send response to client
		}
	}

	response["success"] = true
	response["error"] = nil

	controller.SendResponse(w, r, "settings", response)
}

func (controller *Controller) SettingsAccountsGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	loggedIn, err := controller.Session.IsLoggedIn(w, r)
	if err != nil {
		response["success"] = false
		response["error"] = err

		// TODO Send response to client
	}

	response["loggedIn"] = loggedIn
	response["pageType"] = 4
	response["pageTitle"] = "Accounts"

	if !loggedIn {
		err = controller.Session.SetLoginRedirect(w, r, "/settings/accounts")
		if err != nil {
			response["success"] = false
			response["error"] = err

			// TODO Send response to client
		}
	}

	response["success"] = true
	response["error"] = nil

	controller.SendResponse(w, r, "settingsaccounts", response)
}

func (controller *Controller) SettingsAPIKeysGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	loggedIn, err := controller.Session.IsLoggedIn(w, r)
	if err != nil {
		response["success"] = false
		response["error"] = err

		// TODO Send response to client
	}

	response["loggedIn"] = loggedIn
	response["pageType"] = 4
	response["pageTitle"] = "API Keys"

	if !loggedIn {
		err = controller.Session.SetLoginRedirect(w, r, "/settings/apikeys")
		if err != nil {
			response["success"] = false
			response["error"] = err

			// TODO Send response to client
		}
	}

	response["success"] = true
	response["error"] = nil

	controller.SendResponse(w, r, "settingsapikeys", response)
}

func (controller *Controller) SettingsCharactersGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	loggedIn, err := controller.Session.IsLoggedIn(w, r)
	if err != nil {
		response["success"] = false
		response["error"] = err

		// TODO Send response to client
	}

	response["loggedIn"] = loggedIn
	response["pageType"] = 4
	response["pageTitle"] = "Characters"

	if !loggedIn {
		err = controller.Session.SetLoginRedirect(w, r, "/settings/characters")
		if err != nil {
			response["success"] = false
			response["error"] = err

			// TODO Send response to client
		}
	}

	response["success"] = true
	response["error"] = nil

	controller.SendResponse(w, r, "settingscharacters", response)
}

func (controller *Controller) LegalGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	loggedIn, err := controller.Session.IsLoggedIn(w, r)
	if err != nil {
		response["success"] = false
		response["error"] = err

		// TODO Send response to client
	}

	response["loggedIn"] = loggedIn
	response["pageType"] = 5
	response["pageTitle"] = "Legal"
	response["success"] = true
	response["error"] = nil

	controller.SendResponse(w, r, "legal", response)
}
