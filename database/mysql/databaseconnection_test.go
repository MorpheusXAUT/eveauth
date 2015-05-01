package mysql

import (
	"fmt"
	"os"
	"testing"

	"github.com/morpheusxaut/eveauth/misc"
	"github.com/morpheusxaut/eveauth/models"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/guregu/null.v2/zero"
)

var (
	database *DatabaseConnection
)

func createMySQLConnection() (*DatabaseConnection, error) {
	if database == nil {
		databaseHost := "localhost:3306"
		if len(os.Getenv("DATABASE_HOST")) > 0 {
			databaseHost = os.Getenv("DATABASE_HOST")
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
			DatabaseType:     1,
			DatabaseHost:     databaseHost,
			DatabaseSchema:   databaseSchema,
			DatabaseUser:     databaseUser,
			DatabasePassword: databasePassword,
			DebugLevel:       1,
			HTTPHost:         "localhost:5000",
		}

		db := &DatabaseConnection{
			Config: config,
		}

		err := db.Connect()
		if err != nil {
			return nil, err
		}

		database = db
	}

	return database, nil
}

func TestDatabaseConnectionConnect(t *testing.T) {
	Convey("Connecting to a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
		})
	})
}

func TestDatabaseConnectionInvalidConnect(t *testing.T) {
	Convey("Connecting to a MySQL database with an invalid configuration", t, func() {
		config := &misc.Configuration{
			DatabaseType:     1,
			DatabaseHost:     "does.not.exist:123456",
			DatabaseSchema:   "nonexistent",
			DatabaseUser:     "nonexistent",
			DatabasePassword: "nonexistent",
			DebugLevel:       1,
			HTTPHost:         "localhost:5000",
		}

		db := &DatabaseConnection{
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

func TestDatabaseConnectionRawQuery(t *testing.T) {
	Convey("Performing a raw query at a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
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
					Convey(fmt.Sprintf("The raw data table for key %v should have 6 entries", key), func() {
						So(len(value), ShouldBeGreaterThan, 0)
						So(len(value), ShouldEqual, 6)
					})
				}
			})
		})
	})
}

func TestDatabaseConnectionRawInvalidQuery(t *testing.T) {
	Convey("Performing a raw invalid query at a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
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
}

func TestDatabaseConnectionLoadAllAccounts(t *testing.T) {
	Convey("Loading all accounts from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
		})

		accounts, err := db.LoadAllAccounts()

		Convey("Loading all accounts should return no error", func() {
			So(err, ShouldBeNil)

			Convey("The returned slice should not be nil", func() {
				So(accounts, ShouldNotBeNil)
			})

			Convey("The length of the returned slice should be 6", func() {
				So(len(accounts), ShouldBeGreaterThan, 0)
				So(len(accounts), ShouldEqual, 6)
			})

			Convey("The returned accounts should match the test data set", func() {
				for index, account := range accounts {
					Convey(fmt.Sprintf("Verifying entry #%d", index), func() {
						So(account, ShouldResemble, testAccounts[index+1])
					})
				}
			})
		})
	})
}

