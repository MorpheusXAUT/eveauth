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
	// DatabaseHost represents the hostname:port of the database backend
	DatabaseHost string
	// DatabaseSchema represents the schema/collection of the database backend
	DatabaseSchema string
	// DatabaseUser represents the username used to authenticate with the database backend
	DatabaseUser string
	// DatabasePassword represents the password used to authenticate with the database backend
	DatabasePassword string
	// RedisHost represents the hostname:port of the Redis data store
	RedisHost string
	// RedisPassword represents the password used to authenticate with the Redis data store
	RedisPassword string
	// RedisDB represents the database to select for the session store
	RedisDB string
	// SMTPHost represents the hostname:port of the SMTP server used for sending mails
	SMTPHost string
	// SMTPStartTLS indicates whether the SMTP connection should use the StartTLS command to secure communications
	SMTPStartTLS bool
	// SMTPUser represents the username used to authenticate with the SMTP server
	SMTPUser string
	// SMTPPassword represents the password used to authenticate with the SMTP server
	SMTPPassword string
	// SMTPSender represents the email address set as the sender of all outgoing emails
	SMTPSender string
	// DebugLevel represents the debug level for log messages
	DebugLevel int
	// DebugTemplates toggles the reloading of all templates for every request
	DebugTemplates bool
	// HTTPHost represents the hostname:port the application should listen to for requests
	HTTPHost string
	// HTTPPublicURL represents the public URL the eveauth app is reachable at
	HTTPPublicURL string
}

// LoadConfig creates a Configuration by either using commandline flags or a configuration file, returning an error if the parsing failed
func LoadConfig() (*Configuration, error) {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: eveauth [options]\n")
		flag.PrintDefaults()
		os.Exit(2)
	}

	flag.Parse()

	config, err := ParseJSONConfig(*configFileFlag)
	config = ParseCommandlineFlags(config)

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
