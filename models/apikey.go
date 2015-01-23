package models

import (
	"encoding/json"
)

// APIKey represents an EVE Online API key associated with a user
type APIKey struct {
	// ID represents the database ID of the API key
	ID int64 `json:"id"`
	// UserID represents the database ID of the user the API key is assigned to
	UserID int64 `json:"userID"`
	// APIKeyID represents the EVE Online API key ID
	APIKeyID int64 `json:"apiKeyID"`
	// APIvCode represents the EVE Online API verification code
	APIvCode string `json:"apivCode"`
	// Active indicates whether the API key is set as active
	Active bool `json:"active"`
}

// String represents a JSON encoded representation of the API key
func (apiKey *APIKey) String() string {
	jsonContent, err := json.Marshal(apiKey)
	if err != nil {
		return ""
	}

	return string(jsonContent)
}
