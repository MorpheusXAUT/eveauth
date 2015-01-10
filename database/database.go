package database

var (
	Database *DatabaseConnection
)

type DatabaseConnection interface {
	Connect() error
	RawQuery(query string, v ...interface{}) ([]interface{}, error)
}

func SetupDatabase() {

}
