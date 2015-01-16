package mysql

import (
	"fmt"
	"github.com/morpheusxaut/eveauth/misc"
	"github.com/morpheusxaut/eveauth/models"
	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
	"testing"
)

var testAPIKeys map[int]*models.APIKey = map[int]*models.APIKey{
	1: &models.APIKey{
		ID:       1,
		UserID:   1,
		APIKeyID: 1,
		APIvCode: "a",
		Active:   true,
	},
	2: &models.APIKey{
		ID:       2,
		UserID:   2,
		APIKeyID: 2,
		APIvCode: "b",
		Active:   false,
	},
	3: &models.APIKey{
		ID:       3,
		UserID:   3,
		APIKeyID: 3,
		APIvCode: "c",
		Active:   true,
	},
	4: &models.APIKey{
		ID:       4,
		UserID:   3,
		APIKeyID: 4,
		APIvCode: "d",
		Active:   true,
	},
	5: &models.APIKey{
		ID:       5,
		UserID:   4,
		APIKeyID: 5,
		APIvCode: "e",
		Active:   false,
	},
	6: &models.APIKey{
		ID:       6,
		UserID:   4,
		APIKeyID: 6,
		APIvCode: "f",
		Active:   false,
	},
}

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
	Convey("Connecting to a MySQL database", t, func() {
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
	Convey("Connecting to a MySQL database with an invalid configuration", t, func() {
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
	Convey("Performing a raw query at a MySQL database", t, func() {
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
	Convey("Performing a raw invalid query at a MySQL database", t, func() {
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
	Convey("Loading all API keys from a MySQL database", t, func() {
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
					for index, apiKey := range apiKeys {
						Convey(fmt.Sprintf("Verifying entry #%d", index), func() {
							So(apiKey, ShouldResemble, testAPIKeys[index+1])
						})
					}
				})
			})
		})
	})
}

func TestMySQLDatabaseConnectionLoadAllCorporations(t *testing.T) {
	Convey("Loading all corporations from a MySQL database", t, func() {
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
	Convey("Loading all characters from a MySQL database", t, func() {
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

func TestMySQLDatabaseConnectionLoadRoles(t *testing.T) {
	Convey("Loading all roles from a MySQL database", t, func() {
		db := createConnection()

		Convey("Connecting to the database", func() {
			err := db.Connect()

			Convey("The returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			roles, err := db.LoadAllRoles()

			Convey("Loading all roles should return no error", func() {
				So(err, ShouldBeNil)

				Convey("The returned slice should not be nil", func() {
					So(roles, ShouldNotBeNil)
				})

				Convey("The length of the returned slice should be 4", func() {
					So(len(roles), ShouldBeGreaterThan, 0)
					So(len(roles), ShouldEqual, 4)
				})

				Convey("The returned roles should match the test data set", func() {
					Convey("Verifying entry #1", func() {
						So(roles[0].ID, ShouldEqual, 1)
						So(roles[0].Name, ShouldEqual, "ping.all")
						So(roles[0].Active, ShouldBeTrue)
					})

					Convey("Verifying entry #2", func() {
						So(roles[1].ID, ShouldEqual, 2)
						So(roles[1].Name, ShouldEqual, "destroy.world")
						So(roles[1].Active, ShouldBeFalse)
					})

					Convey("Verifying entry #3", func() {
						So(roles[2].ID, ShouldEqual, 3)
						So(roles[2].Name, ShouldEqual, "logistics.read")
						So(roles[2].Active, ShouldBeTrue)
					})

					Convey("Verifying entry #4", func() {
						So(roles[3].ID, ShouldEqual, 4)
						So(roles[3].Name, ShouldEqual, "logistics.write")
						So(roles[3].Active, ShouldBeTrue)
					})
				})
			})
		})
	})
}

func TestMySQLDatabaseConnectionLoadAllGroupRoles(t *testing.T) {
	Convey("Loading all group roles from a MySQL database", t, func() {
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
						So(groupRoles[0].Role, ShouldNotBeNil)
						So(groupRoles[0].Role.ID, ShouldEqual, 1)
						So(groupRoles[0].Role.Name, ShouldEqual, "ping.all")
						So(groupRoles[0].Role.Active, ShouldBeTrue)
						So(groupRoles[0].AutoAdded, ShouldBeTrue)
						So(groupRoles[0].Granted, ShouldBeTrue)
					})

					Convey("Verifying entry #2", func() {
						So(groupRoles[1].ID, ShouldEqual, 2)
						So(groupRoles[1].GroupID, ShouldEqual, 1)
						So(groupRoles[1].Role, ShouldNotBeNil)
						So(groupRoles[1].Role.ID, ShouldEqual, 3)
						So(groupRoles[1].Role.Name, ShouldEqual, "logistics.read")
						So(groupRoles[1].Role.Active, ShouldBeTrue)
						So(groupRoles[1].AutoAdded, ShouldBeFalse)
						So(groupRoles[1].Granted, ShouldBeTrue)
					})

					Convey("Verifying entry #3", func() {
						So(groupRoles[2].ID, ShouldEqual, 3)
						So(groupRoles[2].GroupID, ShouldEqual, 2)
						So(groupRoles[2].Role, ShouldNotBeNil)
						So(groupRoles[2].Role.ID, ShouldEqual, 2)
						So(groupRoles[2].Role.Name, ShouldEqual, "destroy.world")
						So(groupRoles[2].Role.Active, ShouldBeFalse)
						So(groupRoles[2].AutoAdded, ShouldBeFalse)
						So(groupRoles[2].Granted, ShouldBeFalse)
					})

					Convey("Verifying entry #4", func() {
						So(groupRoles[3].ID, ShouldEqual, 4)
						So(groupRoles[3].GroupID, ShouldEqual, 2)
						So(groupRoles[3].Role, ShouldNotBeNil)
						So(groupRoles[3].Role.ID, ShouldEqual, 4)
						So(groupRoles[3].Role.Name, ShouldEqual, "logistics.write")
						So(groupRoles[3].Role.Active, ShouldBeTrue)
						So(groupRoles[3].AutoAdded, ShouldBeTrue)
						So(groupRoles[3].Granted, ShouldBeFalse)
					})
				})
			})
		})
	})
}

