package models

import (
	"encoding/json"
)

// Character represents an EVE Online character associated with a user
type Character struct {
	// ID represents the database ID of the Character
	ID int64 `json:"id"`
	// AccountID represents the database ID of the user the Character is assigned to
	AccountID int64 `json:"accountID"`
	// CorporationID represents the ID of the corporation the Character is in
	CorporationID int64 `json:"corporationID"`
	// Name represents the ingame name of the Character
	Name string `json:"name"`
	// EVECharacterID represents the ingame character ID of the Character
	EVECharacterID int64 `json:"eveCharacterID"`
	// Active indicates whether the Character is set as active
	Active bool `json:"active"`
}

// NewCharacter creates a new character with the given information
func NewCharacter(accountID int64, corporationID int64, name string, eveCharacterID int64, active bool) *Character {
	character := &Character{
		ID:             -1,
		AccountID:      accountID,
		CorporationID:  corporationID,
		Name:           name,
		EVECharacterID: eveCharacterID,
		Active:         active,
	}

	return character
}

// String represents a JSON encoded representation of the character
func (character *Character) String() string {
	jsonContent, err := json.Marshal(character)
	if err != nil {
		return ""
	}

	return string(jsonContent)
}
