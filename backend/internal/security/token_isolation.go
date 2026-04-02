package security

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// sensitivePatterns matches patterns that should never appear in client-distributed code.
var sensitivePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)CLOUDFLARE_API_TOKEN`),
	regexp.MustCompile(`(?i)CLOUDFLARE_TUNNEL_TOKEN`),
	regexp.MustCompile(`(?i)CLOUDFLARE_API_KEY`),
	regexp.MustCompile(`(?i)CF_API_TOKEN`),
	regexp.MustCompile(`(?i)CF_TUNNEL_TOKEN`),
	regexp.MustCompile(`(?i)cloudflared.*credentials`),
	regexp.MustCompile(`(?i)CERTPEM`),
	regexp.MustCompile(`(?i)origin-key\.pem`),
	regexp.MustCompile(`(?i)origin\.pem`),
	regexp.MustCompile(`(?i)credentials\.json`),
}

// hardcodedSecretPatterns matches common secret formats that may be hardcoded.
var hardcodedSecretPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(api[_-]?key|secret|token|password)\s*[:=]\s*["'][A-Za-z0-9_\-]{20,}["']`),
	regexp.MustCompile(`-----BEGIN (RSA |EC )?PRIVATE KEY-----`),
	regexp.MustCompile(`eyJ[A-Za-z0-9_-]{20,}\.eyJ[A-Za-z0-9_-]{20,}`), // JWT-like
}

// VerifyTokenIsolation scans the codebase at configPath and verifies that no
// Cloudflare tokens, API keys or other server-side secrets leak into
// client-distributable code.
func VerifyTokenIsolation(configPath string) error {
	var violations []string

	// 1. Scan .env files for CLOUDFLARE_* variables (must not be in client dist)
	if err := scanEnvFiles(configPath, &violations); err != nil {
		return fmt.Errorf("env scan: %w", err)
	}

	// 2. Scan Go source for hardcoded secrets
	if err := scanGoSource(configPath, &violations); err != nil {
		return fmt.Errorf("source scan: %w", err)
	}

	// 3. Scan for credential files that should not be distributed
	if err := scanCredentialFiles(configPath, &violations); err != nil {
		return fmt.Errorf("credential file scan: %w", err)
	}

	if len(violations) > 0 {
		return fmt.Errorf("token isolation violations found:\n  - %s", strings.Join(violations, "\n  - "))
	}

	return nil
}

// scanEnvFiles checks .env files for CLOUDFLARE_* variables.
func scanEnvFiles(root string, violations *[]string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		// Skip .git, node_modules, vendor
		base := filepath.Base(path)
		if d.IsDir() && (base == ".git" || base == "node_modules" || base == "vendor" || base == "bin") {
			return filepath.SkipDir
		}

		if !d.IsDir() && (strings.HasPrefix(base, ".env") || strings.HasSuffix(base, ".env")) {
			if err := checkEnvFile(path, violations); err != nil {
				return err
			}
		}
		return nil
	})
}

func checkEnvFile(path string, violations *[]string) error {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		for _, pat := range sensitivePatterns {
			if pat.MatchString(line) {
				*violations = append(*violations,
					fmt.Sprintf("%s:%d — sensitive env variable detected: matches %s", path, lineNum, pat.String()))
			}
		}
	}
	return scanner.Err()
}

// scanGoSource walks Go files and checks for hardcoded secrets.
func scanGoSource(root string, violations *[]string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		base := filepath.Base(path)
		if d.IsDir() && (base == ".git" || base == "node_modules" || base == "vendor" || base == "bin") {
			return filepath.SkipDir
		}

		if !d.IsDir() && strings.HasSuffix(base, ".go") {
			if err := checkGoFile(path, violations); err != nil {
				return err
			}
		}
		return nil
	})
}

func checkGoFile(path string, violations *[]string) error {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		for _, pat := range hardcodedSecretPatterns {
			if pat.MatchString(line) {
				*violations = append(*violations,
					fmt.Sprintf("%s:%d — potential hardcoded secret: matches %s", path, lineNum, pat.String()))
			}
		}
	}
	return scanner.Err()
}

// scanCredentialFiles checks for credential files that should not be in the repo.
func scanCredentialFiles(root string, violations *[]string) error {
	dangerousFiles := []string{
		"credentials.json",
		"cert.pem",
		"origin.pem",
		"origin-key.pem",
	}

	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		base := filepath.Base(path)
		if d.IsDir() && (base == ".git" || base == "node_modules" || base == "vendor") {
			return filepath.SkipDir
		}
		if d.IsDir() {
			return nil
		}
		for _, dangerous := range dangerousFiles {
			if base == dangerous {
				*violations = append(*violations,
					fmt.Sprintf("%s — credential file should not be in source tree", path))
			}
		}
		return nil
	})
}
