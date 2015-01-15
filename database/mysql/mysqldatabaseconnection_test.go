package mysql

import (
	"fmt"
	"github.com/morpheusxaut/eveauth/misc"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
	"testing"
)

func createConnection() *MySQLDatabaseConnection {
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
		DatabaseType:     1,
		DatabaseHost:     mysqlHost,
		DatabasePort:     mysqlPort,
		DatabaseSchema:   mysqlSchema,
		DatabaseUser:     mysqlUser,
		DatabasePassword: mysqlPassword,
		DebugLevel:       1,
		HTTPHost:         "localhost",
		HTTPPort:         5000,
	}

	db := &MySQLDatabaseConnection{
		Config: config,
	}

	return db
}

func TestMySQLDatabaseConnectionConnect(t *testing.T) {
	Convey("Connecting to a MySQL database with a MySQL configuration", t, func() {
		db := createConnection()

		Convey("Connecting to the database", func() {
			err := db.Connect()

			Convey("The returned error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestMySQLDatabaseConnectionInvalidConnect(t *testing.T) {
	Convey("Connecting to a MySQL database with an invalid MySQL configuration", t, func() {
		config := &misc.Configuration{
			DatabaseType:     1,
			DatabaseHost:     "does.not.exist",
			DatabasePort:     123456,
			DatabaseSchema:   "nonexistent",
			DatabaseUser:     "nonexistent",
			DatabasePassword: "nonexistent",
			DebugLevel:       1,
			HTTPHost:         "localhost",
			HTTPPort:         5000,
		}

		db := &MySQLDatabaseConnection{
			Config: config,
		}

		Convey("Connecting to the database", func() {
			err := db.Connect()

			Convey("The returned error should be not be", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}

func TestMySQLDatabaseConnectionRawQuery(t *testing.T) {
	Convey("Performing a raw query at a MySQL database with a MySQL configuration", t, func() {
		db := createConnection()

		Convey("Connecting to the database", func() {
			err := db.Connect()

			Convey("The returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Performing a raw query of the users table", func() {
				result, err := db.RawQuery("SELECT * FROM users;")

				Convey("The returned error should be nil", func() {
					So(err, ShouldBeNil)
				})

				Convey("The returned map should not be nil", func() {
					So(result, ShouldNotBeNil)
				})

				Convey("The returned map should have 4 entries", func() {
					So(len(result), ShouldBeGreaterThan, 0)
					So(len(result), ShouldEqual, 4)
				})

				Convey("Iterating over the result map", func() {
					for key, value := range result {
						Convey(fmt.Sprintf("The raw data table for key %v should have 4 entries", key), func() {
							So(len(value), ShouldBeGreaterThan, 0)
							So(len(value), ShouldEqual, 4)
						})
					}
				})
			})
		})
	})
}

func TestMySQLDatabaseConnectionRawInvalidQuery(t *testing.T) {
	Convey("Performing a raw invalid query at a MySQL database with a MySQL configuration", t, func() {
		db := createConnection()

		Convey("Connecting to the database", func() {
			err := db.Connect()

			Convey("The returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			Convey("Performing a raw invalid query of the users table", func() {
				result, err := db.RawQuery("SELECT nonexistent FROM users;")

				Convey("The returned error should not be nil", func() {
					So(err, ShouldNotBeNil)
				})

				Convey("The returned map should be nil", func() {
					So(result, ShouldBeNil)
				})
			})
		})
	})
}

func TestMySQLDatabaseConnectionLoadAllAPIKeys(t *testing.T) {
	Convey("Loading all API keys from a MySQL database with a MySQL configuration", t, func() {
		db := createConnection()

		Convey("Connecting to the database", func() {
			err := db.Connect()

			Convey("The returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

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

				Convey("The returned API keys should match the test data set", func() {
					Convey("Verifying entry #1", func() {
						So(apiKeys[0].ID, ShouldEqual, 1)
						So(apiKeys[0].UserID, ShouldEqual, 1)
						So(apiKeys[0].APIKeyID, ShouldEqual, 1)
						So(apiKeys[0].APIvCode, ShouldEqual, "a")
						So(apiKeys[0].Active, ShouldBeTrue)
					})

					Convey("Verifying entry #2", func() {
						So(apiKeys[1].ID, ShouldEqual, 2)
						So(apiKeys[1].UserID, ShouldEqual, 2)
						So(apiKeys[1].APIKeyID, ShouldEqual, 2)
						So(apiKeys[1].APIvCode, ShouldEqual, "b")
						So(apiKeys[1].Active, ShouldBeFalse)
					})

					Convey("Verifying entry #3", func() {
						So(apiKeys[2].ID, ShouldEqual, 3)
						So(apiKeys[2].UserID, ShouldEqual, 3)
						So(apiKeys[2].APIKeyID, ShouldEqual, 3)
						So(apiKeys[2].APIvCode, ShouldEqual, "c")
						So(apiKeys[2].Active, ShouldBeTrue)
					})

					Convey("Verifying entry #4", func() {
						So(apiKeys[3].ID, ShouldEqual, 4)
						So(apiKeys[3].UserID, ShouldEqual, 3)
						So(apiKeys[3].APIKeyID, ShouldEqual, 4)
						So(apiKeys[3].APIvCode, ShouldEqual, "d")
						So(apiKeys[3].Active, ShouldBeTrue)
					})

					Convey("Verifying entry #5", func() {
						So(apiKeys[4].ID, ShouldEqual, 5)
						So(apiKeys[4].UserID, ShouldEqual, 4)
						So(apiKeys[4].APIKeyID, ShouldEqual, 5)
						So(apiKeys[4].APIvCode, ShouldEqual, "e")
						So(apiKeys[4].Active, ShouldBeFalse)
					})

					Convey("Verifying entry #6", func() {
						So(apiKeys[5].ID, ShouldEqual, 6)
						So(apiKeys[5].UserID, ShouldEqual, 4)
						So(apiKeys[5].APIKeyID, ShouldEqual, 6)
						So(apiKeys[5].APIvCode, ShouldEqual, "f")
						So(apiKeys[5].Active, ShouldBeFalse)
					})
				})
			})
		})
	})
}

func TestMySQLDatabaseConnectionLoadAllCorporations(t *testing.T) {
	Convey("Loading all corporations from a MySQL database with a MySQL configuration", t, func() {
		db := createConnection()

		Convey("Connecting to the database", func() {
			err := db.Connect()

			Convey("The returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			corporations, err := db.LoadAllCorporations()

			Convey("Loading all corporations should return no error", func() {
				So(err, ShouldBeNil)

				Convey("The returned slice should not be nil", func() {
					So(corporations, ShouldNotBeNil)
				})

				Convey("The length of the returned slice should be 2", func() {
					So(len(corporations), ShouldBeGreaterThan, 0)
					So(len(corporations), ShouldEqual, 2)
				})

				Convey("The returned corporations should match the test data set", func() {
					Convey("Verifying entry #1", func() {
						So(corporations[0].ID, ShouldEqual, 1)
						So(corporations[0].Name, ShouldEqual, "Test Corp Please Ignore")
						So(corporations[0].Ticker, ShouldEqual, "TEST")
						So(corporations[0].EVECorporationID, ShouldEqual, 1)
						So(corporations[0].APIKeyID.IsZero(), ShouldBeFalse)
						So(corporations[0].APIKeyID.Valid, ShouldBeTrue)
						So(corporations[0].APIKeyID.Int64, ShouldEqual, 1)
						So(corporations[0].APIvCode.IsZero(), ShouldBeFalse)
						So(corporations[0].APIvCode.Valid, ShouldBeTrue)
						So(corporations[0].APIvCode.String, ShouldEqual, "a")
						So(corporations[0].Active, ShouldBeTrue)
					})

					Convey("Verifying entry #2", func() {
						So(corporations[1].ID, ShouldEqual, 2)
						So(corporations[1].Name, ShouldEqual, "Corp Test Ignore Please")
						So(corporations[1].Ticker, ShouldEqual, "CORP")
						So(corporations[1].EVECorporationID, ShouldEqual, 2)
						So(corporations[1].APIKeyID.IsZero(), ShouldBeTrue)
						So(corporations[1].APIKeyID.Valid, ShouldBeFalse)
						So(corporations[1].APIKeyID.Int64, ShouldEqual, 0)
						So(corporations[1].APIvCode.IsZero(), ShouldBeTrue)
						So(corporations[1].APIvCode.Valid, ShouldBeFalse)
						So(corporations[1].APIvCode.String, ShouldEqual, "")
						So(corporations[1].Active, ShouldBeFalse)
					})
				})
			})
		})
	})
}

func TestMySQLDatabaseConnectionLoadAllCharacters(t *testing.T) {
	Convey("Loading all characters from a MySQL database with a MySQL configuration", t, func() {
		db := createConnection()

		Convey("Connecting to the database", func() {
			err := db.Connect()

			Convey("The returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			characters, err := db.LoadAllCharacters()

			Convey("Loading all characters should return no error", func() {
				So(err, ShouldBeNil)

				Convey("The returned slice should not be nil", func() {
					So(characters, ShouldNotBeNil)
				})

				Convey("The length of the returned slice should be 6", func() {
					So(len(characters), ShouldBeGreaterThan, 0)
					So(len(characters), ShouldEqual, 6)
				})

				Convey("The returned characters should match the test data set", func() {
					Convey("Verifying entry #1", func() {
						So(characters[0].ID, ShouldEqual, 1)
						So(characters[0].UserID, ShouldEqual, 1)
						So(characters[0].CorporationID, ShouldEqual, 1)
						So(characters[0].Name, ShouldEqual, "Test Character")
						So(characters[0].EVECharacterID, ShouldEqual, 1)
						So(characters[0].Active, ShouldBeTrue)
					})

					Convey("Verifying entry #2", func() {
						So(characters[1].ID, ShouldEqual, 2)
						So(characters[1].UserID, ShouldEqual, 2)
						So(characters[1].CorporationID, ShouldEqual, 2)
						So(characters[1].Name, ShouldEqual, "Please Ignore")
						So(characters[1].EVECharacterID, ShouldEqual, 2)
						So(characters[1].Active, ShouldBeTrue)
					})

					Convey("Verifying entry #3", func() {
						So(characters[2].ID, ShouldEqual, 3)
						So(characters[2].UserID, ShouldEqual, 3)
						So(characters[2].CorporationID, ShouldEqual, 1)
						So(characters[2].Name, ShouldEqual, "Herp")
						So(characters[2].EVECharacterID, ShouldEqual, 3)
						So(characters[2].Active, ShouldBeTrue)
					})

					Convey("Verifying entry #4", func() {
						So(characters[3].ID, ShouldEqual, 4)
						So(characters[3].UserID, ShouldEqual, 3)
						So(characters[3].CorporationID, ShouldEqual, 1)
						So(characters[3].Name, ShouldEqual, "Derp")
						So(characters[3].EVECharacterID, ShouldEqual, 4)
						So(characters[3].Active, ShouldBeTrue)
					})

					Convey("Verifying entry #5", func() {
						So(characters[4].ID, ShouldEqual, 5)
						So(characters[4].UserID, ShouldEqual, 4)
						So(characters[4].CorporationID, ShouldEqual, 2)
						So(characters[4].Name, ShouldEqual, "Spai")
						So(characters[4].EVECharacterID, ShouldEqual, 5)
						So(characters[4].Active, ShouldBeFalse)
					})

					Convey("Verifying entry #6", func() {
						So(characters[5].ID, ShouldEqual, 6)
						So(characters[5].UserID, ShouldEqual, 4)
						So(characters[5].CorporationID, ShouldEqual, 2)
						So(characters[5].Name, ShouldEqual, "NoSpai")
						So(characters[5].EVECharacterID, ShouldEqual, 6)
						So(characters[5].Active, ShouldBeFalse)
					})
				})
			})
		})
	})
}

func TestMySQLDatabaseConnectionLoadAllGroupRoles(t *testing.T) {
	Convey("Loading all group roles from a MySQL database with a MySQL configuration", t, func() {
		db := createConnection()

		Convey("Connecting to the database", func() {
			err := db.Connect()

			Convey("The returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			groupRoles, err := db.LoadAllGroupRoles()

			Convey("Loading all group roles should return no error", func() {
				So(err, ShouldBeNil)

				Convey("The returned slice should not be nil", func() {
					So(groupRoles, ShouldNotBeNil)
				})

				Convey("The length of the returned slice should be 4", func() {
					So(len(groupRoles), ShouldBeGreaterThan, 0)
					So(len(groupRoles), ShouldEqual, 4)
				})

				Convey("The returned group roles should match the test data set", func() {
					Convey("Verifying entry #1", func() {
						So(groupRoles[0].ID, ShouldEqual, 1)
						So(groupRoles[0].GroupID, ShouldEqual, 1)
						// TODO verify actual role
						So(groupRoles[0].AutoAdded, ShouldBeTrue)
						So(groupRoles[0].Granted, ShouldBeTrue)
					})

					Convey("Verifying entry #2", func() {
						So(groupRoles[1].ID, ShouldEqual, 2)
						So(groupRoles[1].GroupID, ShouldEqual, 1)
						// TODO verify actual role
						So(groupRoles[1].AutoAdded, ShouldBeFalse)
						So(groupRoles[1].Granted, ShouldBeTrue)
					})

					Convey("Verifying entry #3", func() {
						So(groupRoles[2].ID, ShouldEqual, 3)
						So(groupRoles[2].GroupID, ShouldEqual, 2)
						// TODO verify actual role
						So(groupRoles[2].AutoAdded, ShouldBeFalse)
						So(groupRoles[2].Granted, ShouldBeFalse)
					})

					Convey("Verifying entry #4", func() {
						So(groupRoles[3].ID, ShouldEqual, 4)
						So(groupRoles[3].GroupID, ShouldEqual, 2)
						// TODO verify actual role
						So(groupRoles[3].AutoAdded, ShouldBeTrue)
						So(groupRoles[3].Granted, ShouldBeFalse)
					})
				})
			})
		})
	})
}

func TestMySQLDatabaseConnectionLoadAllUserRoles(t *testing.T) {
	Convey("Loading all user roles from a MySQL database with a MySQL configuration", t, func() {
		db := createConnection()

		Convey("Connecting to the database", func() {
			err := db.Connect()

			Convey("The returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			userRoles, err := db.LoadAllUserRoles()

			Convey("Loading all user roles should return no error", func() {
				So(err, ShouldBeNil)

				Convey("The returned slice should not be nil", func() {
					So(userRoles, ShouldNotBeNil)
				})

				Convey("The length of the returned slice should be 2", func() {
					So(len(userRoles), ShouldBeGreaterThan, 0)
					So(len(userRoles), ShouldEqual, 2)
				})

				Convey("The returned user roles should match the test data set", func() {
					Convey("Verifying entry #1", func() {
						So(userRoles[0].ID, ShouldEqual, 1)
						So(userRoles[0].UserID, ShouldEqual, 1)
						// TODO verify actual role
						So(userRoles[0].AutoAdded, ShouldBeFalse)
						So(userRoles[0].Granted, ShouldBeFalse)
					})

					Convey("Verifying entry #2", func() {
						So(userRoles[1].ID, ShouldEqual, 2)
						So(userRoles[1].UserID, ShouldEqual, 3)
						// TODO verify actual role
						So(userRoles[1].AutoAdded, ShouldBeTrue)
						So(userRoles[1].Granted, ShouldBeTrue)
					})
				})
			})
		})
	})
}

func TestMySQLDatabaseConnectionLoadAllGroups(t *testing.T) {
	Convey("Loading all groups from a MySQL database with a MySQL configuration", t, func() {
		db := createConnection()

		Convey("Connecting to the database", func() {
			err := db.Connect()

			Convey("The returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			groups, err := db.LoadAllGroups()

			Convey("Loading all groups should return no error", func() {
				So(err, ShouldBeNil)

				Convey("The returned slice should not be nil", func() {
					So(groups, ShouldNotBeNil)
				})

				Convey("The length of the returned slice should be 2", func() {
					So(len(groups), ShouldBeGreaterThan, 0)
					So(len(groups), ShouldEqual, 2)
				})

				Convey("The returned groups should match the test data set", func() {
					Convey("Verifying entry #1", func() {
						So(groups[0].ID, ShouldEqual, 1)
						So(groups[0].Name, ShouldEqual, "Test Group")
						So(groups[0].Active, ShouldBeTrue)
					})

					Convey("Verifying entry #2", func() {
						So(groups[1].ID, ShouldEqual, 2)
						So(groups[1].Name, ShouldEqual, "Dank Access")
						So(groups[1].Active, ShouldBeFalse)
					})
				})
			})
		})
	})
}

func TestMySQLDatabaseConnectionLoadAllUsers(t *testing.T) {
	Convey("Loading all users from a MySQL database with a MySQL configuration", t, func() {
		db := createConnection()

		Convey("Connecting to the database", func() {
			err := db.Connect()

			Convey("The returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			users, err := db.LoadAllUsers()

			Convey("Loading all users should return no error", func() {
				So(err, ShouldBeNil)

				Convey("The returned slice should not be nil", func() {
					So(users, ShouldNotBeNil)
				})

				Convey("The length of the returned slice should be 4", func() {
					So(len(users), ShouldBeGreaterThan, 0)
					So(len(users), ShouldEqual, 4)
				})

				Convey("The returned users should match the test data set", func() {
					Convey("Verifying entry #1", func() {
						So(users[0].ID, ShouldEqual, 1)
						So(users[0].Username, ShouldEqual, "test1")
						So(users[0].Password.IsZero(), ShouldBeTrue)
						So(users[0].Password.Valid, ShouldBeFalse)
						So(users[0].Password.String, ShouldEqual, "")
						So(users[0].Active, ShouldBeTrue)
					})

					Convey("Verifying entry #2", func() {
						So(users[1].ID, ShouldEqual, 2)
						So(users[1].Username, ShouldEqual, "test2")
						So(users[1].Password.IsZero(), ShouldBeTrue)
						So(users[1].Password.Valid, ShouldBeFalse)
						So(users[1].Password.String, ShouldEqual, "")
						So(users[1].Active, ShouldBeFalse)
					})

					Convey("Verifying entry #3", func() {
						So(users[2].ID, ShouldEqual, 3)
						So(users[2].Username, ShouldEqual, "test3")
						So(users[2].Password.IsZero(), ShouldBeFalse)
						So(users[2].Password.Valid, ShouldBeTrue)
						So(bcrypt.CompareHashAndPassword([]byte(users[2].Password.String), []byte("test3")), ShouldBeNil)
						So(users[2].Active, ShouldBeTrue)
					})

					Convey("Verifying entry #4", func() {
						So(users[3].ID, ShouldEqual, 4)
						So(users[3].Username, ShouldEqual, "test4")
						So(users[3].Password.IsZero(), ShouldBeFalse)
						So(users[3].Password.Valid, ShouldBeTrue)
						So(bcrypt.CompareHashAndPassword([]byte(users[3].Password.String), []byte("test4")), ShouldBeNil)
						So(users[3].Active, ShouldBeFalse)
					})
				})
			})
		})
	})
}