func TestMySQLDatabaseConnectionLoadAllUserRoles(t *testing.T) {
	Convey("Loading all user roles from a MySQL database", t, func() {
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
						So(userRoles[0].Role, ShouldNotBeNil)
						So(userRoles[0].Role.ID, ShouldEqual, 1)
						So(userRoles[0].Role.Name, ShouldEqual, "ping.all")
						So(userRoles[0].Role.Active, ShouldBeTrue)
						So(userRoles[0].AutoAdded, ShouldBeFalse)
						So(userRoles[0].Granted, ShouldBeFalse)
					})

					Convey("Verifying entry #2", func() {
						So(userRoles[1].ID, ShouldEqual, 2)
						So(userRoles[1].UserID, ShouldEqual, 3)
						So(userRoles[1].Role, ShouldNotBeNil)
						So(userRoles[1].Role.ID, ShouldEqual, 2)
						So(userRoles[1].Role.Name, ShouldEqual, "destroy.world")
						So(userRoles[1].Role.Active, ShouldBeFalse)
						So(userRoles[1].AutoAdded, ShouldBeTrue)
						So(userRoles[1].Granted, ShouldBeTrue)
					})
				})
			})
		})
	})
}

func TestMySQLDatabaseConnectionLoadAllGroups(t *testing.T) {
	Convey("Loading all groups from a MySQL database", t, func() {
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
						So(groups[0].GroupRoles, ShouldNotBeNil)
						So(len(groups[0].GroupRoles), ShouldEqual, 2)
						So(groups[0].GroupRoles[0].ID, ShouldEqual, 1)
						So(groups[0].GroupRoles[0].GroupID, ShouldEqual, 1)
						So(groups[0].GroupRoles[0].Role, ShouldNotBeNil)
						So(groups[0].GroupRoles[0].Role.ID, ShouldEqual, 1)
						So(groups[0].GroupRoles[0].Role.Name, ShouldEqual, "ping.all")
						So(groups[0].GroupRoles[0].Role.Active, ShouldBeTrue)
						So(groups[0].GroupRoles[0].AutoAdded, ShouldBeTrue)
						So(groups[0].GroupRoles[0].Granted, ShouldBeTrue)
						So(groups[0].GroupRoles[1].ID, ShouldEqual, 2)
						So(groups[0].GroupRoles[1].GroupID, ShouldEqual, 1)
						So(groups[0].GroupRoles[1].Role, ShouldNotBeNil)
						So(groups[0].GroupRoles[1].Role.ID, ShouldEqual, 3)
						So(groups[0].GroupRoles[1].Role.Name, ShouldEqual, "logistics.read")
						So(groups[0].GroupRoles[1].Role.Active, ShouldBeTrue)
						So(groups[0].GroupRoles[1].AutoAdded, ShouldBeFalse)
						So(groups[0].GroupRoles[1].Granted, ShouldBeTrue)
					})

					Convey("Verifying entry #2", func() {
						So(groups[1].ID, ShouldEqual, 2)
						So(groups[1].Name, ShouldEqual, "Dank Access")
						So(groups[1].Active, ShouldBeFalse)
						So(groups[1].GroupRoles, ShouldNotBeNil)
						So(len(groups[1].GroupRoles), ShouldEqual, 2)
						So(groups[1].GroupRoles[0].ID, ShouldEqual, 3)
						So(groups[1].GroupRoles[0].GroupID, ShouldEqual, 2)
						So(groups[1].GroupRoles[0].Role, ShouldNotBeNil)
						So(groups[1].GroupRoles[0].Role.ID, ShouldEqual, 2)
						So(groups[1].GroupRoles[0].Role.Name, ShouldEqual, "destroy.world")
						So(groups[1].GroupRoles[0].Role.Active, ShouldBeFalse)
						So(groups[1].GroupRoles[0].AutoAdded, ShouldBeFalse)
						So(groups[1].GroupRoles[0].Granted, ShouldBeFalse)
						So(groups[1].GroupRoles[1].ID, ShouldEqual, 4)
						So(groups[1].GroupRoles[1].GroupID, ShouldEqual, 2)
						So(groups[1].GroupRoles[1].Role, ShouldNotBeNil)
						So(groups[1].GroupRoles[1].Role.ID, ShouldEqual, 4)
						So(groups[1].GroupRoles[1].Role.Name, ShouldEqual, "logistics.write")
						So(groups[1].GroupRoles[1].Role.Active, ShouldBeTrue)
						So(groups[1].GroupRoles[1].AutoAdded, ShouldBeTrue)
						So(groups[1].GroupRoles[1].Granted, ShouldBeFalse)
					})
				})
			})
		})
	})
}

