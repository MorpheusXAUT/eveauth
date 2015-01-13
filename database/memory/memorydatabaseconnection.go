package memory

import (
	"github.com/morpheusxaut/eveauth/misc"
	"github.com/morpheusxaut/eveauth/models"
)

type MemoryDatabaseConnection struct {
	Config *misc.Configuration
}

func (c *MemoryDatabaseConnection) Connect() error {
	return misc.ErrNotImplemented
}

func (c *MemoryDatabaseConnection) RawQuery(query string, v ...interface{}) ([]map[string]interface{}, error) {
	return nil, misc.ErrNotImplemented
}

func (c *MemoryDatabaseConnection) LoadAllAPIKeys() ([]*models.APIKey, error) {
	return nil, misc.ErrNotImplemented
}

func (c *MemoryDatabaseConnection) LoadAllCorporations() ([]*models.Corporation, error) {
	return nil, misc.ErrNotImplemented
}

func (c *MemoryDatabaseConnection) LoadAllCharacters() ([]*models.Character, error) {
	return nil, misc.ErrNotImplemented
}

func (c *MemoryDatabaseConnection) LoadAllGroupRoles() ([]*models.Role, error) {
	return nil, misc.ErrNotImplemented
}

func (c *MemoryDatabaseConnection) LoadAllUserRoles() ([]*models.Role, error) {
	return nil, misc.ErrNotImplemented
}

func (c *MemoryDatabaseConnection) LoadAllGroups() ([]*models.Group, error) {
	return nil, misc.ErrNotImplemented
}

func (c *MemoryDatabaseConnection) LoadAllUsers() ([]*models.User, error) {
	return nil, misc.ErrNotImplemented
}
