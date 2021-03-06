package mysql

import (
	"fmt"

	"github.com/morpheusxaut/eveauth/misc"
	"github.com/morpheusxaut/eveauth/models"

	// Blank import of the MySQL driver to use with sqlx
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// DatabaseConnection provides an implementation of the Connection interface using a MySQL database
type DatabaseConnection struct {
	// Config stores the current configuration values being used
	Config *misc.Configuration

	conn *sqlx.DB
}

// Connect tries to establish a connection to the MySQL backend, returning an error if the attempt failed
func (c *DatabaseConnection) Connect() error {
	conn, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", c.Config.DatabaseUser, c.Config.DatabasePassword, c.Config.DatabaseHost, c.Config.DatabaseSchema))
	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}

// RawQuery performs a raw MySQL query and returns a map of interfaces containing the retrieve data. An error is returned if the query failed
func (c *DatabaseConnection) RawQuery(query string, v ...interface{}) ([]map[string]interface{}, error) {
	rows, err := c.conn.Query(query, v...)
	if err != nil {
		return nil, err
	}

	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	var results []map[string]interface{}

	for rows.Next() {
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		rows.Scan(valuePtrs...)

		resultRow := make(map[string]interface{})

		for i, col := range columns {
			resultRow[col] = values[i]
		}

		results = append(results, resultRow)
	}

	return results, nil
}

// LoadAllAccounts retrieves all accounts from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadAllAccounts() ([]*models.Account, error) {
	var accounts []*models.Account

	err := c.conn.Select(&accounts, "SELECT id, userid, apikeyid, apivcode, apiaccessmask, active FROM accounts")
	if err != nil {
		return nil, err
	}

	for _, account := range accounts {
		characters, err := c.LoadAllCharactersForAccount(account.ID)
		if err != nil {
			return nil, err
		}

		account.Characters = characters
	}

	return accounts, nil
}

// LoadAllCorporations retrieves all corporations from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadAllCorporations() ([]*models.Corporation, error) {
	var corporations []*models.Corporation

	err := c.conn.Select(&corporations, "SELECT id, name, ticker, evecorporationid, ceoid, apikeyid, apivcode, active FROM corporations")
	if err != nil {
		return nil, err
	}

	return corporations, nil
}

// LoadAllCharacters retrieves all characters from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadAllCharacters() ([]*models.Character, error) {
	var characters []*models.Character

	err := c.conn.Select(&characters, "SELECT id, accountid, corporationid, name, evecharacterid, defaultcharacter, active FROM characters")
	if err != nil {
		return nil, err
	}

	return characters, nil
}

// LoadAllRoles retrieves all roles from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadAllRoles() ([]*models.Role, error) {
	var roles []*models.Role

	err := c.conn.Select(&roles, "SELECT id, name, active, locked FROM roles")
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// LoadAllGroupRoles retrieves all group roles (and their associated roles) from the database, returning an error if the query failed
func (c *DatabaseConnection) LoadAllGroupRoles() ([]*models.GroupRole, error) {
	var groupRoles []*models.GroupRole

	rows, err := c.conn.Queryx("SELECT id, groupid, roleid, autoadded, granted FROM grouproles")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id, groupID, roleID int64
		var autoadded, granted int

		err = rows.Scan(&id, &groupID, &roleID, &autoadded, &granted)
		if err != nil {
			return nil, err
		}

		role, err := c.LoadRole(roleID)
		if err != nil {
			return nil, err
		}

		groupRole := &models.GroupRole{
			ID:        id,
			GroupID:   groupID,
			Role:      role,
			AutoAdded: (autoadded != 0),
			Granted:   (granted != 0),
		}

		groupRoles = append(groupRoles, groupRole)
	}

	return groupRoles, nil
}

// LoadAllUserRoles retrieves all user roles (and their associated roles) from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadAllUserRoles() ([]*models.UserRole, error) {
	var userRoles []*models.UserRole

	rows, err := c.conn.Queryx("SELECT id, userid, roleid, autoadded, granted FROM userroles")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id, userID, roleID int64
		var autoadded, granted int

		err = rows.Scan(&id, &userID, &roleID, &autoadded, &granted)
		if err != nil {
			return nil, err
		}

		role, err := c.LoadRole(roleID)
		if err != nil {
			return nil, err
		}

		userRole := &models.UserRole{
			ID:        id,
			UserID:    userID,
			Role:      role,
			AutoAdded: (autoadded != 0),
			Granted:   (granted != 0),
		}

		userRoles = append(userRoles, userRole)
	}

	return userRoles, nil
}

// LoadAllGroups retrieves all groups (and their associated group roles) from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadAllGroups() ([]*models.Group, error) {
	var groups []*models.Group

	err := c.conn.Select(&groups, "SELECT id, name, active FROM groups")
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		groupRoles, err := c.LoadAllGroupRolesForGroup(group.ID)
		if err != nil {
			return nil, err
		}

		group.GroupRoles = groupRoles
	}

	return groups, nil
}

