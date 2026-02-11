package scanner

import (
	"context"
	"fmt"
	"io"
	"net"
	"sort"
	"sync"
	"time"
)

func scanPort(ctx context.Context, target string, port int, timeout time.Duration, grabBanner bool) PortResult {
	addr := net.JoinHostPort(target, fmt.Sprintf("%d", port))
	started := time.Now()
	dialer := net.Dialer{Timeout: timeout}

	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return PortResult{
			Port:      port,
			State:     "closed",
			RTTMillis: time.Since(started).Milliseconds(),
		}
	}
	defer conn.Close()

	result := PortResult{
		Port:      port,
		State:     "open",
		RTTMillis: time.Since(started).Milliseconds(),
	}

	if grabBanner {
		_ = conn.SetDeadline(time.Now().Add(timeout))
		_, _ = conn.Write([]byte("\r\n"))
		buf := make([]byte, 128)
		n, readErr := conn.Read(buf)
		if readErr == nil || readErr == io.EOF {
			result.Banner = sanitizeBanner(string(buf[:n]))
		}
	}

	result.Service = guessService(port, result.Banner)
	return result
}

func sanitizeBanner(input string) string {
	if len(input) == 0 {
		return ""
	}
	out := make([]rune, 0, len(input))
	for _, ch := range input {
		if ch >= 32 && ch < 127 {
			out = append(out, ch)
		}
	}
	if len(out) > 80 {
		out = out[:80]
	}
	return string(out)
}

func Scan(ctx context.Context, opts Options) (Report, error) {
	if opts.Target == "" {
		return Report{}, fmt.Errorf("target is required")
	}
	if len(opts.Ports) == 0 {
		return Report{}, fmt.Errorf("at least one port must be provided")
	}
	if opts.Timeout <= 0 {
		opts.Timeout = 400 * time.Millisecond
	}
	if opts.Concurrency <= 0 {
		opts.Concurrency = 200
	}

	start := time.Now()
	jobs := make(chan int)
	results := make(chan PortResult, len(opts.Ports))

	var wg sync.WaitGroup
	for i := 0; i < opts.Concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for port := range jobs {
				results <- scanPort(ctx, opts.Target, port, opts.Timeout, opts.GrabBanner)
			}
		}()
	}

	for _, port := range opts.Ports {
		jobs <- port
	}
	close(jobs)

	wg.Wait()
	close(results)

	report := Report{
		StartedAt: start.UTC().Format(time.RFC3339),
	}
	for item := range results {
		report.Results = append(report.Results, item)
	}

	sort.Slice(report.Results, func(i, j int) bool {
		return report.Results[i].Port < report.Results[j].Port
	})

	openPorts := 0
	for _, item := range report.Results {
		if item.State == "open" {
			openPorts++
		}
	}

	finished := time.Now()
	report.FinishedAt = finished.UTC().Format(time.RFC3339)
	report.Summary = Summary{
		Target:         opts.Target,
		PortsScanned:   len(report.Results),
		OpenPorts:      openPorts,
		ClosedPorts:    len(report.Results) - openPorts,
		DurationMillis: finished.Sub(start).Milliseconds(),
	}

	return report, nil
}
