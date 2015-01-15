package models

import (
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
