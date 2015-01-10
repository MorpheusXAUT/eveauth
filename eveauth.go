// eveauth provides an authentication backend used for EVE Online services.
// Beside providing basic user authentication, the application can provide access using different permissions based on access rights.
package main

import (
	"github.com/morpheusxaut/eveauth/misc"
)

func main() {
	misc.ParseConfigFlags()

	misc.SetupLogger(uint(misc.Config.DebugLevel))

	misc.Logger.Infoln("Logger set up! Continuing initialisation...")
}
