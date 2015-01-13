package database

import (
	"github.com/morpheusxaut/eveauth/database/memory"
	"github.com/morpheusxaut/eveauth/database/mock"
	"github.com/morpheusxaut/eveauth/database/mysql"
	"github.com/morpheusxaut/eveauth/misc"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestMockDatabaseSetup(t *testing.T) {
	Convey("Running the database setup using a mock configuration", t, func() {
		config := &misc.Configuration{
			DatabaseType:     -1,
			DatabaseHost:     "localhost",
			DatabasePort:     3306,
			DatabaseSchema:   "eveauth",
			DatabaseUser:     "eveauth",
			DatabasePassword: "eveauth",
			DebugLevel:       1,
			HTTPHost:         "localhost",
			HTTPPort:         5000,
		}

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
		config := &misc.Configuration{
			DatabaseType:     0,
			DatabaseHost:     "localhost",
			DatabasePort:     3306,
			DatabaseSchema:   "eveauth",
			DatabaseUser:     "eveauth",
			DatabasePassword: "eveauth",
			DebugLevel:       1,
			HTTPHost:         "localhost",
			HTTPPort:         5000,
		}

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
		config := &misc.Configuration{
			DatabaseType:     1,
			DatabaseHost:     "localhost",
			DatabasePort:     3306,
			DatabaseSchema:   "eveauth",
			DatabaseUser:     "eveauth",
			DatabasePassword: "eveauth",
			DebugLevel:       1,
			HTTPHost:         "localhost",
			HTTPPort:         5000,
		}

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
		config := &misc.Configuration{
			DatabaseType:     1337,
			DatabaseHost:     "localhost",
			DatabasePort:     3306,
			DatabaseSchema:   "eveauth",
			DatabaseUser:     "eveauth",
			DatabasePassword: "eveauth",
			DebugLevel:       1,
			HTTPHost:         "localhost",
			HTTPPort:         5000,
		}

		db, err := SetupDatabase(config)

		Convey("The returned error should not be nil", func() {
			So(err, ShouldNotBeNil)
		})

		Convey("The returned DatabaseConnection should be nil", func() {
			So(db, ShouldBeNil)
		})
	})
}
