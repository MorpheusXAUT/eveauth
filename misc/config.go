package misc

var (
	Config *Configuration
)

type Configuration struct {
	DebugLevel int
	HTTPHost   string
	HTTPPort   int
}
