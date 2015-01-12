package misc

import (
	"github.com/kdar/factorlog"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	Logger *factorlog.FactorLog
)

func SetupLogger(debugLevel int) {
	Logger = factorlog.New(os.Stdout, factorlog.NewStdFormatter("[%{Date} %{Time}] {%{SEVERITY}:%{File}/%{Function}:%{Line}} %{SafeMessage}"))
	Logger.SetMinMaxSeverity(factorlog.Severity(1<<uint(debugLevel)), factorlog.PANIC)
}

func WebLogging(inner http.Handler, name string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		remoteAddr := r.Header.Get("X-Forwarded-For")

		if len(remoteAddr) > 0 {
			remoteAddrs := strings.Split(remoteAddr, ", ")
			if len(remoteAddrs) > 1 {
				r.RemoteAddr = remoteAddrs[0]
			} else {
				r.RemoteAddr = remoteAddr
			}
		}

		inner.ServeHTTP(w, r)

		Logger.Debugf("WebLogging: [%s] %s %q {%s} - %s ", r.Method, r.RemoteAddr, r.RequestURI, name, time.Since(start))
	})
}
