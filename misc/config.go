package misc

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

// Configuration stores all configuration values required by the application
type Configuration struct {
	// DatabaseType represents the database type to be used as a backend
	DatabaseType int
	// DatabaseHost represents the hostname of the database backend
	DatabaseHost string
	// DatabasePort represents the port of the database backend
	DatabasePort int
	// DatabaseSchema represents the schema/collection of the database backend
	DatabaseSchema string
	// DatabaseUser represents the username used to authenticate with the database backend
	DatabaseUser string
	// DatabasePassword represents the password used to authenticate with the database backend
	DatabasePassword string
	// DebugLevel represents the debug level for log messages
	DebugLevel int
	// DebugTemplates toggles the reloading of all templates for every request
	DebugTemplates bool
	// HTTPHost represents the hostname/IP the application should listen to for requests
	HTTPHost string
	// HTTPPort represents the port the application should listen to for requests
	HTTPPort int
}

// LoadConfig creates a Configuration by either using commandline flags or a configuration file, returning an error if the parsing failed
func LoadConfig() (*Configuration, error) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: eveauth [options]\n")
		flag.PrintDefaults()
		os.Exit(2)
	}

	flag.Parse()

	var config *Configuration
	var err error

	if len(*configFileFlag) > 0 {
		config, err = ParseJSONConfig(*configFileFlag)
	} else {
		config = ParseCommandlineFlags()
		err = nil
	}

	return config, err
}

// ParseJSONConfig parses a Configuration from a JSON encoded file, returning an error if the process failed
func ParseJSONConfig(path string) (*Configuration, error) {
	configFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var config *Configuration

	err = json.NewDecoder(configFile).Decode(&config)

	return config, err
}
