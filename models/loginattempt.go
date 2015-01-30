package models

import (
	"encoding/json"
	"time"
)

// LoginAttempt represents a login attempt to the auth backend
type LoginAttempt struct {
	ID         int64     `json:"id"`
	Username   string    `json:"username"`
	RemoteAddr string    `json:"remoteAddr"`
	UserAgent  string    `json:"userAgent"`
	Successful bool      `json:"successful"`
	Timestamp  time.Time `json:"timestamp"`
}

// NewLoginAttempt creates a new login attempt with the given information
func NewLoginAttempt(username string, remoteAddr string, userAgent string, successful bool) *LoginAttempt {
	loginAttempt := &LoginAttempt{
		ID:         -1,
		Username:   username,
		RemoteAddr: remoteAddr,
		UserAgent:  userAgent,
		Successful: successful,
		Timestamp:  time.Now(),
	}

	return loginAttempt
}

// String represents a JSON encoded representation of the login attempt
func (loginAttempt *LoginAttempt) String() string {
	jsonContent, err := json.Marshal(loginAttempt)
	if err != nil {
		return ""
	}

	return string(jsonContent)
}
