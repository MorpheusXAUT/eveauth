package models

import (
	"encoding/json"
	"net/http"
	"time"
)

type CSRFFailure struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"userID"`
	Request   string    `json:"request"`
	Timestamp time.Time `json:"timestamp"`
}

func NewCSRFFailure(userID int64, r *http.Request) *CSRFFailure {
	csrfFailure := &CSRFFailure{
		ID:        -1,
		UserID:    userID,
		Timestamp: time.Now(),
	}

	jsonRequest, err := json.Marshal(r)
	if err != nil {
		return nil
	}

	csrfFailure.Request = string(jsonRequest)

	return csrfFailure
}

// String represents a JSON encoded representation of the CSRF failure
func (csrfFailure *CSRFFailure) String() string {
	jsonContent, err := json.Marshal(csrfFailure)
	if err != nil {
		return ""
	}

	return string(jsonContent)
}
