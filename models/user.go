package models

import (
	"encoding/json"
)

// User represents an user within the authentication system
type User struct {
	// ID represents the database ID of the User
	ID int64 `json:"id"`
	// Username represents the username of the User
	Username string `json:"username"`
	// Password represents the bcrypt-hashed password of the User
	Password string `json:"-"`
	// Email represents the email address of the User
	Email string `json:"email"`
	// VerifiedEmail indicates whether the user has verified their email address
	VerifiedEmail bool `json:"verifiedEmail"`
	// Active indicates whether the User is set as active
	Active bool `json:"active"`
	// Accounts contains all accounts associated with the User
	Accounts []*Account `json:"accounts,omitempty"`
	// UserRoles contains all UserRoles associated with the User
	UserRoles []*UserRole `json:"userRoles,omitempty"`
	// Groups contains all Groups associated with the User
	Groups []*Group `json:"groups,omitempty"`
}

// AuthUser represents a user used by the authorization handler to pass required information to apps
type AuthUser struct {
	// ID represents the database ID of the User
	ID int64 `json:"id"`
	// Username represents the username of the User
	Username string `json:"username"`
	// Roles contains the names of all the rules the User has been granted
	Roles []string `json:"roles,omitempty"`
	// Characters contains all AuthCharacters for the User
	Characters []*AuthCharacter `json:"characters,omitempty"`
}

// NewUser creates a new user with the given information
func NewUser(username string, password string, email string, verified bool, active bool) *User {
	user := &User{
		ID:            -1,
		Username:      username,
		Password:      password,
		Email:         email,
		VerifiedEmail: verified,
		Active:        active,
		Accounts:      make([]*Account, 0),
		UserRoles:     make([]*UserRole, 0),
		Groups:        make([]*Group, 0),
	}

	return user
}

// HasRole returns the RoleStatus for the provided role name
func (user *User) HasRole(role string) RoleStatus {
	roleStatus := RoleStatusNonExistent

	for _, group := range user.Groups {
		status := group.HasRole(role)
		if status != RoleStatusNonExistent {
			roleStatus = status
			break
		}
	}

	for _, userRole := range user.UserRoles {
		status := userRole.IsRole(role)
		if status != RoleStatusNonExistent {
			roleStatus = status
			break
		}
	}

	return roleStatus
}

// GetCharacterCount returns the number of characters associated with the current user
func (user *User) GetCharacterCount() int {
	characterCount := 0

	for _, account := range user.Accounts {
		characterCount += len(account.Characters)
	}

	return characterCount
}

// GetRoleCount returns the number of roles associated with the current user
func (user *User) GetRoleCount() int {
	rolesCount := 0

	for _, userRole := range user.UserRoles {
		if userRole.Role.Active {
			rolesCount++
		}
	}

	for _, group := range user.Groups {
		for _, groupRole := range group.GroupRoles {
			if groupRole.Role.Active {
				rolesCount++
			}
		}
	}

	return rolesCount
}

// GetDefaultCharacter returns the Character object set as a default character
func (user *User) GetDefaultCharacter() *Character {
	for _, accounts := range user.Accounts {
		for _, character := range accounts.Characters {
			if character.DefaultCharacter {
				return character
			}
		}
	}

	return nil
}

// GetEffectiveRoles returns a map of roles as defined via granted/denied group and user roles, index by the role ID
func (user *User) GetEffectiveRoles() map[int64]*Role {
	roles := make(map[int64]*Role)

	for _, group := range user.Groups {
		for _, groupRole := range group.GroupRoles {
			if groupRole.Granted {
				roles[groupRole.Role.ID] = groupRole.Role
			}
		}
	}

	for _, userRole := range user.UserRoles {
		if userRole.Granted {
			roles[userRole.Role.ID] = userRole.Role
		} else {
			_, ok := roles[userRole.Role.ID]
			if ok {
				delete(roles, userRole.Role.ID)
			}
		}
	}

	return roles
}

// ToAuthUser converts the given iser to an AuthUser, exporting only the information required by third-party apps
func (user *User) ToAuthUser() *AuthUser {
	authUser := &AuthUser{
		ID:         user.ID,
		Username:   user.Username,
		Roles:      make([]string, 0),
		Characters: make([]*AuthCharacter, 0),
	}

	effectiveRoles := user.GetEffectiveRoles()

	for _, role := range effectiveRoles {
		authUser.Roles = append(authUser.Roles, role.Name)
	}

	for _, account := range user.Accounts {
		for _, character := range account.Characters {
			authUser.Characters = append(authUser.Characters, character.ToAuthCharacter())
		}
	}

	return authUser
}

// String represents a JSON encoded representation of the user
func (user *User) String() string {
	jsonContent, err := json.Marshal(user)
	if err != nil {
		return ""
	}

	return string(jsonContent)
}

// String represents a JSON encoded representation of the auth user
func (authUser *AuthUser) String() string {
	jsonContent, err := json.Marshal(authUser)
	if err != nil {
		return ""
	}

	return string(jsonContent)
}
