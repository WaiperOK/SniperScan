package metrics

import (
	"fmt"
	"sync"
)

type Counters struct {
	mu                  sync.Mutex
	RequestsTotal       int
	ScansTotal          int
	PortsScanned        int
	OpenPortsDiscovered int
	ErrorsTotal         int
}

func (c *Counters) AddScan(ports, open int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.RequestsTotal++
	c.ScansTotal++
	c.PortsScanned += ports
	c.OpenPortsDiscovered += open
}

func (c *Counters) AddError() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.RequestsTotal++
	c.ErrorsTotal++
}

func (c *Counters) RenderPrometheus() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	return fmt.Sprintf(
		"# TYPE sniperscan_requests_total counter\n"+
			"sniperscan_requests_total %d\n"+
			"# TYPE sniperscan_scans_total counter\n"+
			"sniperscan_scans_total %d\n"+
			"# TYPE sniperscan_ports_scanned_total counter\n"+
			"sniperscan_ports_scanned_total %d\n"+
			"# TYPE sniperscan_open_ports_total counter\n"+
			"sniperscan_open_ports_total %d\n"+
			"# TYPE sniperscan_errors_total counter\n"+
			"sniperscan_errors_total %d\n",
		c.RequestsTotal,
		c.ScansTotal,
		c.PortsScanned,
		c.OpenPortsDiscovered,
		c.ErrorsTotal,
	)
}
