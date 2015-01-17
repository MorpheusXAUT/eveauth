package models

import (
	"encoding/json"
)

type APIKey struct {
	ID       int64  `json:"id"`
	UserID   int64  `json:"userID"`
	APIKeyID int64  `json:"apiKeyID"`
	APIvCode string `json:"apivCode"`
	Active   bool   `json:"active"`
}

func (apiKey *APIKey) String() string {
	jsonContent, err := json.Marshal(apiKey)
	if err != nil {
		return ""
	}

	return string(jsonContent)
}