func TestDatabaseConnectionLoadAllCorporations(t *testing.T) {
	Convey("Loading all corporations from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
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
}

func TestDatabaseConnectionLoadAllCharacters(t *testing.T) {
	Convey("Loading all characters from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
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
}

func TestDatabaseConnectionLoadAllRoles(t *testing.T) {
	Convey("Loading all roles from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
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
}

func TestDatabaseConnectionLoadAllGroupRoles(t *testing.T) {
	Convey("Loading all group roles from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
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
}

func TestDatabaseConnectionLoadAllUserRoles(t *testing.T) {
	Convey("Loading all user roles from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
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
}

func TestDatabaseConnectionLoadAllGroups(t *testing.T) {
	Convey("Loading all groups from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
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
}

func TestDatabaseConnectionLoadAllUsers(t *testing.T) {
	Convey("Loading all users from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
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
}

func TestDatabaseConnectionLoadAllApplications(t *testing.T) {
	Convey("Loading all applications from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
		})

		applications, err := db.LoadAllApplications()

		Convey("Loading all applications should return no error", func() {
			So(err, ShouldBeNil)

			Convey("The returned slice should not be nil", func() {
				So(applications, ShouldNotBeNil)
			})

			Convey("The length of the returned slice should be 2", func() {
				So(len(applications), ShouldBeGreaterThan, 0)
				So(len(applications), ShouldEqual, 2)
			})

			Convey("The returned applications should match the test data set", func() {
				for index, application := range applications {
					Convey(fmt.Sprintf("Verifying entry #%d", index), func() {
						So(application, ShouldResemble, testApplications[index+1])
					})
				}
			})
		})
	})
}

func TestDatabaseConnectionLoadAccount(t *testing.T) {
	Convey("Loading account #1 from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
		})

		account, err := db.LoadAccount(1)

		Convey("Loading account #1 should return no error", func() {
			So(err, ShouldBeNil)

			Convey("The result should not be nil", func() {
				So(account, ShouldNotBeNil)
			})

			Convey("The returned account should match the test data set", func() {
				Convey("Verifying entry", func() {
					So(account, ShouldResemble, testAccounts[1])
				})
			})
		})
	})
}

func TestDatabaseConnectionLoadCorporation(t *testing.T) {
	Convey("Loading corporation #1 from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
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
}

func TestDatabaseConnectionLoadCorporationFromEVECorporationID(t *testing.T) {
	Convey("Loading corporation with EVE corporation ID #1 from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
		})

		corporation, err := db.LoadCorporationFromEVECorporationID(1)

		Convey("Loading corporation with EVE corporation ID #1 should return no error", func() {
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
}

func TestDatabaseConnectionLoadCorporationNameFromID(t *testing.T) {
	Convey("Loading corporation name for corporation ID #1 from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
		})

		name, err := db.LoadCorporationNameFromID(1)

		Convey("Loading corporation name for corporation ID #1 should return no error", func() {
			So(err, ShouldBeNil)

			Convey("The result should not be empty", func() {
				So(len(name), ShouldNotEqual, 0)
				So(len(name), ShouldBeGreaterThan, 0)
				So(name, ShouldNotEqual, "")
			})

			Convey("The returned name should match the test data set", func() {
				Convey("Verifying entry", func() {
					So(name, ShouldEqual, testCorporations[1].Name)
				})
			})
		})
	})
}

