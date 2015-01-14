package mock

import (
	"github.com/morpheusxaut/eveauth/misc"
	"github.com/morpheusxaut/eveauth/models"
)

type MockDatabaseConnection struct {
}

func (c *MockDatabaseConnection) Connect() error {
	return misc.ErrNotImplemented
}

func (c *MockDatabaseConnection) RawQuery(query string, v ...interface{}) ([]map[string]interface{}, error) {
	return nil, misc.ErrNotImplemented
}

func (c *MockDatabaseConnection) LoadAllAPIKeys() ([]*models.APIKey, error) {
	return nil, misc.ErrNotImplemented
}

func (c *MockDatabaseConnection) LoadAllCorporations() ([]*models.Corporation, error) {
	return nil, misc.ErrNotImplemented
}

func (c *MockDatabaseConnection) LoadAllCharacters() ([]*models.Character, error) {
	return nil, misc.ErrNotImplemented
}

func (c *MockDatabaseConnection) LoadAllGroupRoles() ([]*models.GroupRole, error) {
	return nil, misc.ErrNotImplemented
}

func (c *MockDatabaseConnection) LoadAllUserRoles() ([]*models.UserRole, error) {
	return nil, misc.ErrNotImplemented
}

func (c *MockDatabaseConnection) LoadAllGroups() ([]*models.Group, error) {
	return nil, misc.ErrNotImplemented
}

func (c *MockDatabaseConnection) LoadAllUsers() ([]*models.User, error) {
	return nil, misc.ErrNotImplemented
}
