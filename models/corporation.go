package models

import (
	"encoding/json"

	"gopkg.in/guregu/null.v2/zero"
)

// Corporation represents an EVE Online corporation
type Corporation struct {
	// ID represents the database ID of the Corporation
	ID int64 `json:"id"`
	// Name represents the ingame name of the Corporation
	Name string `json:"name"`
	// Ticker represents the ingame ticker of the Corporation
	Ticker string `json:"ticker"`
	// EVECorporationID represents the ingame corporation ID of the Corporation
	EVECorporationID int64 `json:"eveCorporationID"`
	// APIKeyID represents the EVE Online API key ID for the Corporation
	APIKeyID zero.Int `json:"apiKeyID"`
	// APIvCode represents the EVE Online API verification code for the Corporation
	APIvCode zero.String `json:"apivCode"`
	// Active indicates whether the Character is set as active
	Active bool `json:"active"`
}

// String represents a JSON encoded representation of the corporation
func (corporation *Corporation) String() string {
	jsonContent, err := json.Marshal(corporation)
	if err != nil {
		return ""
	}

	return string(jsonContent)
}
