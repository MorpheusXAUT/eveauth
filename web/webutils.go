package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/morpheusxaut/eveauth/misc"
)

// SendResponse sends a response to the client by executing the templates and appending the asset checksum data
func (controller *Controller) SendResponse(w http.ResponseWriter, r *http.Request, template string, response map[string]interface{}) {
	response["assetChecksums"] = controller.Checksums

	err := controller.Templates.ExecuteTemplates(w, r, template, response)
	if err != nil {
		misc.Logger.Warnf("Failed to execute template %q: [%v]", template, err)
		controller.SendRawError(w, http.StatusInternalServerError, err)
		return
	}
}

// SendRawError sends a raw error messages with the given HTTP status code to the client
func (controller *Controller) SendRawError(w http.ResponseWriter, statusCode int, err error) {
	errorMessage := []byte(fmt.Sprintf("Received fatal error during operation: [%v]", err))

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(errorMessage)))

	w.WriteHeader(statusCode)

	w.Write(errorMessage)
}

// SendJSONResponse sends the given reponse data as a JSON encoded string to the client
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