// LoadAllUsers retrieves all users (and their associates groups and user roles) from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadAllUsers() ([]*models.User, error) {
	var users []*models.User

	err := c.conn.Select(&users, "SELECT id, username, password, email, verifiedemail, active FROM users")
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		accounts, err := c.LoadAllAccountsForUser(user.ID)
		if err != nil {
			return nil, err
		}

		userRoles, err := c.LoadAllUserRolesForUser(user.ID)
		if err != nil {
			return nil, err
		}

		groups, err := c.LoadAllGroupsForUser(user.ID)
		if err != nil {
			return nil, err
		}

		user.Accounts = accounts
		user.UserRoles = userRoles
		user.Groups = groups
	}

	return users, nil
}

// LoadAllApplications retrieves all applications from the database, returning an error if the query failed
func (c *DatabaseConnection) LoadAllApplications() ([]*models.Application, error) {
	var applications []*models.Application

	err := c.conn.Select(&applications, "SELECT id, name, maintainerid, secret, callback, active FROM applications")
	if err != nil {
		return nil, err
	}

	return applications, nil
}

// LoadAccount retrieves the account with the given ID from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadAccount(accountID int64) (*models.Account, error) {
	account := &models.Account{}

	err := c.conn.Get(account, "SELECT id, userid, apikeyid, apivcode, apiaccessmask, active FROM accounts WHERE id=?", accountID)
	if err != nil {
		return nil, err
	}

	characters, err := c.LoadAllCharactersForAccount(account.ID)
	if err != nil {
		return nil, err
	}

	account.Characters = characters

	return account, nil
}

// LoadCorporation retrieves the corporation with the given ID from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadCorporation(corporationID int64) (*models.Corporation, error) {
	corporation := &models.Corporation{}

	err := c.conn.Get(corporation, "SELECT id, name, ticker, evecorporationid, ceoid, apikeyid, apivcode, active FROM corporations WHERE id=?", corporationID)
	if err != nil {
		return nil, err
	}

	return corporation, nil
}

// LoadCorporationFromEVECorporationID retrieves the corporation with the given EVE Online corporation ID from the database, returning an error if the query failed
func (c *DatabaseConnection) LoadCorporationFromEVECorporationID(eveCorporationID int64) (*models.Corporation, error) {
	corporation := &models.Corporation{}

	err := c.conn.Get(corporation, "SELECT id, name, ticker, evecorporationid, ceoid, apikeyid, apivcode, active FROM corporations WHERE evecorporationid=?", eveCorporationID)
	if err != nil {
		return nil, err
	}

	return corporation, nil
}

// LoadCorporationNameFromID retrieves the name of the corporation with the given ID, returning an error if the query failed
func (c *DatabaseConnection) LoadCorporationNameFromID(corporationID int64) (string, error) {
	var corporationName string

	err := c.conn.Get(&corporationName, "SELECT name FROM corporations WHERE id=?", corporationID)
	if err != nil {
		return "", err
	}

	return corporationName, nil
}

// LoadCharacter retrieves the character with the given ID from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadCharacter(characterID int64) (*models.Character, error) {
	character := &models.Character{}

	err := c.conn.Get(character, "SELECT id, accountid, corporationid, name, evecharacterid, defaultcharacter, active FROM characters WHERE id=?", characterID)
	if err != nil {
		return nil, err
	}

	return character, nil
}

// LoadRole retrieves the role with the given ID from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadRole(roleID int64) (*models.Role, error) {
	role := &models.Role{}

	err := c.conn.Get(role, "SELECT id, name, active, locked FROM roles WHERE id=?", roleID)
	if err != nil {
		return nil, err
	}

	return role, nil
}

// LoadGroupRole retrieves the group role (and its associated role) with the given ID from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadGroupRole(groupRoleID int64) (*models.GroupRole, error) {
	row := c.conn.QueryRowx("SELECT id, groupid, roleid, autoadded, granted FROM grouproles WHERE id=?", groupRoleID)

	var id, groupID, roleID int64
	var autoadded, granted int

	err := row.Scan(&id, &groupID, &roleID, &autoadded, &granted)
	if err != nil {
		return nil, err
	}

	role, err := c.LoadRole(roleID)
	if err != nil {
		return nil, err
	}

	groupRole := &models.GroupRole{
		ID:        id,
		GroupID:   groupID,
		Role:      role,
		AutoAdded: (autoadded != 0),
		Granted:   (granted != 0),
	}

	return groupRole, nil
}

// LoadUserRole retrieves the user role (and its associated role) with the given ID from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadUserRole(userRoleID int64) (*models.UserRole, error) {
	row := c.conn.QueryRowx("SELECT id, userid, roleid, autoadded, granted FROM userroles WHERE id=?", userRoleID)

	var id, userID, roleID int64
	var autoadded, granted int

	err := row.Scan(&id, &userID, &roleID, &autoadded, &granted)
	if err != nil {
		return nil, err
	}

	role, err := c.LoadRole(roleID)
	if err != nil {
		return nil, err
	}

	userRole := &models.UserRole{
		ID:        id,
		UserID:    userID,
		Role:      role,
		AutoAdded: (autoadded != 0),
		Granted:   (granted != 0),
	}

	return userRole, nil
}

