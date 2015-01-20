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
	response["success"] = true
	response["error"] = nil

	// TODO Send response to client
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

func (controller *Controller) LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	loggedIn, err := controller.Session.IsLoggedIn(w, r)
	if err != nil {
		response["success"] = false
		response["error"] = err

		// TODO Send response to client
	}

	response["loggedIn"] = loggedIn

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

func (controller *Controller) AuthorizeGetHandler(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]interface{})

	loggedIn, err := controller.Session.IsLoggedIn(w, r)
	if err != nil {
		response["success"] = false
		response["error"] = err

		// TODO Send response to client
	}

	response["loggedIn"] = loggedIn

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

	// TODO Send response to client
}
