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
	// CEOID represents the EVE character ID of the CEO of the Corporation
	CEOID int64 `json:"ceoID"`
	// APIKeyID represents the EVE Online API key ID for the Corporation
	APIKeyID zero.Int `json:"apiKeyID"`
	// APIvCode represents the EVE Online API verification code for the Corporation
	APIvCode zero.String `json:"apivCode"`
	// Active indicates whether the Character is set as active
	Active bool `json:"active"`
}

// NewCorporation creates a new corporation with the given information
func NewCorporation(name string, ticker string, eveCorporationID int64, ceoID int64, apiKeyID zero.Int, apivCode zero.String, active bool) *Corporation {
	corporation := &Corporation{
		ID:               -1,
		Name:             name,
		Ticker:           ticker,
		EVECorporationID: eveCorporationID,
		CEOID:            ceoID,
		APIKeyID:         apiKeyID,
		APIvCode:         apivCode,
		Active:           active,
	}

	return corporation
}

// String represents a JSON encoded representation of the corporation
func (corporation *Corporation) String() string {
	jsonContent, err := json.Marshal(corporation)
	if err != nil {
		return ""
	}

	return string(jsonContent)
}
