package models

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
