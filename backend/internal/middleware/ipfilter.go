package middleware

import (
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
)

// IPFilterConfig holds settings for IP filtering.
type IPFilterConfig struct {
	// WhitelistCIDRs are CIDR ranges that are always allowed.
	WhitelistCIDRs []string
	// BlacklistCIDRs are CIDR ranges that are always denied.
	BlacklistCIDRs []string
	// DefaultAllow controls whether non-matching IPs are allowed.
	DefaultAllow bool
}

// ipFilter holds compiled CIDR networks.
type ipFilter struct {
	whitelist    []*net.IPNet
	blacklist    []*net.IPNet
	defaultAllow bool
	mu           sync.RWMutex
}

// newIPFilter compiles CIDR strings into net.IPNet slices.
func newIPFilter(cfg *IPFilterConfig) (*ipFilter, error) {
	f := &ipFilter{
		defaultAllow: cfg.DefaultAllow,
	}

	for _, cidr := range cfg.WhitelistCIDRs {
		_, network, err := net.ParseCIDR(normalizeCIDR(cidr))
		if err != nil {
			return nil, err
		}
		f.whitelist = append(f.whitelist, network)
	}

	for _, cidr := range cfg.BlacklistCIDRs {
		_, network, err := net.ParseCIDR(normalizeCIDR(cidr))
		if err != nil {
			return nil, err
		}
		f.blacklist = append(f.blacklist, network)
	}

	return f, nil
}

// normalizeCIDR converts a plain IP to CIDR notation.
func normalizeCIDR(s string) string {
	if strings.Contains(s, "/") {
		return s
	}
	// Plain IP → /32 for IPv4 or /128 for IPv6
	if strings.Contains(s, ":") {
		return s + "/128"
	}
	return s + "/32"
}

// isAllowed checks whether an IP is permitted.
func (f *ipFilter) isAllowed(ip net.IP) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	// Check whitelist first — if set, only whitelisted IPs are allowed.
	if len(f.whitelist) > 0 {
		matched := false
		for _, network := range f.whitelist {
			if network.Contains(ip) {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// Check blacklist — always deny if matched.
	for _, network := range f.blacklist {
		if network.Contains(ip) {
			return false
		}
	}

	return f.defaultAllow
}

// addToBlacklist adds a CIDR to the blacklist at runtime.
func (f *ipFilter) addToBlacklist(cidr string) error {
	_, network, err := net.ParseCIDR(normalizeCIDR(cidr))
	if err != nil {
		return err
	}
	f.mu.Lock()
	f.blacklist = append(f.blacklist, network)
	f.mu.Unlock()
	return nil
}

// addToWhitelist adds a CIDR to the whitelist at runtime.
func (f *ipFilter) addToWhitelist(cidr string) error {
	_, network, err := net.ParseCIDR(normalizeCIDR(cidr))
	if err != nil {
		return err
	}
	f.mu.Lock()
	f.whitelist = append(f.whitelist, network)
	f.mu.Unlock()
	return nil
}

// IPFilter returns a middleware that filters requests by IP.
// IPs in the blacklist are denied; if a whitelist is set, only those IPs are allowed.
func IPFilter(cfg *IPFilterConfig) gin.HandlerFunc {
	filter, err := newIPFilter(cfg)
	if err != nil {
		slog.Error("failed to compile IP filter rules", "error", err)
		// Fail-closed: deny all if config is broken
		return func(c *gin.Context) {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "访问被拒绝",
			})
			c.Abort()
		}
	}

	return func(c *gin.Context) {
		ip := net.ParseIP(c.ClientIP())
		if ip == nil {
			slog.Warn("failed to parse client IP", "ip", c.ClientIP())
			c.JSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "访问被拒绝",
			})
			c.Abort()
			return
		}

		if !filter.isAllowed(ip) {
			slog.Warn("IP blocked by filter", "ip", c.ClientIP(), "path", c.Request.URL.Path)
			c.JSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "访问被拒绝",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
