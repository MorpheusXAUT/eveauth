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
	// DefaultCharacter indicates whether the character is selected as default for the account
	DefaultCharacter bool `json:"defaultCharacter"`
	// Active indicates whether the Character is set as active
	Active bool `json:"active"`
}

// AuthCharacter represents a character used by the authorization handler to pass required information to apps
type AuthCharacter struct {
	// ID represents the database ID of the Character
	ID int64 `json:"id"`
	// CorporationID represents the ID of the corporation the Character is in
	CorporationID int64 `json:"corporationID"`
	// Name represents the ingame name of the Character
	Name string `json:"name"`
	// EVECharacterID represents the ingame character ID of the Character
	EVECharacterID int64 `json:"eveCharacterID"`
	// DefaultCharacter indicates whether the character is selected as default for the account
	DefaultCharacter bool `json:"defaultCharacter"`
}

// NewCharacter creates a new character with the given information
func NewCharacter(accountID int64, corporationID int64, name string, eveCharacterID int64, defaultCharacter bool, active bool) *Character {
	character := &Character{
		ID:               -1,
		AccountID:        accountID,
		CorporationID:    corporationID,
		Name:             name,
		EVECharacterID:   eveCharacterID,
		DefaultCharacter: defaultCharacter,
		Active:           active,
	}

	return character
}

// ToAuthCharacter converts the given character to an AuthCharacter, exporting only the information required by third-party apps
func (character *Character) ToAuthCharacter() *AuthCharacter {
	authCharacter := &AuthCharacter{
		ID:               character.ID,
		CorporationID:    character.CorporationID,
		Name:             character.Name,
		EVECharacterID:   character.EVECharacterID,
		DefaultCharacter: character.DefaultCharacter,
	}

	return authCharacter
}

// String represents a JSON encoded representation of the character
func (character *Character) String() string {
	jsonContent, err := json.Marshal(character)
	if err != nil {
		return ""
	}

	return string(jsonContent)
}

// String represents a JSON encoded representation of the auth character
func (authCharacter *AuthCharacter) String() string {
	jsonContent, err := json.Marshal(authCharacter)
	if err != nil {
		return ""
	}

	return string(jsonContent)
}
