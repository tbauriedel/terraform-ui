package netutils

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

// testTcp tests tcp connections without tls
func testTcp(addr string, timeout time.Duration) (bool, error) {
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err == nil {
		conn.Close()
		return true, nil
	}
	return false, err
}

// testTcpTls tests tcp connections with tls
func testTcpTls(addr string, timeout time.Duration) (bool, error) {
	dialer := &net.Dialer{Timeout: timeout}
	config := &tls.Config{
		InsecureSkipVerify: true,
	}

	conn, err := tls.DialWithDialer(dialer, "tcp", addr, config)
	if err == nil {
		conn.Close()
		return true, nil
	}

	return false, err
}

// WaitForConnection waits for a TCP connection to be established. Each try waits 3 seconds for a successful connection. Retries every 500ms until the timeout is reached
// If timeout is reached, an error is returned.
// Tests tls and non-tls connections
//
// 'addr' is in the format "host:port"
func WaitForConnection(addr string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	connectionTimeout := 3 * time.Second

	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout waiting for '%s'", addr)
		}

		// Test tls connection
		if tlsOK, _ := testTcpTls(addr, connectionTimeout); tlsOK {
			return nil
		}

		// test plain tcp connection
		if tcpOK, _ := testTcp(addr, connectionTimeout); tcpOK {
			return nil
		}

		time.Sleep(500 * time.Millisecond)
	}
}
