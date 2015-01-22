package misc

import (
	"os"

	"github.com/kdar/factorlog"
)

var (
	Logger *factorlog.FactorLog
)

func SetupLogger(debugLevel int) {
	Logger = factorlog.New(os.Stdout, factorlog.NewStdFormatter("[%{Date} %{Time}] {%{SEVERITY}:%{File}/%{PkgFunction}:%{Line}} %{SafeMessage}"))
	Logger.SetMinMaxSeverity(factorlog.Severity(1<<uint(debugLevel)), factorlog.PANIC)
}
