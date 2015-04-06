package models

import (
	"encoding/json"
	"strings"
)

// Role represents a permission/role that can be assigned to a user. Applications can require certain roles for access
type Role struct {
	// ID represents the database ID of the Role
	ID int64 `json:"id"`
	// Name represents the name of the Role, which will also be used to reference it in applications
	Name string `json:"name"`
	// Active indicates whether the Role is set as active
	Active bool `json:"active"`
	// Locked indicates whether the Role is locked in database and cannot be delete
	Locked bool `json:"locked"`
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

// NewRole creates a new role with the given information
func NewRole(name string, active bool, locked bool) *Role {
	role := &Role{
		ID:     -1,
		Name:   name,
		Active: active,
		Locked: locked,
	}

	return role
}

// NewGroupRole creates a new group role with the given information
func NewGroupRole(groupID int64, role *Role, autoAdded bool, granted bool) *GroupRole {
	groupRole := &GroupRole{
		ID:        -1,
		GroupID:   groupID,
		Role:      role,
		AutoAdded: autoAdded,
		Granted:   granted,
	}

	return groupRole
}

// NewUserRole creates a new user role with the given information
func NewUserRole(userID int64, role *Role, autoAdded bool, granted bool) *UserRole {
	userRole := &UserRole{
		ID:        -1,
		UserID:    userID,
		Role:      role,
		AutoAdded: autoAdded,
		Granted:   granted,
	}

	return userRole
}

// IsActiveRole checks whether the current role is active and the role name matches the provided value
func (role *Role) IsActiveRole(r string) bool {
	return role.Active && strings.EqualFold(role.Name, r)
}

// IsRole checks whether the current GroupRole has the same role name and returns the appropriate RoleStatus
func (groupRole *GroupRole) IsRole(r string) RoleStatus {
	if !groupRole.Role.IsActiveRole(r) {
		return RoleStatusNonExistent
	}

	if groupRole.Granted {
		return RoleStatusGranted
	}

	return RoleStatusDenied
}

// IsRole checks whether the current UserRole has the same role name and returns the appropriate RoleStatus
func (userRole *UserRole) IsRole(r string) RoleStatus {
	if !userRole.Role.IsActiveRole(r) {
		return RoleStatusNonExistent
	}

	if userRole.Granted {
		return RoleStatusGranted
	}

	return RoleStatusDenied
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

// RoleStatus indicates whether the queried role existed or was granted/denied
type RoleStatus int

const (
	// RoleStatusNonExistent indicates that the role names didn't match
	RoleStatusNonExistent RoleStatus = iota
	// RoleStatusDenied indicates that the role names matched, but the role was not granted
	RoleStatusDenied
	// RoleStatusGranted indicates that the role names matched and the role was granted
	RoleStatusGranted
)

// String returns an easily readable string representation of the role status
func (roleStatus RoleStatus) String() string {
	switch roleStatus {
	case RoleStatusNonExistent:
		return "non-existent"
	case RoleStatusDenied:
		return "denied"
	case RoleStatusGranted:
		return "granted"
	}

	return "unknown"
}