// LoadGroup retrieves the group (and its associated group roles) with the given ID from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadGroup(groupID int64) (*models.Group, error) {
	group := &models.Group{}

	err := c.conn.Get(group, "SELECT id, name, active FROM groups WHERE id=?", groupID)
	if err != nil {
		return nil, err
	}

	groupRoles, err := c.LoadAllGroupRolesForGroup(group.ID)
	if err != nil {
		return nil, err
	}

	group.GroupRoles = groupRoles

	return group, nil
}

// LoadUser retrieves the user (and its associated groups and user roles) with the given ID from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadUser(userID int64) (*models.User, error) {
	user := &models.User{}

	err := c.conn.Get(user, "SELECT id, username, password, email, verifiedemail, active FROM users WHERE id=?", userID)
	if err != nil {
		return nil, err
	}

	accounts, err := c.LoadAllAccountsForUser(user.ID)
	if err != nil {
		return nil, err
	}

	userRoles, err := c.LoadAllUserRolesForUser(user.ID)
	if err != nil {
		return nil, err
	}

	groups, err := c.LoadAllGroupsForUser(user.ID)
	if err != nil {
		return nil, err
	}

	user.Accounts = accounts
	user.UserRoles = userRoles
	user.Groups = groups

	return user, nil
}

// LoadUserFromUsername retrieves the user (and its associated groups and user roles) with the given username from the database, returning an error if the query failed
func (c *DatabaseConnection) LoadUserFromUsername(username string) (*models.User, error) {
	user := &models.User{}

	err := c.conn.Get(user, "SELECT id, username, password, email, verifiedemail, active FROM users WHERE username LIKE ?", username)
	if err != nil {
		return nil, err
	}

	accounts, err := c.LoadAllAccountsForUser(user.ID)
	if err != nil {
		return nil, err
	}

	userRoles, err := c.LoadAllUserRolesForUser(user.ID)
	if err != nil {
		return nil, err
	}

	groups, err := c.LoadAllGroupsForUser(user.ID)
	if err != nil {
		return nil, err
	}

	user.Accounts = accounts
	user.UserRoles = userRoles
	user.Groups = groups

	return user, nil
}

// LoadApplication retrieves the application with the given application ID from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadApplication(applicationID int64) (*models.Application, error) {
	application := &models.Application{}

	err := c.conn.Get(application, "SELECT id, name, maintainerid, secret, callback, active FROM applications WHERE id=?", applicationID)
	if err != nil {
		return nil, err
	}

	return application, nil
}

// LoadAllAccountsForUser retrieves all accounts associated with the given user from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadAllAccountsForUser(userID int64) ([]*models.Account, error) {
	var accounts []*models.Account

	err := c.conn.Select(&accounts, "SELECT id, userid, apikeyid, apivcode, apiaccessmask, active FROM accounts WHERE userid=?", userID)
	if err != nil {
		return nil, err
	}

	for _, account := range accounts {
		characters, err := c.LoadAllCharactersForAccount(account.ID)
		if err != nil {
			return nil, err
		}

		account.Characters = characters
	}

	return accounts, nil
}

// LoadAllCharactersForAccount retrieves all characters associated with the given account from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadAllCharactersForAccount(accountID int64) ([]*models.Character, error) {
	var characters []*models.Character

	err := c.conn.Select(&characters, "SELECT id, accountid, corporationid, name, evecharacterid, defaultcharacter, active FROM characters WHERE accountid=?", accountID)
	if err != nil {
		return nil, err
	}

	return characters, nil
}

// LoadAllGroupRolesForGroup retrieves all group roles (and their associated roles) associated with the given group from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadAllGroupRolesForGroup(groupID int64) ([]*models.GroupRole, error) {
	var groupRoles []*models.GroupRole

	rows, err := c.conn.Queryx("SELECT id, groupid, roleid, autoadded, granted FROM grouproles WHERE groupid=?", groupID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id, grID, roleID int64
		var autoadded, granted int

		err = rows.Scan(&id, &grID, &roleID, &autoadded, &granted)
		if err != nil {
			return nil, err
		}

		role, err := c.LoadRole(roleID)
		if err != nil {
			return nil, err
		}

		groupRole := &models.GroupRole{
			ID:        id,
			GroupID:   grID,
			Role:      role,
			AutoAdded: (autoadded != 0),
			Granted:   (granted != 0),
		}

		groupRoles = append(groupRoles, groupRole)
	}

	return groupRoles, nil
}

