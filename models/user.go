package models

import (
	"database/sql"
)

type User struct {
	ID         int64          `json:"id"`
	Username   string         `json:"username"`
	Password   sql.NullString `json:"password"`
	Active     bool           `json:"active"`
	UserRoles  []*UserRole    `json:"userRoles,omitempty"`
	GroupRoles []*GroupRole   `json:"groupRoles,omitempty"`
	Characters []*Character   `json:"characters,omitempty"`
}