func TestMySQLDatabaseConnectionLoadAllUsers(t *testing.T) {
	Convey("Loading all users from a MySQL database", t, func() {
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

func TestMySQLDatabaseConnectionLoadAPIKey(t *testing.T) {
	Convey("Loading API key #1 from a MySQL database", t, func() {
		db := createConnection()

		Convey("Connecting to the database", func() {
			err := db.Connect()

			Convey("The returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			apiKey, err := db.LoadAPIKey(1)

			Convey("Loading API key #1 should return no error", func() {
				So(err, ShouldBeNil)

				Convey("The result should not be nil", func() {
					So(apiKey, ShouldNotBeNil)
				})

				Convey("The returned role should match the test data set", func() {
					Convey("Verifying entry", func() {
						So(apiKey, ShouldResemble, testAPIKeys[1])
					})
				})
			})
		})
	})
}

func TestMySQLDatabaseConnectionLoadRole(t *testing.T) {
	Convey("Loading role #1 from a MySQL database", t, func() {
		db := createConnection()

		Convey("Connecting to the database", func() {
			err := db.Connect()

			Convey("The returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			role, err := db.LoadRole(1)

			Convey("Loading role #1 should return no error", func() {
				So(err, ShouldBeNil)

				Convey("The result should not be nil", func() {
					So(role, ShouldNotBeNil)
				})

				Convey("The returned role should match the test data set", func() {
					Convey("Verifying entry", func() {
						So(role.ID, ShouldEqual, 1)
						So(role.Name, ShouldEqual, "ping.all")
						So(role.Active, ShouldBeTrue)
					})
				})
			})
		})
	})
}
