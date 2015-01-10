// eveauth provides an authentication backend used for EVE Online services.
// Beside providing basic user authentication, the application can provide access using different permissions based on access rights.
package main

import (
	"github.com/morpheusxaut/eveauth/database"
	"github.com/morpheusxaut/eveauth/misc"
	"github.com/morpheusxaut/eveauth/web"
)

func main() {
	misc.ParseConfigFlags()

	misc.SetupLogger()

	database.SetupDatabase()

	web.SetupRouter()

	web.Router.HandleRequests()
}