// LoadAllUserRolesForUser retrieves all user roles (and their associated roles) associated with the given user from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadAllUserRolesForUser(userID int64) ([]*models.UserRole, error) {
	// For whatever weird reason, only using "var userRoles []*models.UserRole" does not work in this case and throws an error...
	var userRoles []*models.UserRole
	userRoles = make([]*models.UserRole, 0)

	rows, err := c.conn.Queryx("SELECT id, userid, roleid, autoadded, granted FROM userroles WHERE userid=?", userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var id, uID, roleID int64
		var autoadded, granted int

		err = rows.Scan(&id, &uID, &roleID, &autoadded, &granted)
		if err != nil {
			return nil, err
		}

		role, err := c.LoadRole(roleID)
		if err != nil {
			return nil, err
		}

		userRole := &models.UserRole{
			ID:        id,
			UserID:    uID,
			Role:      role,
			AutoAdded: (autoadded != 0),
			Granted:   (granted != 0),
		}

		userRoles = append(userRoles, userRole)
	}

	return userRoles, nil
}

// LoadAllGroupsForUser retrieves all groups (and their associated group roles) associated with the given user from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadAllGroupsForUser(userID int64) ([]*models.Group, error) {
	// For whatever weird reason, only using "var groups []*models.Group" does not work in this case and throws an error...
	var groups []*models.Group
	groups = make([]*models.Group, 0)

	err := c.conn.Select(&groups, "SELECT g.id, g.name, g.active FROM groups AS g INNER JOIN usergroups AS ug ON (g.id = ug.groupid) WHERE ug.active=1 AND ug.userid=? GROUP BY g.id", userID)
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		groupRoles, err := c.LoadAllGroupRolesForGroup(group.ID)
		if err != nil {
			return nil, err
		}

		group.GroupRoles = groupRoles
	}

	return groups, nil
}

// LoadAvailableGroupsForUser retrieves all available groups (and their associated group roles) associated with the given user from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadAvailableGroupsForUser(userID int64) ([]*models.Group, error) {
	// For whatever weird reason, only using "var groups []*models.Group" does not work in this case and throws an error...
	var groups []*models.Group
	groups = make([]*models.Group, 0)

	err := c.conn.Select(&groups, "SELECT g.id, g.name, g.active FROM groups AS g WHERE g.id NOT IN (SELECT gi.id FROM groups AS gi INNER JOIN usergroups AS ug ON (gi.id = ug.groupid) WHERE ug.active=1 AND ug.userid=?) GROUP BY g.id ORDER BY g.name", userID)
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		groupRoles, err := c.LoadAllGroupRolesForGroup(group.ID)
		if err != nil {
			return nil, err
		}

		group.GroupRoles = groupRoles
	}

	return groups, nil
}

// LoadAvailableUserRolesForUser retrieves all available user roles for the given user from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadAvailableUserRolesForUser(userID int64) ([]*models.Role, error) {
	// For whatever weird reason, only using "var roles []*models.Role" does not work in this case and throws an error...
	var roles []*models.Role
	roles = make([]*models.Role, 0)

	err := c.conn.Select(&roles, "SELECT r.id, r.name, r.active, r.locked FROM roles AS r WHERE r.id NOT IN (SELECT ur.roleid FROM userroles AS ur WHERE ur.userid=?) GROUP BY r.id ORDER BY r.name", userID)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// LoadAvailableGroupRolesForGroup retrieves all available group roles for the given group from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadAvailableGroupRolesForGroup(groupID int64) ([]*models.Role, error) {
	// For whatever weird reason, only using "var roles []*models.Role" does not work in this case and throws an error...
	var roles []*models.Role
	roles = make([]*models.Role, 0)

	err := c.conn.Select(&roles, "SELECT r.id, r.name, r.active, r.locked FROM roles AS r WHERE r.id NOT IN (SELECT gr.roleid FROM grouproles AS gr WHERE gr.groupid=?) GROUP BY r.id ORDER BY r.name", groupID)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

// LoadAllApplicationsForUser retrieves all applications associated with the given user from the database, returning an error if the query failed
func (c *DatabaseConnection) LoadAllApplicationsForUser(userID int64) ([]*models.Application, error) {
	var applications []*models.Application

	err := c.conn.Select(&applications, "SELECT id, name, maintainerid, secret, callback, active FROM applications WHERE maintainerid=?", userID)
	if err != nil {
		return nil, err
	}

	return applications, nil
}

// LoadPasswordForUser retrieves the password associated with the given username from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadPasswordForUser(username string) (string, error) {
	row := c.conn.QueryRowx("SELECT password FROM users WHERE username LIKE ?", username)

	var password string

	err := row.Scan(&password)
	if err != nil {
		return "", err
	}

	return password, nil
}

// QueryUserIDExists checks whether a user with the given user ID exists in the database, returning an error if the query failed
func (c *DatabaseConnection) QueryUserIDExists(userID int64) (bool, error) {
	row := c.conn.QueryRowx("SELECT COUNT(id) AS count FROM users WHERE id=?", userID)

	var count int

	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	return (count > 0), nil
}

