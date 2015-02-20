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

// ToAuthUser converts the given iser to an AuthUser, exporting only the information required by third-party apps
func (user *User) ToAuthUser() *AuthUser {
	authUser := &AuthUser{
		ID:         user.ID,
		Username:   user.Username,
		Roles:      make([]string, 0),
		Characters: make([]*AuthCharacter, 0),
	}

	for _, userRole := range user.UserRoles {
		if userRole.Granted {
			authUser.Roles = append(authUser.Roles, userRole.Role.Name)
		}
	}

	for _, group := range user.Groups {
		for _, groupRole := range group.GroupRoles {
			if groupRole.Granted {
				authUser.Roles = append(authUser.Roles, groupRole.Role.Name)
			}
		}
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
