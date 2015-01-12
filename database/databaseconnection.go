package database

import (
	"fmt"
	"github.com/morpheusxaut/eveauth/database/memory"
	"github.com/morpheusxaut/eveauth/database/mock"
	"github.com/morpheusxaut/eveauth/database/mysql"
	"github.com/morpheusxaut/eveauth/misc"
)

type DatabaseConnection interface {
	Connect() error
	RawQuery(query string, v ...interface{}) ([]interface{}, error)
}

func SetupDatabase(conf *misc.Configuration) (DatabaseConnection, error) {
	var database DatabaseConnection

	switch DatabaseType(conf.DatabaseType) {
	case DatabaseTypeMock:
		database = &mock.MockDatabaseConnection{}
		break
	case DatabaseTypeMemory:
		database = &memory.MemoryDatabaseConnection{
			Config: conf,
		}
		break
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
