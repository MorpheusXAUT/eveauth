package database

import (
	"fmt"

	"github.com/morpheusxaut/eveauth/database/mysql"
	"github.com/morpheusxaut/eveauth/misc"
	"github.com/morpheusxaut/eveauth/models"
)

// Connection provides an interface for communicating with a database backend in order to retrieve and persist the needed information
type Connection interface {
	// Connect tries to establish a connection to the database backend, returning an error if the attempt failed
	Connect() error

	// RawQuery performs a raw database query and returns a map of interfaces containing the retrieve data. An error is returned if the query failed
	RawQuery(query string, v ...interface{}) ([]map[string]interface{}, error)

	// LoadAllAccounts retrieves all accounts from the database, returning an error if the query failed
	LoadAllAccounts() ([]*models.Account, error)
	// LoadAllCorporations retrieves all corporations from the database, returning an error if the query failed
	LoadAllCorporations() ([]*models.Corporation, error)
	// LoadAllCharacters retrieves all characters from the database, returning an error if the query failed
	LoadAllCharacters() ([]*models.Character, error)
	// LoadAllRoles retrieves all roles from the database, returning an error if the query failed
	LoadAllRoles() ([]*models.Role, error)
	// LoadAllGroupRoles retrieves all group roles (and their associated roles) from the database, returning an error if the query failed
	LoadAllGroupRoles() ([]*models.GroupRole, error)
	// LoadAllUserRoles retrieves all user roles (and their associated roles) from the database, returning an error if the query failed
	LoadAllUserRoles() ([]*models.UserRole, error)
	// LoadAllGroups retrieves all groups (and their associated group roles) from the database, returning an error if the query failed
	LoadAllGroups() ([]*models.Group, error)
	// LoadAllUsers retrieves all users (and their associates groups and user roles) from the database, returning an error if the query failed
	LoadAllUsers() ([]*models.User, error)

	// LoadAccount retrieves the account with the given ID from the database, returning an error if the query failed
	LoadAccount(accountID int64) (*models.Account, error)
	// LoadCorporation retrieves the corporation with the given ID from the database, returning an error if the query failed
	LoadCorporation(corporationID int64) (*models.Corporation, error)
	// LoadCharacter retrieves the character with the given ID from the database, returning an error if the query failed
	LoadCharacter(characterID int64) (*models.Character, error)
	// LoadRole retrieves the role with the given ID from the database, returning an error if the query failed
	LoadRole(roleID int64) (*models.Role, error)
	// LoadGroupRole retrieves the group role (and its associated role) with the given ID from the database, returning an error if the query failed
	LoadGroupRole(groupRoleID int64) (*models.GroupRole, error)
	// LoadUserRole retrieves the user role (and its associated role) with the given ID from the database, returning an error if the query failed
	LoadUserRole(userRoleID int64) (*models.UserRole, error)
	// LoadGroup retrieves the group (and its associated group roles) with the given ID from the database, returning an error if the query failed
	LoadGroup(groupID int64) (*models.Group, error)
	// LoadUser retrieves the user (and its associated groups and user roles) with the given ID from the database, returning an error if the query failed
	LoadUser(userID int64) (*models.User, error)

	// LoadAllAccountsForUser retrieves all accounts associated with the given user from the database, returning an error if the query failed
	LoadAllAccountsForUser(userID int64) ([]*models.Account, error)
	// LoadAllCharactersForAccount retrieves all characters associated with the given account from the database, returning an error if the query failed
	LoadAllCharactersForAccount(accountID int64) ([]*models.Character, error)
	// LoadAllGroupRolesForGroup retrieves all group roles (and their associated roles) associated with the given group from the database, returning an error if the query failed
	LoadAllGroupRolesForGroup(groupID int64) ([]*models.GroupRole, error)
	// LoadAllUserRolesForUser retrieves all user roles (and their associated roles) associated with the given user from the database, returning an error if the query failed
	LoadAllUserRolesForUser(userID int64) ([]*models.UserRole, error)
	// LoadAllGroupsForUser retrieves all groups (and their associated group roles) associated with the given user from the database, returning an error if the query failed
	LoadAllGroupsForUser(userID int64) ([]*models.Group, error)
}

// SetupDatabase parses the database type set in the configuration and returns an appropriate database implementation or an error if the type is unknown
func SetupDatabase(conf *misc.Configuration) (Connection, error) {
	var database Connection

	switch Type(conf.DatabaseType) {
	case TypeMySQL:
		database = &mysql.DatabaseConnection{
			Config: conf,
		}
		break
	default:
		return nil, fmt.Errorf("Unknown type #%d", conf.DatabaseType)
	}

	return database, nil
}