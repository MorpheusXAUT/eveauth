package models

import (
	"encoding/json"
)

type Group struct {
	ID         int64        `json:"id"`
	Name       string       `json:"name"`
	Active     bool         `json:"active"`
	GroupRoles []*GroupRole `json:"groupRoles,omitempty"`
}

func (group *Group) String() string {
	jsonContent, err := json.Marshal(group)
	if err != nil {
		return ""
	}

	return string(jsonContent)
}
