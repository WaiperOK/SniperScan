package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestScanEndpoint(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen failed: %v", err)
	}
	defer listener.Close()
	port := listener.Addr().(*net.TCPAddr).Port

	go func() {
		conn, acceptErr := listener.Accept()
		if acceptErr == nil {
			_ = conn.Close()
		}
	}()

	handler := NewHandler(150*time.Millisecond, 16)
	ts := httptest.NewServer(handler)
	defer ts.Close()

	body, _ := json.Marshal(map[string]any{
		"target":      "127.0.0.1",
		"ports":       fmt.Sprintf("%d", port),
		"timeout_ms":  200,
		"concurrency": 8,
		"banner":      false,
	})

	resp, err := http.Post(ts.URL+"/v1/scan", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("post failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}

func TestHealthEndpoint(t *testing.T) {
	ts := httptest.NewServer(NewHandler(200*time.Millisecond, 10))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/healthz")
	if err != nil {
		t.Fatalf("health call failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}
}
