package misc

var (
	Config *Configuration
)

type Configuration struct {
	DatabaseType     int
	DatabaseHost     string
	DatabasePort     int
	DatabaseSchema   string
	DatabaseUser     string
	DatabasePassword string
	DebugLevel       int
	HTTPHost         string
	HTTPPort         int
}
