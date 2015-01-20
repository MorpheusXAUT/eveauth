package web

import (
	"encoding/json"
	"github.com/morpheusxaut/eveauth/misc"
	"net/http"
	"strconv"
)

func (controller *Controller) SendResponse(w http.ResponseWriter, r *http.Request, template string, response map[string]interface{}) {
	err := controller.Templates.ExecuteTemplates(w, r, template, response)
	if err != nil {
		// TODO Display error to client
		misc.Logger.Warnf("Failed to execute template %q: [%v]", template, err)
		return
	}
}

func (controller *Controller) SendJSONResponse(w http.ResponseWriter, r *http.Request, response map[string]interface{}) {
	responseContent, err := json.Marshal(response)
	if err != nil {
		// TODO Display error to client
		misc.Logger.Warnf("Failed to marshal JSON response: [%v]", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(responseContent)))

	w.WriteHeader(http.StatusOK)

	w.Write(responseContent)
}
