package misc

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type Configuration struct {
	DatabaseType       int
	DatabaseHost       string
	DatabasePort       int
	DatabaseSchema     string
	DatabaseUser       string
	DatabasePassword   string
	DebugLevel         int
	DebugTemplates     bool
	HTTPHost           string
	HTTPPort           int
	EVESSOClientID     string
	EVESSOClientSecret string
	EVESSOCallbackURL  string
}

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

func ParseJSONConfig(path string) (*Configuration, error) {
	configFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	var config *Configuration

	err = json.NewDecoder(configFile).Decode(&config)

	return config, err
}
