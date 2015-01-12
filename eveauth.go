// eveauth provides an authentication backend used for EVE Online services.
// Beside providing basic user authentication, the application can provide access using different permissions based on access rights.
package main

import (
	"github.com/morpheusxaut/eveauth/database"
	"github.com/morpheusxaut/eveauth/misc"
	"github.com/morpheusxaut/eveauth/web"
	"os"
)

func main() {
	config, err := misc.LoadConfig()
	if err != nil {
		misc.Logger.Criticalf("Failed to load config: [%v]", err)
		os.Exit(2)
	}

	misc.SetupLogger(config.DebugLevel)

	db, err := database.SetupDatabase(config)
	if err != nil {
		misc.Logger.Criticalf("Failed to set up database: [%v]", err)
		os.Exit(2)
	}

	err = db.Connect()
	if err != nil {
		misc.Logger.Criticalf("Failed to connect to database: [%v]", err)
		os.Exit(2)
	}

	controller := web.SetupController(config, db)

	controller.HandleRequests()
}
