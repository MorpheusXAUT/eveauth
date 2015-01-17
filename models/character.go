package models

import (
	"encoding/json"
)

type Character struct {
	ID             int64  `json:"id"`
	UserID         int64  `json:"userID"`
	CorporationID  int64  `json:"corporationID"`
	Name           string `json:"name"`
	EVECharacterID int64  `json:"eveCharacterID"`
	Active         bool   `json:"active"`
}

func (character *Character) String() string {
	jsonContent, err := json.Marshal(character)
	if err != nil {
		return ""
	}

	return string(jsonContent)
}
