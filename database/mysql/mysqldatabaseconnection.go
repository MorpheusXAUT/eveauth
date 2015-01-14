package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/morpheusxaut/eveauth/misc"
	"github.com/morpheusxaut/eveauth/models"
	"net"
	"strconv"
)

type MySQLDatabaseConnection struct {
	Config *misc.Configuration

	conn *sqlx.DB
}

func (c *MySQLDatabaseConnection) Connect() error {
	conn, err := sqlx.Connect("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", c.Config.DatabaseUser, c.Config.DatabasePassword, net.JoinHostPort(c.Config.DatabaseHost, strconv.Itoa(c.Config.DatabasePort)), c.Config.DatabaseSchema))
	if err != nil {
		return err
	}

	c.conn = conn

	err = c.conn.Ping()
	if err != nil {
		return err
	}

	return nil
}

func (c *MySQLDatabaseConnection) RawQuery(query string, v ...interface{}) ([]map[string]interface{}, error) {
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
		for i, _ := range columns {
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

func (c *MySQLDatabaseConnection) LoadAllAPIKeys() ([]*models.APIKey, error) {
	apiKeys := make([]*models.APIKey, 0)

	err := c.conn.Select(&apiKeys, "SELECT id, userid, apikeyid, apivcode, active FROM userapikeys;")
	if err != nil {
		return nil, err
	}

	return apiKeys, nil
}

func (c *MySQLDatabaseConnection) LoadAllCorporations() ([]*models.Corporation, error) {
	corporations := make([]*models.Corporation, 0)

	err := c.conn.Select(&corporations, "SELECT id, name, ticker, evecorporationid, apikeyid, apivcode, active FROM corporations;")
	if err != nil {
		return nil, err
	}

	return corporations, nil
}

func (c *MySQLDatabaseConnection) LoadAllCharacters() ([]*models.Character, error) {
	characters := make([]*models.Character, 0)

	err := c.conn.Select(&characters, "SELECT id, userid, corporationid, name, evecharacterid, active FROM characters;")
	if err != nil {
		return nil, err
	}

	return characters, nil
}

func (c *MySQLDatabaseConnection) LoadAllGroupRoles() ([]*models.GroupRole, error) {
	groupRoles := make([]*models.GroupRole, 0)

	rows, err := c.conn.Queryx("SELECT id, groupid, roleid, autoadded, granted FROM grouproles;")
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

		// TODO load role from database

		groupRole := &models.GroupRole{
			ID:        id,
			GroupID:   groupID,
			Role:      nil,
			AutoAdded: (autoadded != 0),
			Granted:   (granted != 0),
		}

		groupRoles = append(groupRoles, groupRole)
	}

	return groupRoles, nil
}

func (c *MySQLDatabaseConnection) LoadAllUserRoles() ([]*models.UserRole, error) {
	userRoles := make([]*models.UserRole, 0)

	rows, err := c.conn.Queryx("SELECT id, userid, roleid, autoadded, granted FROM userroles;")
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

		// TODO load role from database

		userRole := &models.UserRole{
			ID:        id,
			UserID:    userID,
			Role:      nil,
			AutoAdded: (autoadded != 0),
			Granted:   (granted != 0),
		}

		userRoles = append(userRoles, userRole)
	}

	return userRoles, nil
}

func (c *MySQLDatabaseConnection) LoadAllGroups() ([]*models.Group, error) {
	groups := make([]*models.Group, 0)

	err := c.conn.Select(&groups, "SELECT id, name, active FROM groups;")
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (c *MySQLDatabaseConnection) LoadAllUsers() ([]*models.User, error) {
	users := make([]*models.User, 0)

	err := c.conn.Select(&users, "SELECT id, username, password, active FROM users;")
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		user.UserRoles = make([]*models.UserRole, 0)
		user.GroupRoles = make([]*models.GroupRole, 0)
		user.Characters = make([]*models.Character, 0)
	}

	return users, nil
}
