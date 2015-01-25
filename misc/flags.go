package misc

import (
	"flag"
)

var (
	databaseTypeFlag     = flag.Int("dbtype", 1, "Selects the type of database to use for backend connections")
	databaseHostFlag     = flag.String("dbhost", "", "Hostname/path of the database to use for backend connections")
	databasePortFlag     = flag.Int("dbport", 0, "Port of the database to use for backend connections")
	databaseSchemaFlag   = flag.String("dbschema", "", "Name of the database schema to use for backend connections")
	databaseUserFlag     = flag.String("dbuser", "", "Username for the database to use for backend connections")
	databasePasswordFlag = flag.String("dbpassword", "", "Password for the database to use for backend connections")
	debugLevelFlag       = flag.Int("debug", 3, "Sets the debug level (0-9), lower number displays more messages")
	debugTemplatesFlag   = flag.Bool("debugtemplates", false, "Enables reloading of HTML templates on every request to help development")
	httpHostFlag         = flag.String("host", "0.0.0.0", "Hostname for the webserver to bind to")
	httpPortFlag         = flag.Int("port", 5000, "Port for the webserver to bind to")
	smtpHostFlag         = flag.String("smtphost", "localhost", "Hostname of the SMTP server to use")
	smtpPortFlag         = flag.Int("smtpport", 587, "Port of the SMTP server to use")
	smtpUserFlag         = flag.String("smtpuser", "", "Username used to authenticate with the SMTP server")
	smtpPasswordFlag     = flag.String("smtppassword", "", "Hostname of the SMTP server to use")
	smtpSenderFlag       = flag.String("smtpsender", "", "Address to send emails as")
	configFileFlag       = flag.String("config", "", "Path to the config file to parse")
)

// ParseCommandlineFlags parses the command line flags used with the application
func ParseCommandlineFlags() *Configuration {
	config := &Configuration{
		DatabaseType:     *databaseTypeFlag,
		DatabaseHost:     *databaseHostFlag,
		DatabasePort:     *databasePortFlag,
		DatabaseSchema:   *databaseSchemaFlag,
		DatabaseUser:     *databaseUserFlag,
		DatabasePassword: *databasePasswordFlag,
		DebugLevel:       *debugLevelFlag,
		DebugTemplates:   *debugTemplatesFlag,
		HTTPHost:         *httpHostFlag,
		HTTPPort:         *httpPortFlag,
		SMTPHost:         *smtpHostFlag,
		SMTPPort:         *smtpPortFlag,
		SMTPUser:         *smtpUserFlag,
		SMTPPassword:     *smtpPasswordFlag,
		SMTPSender:       *smtpSenderFlag,
	}

	return config
}
