package models

import (
	"encoding/json"

	"gopkg.in/guregu/null.v2/zero"
)

// User represents an user within the authentication system
type User struct {
	// ID represents the database ID of the User
	ID int64 `json:"id"`
	// Username represents the username of the User
	Username string `json:"username"`
	// Password represents the bcrypt-hashed password of the User
	Password zero.String `json:"password"`
	// Active indicates whether the User is set as active
	Active bool `json:"active"`
	// Characters contains all characters associated with the User
	Characters []*Character `json:"characters,omitempty"`
	// APIKeys contains all API keys associated with the User
	APIKeys []*APIKey `json:"apiKeys,omitempty"`
	// UserRoles contains all UserRoles associated with the User
	UserRoles []*UserRole `json:"userRoles,omitempty"`
	// Groups contains all Groups associated with the User
	Groups []*Group `json:"groups,omitempty"`
}

// String represents a JSON encoded representation of the user
func (user *User) String() string {
	jsonContent, err := json.Marshal(user)
	if err != nil {
		return ""
	}

	return string(jsonContent)
}
