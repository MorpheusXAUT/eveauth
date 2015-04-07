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
	// LoadAllApplications retrieves all applications from the database, returning an error if the query failed
	LoadAllApplications() ([]*models.Application, error)

	// LoadAccount retrieves the account with the given ID from the database, returning an error if the query failed
	LoadAccount(accountID int64) (*models.Account, error)
	// LoadCorporation retrieves the corporation with the given ID from the database, returning an error if the query failed
	LoadCorporation(corporationID int64) (*models.Corporation, error)
	// LoadCorporationFromEVECorporationID retrieves the corporation with the given EVE Online corporation ID from the database, returning an error if the query failed
	LoadCorporationFromEVECorporationID(eveCorporationID int64) (*models.Corporation, error)
	// LoadCorporationNameFromID retrieves the name of the corporation with the given ID, returning an error if the query failed
	LoadCorporationNameFromID(corporationID int64) (string, error)
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
	// LoadUserFromUsername retrieves the user (and its associated groups and user roles) with the given username from the database, returning an error if the query failed
	LoadUserFromUsername(username string) (*models.User, error)
	// LoadApplication retrieves the application with the given application ID from the database, returning an error if the query failed
	LoadApplication(applicationID int64) (*models.Application, error)

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
	// LoadAvailableGroupsForUser retrieves all available groups (and their associated group roles) associated with the given user from the MySQL database, returning an error if the query failed
	LoadAvailableGroupsForUser(userID int64) ([]*models.Group, error)
	// LoadAvailableUserRolesForUser retrieves all available user roles for the given user from the MySQL database, returning an error if the query failed
	LoadAvailableUserRolesForUser(userID int64) ([]*models.Role, error)
	// LoadAvailableGroupRolesForGroup retrieves all available group roles for the given group from the MySQL database, returning an error if the query failed
	LoadAvailableGroupRolesForGroup(groupID int64) ([]*models.Role, error)

	// LoadPasswordForUser retrieves the password associated with the given username from the database, returning an error if the query failed
	LoadPasswordForUser(username string) (string, error)

	// QueryUserIDExists checks whether a user with the given user ID exists in the database, returning an error if the query failed
	QueryUserIDExists(userID int64) (bool, error)
	// QueryUserNameEmailExists checks whether a user with the given username or email address exists in the database, returning an error if the query failed
	QueryUserNameEmailExists(username string, email string) (bool, error)

	// SaveAccount saves an account to the database, returning the updated model or an error if the query failed
	SaveAccount(account *models.Account) (*models.Account, error)
	// SaveCorporation saves a corporation to the database, returning the updated model or an error if the query failed
	SaveCorporation(corporation *models.Corporation) (*models.Corporation, error)
	// SaveCharacter saves a character to the database, returning the updated model or an error if the query failed
	SaveCharacter(character *models.Character) (*models.Character, error)
	// SaveRole saves a role to the database, returning the updated model or an error if the query failed
	SaveRole(role *models.Role) (*models.Role, error)
	// SaveGroupRole saves a group role to the database, returning the updated model or an error if the query failed
	SaveGroupRole(groupRole *models.GroupRole) (*models.GroupRole, error)
	// SaveUserRole saves a user role to the database, returning the updated model or an error if the query failed
	SaveUserRole(userRole *models.UserRole) (*models.UserRole, error)
	// SaveGroup saves a group to the database, returning the updated model or an error if the query failed
	SaveGroup(group *models.Group) (*models.Group, error)
	// SaveUser saves a user to the database, returning the updated model or an error if the query failed
	SaveUser(user *models.User) (*models.User, error)
	// SaveLoginAttempt saves a login attempt to the database, returning an error if the query failed
	SaveLoginAttempt(loginAttempt *models.LoginAttempt) error
	// SaveCSRFFailure saves a CSRF failure to the database, returning an error if the query failed
	SaveCSRFFailure(csrfFailure *models.CSRFFailure) error

	// RemoveUserFromGroup removes a user from the given group, updates the database and returns the updated model
	RemoveUserFromGroup(user *models.User, groupID int64) (*models.User, error)
	// RemoveAPIKeyFromUser removes an API key from the given user, updates the database and returns the updated model
	RemoveAPIKeyFromUser(user *models.User, apiKeyID int64) (*models.User, error)
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
