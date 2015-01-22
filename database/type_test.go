package database

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDatabaseTypeString(t *testing.T) {
	Convey("Printing the string representation of the different DatabaseTypes", t, func() {
		Convey("The DatabaseTypeMySQL should print \"MySQL\"", func() {
			So(fmt.Sprintf("%v", TypeMySQL), ShouldEqual, "MySQL")
			So(fmt.Sprintf("%v", Type(1)), ShouldEqual, "MySQL")
		})

		Convey("Any other Type should print \"Unknown\"", func() {
			So(fmt.Sprintf("%v", Type(0)), ShouldEqual, "Unknown")
			So(fmt.Sprintf("%v", Type(2)), ShouldEqual, "Unknown")
			So(fmt.Sprintf("%v", Type(3)), ShouldEqual, "Unknown")
			So(fmt.Sprintf("%v", Type(1337)), ShouldEqual, "Unknown")
		})
	})
}
