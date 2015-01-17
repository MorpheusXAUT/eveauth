package models

import (
	"encoding/json"
	"gopkg.in/guregu/null.v2/zero"
)

type Corporation struct {
	ID               int64       `json:"id"`
	Name             string      `json:"name"`
	Ticker           string      `json:"ticker"`
	EVECorporationID int64       `json:"eveCorporationID"`
	APIKeyID         zero.Int    `json:"apiKeyID"`
	APIvCode         zero.String `json:"apivCode"`
	Active           bool        `json:"active"`
}

func (corporation *Corporation) String() string {
	jsonContent, err := json.Marshal(corporation)
	if err != nil {
		return ""
	}

	return string(jsonContent)
}
