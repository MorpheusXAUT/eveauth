package mock

import (
	"github.com/morpheusxaut/eveauth/misc"
)

type MockDatabaseConnection struct {
}

func (connection *MockDatabaseConnection) Connect() error {
	return misc.ErrNotImplemented
}

func (connection *MockDatabaseConnection) RawQuery(query string, v ...interface{}) ([]interface{}, error) {
	return nil, misc.ErrNotImplemented
}
