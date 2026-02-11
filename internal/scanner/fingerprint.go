package scanner

import (
	"strings"
)

func guessService(port int, banner string) string {
	b := strings.ToLower(banner)
	switch {
	case strings.Contains(b, "ssh") || port == 22:
		return "ssh"
	case strings.Contains(b, "smtp") || port == 25:
		return "smtp"
	case strings.Contains(b, "http") || port == 80 || port == 8080 || port == 443:
		return "http"
	case strings.Contains(b, "mysql") || port == 3306:
		return "mysql"
	case strings.Contains(b, "redis") || port == 6379:
		return "redis"
	default:
		return "unknown"
	}
}
