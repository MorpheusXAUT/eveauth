package database

var (
	Database *DatabaseConnection
)

type DatabaseConnection interface {
}

func SetupDatabase(dbType DatabaseType) {

}
