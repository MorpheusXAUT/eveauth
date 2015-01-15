package models

import (
	"gopkg.in/guregu/null.v2/zero"
)

type User struct {
	ID         int64        `json:"id"`
	Username   string       `json:"username"`
	Password   zero.String  `json:"password"`
	Active     bool         `json:"active"`
	UserRoles  []*UserRole  `json:"userRoles,omitempty"`
	GroupRoles []*GroupRole `json:"groupRoles,omitempty"`
	Characters []*Character `json:"characters,omitempty"`
}
