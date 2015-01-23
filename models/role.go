package models

import (
	"encoding/json"
)

// Role represents a permission/role that can be assigned to a user. Applications can require certain roles for access
type Role struct {
	// ID represents the database ID of the Role
	ID int64 `json:"id"`
	// Name represents the name of the Role, which will also be used to reference it in applications
	Name string `json:"name"`
	// Active indicates whether the Role is set as active
	Active bool `json:"active"`
}

// GroupRole represents a role assigned to a Group. Group permissions affect all people within the group
type GroupRole struct {
	// ID represents the database ID of the GroupRole
	ID int64 `json:"id"`
	// GroupID represents the database ID of the group the GroupRole is assigned to
	GroupID int64 `json:"groupID"`
	// Role represents the actual Role that is associated to the GroupRole
	Role *Role `json:"role"`
	// AutoAdded indicates whether the GroupRole was added automatically due to dependencies or titles
	AutoAdded bool `json:"autoAdded"`
	// Granted indicates whether the GroupRole is currently set as granted
	Granted bool `json:"granted"`
}

// UserRole represents a role assigned to a single User
type UserRole struct {
	// ID represents the database ID of the UserRole
	ID int64 `json:"id"`
	// UserID represents the database ID of the user the UserRole is assigned to
	UserID int64 `json:"userID"`
	// Role represents the actual Role that is associated to the UserRoles
	Role *Role `json:"role"`
	// AutoAdded indicates whether the UserRole was added automatically due to dependencies or titles
	AutoAdded bool `json:"autoAdded"`
	// Granted indicates whether the UserRole is currently set as granted
	Granted bool `json:"granted"`
}

// String represents a JSON encoded representation of the role
func (role *Role) String() string {
	jsonContent, err := json.Marshal(role)
	if err != nil {
		return ""
	}

	return string(jsonContent)
}

// String represents a JSON encoded representation of the group role
func (groupRole *GroupRole) String() string {
	jsonContent, err := json.Marshal(groupRole)
	if err != nil {
		return ""
	}

	return string(jsonContent)
}

// String represents a JSON encoded representation of the user role
func (userRole *UserRole) String() string {
	jsonContent, err := json.Marshal(userRole)
	if err != nil {
		return ""
	}

	return string(jsonContent)
}