func TestDatabaseConnectionLoadCharacter(t *testing.T) {
	Convey("Loading character #1 from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
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
}

func TestDatabaseConnectionLoadRole(t *testing.T) {
	Convey("Loading role #1 from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
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
}

func TestDatabaseConnectionLoadGroupRole(t *testing.T) {
	Convey("Loading group role #1 from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
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
}

func TestDatabaseConnectionLoadUserRole(t *testing.T) {
	Convey("Loading user role #1 from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
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
}

func TestDatabaseConnectionLoadGroup(t *testing.T) {
	Convey("Loading group #1 from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
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
}

func TestDatabaseConnectionLoadUser(t *testing.T) {
	Convey("Loading user #1 from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
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
}

func TestDatabaseConnectionLoadUserFromUsername(t *testing.T) {
	Convey("Loading user with name test1 from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
		})

		user, err := db.LoadUserFromUsername("test1")

		Convey("Loading user with name test1 should return no error", func() {
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
}

func TestDatabaseConnectionLoadApplication(t *testing.T) {
	Convey("Loading application #1 from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
		})

		application, err := db.LoadApplication(1)

		Convey("Loading application #1 should return no error", func() {
			So(err, ShouldBeNil)

			Convey("The result should not be nil", func() {
				So(application, ShouldNotBeNil)
			})

			Convey("The returned application should match the test data set", func() {
				Convey("Verifying entry", func() {
					So(application, ShouldResemble, testApplications[1])
				})
			})
		})
	})
}

func TestDatabaseConnectionLoadAvailableGroupsForUser(t *testing.T) {
	Convey("Loading available groups for user #1 from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
		})

		groups, err := db.LoadAvailableGroupsForUser(1)

		Convey("Loading available groups for user #1 should return no error", func() {
			So(err, ShouldBeNil)

			Convey("The result should not be nil", func() {
				So(groups, ShouldNotBeNil)
			})

			Convey("The returned groups should match the test data set", func() {
				Convey("Verifying entry", func() {
					So(groups[0], ShouldResemble, testGroups[2])
				})
			})
		})
	})
}

func TestDatabaseConnectionLoadAvailableUserRolesForUser(t *testing.T) {
	Convey("Loading available user roles for user #1 from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
		})

		roles, err := db.LoadAvailableUserRolesForUser(1)

		Convey("Loading available user roles for user #1 should return no error", func() {
			So(err, ShouldBeNil)

			Convey("The result should not be nil", func() {
				So(roles, ShouldNotBeNil)
			})

			Convey("The returned roles should match the test data set", func() {
				Convey("Verifying entry", func() {
					So(len(roles), ShouldEqual, 3)
					// ShouldContain doesn't work, skipping for now
					SkipSo(roles, ShouldContain, testRoles[2])
					SkipSo(roles, ShouldContain, testRoles[3])
					SkipSo(roles, ShouldContain, testRoles[4])
				})
			})
		})
	})
}

func TestDatabaseConnectionLoadAvailableGroupRolesForGroup(t *testing.T) {
	Convey("Loading available group roles for group #1 from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
		})

		roles, err := db.LoadAvailableGroupRolesForGroup(1)

		Convey("Loading available group roles for group #1 should return no error", func() {
			So(err, ShouldBeNil)

			Convey("The result should not be nil", func() {
				So(roles, ShouldNotBeNil)
			})

			Convey("The returned roles should match the test data set", func() {
				Convey("Verifying entry", func() {
					So(len(roles), ShouldEqual, 2)
					// ShouldContain doesn't work, skipping for now
					SkipSo(roles, ShouldContain, testRoles[2])
					SkipSo(roles, ShouldContain, testRoles[4])
				})
			})
		})
	})
}

func TestDatabaseConnectionLoadAllApplicationsForUser(t *testing.T) {
	Convey("Loading all applications for user #1 from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
		})

		applications, err := db.LoadAllApplicationsForUser(1)

		Convey("Loading all applications for user #1 should return no error", func() {
			So(err, ShouldBeNil)

			Convey("The result should not be nil", func() {
				So(applications, ShouldNotBeNil)
			})

			Convey("The returned applications should match the test data set", func() {
				Convey("Verifying entry", func() {
					So(applications[0], ShouldResemble, testApplications[1])
				})
			})
		})
	})
}

func TestDatabaseConnectionLoadPasswordForUser(t *testing.T) {
	Convey("Loading password for user test1 from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
		})

		password, err := db.LoadPasswordForUser("test1")

		Convey("Loading password for user test1 should return no error", func() {
			So(err, ShouldBeNil)

			Convey("The result should have length 60", func() {
				So(len(password), ShouldNotEqual, 0)
				So(len(password), ShouldEqual, 60)
			})

			Convey("The returned password hash should match the test data set", func() {
				Convey("Verifying entry", func() {
					So(password, ShouldEqual, testUsers[1].Password)
				})
			})
		})
	})
}

func TestDatabaseConnectionQueryUserIDExists(t *testing.T) {
	Convey("Querying whether a user ID exists in a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
		})

		exists, err := db.QueryUserIDExists(1)

		Convey("Querying whether user ID #1 exists should return no error", func() {
			So(err, ShouldBeNil)

			Convey("The queried user ID should exist", func() {
				So(exists, ShouldBeTrue)
			})
		})

		exists, err = db.QueryUserIDExists(-1)

		Convey("Querying whether user ID #-1 exists should return no error", func() {
			So(err, ShouldBeNil)

			Convey("The queried user ID should not exist", func() {
				So(exists, ShouldBeFalse)
			})
		})
	})
}

