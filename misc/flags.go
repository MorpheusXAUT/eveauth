package misc

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func ParseConfigFlags() {
	databaseTypeFlag := flag.Int("dbtype", 1, "Selects the type of database to use for backend connections")
	databaseHostFlag := flag.String("dbhost", "", "Hostname/path of the database to use for backend connections")
	databasePortFlag := flag.Int("dbport", 0, "Port of the database to use for backend connections")
	databaseSchemaFlag := flag.String("dbschema", "", "Name of the database schema to use for backend connections")
	databaseUserFlag := flag.String("dbuser", "", "Username for the database to use for backend connections")
	databasePasswordFlag := flag.String("dbpassword", "", "Password for the database to use for backend connections")
	debugLevelFlag := flag.Int("debug", 3, "Sets the debug level (0-9), lower number displays more messages")
	httpHostFlag := flag.String("host", "0.0.0.0", "Hostname for the webserver to bind to")
	httpPortFlag := flag.Int("port", 5000, "Port for the webserver to bind to")
	configFileFlag := flag.String("config", "", "Path to the config file to parse")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: eveauth [options]\n")
		flag.PrintDefaults()
		os.Exit(2)
	}

	flag.Parse()

	if len(*configFileFlag) > 0 {
		configFile, err := os.Open(*configFileFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read config file %q: [%v]\n", *configFileFlag, err)
			os.Exit(2)
		}

		decoder := json.NewDecoder(configFile)

		Config = &Configuration{}

		err = decoder.Decode(&Config)
		if err != nil {
			return
		}
	} else {
		Config = &Configuration{
			DatabaseType:     *databaseTypeFlag,
			DatabaseHost:     *databaseHostFlag,
			DatabasePort:     *databasePortFlag,
			DatabaseSchema:   *databaseSchemaFlag,
			DatabaseUser:     *databaseUserFlag,
			DatabasePassword: *databasePasswordFlag,
			DebugLevel:       *debugLevelFlag,
			HTTPHost:         *httpHostFlag,
			HTTPPort:         *httpPortFlag,
		}
	}
}