// QueryUserNameEmailExists checks whether a user with the given username or email address exists in the database, returning an error if the query failed
func (c *DatabaseConnection) QueryUserNameEmailExists(username string, email string) (bool, error) {
	row := c.conn.QueryRowx("SELECT COUNT(username) AS count FROM users WHERE username LIKE ? OR email LIKE ?", username, email)

	var count int

	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	return (count > 0), nil
}

// SaveAccount saves an account to the MySQL database, returning the updated model or an error if the query failed
func (c *DatabaseConnection) SaveAccount(account *models.Account) (*models.Account, error) {
	if account.ID > 0 {
		for _, character := range account.Characters {
			char, err := c.SaveCharacter(character)
			if err != nil {
				return nil, err
			}

			character = char
		}

		_, err := c.conn.Exec("UPDATE accounts SET userid=?, apikeyid=?, apivcode=?, apiaccessmask=?, active=? WHERE id=?", account.UserID, account.APIKeyID, account.APIvCode, account.APIAccessMask, account.Active, account.ID)
		if err != nil {
			return nil, err
		}
	} else {
		resp, err := c.conn.Exec("INSERT INTO accounts(userid, apikeyid, apivcode, apiaccessmask, active) VALUES(?, ?, ?, ?, ?)", account.UserID, account.APIKeyID, account.APIvCode, account.APIAccessMask, account.Active)
		if err != nil {
			return nil, err
		}

		lastInsertedID, err := resp.LastInsertId()
		if err != nil {
			return nil, err
		}

		account.ID = lastInsertedID

		for _, character := range account.Characters {
			character.AccountID = account.ID

			char, err := c.SaveCharacter(character)
			if err != nil {
				return nil, err
			}

			character = char
		}
	}

	return account, nil
}

// SaveCorporation saves a corporation to the MySQL database, returning the updated model or an error if the query failed
func (c *DatabaseConnection) SaveCorporation(corporation *models.Corporation) (*models.Corporation, error) {
	if corporation.ID > 0 {
		_, err := c.conn.Exec("UPDATE corporations SET name=?, ticker=?, evecorporationid=?, ceoid=?, apikeyid=?, apivcode=?, active=? WHERE id=?", corporation.Name, corporation.Ticker, corporation.EVECorporationID, corporation.CEOID, corporation.APIKeyID, corporation.APIvCode, corporation.Active, corporation.ID)
		if err != nil {
			return nil, err
		}
	} else {
		resp, err := c.conn.Exec("INSERT INTO corporations(name, ticker, evecorporationid, ceoid, apikeyid, apivcode, active) VALUES(?, ?, ?, ?, ?, ?, ?)", corporation.Name, corporation.Ticker, corporation.EVECorporationID, corporation.CEOID, corporation.APIKeyID, corporation.APIvCode, corporation.Active)
		if err != nil {
			return nil, err
		}

		lastInsertedID, err := resp.LastInsertId()
		if err != nil {
			return nil, err
		}

		corporation.ID = lastInsertedID
	}

	return corporation, nil
}

// SaveCharacter saves a character to the MySQL database, returning the updated model or an error if the query failed
func (c *DatabaseConnection) SaveCharacter(character *models.Character) (*models.Character, error) {
	if character.ID > 0 {
		_, err := c.conn.Exec("UPDATE characters SET accountid=?, corporationid=?, name=?, evecharacterid=?, defaultcharacter=?, active=? WHERE id=?", character.AccountID, character.CorporationID, character.Name, character.EVECharacterID, character.DefaultCharacter, character.Active, character.ID)
		if err != nil {
			return nil, err
		}
	} else {
		resp, err := c.conn.Exec("INSERT INTO characters(accountid, corporationid, name, evecharacterid, defaultcharacter, active) VALUES(?, ?, ?, ?, ?, ?)", character.AccountID, character.CorporationID, character.Name, character.EVECharacterID, character.DefaultCharacter, character.Active)
		if err != nil {
			return nil, err
		}

		lastInsertedID, err := resp.LastInsertId()
		if err != nil {
			return nil, err
		}

		character.ID = lastInsertedID
	}

	return character, nil
}

// SaveRole saves a role to the MySQL database, returning the updated model or an error if the query failed
func (c *DatabaseConnection) SaveRole(role *models.Role) (*models.Role, error) {
	if role.ID > 0 {
		_, err := c.conn.Exec("UPDATE roles SET name=?, active=?, locked=? WHERE id=?", role.Name, role.Active, role.Locked, role.ID)
		if err != nil {
			return nil, err
		}
	} else {
		resp, err := c.conn.Exec("INSERT INTO roles(name, active, locked) VALUES(?, ?, ?)", role.Name, role.Active, role.Locked)
		if err != nil {
			return nil, err
		}

		lastInsertedID, err := resp.LastInsertId()
		if err != nil {
			return nil, err
		}

		role.ID = lastInsertedID
	}

	return role, nil
}

