package models

import (
	"database/sql"
)

type Corporation struct {
	ID               int64          `json:"id"`
	Name             string         `json:"name"`
	Ticker           string         `json:"ticker"`
	EVECorporationID int64          `json:"eveCorporationID"`
	APIKeyID         sql.NullInt64  `json:"apiKeyID"`
	APIvCode         sql.NullString `json:"apivCode"`
	Active           bool           `json:"active"`
}
