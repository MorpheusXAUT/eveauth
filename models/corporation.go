package models

type Corporation struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	Ticker           string `json:"ticker"`
	EVECorporationID int64  `json:"eveCorporationID"`
	APIKeyID         int64  `json:"apiKeyID"`
	APIvCode         string `json:"apivCode"`
	Active           bool   `json:"active"`
}
