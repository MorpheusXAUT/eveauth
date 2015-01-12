package mysql

import (
	"github.com/morpheusxaut/eveauth/misc"
)

type MySQLDatabaseConnection struct {
	Config *misc.Configuration
}

func (connection *MySQLDatabaseConnection) Connect() error {
	return misc.ErrNotImplemented
}

func (connection *MySQLDatabaseConnection) RawQuery(query string, v ...interface{}) ([]interface{}, error) {
	return nil, misc.ErrNotImplemented
}
