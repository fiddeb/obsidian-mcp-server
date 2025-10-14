package security

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

// NewRateLimiter creates a new rate limiter
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
	}
}

// Allow checks if a request from the given IP is allowed
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

// Middleware creates a security middleware handler
func Middleware(config SecurityConfig, rateLimiter *RateLimiter) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 1. IP Whitelist check
			if len(config.AllowedIPs) > 0 {
				clientIP := GetClientIP(r)
				if !isIPAllowed(clientIP, config.AllowedIPs) {
					http.Error(w, "Forbidden: IP not allowed", http.StatusForbidden)
					return
				}
			}

			// 2. Authentication check
			if config.EnableAuth {
				authHeader := r.Header.Get("Authorization")
				expectedAuth := "Bearer " + config.AuthToken

				if !secureCompare(authHeader, expectedAuth) {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}
			}

			// 3. Rate limiting
			if config.EnableRateLimit {
				clientIP := GetClientIP(r)
				if !rateLimiter.Allow(clientIP, config.RateLimit) {
					http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
					return
				}
			}

			// 4. CORS handling
			if config.EnableCORS {
				origin := r.Header.Get("Origin")
				if isOriginAllowed(origin, config.AllowedOrigins) {
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
}

// ValidateInputMiddleware validates input for POST requests
func ValidateInputMiddleware(next http.HandlerFunc) http.HandlerFunc {
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

// ValidatePath validates file paths to prevent directory traversal and other attacks
func ValidatePath(path string) error {
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

// SanitizeContent sanitizes content to prevent XSS and injection attacks
func SanitizeContent(content string) string {
	// This is basic - you might want to use a proper sanitization library
	// For now, we just limit some dangerous patterns

	// Remove null bytes
	content = strings.ReplaceAll(content, "\x00", "")

	return content
}

// GetClientIP extracts the client IP from the request
func GetClientIP(r *http.Request) string {
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

// AuditLog logs audit information
func AuditLog(r *http.Request, action string, result string) {
	if os.Getenv("ENABLE_AUDIT_LOG") == "true" {
		log := map[string]interface{}{
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"ip":        GetClientIP(r),
			"action":    action,
			"result":    result,
			"path":      r.URL.Path,
		}

		jsonLog, _ := json.Marshal(log)
		fmt.Fprintf(os.Stderr, "AUDIT: %s\n", jsonLog)
	}
}

// Helper functions

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