// SaveGroupRole saves a group role to the MySQL database, returning the updated model or an error if the query failed
func (c *DatabaseConnection) SaveGroupRole(groupRole *models.GroupRole) (*models.GroupRole, error) {
	role, err := c.SaveRole(groupRole.Role)
	if err != nil {
		return nil, err
	}

	groupRole.Role = role

	if groupRole.ID > 0 {
		_, err = c.conn.Exec("UPDATE grouproles SET groupid=?, roleid=?, autoadded=?, granted=? WHERE id=?", groupRole.GroupID, groupRole.Role.ID, groupRole.AutoAdded, groupRole.Granted, groupRole.ID)
		if err != nil {
			return nil, err
		}
	} else {
		resp, err := c.conn.Exec("INSERT INTO grouproles(groupid, roleid, autoadded, granted) VALUES(?, ?, ?, ?)", groupRole.GroupID, groupRole.Role.ID, groupRole.AutoAdded, groupRole.Granted)
		if err != nil {
			return nil, err
		}

		lastInsertedID, err := resp.LastInsertId()
		if err != nil {
			return nil, err
		}

		groupRole.ID = lastInsertedID
	}

	return groupRole, nil
}

// SaveUserRole saves a user role to the MySQL database, returning the updated model or an error if the query failed
func (c *DatabaseConnection) SaveUserRole(userRole *models.UserRole) (*models.UserRole, error) {
	role, err := c.SaveRole(userRole.Role)
	if err != nil {
		return nil, err
	}

	userRole.Role = role

	if userRole.ID > 0 {
		_, err = c.conn.Exec("UPDATE userroles SET userid=?, roleid=?, autoadded=?, granted=? WHERE id=?", userRole.UserID, userRole.Role.ID, userRole.AutoAdded, userRole.Granted, userRole.ID)
		if err != nil {
			return nil, err
		}
	} else {
		resp, err := c.conn.Exec("INSERT INTO userroles(userid, roleid, autoadded, granted) VALUES(?, ?, ?, ?)", userRole.UserID, userRole.Role.ID, userRole.AutoAdded, userRole.Granted)
		if err != nil {
			return nil, err
		}

		lastInsertedID, err := resp.LastInsertId()
		if err != nil {
			return nil, err
		}

		userRole.ID = lastInsertedID
	}

	return userRole, nil
}

// SaveGroup saves a group to the MySQL database, returning the updated model or an error if the query failed
func (c *DatabaseConnection) SaveGroup(group *models.Group) (*models.Group, error) {
	if group.ID > 0 {
		for _, groupRole := range group.GroupRoles {
			role, err := c.SaveGroupRole(groupRole)
			if err != nil {
				return nil, err
			}

			groupRole = role
		}

		_, err := c.conn.Exec("UPDATE groups SET name=?, active=? WHERE id=?", group.Name, group.Active, group.ID)
		if err != nil {
			return nil, err
		}
	} else {
		resp, err := c.conn.Exec("INSERT INTO groups(name, active) VALUES(?, ?)", group.Name, group.Active)
		if err != nil {
			return nil, err
		}

		lastInsertedID, err := resp.LastInsertId()
		if err != nil {
			return nil, err
		}

		group.ID = lastInsertedID

		for _, groupRole := range group.GroupRoles {
			groupRole.GroupID = group.ID

			role, err := c.SaveGroupRole(groupRole)
			if err != nil {
				return nil, err
			}

			groupRole = role
		}
	}

	return group, nil
}

// SaveUser saves a user to the MySQL database, returning the updated model or an error if the query failed
func (c *DatabaseConnection) SaveUser(user *models.User) (*models.User, error) {
	if user.ID > 0 {
		for _, account := range user.Accounts {
			acc, err := c.SaveAccount(account)
			if err != nil {
				return nil, err
			}

			account = acc
		}

		for _, userRole := range user.UserRoles {
			role, err := c.SaveUserRole(userRole)
			if err != nil {
				return nil, err
			}

			userRole = role
		}

		groups, err := c.SaveAllGroupsForUser(user.ID, user.Groups)
		if err != nil {
			return nil, err
		}

		user.Groups = groups

		_, err = c.conn.Exec("UPDATE users SET username=?, password=?, email=?, verifiedemail=?, active=? WHERE id=?", user.Username, user.Password, user.Email, user.VerifiedEmail, user.Active, user.ID)
		if err != nil {
			return nil, err
		}
	} else {
		resp, err := c.conn.Exec("INSERT INTO users(username, password, email, verifiedemail, active) VALUES(?, ?, ?, ?, ?)", user.Username, user.Password, user.Email, user.VerifiedEmail, user.Active)
		if err != nil {
			return nil, err
		}

		lastInsertedID, err := resp.LastInsertId()
		if err != nil {
			return nil, err
		}

		user.ID = lastInsertedID

		for _, account := range user.Accounts {
			account.UserID = user.ID

			acc, err := c.SaveAccount(account)
			if err != nil {
				return nil, err
			}

			account = acc
		}

		for _, userRole := range user.UserRoles {
			userRole.UserID = user.ID

			role, err := c.SaveUserRole(userRole)
			if err != nil {
				return nil, err
			}

			userRole = role
		}

		groups, err := c.SaveAllGroupsForUser(user.ID, user.Groups)
		if err != nil {
			return nil, err
		}

		user.Groups = groups
	}

	return user, nil
}

