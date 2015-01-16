package models

import (
	"gopkg.in/guregu/null.v2/zero"
)

type User struct {
	ID         int64        `json:"id"`
	Username   string       `json:"username"`
	Password   zero.String  `json:"password"`
	Active     bool         `json:"active"`
	Characters []*Character `json:"characters,omitempty"`
	APIKeys    []*APIKey    `json:"apiKeys,omitempty"`
	UserRoles  []*UserRole  `json:"userRoles,omitempty"`
	Groups     []*Group     `json:"groups,omitempty"`
}
