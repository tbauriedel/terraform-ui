package listener

import (
	"testing"
	"time"

	"github.com/tbauriedel/resource-nexus-core/internal/config"
	"github.com/tbauriedel/resource-nexus-core/internal/logging"
)

func TestWithMiddleWare(t *testing.T) {
	log := logging.NewLoggerStdout(config.Logger{Type: "stdout", Level: "debug"})

	l := NewListener(
		config.Listener{
			ListenAddr:  "localhost:0",
			ReadTimeout: 30 * time.Second,
			IdleTimeout: 120 * time.Second,
			TLSEnabled:  false,
		},
		nil,
		WithMiddleWare(MiddlewareLogging(log)),
	)

	if l.middlewares == nil {
		t.Fatal("middleware should be set")
	}

	if len(l.middlewares) != 1 {
		t.Fatal("applied middlewares should be equal to one")
	}
}
