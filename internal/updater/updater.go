package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
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

type githubRelease struct {
	TagName string `json:"tag_name"`
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
		latest, err := resolveLatestVersion()
		if err != nil || latest == "" {
			ch <- ""
			return
		}

		current := strings.TrimPrefix(currentVersion, "v")
		latestClean := strings.TrimPrefix(latest, "v")

		if current == latestClean || !isNewer(latestClean, current) {
			ch <- ""
			return
		}

		// "  ║  Update available: v{current} → v{latest}{pad}║"
		// Fixed content length: "  ║  Update available: v" = 22, " → v" = 4, "{pad}║" at end
		// Total inner width = 58 chars (between the ║ borders)
		inner := 58
		updateLine := fmt.Sprintf("Update available: v%s → v%s", current, latestClean)
		pad := strings.Repeat(" ", max(0, inner-2-len(updateLine)))
		msg := fmt.Sprintf(
			"\n  ╔══════════════════════════════════════════════════════════╗\n"+
				"  ║  %s%s║\n"+
				"  ║  Run: npm update -g agentic-doc-parse-and-extract-cli   ║\n"+
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
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)

	client := &http.Client{Timeout: requestTimeout}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "adp-cli-updater/1.0 ("+runtime.GOOS+")")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", nil // private repo or no releases yet, treat as "no update"
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github api returned %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var release githubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return "", err
	}

	return release.TagName, nil
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
