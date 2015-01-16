package database

import (
	"fmt"
	"github.com/morpheusxaut/eveauth/database/mysql"
	"github.com/morpheusxaut/eveauth/misc"
	"github.com/morpheusxaut/eveauth/models"
)

type DatabaseConnection interface {
	Connect() error

	RawQuery(query string, v ...interface{}) ([]map[string]interface{}, error)

	LoadAllAPIKeys() ([]*models.APIKey, error)
	LoadAllCorporations() ([]*models.Corporation, error)
	LoadAllCharacters() ([]*models.Character, error)
	LoadAllRoles() ([]*models.Role, error)
	LoadAllGroupRoles() ([]*models.GroupRole, error)
	LoadAllUserRoles() ([]*models.UserRole, error)
	LoadAllGroups() ([]*models.Group, error)
	LoadAllUsers() ([]*models.User, error)

	LoadAPIKey(apiKeyID int64) (*models.APIKey, error)
	LoadCorporation(corporationID int64) (*models.Corporation, error)
	LoadCharacter(characterID int64) (*models.Character, error)
	LoadRole(roleID int64) (*models.Role, error)
	LoadGroupRole(groupRoleID int64) (*models.GroupRole, error)
	LoadUserRole(userRoleID int64) (*models.UserRole, error)
	LoadGroup(groupID int64) (*models.Group, error)
	LoadUser(userID int64) (*models.User, error)

	LoadAllAPIKeysForUser(userID int64) ([]*models.APIKey, error)
	LoadAllCharactersForUser(userID int64) ([]*models.Character, error)
	LoadAllGroupRolesForGroup(groupID int64) ([]*models.GroupRole, error)
	LoadAllUserRolesForUser(userID int64) ([]*models.UserRole, error)
	LoadAllGroupsForUser(userID int64) ([]*models.Group, error)
}

func SetupDatabase(conf *misc.Configuration) (DatabaseConnection, error) {
	var database DatabaseConnection

	switch DatabaseType(conf.DatabaseType) {
	case DatabaseTypeMySQL:
		database = &mysql.MySQLDatabaseConnection{
			Config: conf,
		}
		break
	default:
		return nil, fmt.Errorf("Unknown DatabaseType #%d", conf.DatabaseType)
	}

	return database, nil
}
