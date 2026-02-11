package scanner

import "time"

type PortResult struct {
	Port      int    `json:"port"`
	State     string `json:"state"`
	Banner    string `json:"banner,omitempty"`
	Service   string `json:"service,omitempty"`
	RTTMillis int64  `json:"rtt_ms"`
}

type Summary struct {
	Target         string `json:"target"`
	PortsScanned   int    `json:"ports_scanned"`
	OpenPorts      int    `json:"open_ports"`
	ClosedPorts    int    `json:"closed_ports"`
	DurationMillis int64  `json:"duration_ms"`
}

type Report struct {
	StartedAt  string       `json:"started_at"`
	FinishedAt string       `json:"finished_at"`
	Summary    Summary      `json:"summary"`
	Results    []PortResult `json:"results"`
}

type Options struct {
	Target      string
	Ports       []int
	Timeout     time.Duration
	Concurrency int
	GrabBanner  bool
}
