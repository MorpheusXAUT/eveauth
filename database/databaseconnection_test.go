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

			Convey("Performing some queries", func() {
				apiKeys, err := db.LoadAllAPIKeys()

				Convey("Loading all API keys should return no error", func() {
					So(err, ShouldBeNil)

					Convey("The returned slice should not be nil", func() {
						So(apiKeys, ShouldNotBeNil)
					})

					Convey("The length of the returned slice should be 6", func() {
						So(len(apiKeys), ShouldBeGreaterThan, 0)
						So(len(apiKeys), ShouldEqual, 6)
					})
				})

				corporations, err := db.LoadAllCorporations()

				Convey("Loading all corporations should return ErrNotImplemented", func() {
					So(err, ShouldEqual, misc.ErrNotImplemented)

					Convey("The returned slice should be nil", func() {
						So(corporations, ShouldBeNil)
					})
				})

				characters, err := db.LoadAllCharacters()

				Convey("Loading all characters should return ErrNotImplemented", func() {
					So(err, ShouldEqual, misc.ErrNotImplemented)

					Convey("The returned slice should be nil", func() {
						So(characters, ShouldBeNil)
					})
				})

				groupRoles, err := db.LoadAllGroupRoles()

				Convey("Loading all group roles should return ErrNotImplemented", func() {
					So(err, ShouldEqual, misc.ErrNotImplemented)

					Convey("The returned slice should be nil", func() {
						So(groupRoles, ShouldBeNil)
					})
				})

				userRoles, err := db.LoadAllUserRoles()

				Convey("Loading all user roles should return ErrNotImplemented", func() {
					So(err, ShouldEqual, misc.ErrNotImplemented)

					Convey("The returned slice should be nil", func() {
						So(userRoles, ShouldBeNil)
					})
				})

				groups, err := db.LoadAllGroups()

				Convey("Loading all groups should return ErrNotImplemented", func() {
					So(err, ShouldEqual, misc.ErrNotImplemented)

					Convey("The returned slice should be nil", func() {
						So(groups, ShouldBeNil)
					})
				})

				users, err := db.LoadAllUsers()

				Convey("Loading all users should return ErrNotImplemented", func() {
					So(err, ShouldEqual, misc.ErrNotImplemented)

					Convey("The returned slice should be nil", func() {
						So(users, ShouldBeNil)
					})
				})
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
