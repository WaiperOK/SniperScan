package server

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/WaiperOK/SniperScan/internal/metrics"
	"github.com/WaiperOK/SniperScan/internal/scanner"
)

type ScanRequest struct {
	Target      string `json:"target"`
	Ports       string `json:"ports"`
	TimeoutMS   int    `json:"timeout_ms"`
	Concurrency int    `json:"concurrency"`
	Banner      bool   `json:"banner"`
}

func NewHandler(defaultTimeout time.Duration, defaultConcurrency int) http.Handler {
	mux := http.NewServeMux()
	counters := &metrics.Counters{}

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/metrics", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; version=0.0.4")
		_, _ = w.Write([]byte(counters.RenderPrometheus()))
	})

	mux.HandleFunc("/v1/scan", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method_not_allowed"})
			return
		}

		var req ScanRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			counters.AddError()
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid_json"})
			return
		}

		report, err := runScan(req, defaultTimeout, defaultConcurrency)
		if err != nil {
			counters.AddError()
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		counters.AddScan(report.Summary.PortsScanned, report.Summary.OpenPorts)
		writeJSON(w, http.StatusOK, report)
	})

	return mux
}

func runScan(req ScanRequest, defaultTimeout time.Duration, defaultConcurrency int) (scanner.Report, error) {
	if req.Target == "" {
		return scanner.Report{}, errors.New("target is required")
	}

	ports, err := scanner.ParsePorts(req.Ports)
	if err != nil {
		return scanner.Report{}, err
	}

	timeout := defaultTimeout
	if req.TimeoutMS > 0 {
		timeout = time.Duration(req.TimeoutMS) * time.Millisecond
	}

	concurrency := defaultConcurrency
	if req.Concurrency > 0 {
		concurrency = req.Concurrency
	}

	return scanner.Scan(context.Background(), scanner.Options{
		Target:      req.Target,
		Ports:       ports,
		Timeout:     timeout,
		Concurrency: concurrency,
		GrabBanner:  req.Banner,
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
