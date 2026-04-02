package security

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"strings"
)

// TLSConfig holds TLS verification settings.
type TLSConfig struct {
	// MinTLSVersion is the minimum TLS version (e.g. "1.2", "1.3").
	MinTLSVersion string
	// RequireFullStrict enables Cloudflare Full (Strict) mode verification.
	RequireFullStrict bool
	// WebSocketWSS forces wss:// for all WebSocket connections.
	WebSocketWSS bool
	// GRPCUnixSocket confirms gRPC uses Unix domain sockets (no network exposure).
	GRPCUnixSocket bool
}

// DefaultTLSConfig returns the recommended production TLS configuration.
func DefaultTLSConfig() *TLSConfig {
	return &TLSConfig{
		MinTLSVersion:    "1.3",
		RequireFullStrict: true,
		WebSocketWSS:     true,
		GRPCUnixSocket:   true,
	}
}

// VerifyTLSConfig validates the TLS configuration and checks that all
// transport channels meet security requirements.
func VerifyTLSConfig(cfg *TLSConfig, wsURL, grpcAddr string) error {
	var errors []string

	// 1. TLS version check
	if cfg.MinTLSVersion != "1.2" && cfg.MinTLSVersion != "1.3" {
		errors = append(errors, fmt.Sprintf("unsupported TLS version: %s (must be 1.2 or 1.3)", cfg.MinTLSVersion))
	}

	// 2. WebSocket must use wss://
	if cfg.WebSocketWSS && wsURL != "" {
		if err := verifyWSS(wsURL); err != nil {
			errors = append(errors, err.Error())
		}
	}

	// 3. gRPC must use Unix socket (not network address)
	if cfg.GRPCUnixSocket && grpcAddr != "" {
		if err := verifyUnixSocket(grpcAddr); err != nil {
			errors = append(errors, err.Error())
		}
	}

	// 4. Verify origin certificate files exist (Full Strict mode)
	if cfg.RequireFullStrict {
		if err := verifyOriginCerts(); err != nil {
			errors = append(errors, err.Error())
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("TLS config verification failed:\n  - %s", strings.Join(errors, "\n  - "))
	}

	return nil
}

// verifyWSS ensures the WebSocket URL uses wss:// scheme.
func verifyWSS(wsURL string) error {
	u, err := url.Parse(wsURL)
	if err != nil {
		return fmt.Errorf("invalid WebSocket URL %q: %w", wsURL, err)
	}

	switch u.Scheme {
	case "wss":
		return nil // good
	case "ws":
		return fmt.Errorf("WebSocket URL %q uses insecure ws:// scheme — must use wss://", wsURL)
	default:
		return fmt.Errorf("WebSocket URL %q has unexpected scheme %q — expected wss://", wsURL, u.Scheme)
	}
}

// verifyUnixSocket ensures the gRPC address is a Unix domain socket.
func verifyUnixSocket(addr string) error {
	// Check for unix:// prefix
	if strings.HasPrefix(addr, "unix://") {
		socketPath := strings.TrimPrefix(addr, "unix://")
		if _, err := os.Stat(socketPath); err != nil {
			return fmt.Errorf("Unix socket %q does not exist", socketPath)
		}
		return nil
	}

	// Check if it's a file path (no port, no colon)
	if !strings.Contains(addr, ":") {
		if _, err := os.Stat(addr); err != nil {
			return fmt.Errorf("Unix socket %q does not exist", addr)
		}
		return nil
	}

	// It's a network address (host:port) — verify it's localhost
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return fmt.Errorf("gRPC address %q is not a Unix socket and cannot be parsed as host:port", addr)
	}

	// Only loopback addresses are acceptable as fallback
	switch host {
	case "localhost", "127.0.0.1", "::1":
		return nil
	default:
		return fmt.Errorf("gRPC address %q exposes gRPC to network (must use Unix socket)", addr)
	}
}

// verifyOriginCerts checks that origin certificate files exist for Full Strict mode.
func verifyOriginCerts() error {
	certPaths := []string{
		"/etc/ssl/cinaseek/origin.pem",
		"/etc/ssl/cinaseek/origin-key.pem",
	}

	var missing []string
	for _, p := range certPaths {
		if _, err := os.Stat(p); err != nil {
			missing = append(missing, p)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("origin certificates not found (required for Full Strict mode): %s — run deploy/cloudflare/SSL.md setup first", strings.Join(missing, ", "))
	}

	return nil
}
