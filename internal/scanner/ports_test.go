package scanner

import "testing"

func TestParsePorts(t *testing.T) {
	ports, err := ParsePorts("22,80,443,8000-8002,443")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := []int{22, 80, 443, 8000, 8001, 8002}
	if len(ports) != len(expected) {
		t.Fatalf("expected %d ports, got %d", len(expected), len(ports))
	}
	for i := range expected {
		if ports[i] != expected[i] {
			t.Fatalf("expected port %d at index %d, got %d", expected[i], i, ports[i])
		}
	}
}

func TestParsePortsRejectsInvalidRange(t *testing.T) {
	_, err := ParsePorts("9000-8000")
	if err == nil {
		t.Fatal("expected error for invalid range")
	}
}