func TestDatabaseConnectionQueryUserNameEmailExists(t *testing.T) {
	Convey("Querying whether a username or email exists in a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
		})

		exists, err := db.QueryUserNameEmailExists("test1", "test1@example.com")

		Convey("Querying whether username test1 or email test1@example.com exists should return no error", func() {
			So(err, ShouldBeNil)

			Convey("The queried username or email should exist", func() {
				So(exists, ShouldBeTrue)
			})
		})

		exists, err = db.QueryUserNameEmailExists("test1", "does.not@exist.com")

		Convey("Querying whether username test1 or email does.not@exist.com exists should return no error", func() {
			So(err, ShouldBeNil)

			Convey("The queried username or email should exist", func() {
				So(exists, ShouldBeTrue)
			})
		})

		exists, err = db.QueryUserNameEmailExists("does.not.exist", "test1@example.com")

		Convey("Querying whether username does.not.exist or email test1@example.com exists should return no error", func() {
			So(err, ShouldBeNil)

			Convey("The queried username or email should exist", func() {
				So(exists, ShouldBeTrue)
			})
		})

		exists, err = db.QueryUserNameEmailExists("does.not.exist", "does.not@exist.com")

		Convey("Querying whether username does.not.exist or email does.not@exist.com exists should return no error", func() {
			So(err, ShouldBeNil)

			Convey("The queried username or email should not exist", func() {
				So(exists, ShouldBeFalse)
			})
		})
	})
}

func TestDatabaseConnectionLoadInvalidAccount(t *testing.T) {
	Convey("Loading invalid account from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
		})

		account, err := db.LoadAccount(-1)

		Convey("Loading an invalid account should return an error", func() {
			So(err, ShouldNotBeNil)

			Convey("The result should be nil", func() {
				So(account, ShouldBeNil)
			})
		})
	})
}

func TestDatabaseConnectionLoadInvalidCorporation(t *testing.T) {
	Convey("Loading invalid corporation from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
		})

		corporation, err := db.LoadCorporation(-1)

		Convey("Loading an invalid corporation should return an error", func() {
			So(err, ShouldNotBeNil)

			Convey("The result should be nil", func() {
				So(corporation, ShouldBeNil)
			})
		})
	})
}

func TestDatabaseConnectionLoadInvalidCharacter(t *testing.T) {
	Convey("Loading invalid character from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
		})

		character, err := db.LoadCharacter(-1)

		Convey("Loading an invalid character should return an error", func() {
			So(err, ShouldNotBeNil)

			Convey("The result should be nil", func() {
				So(character, ShouldBeNil)
			})
		})
	})
}

func TestDatabaseConnectionLoadInvalidRole(t *testing.T) {
	Convey("Loading invalid role from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
		})

		role, err := db.LoadRole(-1)

		Convey("Loading an role should return an error", func() {
			So(err, ShouldNotBeNil)

			Convey("The result should be nil", func() {
				So(role, ShouldBeNil)
			})
		})
	})
}

func TestDatabaseConnectionLoadInvalidGroupRole(t *testing.T) {
	Convey("Loading invalid group role from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
		})

		groupRole, err := db.LoadGroupRole(-1)

		Convey("Loading an invalid group role should return an error", func() {
			So(err, ShouldNotBeNil)

			Convey("The result should be nil", func() {
				So(groupRole, ShouldBeNil)
			})
		})
	})
}

func TestDatabaseConnectionLoadInvalidUserRole(t *testing.T) {
	Convey("Loading invalid user role from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
		})

		userRole, err := db.LoadUserRole(-1)

		Convey("Loading an invalid user role should return an error", func() {
			So(err, ShouldNotBeNil)

			Convey("The result should be nil", func() {
				So(userRole, ShouldBeNil)
			})
		})
	})
}

