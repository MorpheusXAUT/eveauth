package models

import (
	"encoding/json"
)

// Group represents a group of users and permissions
type Group struct {
	// ID represents the database ID of the Group
	ID int64 `json:"id"`
	// Name represents the name of the Group
	Name string `json:"name"`
	// Active indicates whether the Group is set as active
	Active bool `json:"active"`
	// GroupRoles stores all the roles associated with the Group
	GroupRoles []*GroupRole `json:"groupRoles,omitempty"`
}

// String represents a JSON encoded representation of the character
func (group *Group) String() string {
	jsonContent, err := json.Marshal(group)
	if err != nil {
		return ""
	}

	return string(jsonContent)
}
