package middleware

import (
	"net/http"
	"strings"
)

// ClientIP retrieves client's real IP from request headers
func ClientIP(r *http.Request) string {
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	} else {
		ip = strings.Split(ip, ",")[0]
	}
	ip = strings.Split(ip, ":")[0]
	return strings.TrimSpace(ip)
}
