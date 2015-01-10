package misc

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func ParseConfigFlags() {
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
			DebugLevel: *debugLevelFlag,
			HTTPHost:   *httpHostFlag,
			HTTPPort:   *httpPortFlag,
		}
	}
}
