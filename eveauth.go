// eveauth provides an authentication backend used for EVE Online services.
// Beside providing basic user authentication, the application can provide access using different permissions based on access rights.
package main

import (
	"log"
	"os"
	"runtime"

	"github.com/morpheusxaut/eveauth/database"
	"github.com/morpheusxaut/eveauth/mail"
	"github.com/morpheusxaut/eveauth/misc"
	"github.com/morpheusxaut/eveauth/session"
	"github.com/morpheusxaut/eveauth/web"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	config, err := misc.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: [%v]", err)
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

	mailer := mail.SetupMailController(config, db)

	sessionController, err := session.SetupSessionController(config, db, mailer)
	if err != nil {
		misc.Logger.Criticalf("Failed to set up session controller: [%v]", err)
		os.Exit(2)
	}

	templates := web.SetupTemplates(sessionController)

	checksums, err := web.SetupAssetChecksums()
	if err != nil {
		misc.Logger.Criticalf("Failed to calculate asset checkums: [%v]", err)
		os.Exit(2)
	}

	controller := web.SetupController(config, db, sessionController, templates, checksums)

	controller.HandleRequests()
}
