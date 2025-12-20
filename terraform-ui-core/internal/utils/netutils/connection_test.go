package netutils

import (
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"testing"
	"time"
)

func TestTestTcp(t *testing.T) {
	// build listener
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	// start fake listener
	go func() { l.Accept() }()

	// wait a bit for the listener to start
	time.Sleep(500 * time.Millisecond)

	addr := l.Addr().String()

	// test successful connection
	ok, err := testTcp(addr, 1*time.Second)
	if ok != true || err != nil {
		t.Fatal(err)
	}

	// test failed connection
	ok, err = testTcp("localhost:4891", 1*time.Second)
	if ok != false || err == nil {
		t.Fatal(err)
	}
}

func TestTestTcpTls(t *testing.T) {
	// load tls certificate
	cer, err := tls.LoadX509KeyPair("../../../test/testdata/config/dummy-cert.pem", "../../../test/testdata/config/dummy-key.pem")
	if err != nil {
		t.Fatal(err)
	}

	config := &tls.Config{Certificates: []tls.Certificate{cer}}

	// build tls listener
	l, err := tls.Listen("tcp", "localhost:0", config)
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	go func() {
		conn, err := l.Accept()
		if err != nil {
			t.Fatal(err)
			return
		}

		defer conn.Close()

		// bidirectional echo loop to complete TLS handshake!
		io.Copy(conn, conn)
	}()

	// wait a bit for the listener to start
	time.Sleep(500 * time.Millisecond)

	addr := l.Addr().String()

	// test successful tcp connection
	ok, err := testTcpTls(addr, 3*time.Second)
	if !ok || err != nil {
		t.Fatal(err)
	}

	tcpAddr := l.Addr().(*net.TCPAddr)
	port := tcpAddr.Port

	// test failed tcp connection
	ok, err = testTcpTls(fmt.Sprintf("localhost:%s", string(rune(port+1))), 3*time.Second)
	if ok || err == nil {
		t.Fatal(err)
	}
}

func TestWaitForConnection(t *testing.T) {
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	addr := l.Addr().String()

	err = WaitForConnection(addr, 1*time.Second)
	if err != nil {
		t.Fatal(err)
	}
}
