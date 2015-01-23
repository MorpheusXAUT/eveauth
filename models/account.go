package models

import (
	"encoding/json"

	"gopkg.in/guregu/null.v2/zero"
)

// Account represents an EVE Online account linked with an API key and a User
type Account struct {
	// ID represents the database ID of the Account
	ID int64 `json:"id"`
	// UserID represents the database ID of the user the Account is assigned to
	UserID int64 `json:"userID"`
	// APIKeyID represents the EVE Online API key ID
	APIKeyID zero.Int `json:"apiKeyID"`
	// APIvCode represents the EVE Online API verification code
	APIvCode zero.String `json:"apivCode"`
	// Active indicates whether the Account is set as active
	Active bool `json:"active"`
	// Characters stores all characters associated to this Account
	Characters []*Character `json:"characters,omitempty"`
}

// String represents a JSON encoded representation of the account
func (account *Account) String() string {
	jsonContent, err := json.Marshal(account)
	if err != nil {
		return ""
	}

	return string(jsonContent)
}
