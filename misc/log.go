package misc

import (
	"github.com/kdar/factorlog"
	"os"
)

var (
	Logger *factorlog.FactorLog
)

func SetupLogger(debugLevel int) {
	Logger = factorlog.New(os.Stdout, factorlog.NewStdFormatter("[%{Date} %{Time}] {%{SEVERITY}:%{File}/%{Function}:%{Line}} %{SafeMessage}"))
	Logger.SetMinMaxSeverity(factorlog.Severity(1<<uint(debugLevel)), factorlog.PANIC)
}
