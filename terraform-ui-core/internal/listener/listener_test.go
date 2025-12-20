package listener

import (
	"context"
	"testing"
	"time"

	"github.com/tbauriedel/terraform-ui-core/internal/config"
	"github.com/tbauriedel/terraform-ui-core/internal/logging"
)

func TestNewListener(t *testing.T) {
	c := config.Listener{}
	ctx := context.TODO()
	log := logging.NewLoggerStdout(config.Logger{Type: "stdout", Level: "debug"})

	// Just returns the struct, nothing to test
	_ = NewListener(c, ctx, log)
}

func getListener(ctx context.Context) Listener {
	c := config.Listener{
		ListenAddr:  ":0",
		ReadTimeout: 10 * time.Second,
		IdleTimeout: 120 * time.Second,
	}

	log := logging.NewLoggerStdout(config.Logger{Type: "stdout", Level: "debug"})

	l := Listener{
		logger:      log,
		config:      c,
		context:     ctx,
		multiplexer: nil,
	}

	return l
}

func TestListenerInsecure(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	l := getListener(ctx)

	// start listener
	go func() {
		err := l.Start()
		if err != nil {
			t.Fatal(err)
		}
	}()

	// wait for the listener to start and let it run
	time.Sleep(1 * time.Second)

	ctx, cancel = context.WithTimeout(context.TODO(), 1*time.Second)
	defer cancel()

	err := l.Shutdown(ctx)
	if err != nil {
		t.Fatal(err)
	}
}

func TestListenerSecure(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	l := getListener(ctx)

	l.config.TlsEnabled = true
	l.config.TlsSkipVerify = true
	l.config.TlsKeyPath = "../../test/testdata/config/dummy-key.pem"
	l.config.TlsCertPath = "../../test/testdata/config/dummy-cert.pem"

	// start listener
	go func() {
		err := l.Start()
		if err != nil {
			t.Fatal(err)
		}
	}()

	// wait for the listener to start and let it run
	time.Sleep(1 * time.Second)

	ctx, cancel = context.WithTimeout(context.TODO(), 1*time.Second)
	defer cancel()

	err := l.Shutdown(ctx)
	if err != nil {
		t.Fatal(err)
	}
}
