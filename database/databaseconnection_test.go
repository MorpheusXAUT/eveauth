package database

import (
	"github.com/morpheusxaut/eveauth/database/memory"
	"github.com/morpheusxaut/eveauth/database/mock"
	"github.com/morpheusxaut/eveauth/database/mysql"
	"github.com/morpheusxaut/eveauth/misc"
	. "github.com/smartystreets/goconvey/convey"
	"os"
	"strconv"
	"testing"
)

func createConfig(databaseType int) *misc.Configuration {
	mysqlHost := "localhost"
	if len(os.Getenv("MYSQL_HOST")) > 0 {
		mysqlHost = os.Getenv("MYSQL_HOST")
	}

	mysqlPort := 3306
	if len(os.Getenv("MYSQL_PORT")) > 0 {
		port, err := strconv.ParseInt(os.Getenv("MYSQL_PORT"), 10, 64)
		if err == nil {
			mysqlPort = int(port)
		}
	}

	mysqlSchema := "eveauth"
	if len(os.Getenv("MYSQL_SCHEMA")) > 0 {
		mysqlSchema = os.Getenv("MYSQL_SCHEMA")
	}

	mysqlUser := "eveauth"
	if len(os.Getenv("MYSQL_USER")) > 0 {
		mysqlUser = os.Getenv("MYSQL_USER")
	}

	mysqlPassword := "eveauth"
	if len(os.Getenv("MYSQL_PASSWORD")) > 0 {
		mysqlPassword = os.Getenv("MYSQL_PASSWORD")
	}

	config := &misc.Configuration{
		DatabaseType:     databaseType,
		DatabaseHost:     mysqlHost,
		DatabasePort:     mysqlPort,
		DatabaseSchema:   mysqlSchema,
		DatabaseUser:     mysqlUser,
		DatabasePassword: mysqlPassword,
		DebugLevel:       1,
		HTTPHost:         "localhost",
		HTTPPort:         5000,
	}

	return config
}

func TestMockDatabaseSetup(t *testing.T) {
	Convey("Running the database setup using a mock configuration", t, func() {
		config := createConfig(-1)

		db, err := SetupDatabase(config)

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The returned DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
		})

		Convey("The returned DatabaseConnection's type should be MockDatabaseConnection", func() {
			mockDbConn := &mock.MockDatabaseConnection{}
			So(db, ShouldHaveSameTypeAs, mockDbConn)
		})

		Convey("Connecting to the database", func() {
			err = db.Connect()

			Convey("The returned error should be ErrNotImplemented", func() {
				So(err, ShouldEqual, misc.ErrNotImplemented)
			})
		})
	})
}

func TestMemoryDatabaseSetup(t *testing.T) {
	Convey("Running the database setup using a memory configuration", t, func() {
		config := createConfig(0)

		db, err := SetupDatabase(config)

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The returned DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
		})

		Convey("The returned DatabaseConnection's type should be MemoryDatabaseConnection", func() {
			memoryDbConn := &memory.MemoryDatabaseConnection{}
			So(db, ShouldHaveSameTypeAs, memoryDbConn)
		})

		Convey("Connecting to the database", func() {
			err = db.Connect()

			Convey("The returned error should be ErrNotImplemented", func() {
				So(err, ShouldEqual, misc.ErrNotImplemented)
			})
		})
	})
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
			mysqlDbConn := &mysql.MySQLDatabaseConnection{}
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
