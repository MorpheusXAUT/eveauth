package database

import (
	"os"
	"strconv"
	"testing"

	"github.com/morpheusxaut/eveauth/database/mysql"
	"github.com/morpheusxaut/eveauth/misc"

	. "github.com/smartystreets/goconvey/convey"
)

func createConfig(databaseType int) *misc.Configuration {
	databaseHost := "localhost"
	if len(os.Getenv("DATABASE_HOST")) > 0 {
		databaseHost = os.Getenv("DATABASE_HOST")
	}

	databasePort := 3306
	if len(os.Getenv("DATABASE_PORT")) > 0 {
		port, err := strconv.ParseInt(os.Getenv("DATABASE_PORT"), 10, 64)
		if err == nil {
			databasePort = int(port)
		}
	}

	databaseSchema := "eveauth"
	if len(os.Getenv("DATABASE_SCHEMA")) > 0 {
		databaseSchema = os.Getenv("DATABASE_SCHEMA")
	}

	databaseUser := "eveauth"
	if len(os.Getenv("DATABASE_USER")) > 0 {
		databaseUser = os.Getenv("DATABASE_USER")
	}

	databasePassword := "eveauth"
	if len(os.Getenv("DATABASE_PASSWORD")) > 0 {
		databasePassword = os.Getenv("DATABASE_PASSWORD")
	}

	config := &misc.Configuration{
		DatabaseType:     databaseType,
		DatabaseHost:     databaseHost,
		DatabasePort:     databasePort,
		DatabaseSchema:   databaseSchema,
		DatabaseUser:     databaseUser,
		DatabasePassword: databasePassword,
		DebugLevel:       1,
		HTTPHost:         "localhost",
		HTTPPort:         5000,
	}

	return config
}

func TestMySQLDatabaseSetup(t *testing.T) {
	Convey("Running the database setup using a MySQL configuration", t, func() {
		config := createConfig(1)

		db, err := SetupDatabase(config)

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The returned DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
		})

		Convey("The returned DatabaseConnection's type should be MySQLDatabaseConnection", func() {
			mysqlDbConn := &mysql.DatabaseConnection{}
			So(db, ShouldHaveSameTypeAs, mysqlDbConn)
		})

		Convey("Connecting to the database", func() {
			err = db.Connect()

			Convey("The returned error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestInvalidDatabaseSetup(t *testing.T) {
	Convey("Running the database setup using an invalid configuration", t, func() {
		config := createConfig(1337)

		db, err := SetupDatabase(config)

		Convey("The returned error should not be nil", func() {
			So(err, ShouldNotBeNil)
		})

		Convey("The returned DatabaseConnection should be nil", func() {
			So(db, ShouldBeNil)
		})
	})
}
