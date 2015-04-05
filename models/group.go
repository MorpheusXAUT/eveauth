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

// NewGroup creates a new group with the given information
func NewGroup(name string, active bool) *Group {
	group := &Group{
		ID:         -1,
		Name:       name,
		Active:     active,
		GroupRoles: make([]*GroupRole, 0),
	}

	return group
}

// HasRole returns the RoleStatus for the provided role name
func (group *Group) HasRole(role string) RoleStatus {
	for _, groupRole := range group.GroupRoles {
		roleStatus := groupRole.IsRole(role)
		if roleStatus != RoleStatusNonExistent {
			return roleStatus
		}
	}

	return RoleStatusNonExistent
}

// GetRoleCount returns the number of roles associated with the current group
func (group *Group) GetRoleCount() int {
	return len(group.GroupRoles)
}

// String represents a JSON encoded representation of the character
func (group *Group) String() string {
	jsonContent, err := json.Marshal(group)
	if err != nil {
		return ""
	}

	return string(jsonContent)
}