// SaveApplication saves an application to the MySQL database, returning the updated model or an error if the query failed
func (c *DatabaseConnection) SaveApplication(application *models.Application) (*models.Application, error) {
	if application.ID > 0 {
		_, err := c.conn.Exec("UPDATE applications SET name=?, maintainerid=?, secret=?, callback=?, active=? WHERE id=?", application.Name, application.MaintainerID, application.Secret, application.Callback, application.Active, application.ID)
		if err != nil {
			return nil, err
		}
	} else {
		resp, err := c.conn.Exec("INSERT INTO applications(name, maintainerid, secret, callback, active) VALUES(?, ?, ?, ?, ?)", application.Name, application.MaintainerID, application.Secret, application.Callback, application.Active)
		if err != nil {
			return nil, err
		}

		lastInsertedID, err := resp.LastInsertId()
		if err != nil {
			return nil, err
		}

		application.ID = lastInsertedID
	}

	return application, nil
}

// SaveLoginAttempt saves a login attempt to the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) SaveLoginAttempt(loginAttempt *models.LoginAttempt) error {
	_, err := c.conn.Exec("INSERT INTO loginattempts(username, remoteaddr, useragent, successful) VALUES(?, ?, ?, ?)", loginAttempt.Username, loginAttempt.RemoteAddr, loginAttempt.UserAgent, loginAttempt.Successful)
	if err != nil {
		return err
	}

	return nil
}

// SaveCSRFFailure saves a CSRF failure to the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) SaveCSRFFailure(csrfFailure *models.CSRFFailure) error {
	_, err := c.conn.Exec("INSERT INTO csrffailures(userid, request) VALUES(?, ?)", csrfFailure.UserID, csrfFailure.Request)
	if err != nil {
		return err
	}

	return nil
}

// SaveAllGroupsForUser saves all group memberships for the user
func (c *DatabaseConnection) SaveAllGroupsForUser(userID int64, groups []*models.Group) ([]*models.Group, error) {
	for _, group := range groups {
		_, err := c.conn.Exec("INSERT INTO usergroups(userid, groupid, active) VALUES(?, ?, ?) ON DUPLICATE KEY UPDATE userid=?, groupid=?, active=?", userID, group.ID, true, userID, group.ID, true)
		if err != nil {
			return nil, err
		}
	}

	return groups, nil
}

