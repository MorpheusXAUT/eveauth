package mysql

import (
	"fmt"
	"github.com/morpheusxaut/eveauth/misc"
	"github.com/morpheusxaut/eveauth/models"
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/guregu/null.v2/zero"
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
					for index, corporation := range corporations {
						Convey(fmt.Sprintf("Verifying entry #%d", index), func() {
							So(corporation, ShouldResemble, testCorporations[index+1])
						})
					}
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
					for index, character := range characters {
						Convey(fmt.Sprintf("Verifying entry #%d", index), func() {
							So(character, ShouldResemble, testCharacters[index+1])
						})
					}
				})
			})
		})
	})
}

func TestMySQLDatabaseConnectionLoadAllRoles(t *testing.T) {
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
					for index, role := range roles {
						Convey(fmt.Sprintf("Verifying entry #%d", index), func() {
							So(role, ShouldResemble, testRoles[index+1])
						})
					}
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
					for index, groupRole := range groupRoles {
						Convey(fmt.Sprintf("Verifying entry #%d", index), func() {
							So(groupRole, ShouldResemble, testGroupRoles[index+1])
						})
					}
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
					for index, userRole := range userRoles {
						Convey(fmt.Sprintf("Verifying entry #%d", index), func() {
							So(userRole, ShouldResemble, testUserRoles[index+1])
						})
					}
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
					for index, group := range groups {
						Convey(fmt.Sprintf("Verifying entry #%d", index), func() {
							So(group, ShouldResemble, testGroups[index+1])
						})
					}
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
					for index, user := range users {
						Convey(fmt.Sprintf("Verifying entry #%d", index), func() {
							So(user, ShouldResemble, testUsers[index+1])
						})
					}
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

func TestMySQLDatabaseConnectionLoadCorporation(t *testing.T) {
	Convey("Loading corporation #1 from a MySQL database", t, func() {
		db := createConnection()

		Convey("Connecting to the database", func() {
			err := db.Connect()

			Convey("The returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			corporation, err := db.LoadCorporation(1)

			Convey("Loading corporation #1 should return no error", func() {
				So(err, ShouldBeNil)

				Convey("The result should not be nil", func() {
					So(corporation, ShouldNotBeNil)
				})

				Convey("The returned corporation should match the test data set", func() {
					Convey("Verifying entry", func() {
						So(corporation, ShouldResemble, testCorporations[1])
					})
				})
			})
		})
	})
}

func TestMySQLDatabaseConnectionLoadCharacter(t *testing.T) {
	Convey("Loading character #1 from a MySQL database", t, func() {
		db := createConnection()

		Convey("Connecting to the database", func() {
			err := db.Connect()

			Convey("The returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			character, err := db.LoadCharacter(1)

			Convey("Loading character #1 should return no error", func() {
				So(err, ShouldBeNil)

				Convey("The result should not be nil", func() {
					So(character, ShouldNotBeNil)
				})

				Convey("The returned character should match the test data set", func() {
					Convey("Verifying entry", func() {
						So(character, ShouldResemble, testCharacters[1])
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
						So(role, ShouldResemble, testRoles[1])
					})
				})
			})
		})
	})
}

func TestMySQLDatabaseConnectionLoadGroupRole(t *testing.T) {
	Convey("Loading group role #1 from a MySQL database", t, func() {
		db := createConnection()

		Convey("Connecting to the database", func() {
			err := db.Connect()

			Convey("The returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			groupRole, err := db.LoadGroupRole(1)

			Convey("Loading group role #1 should return no error", func() {
				So(err, ShouldBeNil)

				Convey("The result should not be nil", func() {
					So(groupRole, ShouldNotBeNil)
				})

				Convey("The returned group role should match the test data set", func() {
					Convey("Verifying entry", func() {
						So(groupRole, ShouldResemble, testGroupRoles[1])
					})
				})
			})
		})
	})
}

func TestMySQLDatabaseConnectionLoadUserRole(t *testing.T) {
	Convey("Loading user role #1 from a MySQL database", t, func() {
		db := createConnection()

		Convey("Connecting to the database", func() {
			err := db.Connect()

			Convey("The returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			userRole, err := db.LoadUserRole(1)

			Convey("Loading user role #1 should return no error", func() {
				So(err, ShouldBeNil)

				Convey("The result should not be nil", func() {
					So(userRole, ShouldNotBeNil)
				})

				Convey("The returned user role should match the test data set", func() {
					Convey("Verifying entry", func() {
						So(userRole, ShouldResemble, testUserRoles[1])
					})
				})
			})
		})
	})
}

func TestMySQLDatabaseConnectionLoadGroup(t *testing.T) {
	Convey("Loading group #1 from a MySQL database", t, func() {
		db := createConnection()

		Convey("Connecting to the database", func() {
			err := db.Connect()

			Convey("The returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			group, err := db.LoadGroup(1)

			Convey("Loading group #1 should return no error", func() {
				So(err, ShouldBeNil)

				Convey("The result should not be nil", func() {
					So(group, ShouldNotBeNil)
				})

				Convey("The returned group should match the test data set", func() {
					Convey("Verifying entry", func() {
						So(group, ShouldResemble, testGroups[1])
					})
				})
			})
		})
	})
}

func TestMySQLDatabaseConnectionLoadUser(t *testing.T) {
	Convey("Loading user #1 from a MySQL database", t, func() {
		db := createConnection()

		Convey("Connecting to the database", func() {
			err := db.Connect()

			Convey("The returned error should be nil", func() {
				So(err, ShouldBeNil)
			})

			user, err := db.LoadUser(1)

			Convey("Loading user #1 should return no error", func() {
				So(err, ShouldBeNil)

				Convey("The result should not be nil", func() {
					So(user, ShouldNotBeNil)
				})

				Convey("The returned user should match the test data set", func() {
					Convey("Verifying entry", func() {
						So(user, ShouldResemble, testUsers[1])
					})
				})
			})
		})
	})
}

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

var testCorporations map[int]*models.Corporation = map[int]*models.Corporation{
	1: &models.Corporation{
		ID:               1,
		Name:             "Test Corp Please Ignore",
		Ticker:           "TEST",
		EVECorporationID: 1,
		APIKeyID:         zero.IntFrom(1),
		APIvCode:         zero.StringFrom("a"),
		Active:           true,
	},
	2: &models.Corporation{
		ID:               2,
		Name:             "Corp Test Ignore Please",
		Ticker:           "CORP",
		EVECorporationID: 2,
		APIKeyID:         zero.NewInt(0, false),
		APIvCode:         zero.NewString("", false),
		Active:           false,
	},
}

var testCharacters map[int]*models.Character = map[int]*models.Character{
	1: &models.Character{
		ID:             1,
		UserID:         1,
		CorporationID:  1,
		Name:           "Test Character",
		EVECharacterID: 1,
		Active:         true,
	},
	2: &models.Character{
		ID:             2,
		UserID:         2,
		CorporationID:  2,
		Name:           "Please Ignore",
		EVECharacterID: 2,
		Active:         true,
	},
	3: &models.Character{
		ID:             3,
		UserID:         3,
		CorporationID:  1,
		Name:           "Herp",
		EVECharacterID: 3,
		Active:         true,
	},
	4: &models.Character{
		ID:             4,
		UserID:         3,
		CorporationID:  1,
		Name:           "Derp",
		EVECharacterID: 4,
		Active:         true,
	},
	5: &models.Character{
		ID:             5,
		UserID:         4,
		CorporationID:  2,
		Name:           "Spai",
		EVECharacterID: 5,
		Active:         false,
	},
	6: &models.Character{
		ID:             6,
		UserID:         4,
		CorporationID:  2,
		Name:           "NoSpai",
		EVECharacterID: 6,
		Active:         false,
	},
}

var testRoles map[int]*models.Role = map[int]*models.Role{
	1: &models.Role{
		ID:     1,
		Name:   "ping.all",
		Active: true,
	},
	2: &models.Role{
		ID:     2,
		Name:   "destroy.world",
		Active: false,
	},
	3: &models.Role{
		ID:     3,
		Name:   "logistics.read",
		Active: true,
	},
	4: &models.Role{
		ID:     4,
		Name:   "logistics.write",
		Active: true,
	},
}

var testGroupRoles map[int]*models.GroupRole = map[int]*models.GroupRole{
	1: &models.GroupRole{
		ID:        1,
		GroupID:   1,
		Role:      testRoles[1],
		AutoAdded: true,
		Granted:   true,
	},
	2: &models.GroupRole{
		ID:        2,
		GroupID:   1,
		Role:      testRoles[3],
		AutoAdded: false,
		Granted:   true,
	},
	3: &models.GroupRole{
		ID:        3,
		GroupID:   2,
		Role:      testRoles[2],
		AutoAdded: false,
		Granted:   false,
	},
	4: &models.GroupRole{
		ID:        4,
		GroupID:   2,
		Role:      testRoles[4],
		AutoAdded: true,
		Granted:   false,
	},
}

var testUserRoles map[int]*models.UserRole = map[int]*models.UserRole{
	1: &models.UserRole{
		ID:        1,
		UserID:    1,
		Role:      testRoles[1],
		AutoAdded: false,
		Granted:   false,
	},
	2: &models.UserRole{
		ID:        2,
		UserID:    3,
		Role:      testRoles[2],
		AutoAdded: true,
		Granted:   true,
	},
}

var testGroups map[int]*models.Group = map[int]*models.Group{
	1: &models.Group{
		ID:     1,
		Name:   "Test Group",
		Active: true,
		GroupRoles: []*models.GroupRole{
			testGroupRoles[1],
			testGroupRoles[2],
		},
	},
	2: &models.Group{
		ID:     2,
		Name:   "Dank Access",
		Active: false,
		GroupRoles: []*models.GroupRole{
			testGroupRoles[3],
			testGroupRoles[4],
		},
	},
}

var testUsers map[int]*models.User = map[int]*models.User{
	1: &models.User{
		ID:       1,
		Username: "test1",
		Password: zero.NewString("", false),
		Active:   true,
		Characters: []*models.Character{
			testCharacters[1],
		},
		APIKeys: []*models.APIKey{
			testAPIKeys[1],
		},
		UserRoles: []*models.UserRole{
			testUserRoles[1],
		},
		Groups: []*models.Group{
			testGroups[1],
		},
	},
	2: &models.User{
		ID:       2,
		Username: "test2",
		Password: zero.NewString("", false),
		Active:   false,
		Characters: []*models.Character{
			testCharacters[2],
		},
		APIKeys: []*models.APIKey{
			testAPIKeys[2],
		},
		UserRoles: []*models.UserRole{},
		Groups:    []*models.Group{},
	},
	3: &models.User{
		ID:       3,
		Username: "test3",
		Password: zero.StringFrom("$2a$10$7Yxm2scdTVpEJpvZAT7tbOFA.G9JfyxtiHbr989iocX6U37C3/j4q"),
		Active:   true,
		Characters: []*models.Character{
			testCharacters[3],
			testCharacters[4],
		},
		APIKeys: []*models.APIKey{
			testAPIKeys[3],
			testAPIKeys[4],
		},
		UserRoles: []*models.UserRole{
			testUserRoles[2],
		},
		Groups: []*models.Group{
			testGroups[1],
			testGroups[2],
		},
	},
	4: &models.User{
		ID:       4,
		Username: "test4",
		Password: zero.StringFrom("$2a$10$WOWTgqaqLKbkb1uhYbtLnOuuYX4kXBC61GVAke7RkjiODoBpgGGzy"),
		Active:   false,
		Characters: []*models.Character{
			testCharacters[5],
			testCharacters[6],
		},
		APIKeys: []*models.APIKey{
			testAPIKeys[5],
			testAPIKeys[6],
		},
		UserRoles: []*models.UserRole{},
		Groups:    []*models.Group{},
	},
}
