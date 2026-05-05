package utils

import (
	"net"
	"net/http"
	"strings"

	"github.com/google/uuid"
)

func GetClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (can contain multiple IPs)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0]) // First IP is the original client
	}

	// Check X-Real-IP header
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return strings.TrimSpace(xrip)
	}

	// Fallback to RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return ip
	}

	return r.RemoteAddr
}

func ValidateUUID(s string) bool {
	if err := uuid.Validate(s); err != nil {
		return false
	}
	return true
}

func GenerateUUID() string {
	return uuid.NewString()
}

func IntPtrValue(val *int, fallback int) int {
	if val == nil {
		return fallback
	}
	return *val
}

func StringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
