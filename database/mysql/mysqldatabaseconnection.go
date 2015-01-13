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
	err := c.conn.Ping()
	if err != nil {
		return nil, err
	}

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
	return nil, misc.ErrNotImplemented
}

func (c *MySQLDatabaseConnection) LoadAllCorporations() ([]*models.Corporation, error) {
	return nil, misc.ErrNotImplemented
}

func (c *MySQLDatabaseConnection) LoadAllCharacters() ([]*models.Character, error) {
	return nil, misc.ErrNotImplemented
}

func (c *MySQLDatabaseConnection) LoadAllGroupRoles() ([]*models.Role, error) {
	return nil, misc.ErrNotImplemented
}

func (c *MySQLDatabaseConnection) LoadAllUserRoles() ([]*models.Role, error) {
	return nil, misc.ErrNotImplemented
}

func (c *MySQLDatabaseConnection) LoadAllGroups() ([]*models.Group, error) {
	return nil, misc.ErrNotImplemented
}

func (c *MySQLDatabaseConnection) LoadAllUsers() ([]*models.User, error) {
	return nil, misc.ErrNotImplemented
}
