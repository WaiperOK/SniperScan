package scanner

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func ParsePorts(spec string) ([]int, error) {
	spec = strings.TrimSpace(spec)
	if spec == "" {
		return nil, fmt.Errorf("ports spec is empty")
	}

	seen := map[int]struct{}{}
	for _, chunk := range strings.Split(spec, ",") {
		chunk = strings.TrimSpace(chunk)
		if chunk == "" {
			continue
		}

		if strings.Contains(chunk, "-") {
			parts := strings.SplitN(chunk, "-", 2)
			start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
			if err != nil {
				return nil, fmt.Errorf("invalid range start %q: %w", parts[0], err)
			}
			end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
			if err != nil {
				return nil, fmt.Errorf("invalid range end %q: %w", parts[1], err)
			}
			if start < 1 || end > 65535 || start > end {
				return nil, fmt.Errorf("invalid range %d-%d", start, end)
			}
			for port := start; port <= end; port++ {
				seen[port] = struct{}{}
			}
			continue
		}

		port, err := strconv.Atoi(chunk)
		if err != nil {
			return nil, fmt.Errorf("invalid port %q: %w", chunk, err)
		}
		if port < 1 || port > 65535 {
			return nil, fmt.Errorf("invalid port %d", port)
		}
		seen[port] = struct{}{}
	}

	if len(seen) == 0 {
		return nil, fmt.Errorf("no valid ports parsed")
	}

	out := make([]int, 0, len(seen))
	for port := range seen {
		out = append(out, port)
	}
	sort.Ints(out)
	return out, nil
}
