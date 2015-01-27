package misc

import (
	"flag"
	"strings"
)

var (
	debugLevelFlag     = flag.Int("debug", 3, "Sets the debug level (0-9), lower number displays more messages")
	debugTemplatesFlag = flag.Bool("templates", false, "Enables reloading of HTML templates on every request to help development")
	httpHostFlag       = flag.String("http", "0.0.0.0:5000", "Hostname:port for the webserver to bind to")
	configFileFlag     = flag.String("config", "config.cfg", "Path to the config file to parse")
)

// ParseCommandlineFlags parses the command line flags used with the application
func ParseCommandlineFlags(config *Configuration) *Configuration {
	if *debugLevelFlag != 3 {
		config.DebugLevel = *debugLevelFlag
	}
	if *debugTemplatesFlag != false {
		config.DebugTemplates = *debugTemplatesFlag
	}
	if !strings.EqualFold(*httpHostFlag, "0.0.0.0:5000") {
		config.HTTPHost = *httpHostFlag
	}

	return config
}
