package scanner

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"
)

func TestScanDetectsOpenPort(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen failed: %v", err)
	}
	defer listener.Close()

	port := listener.Addr().(*net.TCPAddr).Port

	done := make(chan struct{})
	go func() {
		defer close(done)
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			return
		}
		defer conn.Close()
		_, _ = conn.Write([]byte("SSH-2.0-TestServer"))
	}()

	report, err := Scan(context.Background(), Options{
		Target:      "127.0.0.1",
		Ports:       []int{port},
		Timeout:     200 * time.Millisecond,
		Concurrency: 4,
		GrabBanner:  true,
	})
	if err != nil {
		t.Fatalf("scan failed: %v", err)
	}

	if report.Summary.OpenPorts != 1 {
		t.Fatalf("expected 1 open port, got %d", report.Summary.OpenPorts)
	}
	if report.Results[0].Port != port {
		t.Fatalf("expected open port %d, got %d", port, report.Results[0].Port)
	}

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("accept goroutine timeout")
	}
}

func TestSanitizeBanner(t *testing.T) {
	in := fmt.Sprintf("OK\x00\x01%v", string([]byte{7, 8, 9}))
	out := sanitizeBanner(in)
	if out != "OK" {
		t.Fatalf("expected sanitized banner to be OK, got %q", out)
	}
}
