package netutils

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

// testTCP tests tcp connections without tls.
func testTCP(addr string, timeout time.Duration) (bool, error) {
	d := net.Dialer{}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := d.DialContext(ctx, "tcp", addr)
	if err == nil {
		_ = conn.Close()

		return true, nil
	}

	return false, fmt.Errorf("failed to connect to %s: %w", addr, err)
}

// testTCPWithTLS tests tcp connections with tls.
func testTCPWithTLS(addr string, insecure bool, timeout time.Duration) (bool, error) {
	d := &tls.Dialer{}
	config := &tls.Config{
		InsecureSkipVerify: insecure, //nolint:gosec
	}

	d.Config = config

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	conn, err := d.DialContext(ctx, "tcp", addr)
	if err == nil {
		_ = conn.Close()

		return true, nil
	}

	return false, fmt.Errorf("failed to connect to %s: %w", addr, err)
}

// WaitForConnection waits for a TCP connection to be established.
// Each try waits 3 seconds for a successful connection.
// Retries every 500ms until the timeout is reached
// If timeout is reached, an error is returned.
// Tests tls and non-tls connections
//
// 'addr' is in the format "host:port".
func WaitForConnection(addr string, insecure bool, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	connectionTimeout := 3 * time.Second

	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout waiting for '%s'", addr)
		}

		// Test tls connection
		if tlsOK, _ := testTCPWithTLS(addr, insecure, connectionTimeout); tlsOK {
			return nil
		}

		// test plain tcp connection
		if tcpOK, _ := testTCP(addr, connectionTimeout); tcpOK {
			return nil
		}

		time.Sleep(500 * time.Millisecond)
	}
}
