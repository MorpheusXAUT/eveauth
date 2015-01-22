package web

import (
	"encoding/json"
	"fmt"
	"github.com/morpheusxaut/eveauth/misc"
	"net/http"
	"strconv"
)

func (controller *Controller) SendResponse(w http.ResponseWriter, r *http.Request, template string, response map[string]interface{}) {
	response["assetChecksums"] = controller.Checksums

	err := controller.Templates.ExecuteTemplates(w, r, template, response)
	if err != nil {
		misc.Logger.Warnf("Failed to execute template %q: [%v]", template, err)
		controller.SendRawError(w, http.StatusInternalServerError, err)
		return
	}
}

func (controller *Controller) SendRawError(w http.ResponseWriter, statusCode int, err error) {
	errorMessage := []byte(fmt.Sprintf("Received fatal error during operation: [%v]", err))

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(errorMessage)))

	w.WriteHeader(statusCode)

	w.Write(errorMessage)
}

func (controller *Controller) SendJSONResponse(w http.ResponseWriter, r *http.Request, response map[string]interface{}) {
	responseContent, err := json.Marshal(response)
	if err != nil {
		misc.Logger.Warnf("Failed to marshal JSON response: [%v]", err)
		controller.SendRawError(w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(responseContent)))

	w.WriteHeader(http.StatusOK)

	w.Write(responseContent)
}
