package misc

import (
	"os"

	"github.com/kdar/factorlog"
)

var (
	// Logger is a accessible logger, providing formatted log messages with different output levels
	Logger *factorlog.FactorLog
)

// SetupLogger configures the logger, setting a custom format and output level
func SetupLogger(debugLevel int) {
	Logger = factorlog.New(os.Stdout, factorlog.NewStdFormatter("[%{Date} %{Time}] {%{SEVERITY}:%{File}/%{PkgFunction}:%{Line}} %{SafeMessage}"))
	Logger.SetMinMaxSeverity(factorlog.Severity(1<<uint(debugLevel)), factorlog.PANIC)
}
