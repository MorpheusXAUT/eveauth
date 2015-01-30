package models

import (
	"encoding/json"
)

// Application represents an application registered with the auth backend
type Application struct {
	// ID represents the app ID in the database
	ID int64 `json:"id"`
	// Name represents the name of the app
	Name string `json:"name"`
	// Maintainer represents the maintainer of the app
	Maintainer string `json:"maintainer"`
	// Secret represents the application's secret used to "authenticate" with the auth backend
	Secret string `json:"-"`
	// Callback represents the defined callback URL for the app
	Callback string `json:"callback"`
	// Active indicates whether the app is set as active
	Active bool `json:"active"`
}

// NewApplication creates a new application with the given information
func NewApplication(name string, maintainer string, secret string, callback string, active bool) *Application {
	application := &Application{
		ID:         -1,
		Name:       name,
		Maintainer: maintainer,
		Secret:     secret,
		Callback:   callback,
		Active:     active,
	}

	return application
}

// String represents a JSON encoded representation of the app
func (application *Application) String() string {
	jsonContent, err := json.Marshal(application)
	if err != nil {
		return ""
	}

	return string(jsonContent)
}
