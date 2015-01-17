package models

import (
	"encoding/json"
)

type Role struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

type GroupRole struct {
	ID        int64 `json:"id"`
	GroupID   int64 `json:"groupID"`
	Role      *Role `json:"role"`
	AutoAdded bool  `json:"autoAdded"`
	Granted   bool  `json:"granted"`
}

type UserRole struct {
	ID        int64 `json:"id"`
	UserID    int64 `json:"userID"`
	Role      *Role `json:"role"`
	AutoAdded bool  `json:"autoAdded"`
	Granted   bool  `json:"granted"`
}

func (role *Role) String() string {
	jsonContent, err := json.Marshal(role)
	if err != nil {
		return ""
	}

	return string(jsonContent)
}

func (groupRole *GroupRole) String() string {
	jsonContent, err := json.Marshal(groupRole)
	if err != nil {
		return ""
	}

	return string(jsonContent)
}

func (userRole *UserRole) String() string {
	jsonContent, err := json.Marshal(userRole)
	if err != nil {
		return ""
	}

	return string(jsonContent)
}
