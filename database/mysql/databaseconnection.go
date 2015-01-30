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

	err := c.conn.Select(&corporations, "SELECT id, name, ticker, evecorporationid, apikeyid, apivcode, active FROM corporations")
	if err != nil {
		return nil, err
	}

	return corporations, nil
}

// LoadAllCharacters retrieves all characters from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadAllCharacters() ([]*models.Character, error) {
	var characters []*models.Character

	err := c.conn.Select(&characters, "SELECT id, accountid, corporationid, name, evecharacterid, active FROM characters")
	if err != nil {
		return nil, err
	}

	return characters, nil
}

// LoadAllRoles retrieves all roles from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadAllRoles() ([]*models.Role, error) {
	var roles []*models.Role

	err := c.conn.Select(&roles, "SELECT id, name, active FROM roles")
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

	err := c.conn.Select(&users, "SELECT id, username, password, email, active FROM users")
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

	err := c.conn.Get(corporation, "SELECT id, name, ticker, evecorporationid, apikeyid, apivcode, active FROM corporations WHERE id=?", corporationID)
	if err != nil {
		return nil, err
	}

	return corporation, nil
}

// LoadCorporationFromEVECorporationID retrieves the corporation with the given EVE Online corporation ID from the database, returning an error if the query failed
func (c *DatabaseConnection) LoadCorporationFromEVECorporationID(eveCorporationID int64) (*models.Corporation, error) {
	corporation := &models.Corporation{}

	err := c.conn.Get(corporation, "SELECT id, name, ticker, evecorporationid, apikeyid, apivcode, active FROM corporations WHERE evecorporationid=?", eveCorporationID)
	if err != nil {
		return nil, err
	}

	return corporation, nil
}

// LoadCharacter retrieves the character with the given ID from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadCharacter(characterID int64) (*models.Character, error) {
	character := &models.Character{}

	err := c.conn.Get(character, "SELECT id, accountid, corporationid, name, evecharacterid, active FROM characters WHERE id=?", characterID)
	if err != nil {
		return nil, err
	}

	return character, nil
}

// LoadRole retrieves the role with the given ID from the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) LoadRole(roleID int64) (*models.Role, error) {
	role := &models.Role{}

	err := c.conn.Get(role, "SELECT id, name, active FROM roles WHERE id=?", roleID)
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
		UserID:    userRoleID,
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

	err := c.conn.Get(user, "SELECT id, username, password, email, active FROM users WHERE id=?", userID)
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

	err := c.conn.Get(user, "SELECT id, username, password, email, active FROM users WHERE username LIKE ?", username)
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

	err := c.conn.Select(&characters, "SELECT id, accountid, corporationid, name, evecharacterid, active FROM characters WHERE accountid=?", accountID)
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
		_, err := c.conn.Exec("UPDATE corporations SET name=?, ticker=?, evecorporationid=?, apikeyid=?, apivcode=?, active=? WHERE id=?", corporation.Name, corporation.Ticker, corporation.EVECorporationID, corporation.APIKeyID, corporation.APIvCode, corporation.Active, corporation.ID)
		if err != nil {
			return nil, err
		}
	} else {
		resp, err := c.conn.Exec("INSERT INTO corporations(name, ticker, evecorporationid, apikeyid, apivcode, active) VALUES(?, ?, ?, ?, ?, ?)", corporation.Name, corporation.Ticker, corporation.EVECorporationID, corporation.APIKeyID, corporation.APIvCode, corporation.Active)
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
		_, err := c.conn.Exec("UPDATE characters SET accountid=?, corporationid=?, name=?, evecharacterid=?, active=? WHERE id=?", character.AccountID, character.CorporationID, character.Name, character.EVECharacterID, character.Active, character.ID)
		if err != nil {
			return nil, err
		}
	} else {
		resp, err := c.conn.Exec("INSERT INTO characters(accountid, corporationid, name, evecharacterid, active) VALUES(?, ?, ?, ?, ?)", character.AccountID, character.CorporationID, character.Name, character.EVECharacterID, character.Active)
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
		_, err := c.conn.Exec("UPDATE roles SET name=?, active=? WHERE id=?", role.Name, role.Active, role.ID)
		if err != nil {
			return nil, err
		}
	} else {
		resp, err := c.conn.Exec("INSERT INTO roles(name, active) VALUES(?, ?)", role.Name, role.Active)
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

		for _, group := range user.Groups {
			gr, err := c.SaveGroup(group)
			if err != nil {
				return nil, err
			}

			group = gr
		}

		_, err := c.conn.Exec("UPDATE users SET username=?, password=?, email=?, active=? WHERE id=?", user.Username, user.Password, user.Email, user.Active, user.ID)
		if err != nil {
			return nil, err
		}
	} else {
		resp, err := c.conn.Exec("INSERT INTO users(username, password, email, active) VALUES(?, ?, ?, ?)", user.Username, user.Password, user.Email, user.Active)
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

// SaveAuthenticationAttempt saves a login attempt to the MySQL database, returning an error if the query failed
func (c *DatabaseConnection) SaveAuthenticationAttempt(loginAttempt *models.LoginAttempt) error {
	_, err := c.conn.Exec("INSERT INTO loginattempts(username, remoteaddr, useragent, successful) VALUES(?, ?, ?, ?)", loginAttempt.Username, loginAttempt.RemoteAddr, loginAttempt.UserAgent, loginAttempt.Successful)
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

// RemoveUserFromGroup removes a user from the given group, updates the MySQL database and returns the updated model
func (c *DatabaseConnection) RemoveUserFromGroup(user *models.User, groupID int64) (*models.User, error) {
	_, err := c.conn.Exec("DELETE FROM usergroups WHERE userid=? AND groupid=?", user.ID, groupID)
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
