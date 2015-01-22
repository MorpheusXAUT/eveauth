package mysql

import (
	"fmt"
	"net"
	"strconv"

	"github.com/morpheusxaut/eveauth/misc"
	"github.com/morpheusxaut/eveauth/models"

	// Blank import of the MySQL driver to use with sqlx
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type DatabaseConnection struct {
	Config *misc.Configuration

	conn *sqlx.DB
}

func (c *DatabaseConnection) Connect() error {
	conn, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", c.Config.DatabaseUser, c.Config.DatabasePassword, net.JoinHostPort(c.Config.DatabaseHost, strconv.Itoa(c.Config.DatabasePort)), c.Config.DatabaseSchema))
	if err != nil {
		return err
	}

	c.conn = conn

	return nil
}

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

func (c *DatabaseConnection) LoadAllAPIKeys() ([]*models.APIKey, error) {
	var apiKeys []*models.APIKey

	err := c.conn.Select(&apiKeys, "SELECT id, userid, apikeyid, apivcode, active FROM userapikeys")
	if err != nil {
		return nil, err
	}

	return apiKeys, nil
}

func (c *DatabaseConnection) LoadAllCorporations() ([]*models.Corporation, error) {
	var corporations []*models.Corporation

	err := c.conn.Select(&corporations, "SELECT id, name, ticker, evecorporationid, apikeyid, apivcode, active FROM corporations")
	if err != nil {
		return nil, err
	}

	return corporations, nil
}

func (c *DatabaseConnection) LoadAllCharacters() ([]*models.Character, error) {
	var characters []*models.Character

	err := c.conn.Select(&characters, "SELECT id, userid, corporationid, name, evecharacterid, active FROM characters")
	if err != nil {
		return nil, err
	}

	return characters, nil
}

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

func (c *DatabaseConnection) LoadAllRoles() ([]*models.Role, error) {
	var roles []*models.Role

	err := c.conn.Select(&roles, "SELECT id, name, active FROM roles")
	if err != nil {
		return nil, err
	}

	return roles, nil
}

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

func (c *DatabaseConnection) LoadAllUsers() ([]*models.User, error) {
	var users []*models.User

	err := c.conn.Select(&users, "SELECT id, username, password, active FROM users")
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		characters, err := c.LoadAllCharactersForUser(user.ID)
		if err != nil {
			return nil, err
		}

		apiKeys, err := c.LoadAllAPIKeysForUser(user.ID)
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

		user.Characters = characters
		user.APIKeys = apiKeys
		user.UserRoles = userRoles
		user.Groups = groups
	}

	return users, nil
}

func (c *DatabaseConnection) LoadAPIKey(apiKeyID int64) (*models.APIKey, error) {
	apiKey := &models.APIKey{}

	err := c.conn.Get(apiKey, "SELECT id, userid, apikeyid, apivcode, active FROM userapikeys WHERE id=?", apiKeyID)
	if err != nil {
		return nil, err
	}

	return apiKey, nil
}

func (c *DatabaseConnection) LoadCorporation(corporationID int64) (*models.Corporation, error) {
	corporation := &models.Corporation{}

	err := c.conn.Get(corporation, "SELECT id, name, ticker, evecorporationid, apikeyid, apivcode, active FROM corporations WHERE id=?", corporationID)
	if err != nil {
		return nil, err
	}

	return corporation, nil
}

func (c *DatabaseConnection) LoadCharacter(characterID int64) (*models.Character, error) {
	character := &models.Character{}

	err := c.conn.Get(character, "SELECT id, userid, corporationid, name, evecharacterid, active FROM characters WHERE id=?", characterID)
	if err != nil {
		return nil, err
	}

	return character, nil
}

func (c *DatabaseConnection) LoadRole(roleID int64) (*models.Role, error) {
	role := &models.Role{}

	err := c.conn.Get(role, "SELECT id, name, active FROM roles WHERE id=?", roleID)
	if err != nil {
		return nil, err
	}

	return role, nil
}

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

func (c *DatabaseConnection) LoadUser(userID int64) (*models.User, error) {
	user := &models.User{}

	err := c.conn.Get(user, "SELECT id, username, password, active FROM users WHERE id=?", userID)
	if err != nil {
		return nil, err
	}

	characters, err := c.LoadAllCharactersForUser(user.ID)
	if err != nil {
		return nil, err
	}

	apiKeys, err := c.LoadAllAPIKeysForUser(user.ID)
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

	user.Characters = characters
	user.APIKeys = apiKeys
	user.UserRoles = userRoles
	user.Groups = groups

	return user, nil
}

func (c *DatabaseConnection) LoadAllAPIKeysForUser(userID int64) ([]*models.APIKey, error) {
	var apiKeys []*models.APIKey

	err := c.conn.Select(&apiKeys, "SELECT id, userid, apikeyid, apivcode, active FROM userapikeys WHERE userid=?", userID)
	if err != nil {
		return nil, err
	}

	return apiKeys, nil
}

func (c *DatabaseConnection) LoadAllCharactersForUser(userID int64) ([]*models.Character, error) {
	var characters []*models.Character

	err := c.conn.Select(&characters, "SELECT id, userid, corporationid, name, evecharacterid, active FROM characters WHERE userid=?", userID)
	if err != nil {
		return nil, err
	}

	return characters, nil
}

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