func TestDatabaseConnectionLoadInvalidGroup(t *testing.T) {
	Convey("Loading invalid group from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
		})

		group, err := db.LoadGroup(-1)

		Convey("Loading an invalid group should return an error", func() {
			So(err, ShouldNotBeNil)

			Convey("The result should be nil", func() {
				So(group, ShouldBeNil)
			})
		})
	})
}

func TestDatabaseConnectionLoadInvalidUser(t *testing.T) {
	Convey("Loading invalid user from a MySQL database", t, func() {
		db, err := createMySQLConnection()

		Convey("The returned error should be nil", func() {
			So(err, ShouldBeNil)
		})

		Convey("The DatabaseConnection should not be nil", func() {
			So(db, ShouldNotBeNil)
		})

		user, err := db.LoadUser(-1)

		Convey("Loading an invalid user should return an error", func() {
			So(err, ShouldNotBeNil)

			Convey("The result should be nil", func() {
				So(user, ShouldBeNil)
			})
		})
	})
}

var (
	testAccounts = map[int]*models.Account{
		1: &models.Account{
			ID:            1,
			UserID:        1,
			APIKeyID:      1,
			APIvCode:      "a",
			APIAccessMask: 0,
			Active:        true,
			Characters: []*models.Character{
				testCharacters[1],
			},
		},
		2: &models.Account{
			ID:            2,
			UserID:        2,
			APIKeyID:      2,
			APIvCode:      "b",
			APIAccessMask: 0,
			Active:        false,
			Characters: []*models.Character{
				testCharacters[2],
			},
		},
		3: &models.Account{
			ID:            3,
			UserID:        3,
			APIKeyID:      3,
			APIvCode:      "c",
			APIAccessMask: 0,
			Active:        true,
			Characters: []*models.Character{
				testCharacters[3],
				testCharacters[4],
			},
		},
		4: &models.Account{
			ID:            4,
			UserID:        3,
			APIKeyID:      4,
			APIvCode:      "d",
			APIAccessMask: 268435455,
			Active:        true,
			Characters: []*models.Character{
				testCharacters[5],
				testCharacters[6],
			},
		},
		5: &models.Account{
			ID:            5,
			UserID:        4,
			APIKeyID:      5,
			APIvCode:      "e",
			APIAccessMask: 268435455,
			Active:        false,
		},
		6: &models.Account{
			ID:            6,
			UserID:        4,
			APIKeyID:      6,
			APIvCode:      "f",
			APIAccessMask: 268435455,
			Active:        false,
		},
	}

	testApplications = map[int]*models.Application{
		1: &models.Application{
			ID:           1,
			Name:         "Testapp",
			MaintainerID: 1,
			Secret:       "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			Callback:     "http://localhost/callback",
			Active:       true,
		},
		2: &models.Application{
			ID:           2,
			Name:         "Apptest",
			MaintainerID: 2,
			Secret:       "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
			Callback:     "http://example.com/callback",
			Active:       false,
		},
	}

	testCorporations = map[int]*models.Corporation{
		1: &models.Corporation{
			ID:               1,
			Name:             "Test Corp Please Ignore",
			Ticker:           "TEST",
			EVECorporationID: 1,
			CEOID:            1,
			APIKeyID:         zero.IntFrom(1),
			APIvCode:         zero.StringFrom("a"),
			Active:           true,
		},
		2: &models.Corporation{
			ID:               2,
			Name:             "Corp Test Ignore Please",
			Ticker:           "CORP",
			EVECorporationID: 2,
			CEOID:            2,
			APIKeyID:         zero.NewInt(0, false),
			APIvCode:         zero.NewString("", false),
			Active:           false,
		},
	}

	testCharacters = map[int]*models.Character{
		1: &models.Character{
			ID:               1,
			AccountID:        1,
			CorporationID:    1,
			Name:             "Test Character",
			EVECharacterID:   1,
			DefaultCharacter: true,
			Active:           true,
		},
		2: &models.Character{
			ID:               2,
			AccountID:        2,
			CorporationID:    2,
			Name:             "Please Ignore",
			EVECharacterID:   2,
			DefaultCharacter: true,
			Active:           true,
		},
		3: &models.Character{
			ID:               3,
			AccountID:        3,
			CorporationID:    1,
			Name:             "Herp",
			EVECharacterID:   3,
			DefaultCharacter: true,
			Active:           true,
		},
		4: &models.Character{
			ID:               4,
			AccountID:        3,
			CorporationID:    1,
			Name:             "Derp",
			EVECharacterID:   4,
			DefaultCharacter: false,
			Active:           true,
		},
		5: &models.Character{
			ID:               5,
			AccountID:        4,
			CorporationID:    2,
			Name:             "Spai",
			EVECharacterID:   5,
			DefaultCharacter: false,
			Active:           false,
		},
		6: &models.Character{
			ID:               6,
			AccountID:        4,
			CorporationID:    2,
			Name:             "NoSpai",
			EVECharacterID:   6,
			DefaultCharacter: true,
			Active:           false,
		},
	}

	testRoles = map[int]*models.Role{
		1: &models.Role{
			ID:     1,
			Name:   "ping.all",
			Active: true,
			Locked: false,
		},
		2: &models.Role{
			ID:     2,
			Name:   "destroy.world",
			Active: false,
			Locked: true,
		},
		3: &models.Role{
			ID:     3,
			Name:   "logistics.read",
			Active: true,
			Locked: false,
		},
		4: &models.Role{
			ID:     4,
			Name:   "logistics.write",
			Active: true,
			Locked: false,
		},
	}

	testGroupRoles = map[int]*models.GroupRole{
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

	testUserRoles = map[int]*models.UserRole{
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

	testGroups = map[int]*models.Group{
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

	testUsers = map[int]*models.User{
		1: &models.User{
			ID:            1,
			Username:      "test1",
			Password:      "$2a$10$veif8VUZt7lShFhJKD0wGeY1YjCwIuWjYL0vQzlTqu8wNaYQMqzbe",
			Email:         "test1@example.com",
			VerifiedEmail: true,
			Active:        true,
			Accounts: []*models.Account{
				testAccounts[1],
			},
			UserRoles: []*models.UserRole{
				testUserRoles[1],
			},
			Groups: []*models.Group{
				testGroups[1],
			},
		},
		2: &models.User{
			ID:            2,
			Username:      "test2",
			Password:      "$2a$10$95z.WXfIreLKJ9px.3KgpOq4aXTG3DF7/5ehGYzUWALhpN6MMq/aK",
			Email:         "test2@example.com",
			VerifiedEmail: false,
			Active:        false,
			Accounts: []*models.Account{
				testAccounts[2],
			},
			UserRoles: []*models.UserRole{},
			Groups:    []*models.Group{},
		},
		3: &models.User{
			ID:            3,
			Username:      "test3",
			Password:      "$2a$10$7Yxm2scdTVpEJpvZAT7tbOFA.G9JfyxtiHbr989iocX6U37C3/j4q",
			Email:         "test3@example.com",
			VerifiedEmail: false,
			Active:        true,
			Accounts: []*models.Account{
				testAccounts[3],
				testAccounts[4],
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
			ID:            4,
			Username:      "test4",
			Password:      "$2a$10$WOWTgqaqLKbkb1uhYbtLnOuuYX4kXBC61GVAke7RkjiODoBpgGGzy",
			Email:         "test4@example.com",
			VerifiedEmail: true,
			Active:        false,
			Accounts: []*models.Account{
				testAccounts[5],
				testAccounts[6],
			},
			UserRoles: []*models.UserRole{},
			Groups:    []*models.Group{},
		},
	}
)
