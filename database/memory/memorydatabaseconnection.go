package memory

import (
	"github.com/morpheusxaut/eveauth/misc"
)

type MemoryDatabaseConnection struct {
	Config *misc.Configuration
}

func (connection *MemoryDatabaseConnection) Connect() error {
	return misc.ErrNotImplemented
}

func (connection *MemoryDatabaseConnection) RawQuery(query string, v ...interface{}) ([]interface{}, error) {
	return nil, misc.ErrNotImplemented
}
