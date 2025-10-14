package main

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

// SecurityConfig holds security settings
type SecurityConfig struct {
	EnableAuth      bool     `yaml:"enable_auth"`
	AuthToken       string   `yaml:"auth_token"`
	AllowedIPs      []string `yaml:"allowed_ips"`
	EnableRateLimit bool     `yaml:"enable_rate_limit"`
	RateLimit       int      `yaml:"rate_limit"` // requests per minute
	EnableCORS      bool     `yaml:"enable_cors"`
	AllowedOrigins  []string `yaml:"allowed_origins"`
}

// RateLimiter simple rate limiter
type RateLimiter struct {
	requests map[string][]time.Time
}

func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
	}
}

func (rl *RateLimiter) Allow(ip string, limit int) bool {
	now := time.Now()
	cutoff := now.Add(-time.Minute)

	// Clean old requests
	if times, ok := rl.requests[ip]; ok {
		var recent []time.Time
		for _, t := range times {
			if t.After(cutoff) {
				recent = append(recent, t)
			}
		}
		rl.requests[ip] = recent
	}

	// Check limit
	if len(rl.requests[ip]) >= limit {
		return false
	}

	// Add request
	rl.requests[ip] = append(rl.requests[ip], now)
	return true
}

// Security middleware
func (s *MCPServer) securityMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. IP Whitelist check
		if len(s.config.Security.AllowedIPs) > 0 {
			clientIP := getClientIP(r)
			if !isIPAllowed(clientIP, s.config.Security.AllowedIPs) {
				http.Error(w, "Forbidden: IP not allowed", http.StatusForbidden)
				return
			}
		}

		// 2. Authentication check
		if s.config.Security.EnableAuth {
			authHeader := r.Header.Get("Authorization")
			expectedAuth := "Bearer " + s.config.Security.AuthToken

			if !secureCompare(authHeader, expectedAuth) {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}

		// 3. Rate limiting
		if s.config.Security.EnableRateLimit {
			clientIP := getClientIP(r)
			if !s.rateLimiter.Allow(clientIP, s.config.Security.RateLimit) {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}
		}

		// 4. CORS handling
		if s.config.Security.EnableCORS {
			origin := r.Header.Get("Origin")
			if isOriginAllowed(origin, s.config.Security.AllowedOrigins) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		next(w, r)
	}
}

// Input validation middleware
func (s *MCPServer) validateInput(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Check content type
		if r.Method == "POST" && r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
			return
		}

		// Limit request body size (1MB)
		r.Body = http.MaxBytesReader(w, r.Body, 1*1024*1024)

		next(w, r)
	}
}

// Path validation - prevent directory traversal
func validatePath(path string) error {
	// Check for directory traversal attempts
	if strings.Contains(path, "..") {
		return fmt.Errorf("invalid path: directory traversal not allowed")
	}

	// Check for absolute paths
	if strings.HasPrefix(path, "/") {
		return fmt.Errorf("invalid path: absolute paths not allowed")
	}

	// Check for null bytes and other control characters
	if strings.Contains(path, "\x00") {
		return fmt.Errorf("invalid path: contains null byte")
	}

	// Check for common dangerous patterns that could cause issues
	// Note: We allow & and other characters that are valid in Obsidian filenames
	dangerous := []string{"~", "$", "`", "|", ";"}
	for _, d := range dangerous {
		if strings.Contains(path, d) {
			return fmt.Errorf("invalid path: contains dangerous character '%s'", d)
		}
	}

	// Validate that path ends with .md
	if !strings.HasSuffix(path, ".md") && !strings.HasSuffix(path, "/") {
		// Allow paths without extension only if they're folder references
		if strings.Contains(path, ".") {
			return fmt.Errorf("invalid path: only .md files are supported")
		}
	}

	return nil
}

// Sanitize content - prevent XSS and injection
func sanitizeContent(content string) string {
	// This is basic - you might want to use a proper sanitization library
	// For now, we just limit some dangerous patterns

	// Remove null bytes
	content = strings.ReplaceAll(content, "\x00", "")

	return content
}

// Helper functions

func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func isIPAllowed(clientIP string, allowedIPs []string) bool {
	for _, allowed := range allowedIPs {
		if allowed == "*" || allowed == clientIP {
			return true
		}

		// Support CIDR notation
		if strings.Contains(allowed, "/") {
			_, ipNet, err := net.ParseCIDR(allowed)
			if err == nil && ipNet.Contains(net.ParseIP(clientIP)) {
				return true
			}
		}
	}
	return len(allowedIPs) == 0 // Allow all if no IPs specified
}

func isOriginAllowed(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if allowed == "*" || allowed == origin {
			return true
		}
	}
	return len(allowedOrigins) == 0
}

func secureCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// Audit logging
func (s *MCPServer) auditLog(r *http.Request, action string, result string) {
	if os.Getenv("ENABLE_AUDIT_LOG") == "true" {
		log := map[string]interface{}{
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"ip":        getClientIP(r),
			"action":    action,
			"result":    result,
			"path":      r.URL.Path,
		}

		jsonLog, _ := json.Marshal(log)
		fmt.Fprintf(os.Stderr, "AUDIT: %s\n", jsonLog)
	}
}
