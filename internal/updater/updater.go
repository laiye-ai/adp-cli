package updater

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

const (
	repo            = "laiye-ai/adp-cli"
	checkInterval   = 24 * time.Hour
	requestTimeout  = 3 * time.Second
	cacheFileName   = "version_check.json"
)

type cacheFile struct {
	CheckedAt   time.Time `json:"checked_at"`
	LatestVersion string  `json:"latest_version"`
}


// CheckAndNotify runs a version check asynchronously.
// Returns a channel that yields the update message (or empty string if no update).
// Caller should receive from the channel after the main command finishes.
func CheckAndNotify(currentVersion string, quiet bool, jsonMode bool) <-chan string {
	ch := make(chan string, 1)

	if quiet || jsonMode || currentVersion == "dev" {
		ch <- ""
		return ch
	}

	go func() {
		log.Info().Msg("Checking for new version")
		latest, err := resolveLatestVersion()
		if err != nil {
			log.Warn().Err(err).Msg("Failed to check latest version")
			ch <- ""
			return
		}
		if latest == "" {
			log.Debug().Msg("No latest version found")
			ch <- ""
			return
		}

		current := strings.TrimPrefix(currentVersion, "v")
		latestClean := strings.TrimPrefix(latest, "v")

		if current == latestClean || !isNewer(latestClean, current) {
			log.Info().Str("current_version", current).Msg("Already on latest version")
			ch <- ""
			return
		}

		log.Info().Str("current_version", current).Str("latest_version", latestClean).Msg("Update available")

		// "  ║  Update available: v{current} → v{latest}{pad}║"
		// Fixed content length: "  ║  Update available: v" = 22, " → v" = 4, "{pad}║" at end
		// Total inner width = 58 chars (between the ║ borders)
		inner := 58
		updateLine := fmt.Sprintf("Update available: v%s → v%s", current, latestClean)
		pad := strings.Repeat(" ", max(0, inner-2-len(updateLine)))
		msg := fmt.Sprintf(
			"\n  ╔══════════════════════════════════════════════════════════╗\n"+
				"  ║  %s%s║\n"+
				"  ║  Run: npm update -g @laiye-adp/agentic-doc-parse-...  ║\n"+
				"  ╚══════════════════════════════════════════════════════════╝\n",
			updateLine, pad,
		)
		ch <- msg
	}()

	return ch
}

// resolveLatestVersion returns the latest version, using cache if fresh enough.
func resolveLatestVersion() (string, error) {
	cacheDir := getCacheDir()
	cachePath := filepath.Join(cacheDir, cacheFileName)

	// Try reading cache
	if cached, ok := readCache(cachePath); ok {
		return cached.LatestVersion, nil
	}

	// Fetch from GitHub
	latest, err := fetchLatestVersion()
	if err != nil {
		return "", err
	}

	// Write cache even if empty (e.g. 404 / no releases), to avoid hammering the API
	writeCache(cachePath, latest)

	return latest, nil
}

func readCache(path string) (*cacheFile, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, false
	}

	var c cacheFile
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, false
	}

	if time.Since(c.CheckedAt) > checkInterval {
		return nil, false
	}

	return &c, true
}

func writeCache(path string, latestVersion string) {
	c := cacheFile{
		CheckedAt:     time.Now(),
		LatestVersion: latestVersion,
	}
	data, err := json.Marshal(c)
	if err != nil {
		return
	}
	_ = os.MkdirAll(filepath.Dir(path), 0o755)
	_ = os.WriteFile(path, data, 0o644)
}

func fetchLatestVersion() (string, error) {
	// Use releases/latest/download/ and extract version from the 302 redirect URL.
	// This avoids the GitHub API entirely (no rate-limit, no auth needed).
	probeURL := fmt.Sprintf("https://github.com/%s/releases/latest/download/adp-%s-%s%s",
		repo, mapOS(runtime.GOOS), mapArch(runtime.GOARCH), archiveExt())

	client := &http.Client{
		Timeout: requestTimeout,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse // don't follow redirects
		},
	}

	req, err := http.NewRequest("HEAD", probeURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "adp-cli-updater/1.0 ("+runtime.GOOS+")")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", nil
	}

	// 301/302 Location: .../download/v1.2.3/adp-linux-x64.tar.gz
	loc := resp.Header.Get("Location")
	if loc == "" {
		return "", fmt.Errorf("no redirect from releases/latest")
	}

	// Extract version from path: .../download/{version}/...
	parts := strings.Split(loc, "/")
	for i, p := range parts {
		if p == "download" && i+1 < len(parts) {
			return parts[i+1], nil
		}
	}

	return "", fmt.Errorf("could not extract version from redirect URL")
}

func archiveExt() string {
	if runtime.GOOS == "windows" {
		return ".zip"
	}
	return ".tar.gz"
}

func mapArch(goarch string) string {
	switch goarch {
	case "amd64":
		return "x64"
	default:
		return goarch // arm64 stays arm64
	}
}

func mapOS(goos string) string {
	switch goos {
	case "windows":
		return "win32"
	default:
		return goos // linux, darwin stay as-is
	}
}

// isNewer returns true if candidate is strictly newer than current.
// Simple string-based semver comparison (major.minor.patch).
func isNewer(candidate, current string) bool {
	return parseVersion(candidate) > parseVersion(current)
}

// parseVersion converts "1.2.3" into a comparable integer 010203.
func parseVersion(v string) int {
	v = strings.TrimPrefix(v, "v")
	parts := strings.SplitN(v, ".", 3)
	result := 0
	multipliers := []int{1000000, 1000, 1}
	for i, p := range parts {
		if i >= len(multipliers) {
			break
		}
		n := 0
		fmt.Sscanf(p, "%d", &n)
		result += n * multipliers[i]
	}
	return result
}

func getCacheDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".adp")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
