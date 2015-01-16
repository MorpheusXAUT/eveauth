package database

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestDatabaseTypeString(t *testing.T) {
	Convey("Printing the string representation of the different DatabaseTypes", t, func() {
		Convey("The DatabaseTypeMySQL should print \"MySQL\"", func() {
			So(fmt.Sprintf("%v", DatabaseTypeMySQL), ShouldEqual, "MySQL")
			So(fmt.Sprintf("%v", DatabaseType(1)), ShouldEqual, "MySQL")
		})

		Convey("Any other DatabaseType should print \"Unknown\"", func() {
			So(fmt.Sprintf("%v", DatabaseType(0)), ShouldEqual, "Unknown")
			So(fmt.Sprintf("%v", DatabaseType(2)), ShouldEqual, "Unknown")
			So(fmt.Sprintf("%v", DatabaseType(3)), ShouldEqual, "Unknown")
			So(fmt.Sprintf("%v", DatabaseType(1337)), ShouldEqual, "Unknown")
		})
	})
}