// DeleteAccount removes an account and all associated characters from the MySQL database
func (c *DatabaseConnection) DeleteAccount(accountID int64) error {
	_, err := c.conn.Exec("DELETE FROM characters WHERE accountid=?", accountID)
	if err != nil {
		return err
	}

	_, err = c.conn.Exec("DELETE FROM accounts WHERE id=?", accountID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteCharacter removes a character from the MySQL database
func (c *DatabaseConnection) DeleteCharacter(characterID int64) error {
	_, err := c.conn.Exec("DELETE FROM characters WHERE id=?", characterID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteRole removes a role and all user and group roles associated from the MySQL database
func (c *DatabaseConnection) DeleteRole(roleID int64) error {
	_, err := c.conn.Exec("DELETE FROM userroles WHERE roleid=?", roleID)
	if err != nil {
		return err
	}

	_, err = c.conn.Exec("DELETE FROM grouproles WHERE roleid=?", roleID)
	if err != nil {
		return err
	}

	_, err = c.conn.Exec("DELETE FROM roles WHERE id=?", roleID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteGroupRole removes a group role from the MySQL database
func (c *DatabaseConnection) DeleteGroupRole(groupRoleID int64) error {
	_, err := c.conn.Exec("DELETE FROM grouproles WHERE id=?", groupRoleID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUserRole removes a user role from the MySQL database
func (c *DatabaseConnection) DeleteUserRole(userRoleID int64) error {
	_, err := c.conn.Exec("DELETE FROM userroles WHERE id=?", userRoleID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteGroup removes a group and all associated group memberships and roles from the MySQL database
func (c *DatabaseConnection) DeleteGroup(groupID int64) error {
	_, err := c.conn.Exec("DELETE FROM grouproles WHERE groupid=?", groupID)
	if err != nil {
		return err
	}

	_, err = c.conn.Exec("DELETE FROM usergroups WHERE groupid=?", groupID)
	if err != nil {
		return err
	}

	_, err = c.conn.Exec("DELETE FROM groups WHERE id=?", groupID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser removes a user and all assoicated group memberships, roles and accounts from the MySQL database
func (c *DatabaseConnection) DeleteUser(userID int64) error {
	_, err := c.conn.Exec("DELETE FROM usergroups WHERE userid=?", userID)
	if err != nil {
		return err
	}

	_, err = c.conn.Exec("DELETE FROM userroles WHERE userid=?", userID)
	if err != nil {
		return err
	}

	_, err = c.conn.Exec("DELETE FROM characters WHERE accountid IN (SELECT id FROM accounts WHERE userid=?)", userID)
	if err != nil {
		return err
	}

	_, err = c.conn.Exec("DELETE FROM accounts WHERE userid=?", userID)
	if err != nil {
		return err
	}

	_, err = c.conn.Exec("DELETE FROM users WHERE id=?", userID)
	if err != nil {
		return err
	}

	return nil
}

// DeleteApplication remove an application from the MySQL database
func (c *DatabaseConnection) DeleteApplication(appID int64) error {
	_, err := c.conn.Exec("DELETE FROM applications WHERE id=?", appID)
	if err != nil {
		return err
	}

	return nil
}

// RemoveUserFromGroup removes a user from the given group, updates the MySQL database and returns the updated model
func (c *DatabaseConnection) RemoveUserFromGroup(userID int64, groupID int64) (*models.User, error) {
	user, err := c.LoadUser(userID)
	if err != nil {
		return nil, err
	}

	_, err = c.conn.Exec("DELETE FROM usergroups WHERE userid=? AND groupid=?", user.ID, groupID)
	if err != nil {
		return nil, err
	}

	var groups []*models.Group

	for _, group := range user.Groups {
		if group.ID != groupID {
			groups = append(groups, group)
		}
	}

	user.Groups = groups

	return user, nil
}

// RemoveUserRoleFromUser removes a user role from the given user, updates the database and returns the updated model
func (c *DatabaseConnection) RemoveUserRoleFromUser(userID int64, roleID int64) (*models.User, error) {
	user, err := c.LoadUser(userID)
	if err != nil {
		return nil, err
	}

	_, err = c.conn.Exec("DELETE FROM userroles WHERE userid=? AND id=?", user.ID, roleID)
	if err != nil {
		return nil, err
	}

	var userRoles []*models.UserRole

	for _, userRole := range user.UserRoles {
		if userRole.ID != roleID {
			userRoles = append(userRoles, userRole)
		}
	}

	user.UserRoles = userRoles

	return user, nil
}

// RemoveGroupRoleFromGroup removes a group role from the given group, updates the database and returns the updated model
func (c *DatabaseConnection) RemoveGroupRoleFromGroup(groupID int64, roleID int64) (*models.Group, error) {
	group, err := c.LoadGroup(groupID)
	if err != nil {
		return nil, err
	}

	_, err = c.conn.Exec("DELETE FROM grouproles WHERE groupid=? AND id=?", group.ID, roleID)
	if err != nil {
		return nil, err
	}

	var groupRoles []*models.GroupRole

	for _, groupRole := range group.GroupRoles {
		if groupRole.ID != roleID {
			groupRoles = append(groupRoles, groupRole)
		}
	}

	group.GroupRoles = groupRoles

	return group, nil
}

// RemoveAPIKeyFromUser removes an API key from the given user, updates the MySQL database and returns the updated model
func (c *DatabaseConnection) RemoveAPIKeyFromUser(user *models.User, apiKeyID int64) (*models.User, error) {
	for index, account := range user.Accounts {
		if account.APIKeyID == apiKeyID {
			for _, character := range account.Characters {
				_, err := c.conn.Exec("DELETE FROM characters WHERE id=? AND accountid=?", character.ID, account.ID)
				if err != nil {
					return nil, err
				}
			}

			_, err := c.conn.Exec("DELETE FROM accounts WHERE id=? AND apikeyid=?", account.ID, apiKeyID)
			if err != nil {
				return nil, err
			}

			user.Accounts[index], user.Accounts[len(user.Accounts)-1], user.Accounts = user.Accounts[len(user.Accounts)-1], nil, user.Accounts[:len(user.Accounts)-1]

			break
		}
	}

	return user, nil
}

// ToggleUserRoleGranted toggles the granted state of the given user role
func (c *DatabaseConnection) ToggleUserRoleGranted(roleID int64) (*models.UserRole, error) {
	userRole, err := c.LoadUserRole(roleID)
	if err != nil {
		return nil, err
	}

	userRole.Granted = !userRole.Granted

	userRole, err = c.SaveUserRole(userRole)
	if err != nil {
		return nil, err
	}

	return userRole, nil
}

// ToggleGroupRoleGranted toggles the granted state of the given group role
func (c *DatabaseConnection) ToggleGroupRoleGranted(roleID int64) (*models.GroupRole, error) {
	groupRole, err := c.LoadGroupRole(roleID)
	if err != nil {
		return nil, err
	}

	groupRole.Granted = !groupRole.Granted

	groupRole, err = c.SaveGroupRole(groupRole)
	if err != nil {
		return nil, err
	}

	return groupRole, nil
}
