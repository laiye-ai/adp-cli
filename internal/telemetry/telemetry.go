package telemetry

import (
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// Event is the telemetry payload sent to the statistics endpoint.
type Event struct {
	Command  string   `json:"command"`
	Metadata Metadata `json:"metadata"`
}

// Metadata holds detailed telemetry data for a single command invocation.
type Metadata struct {
	Command      string         `json:"command"`
	Params       map[string]any `json:"params"`
	StartedAt    string         `json:"started_at"`
	EndedAt      string         `json:"ended_at"`
	DurationMs   int64          `json:"duration_ms"`
	Status       string         `json:"status"`
	ErrorMessage string         `json:"error_message"`
	CLIVersion   string         `json:"cli_version"`
	OS           string         `json:"os"`
	Arch         string         `json:"arch"`
	Terminal     string         `json:"terminal"`
	Source       string         `json:"source"`
}

// skipCommands lists command paths that should NOT be tracked.
var skipCommands = map[string]bool{
	"adp config":       true,
	"adp config set":   true,
	"adp config get":   true,
	"adp config clear": true,
	"adp version":      true,
	"adp help":         true,
	"adp":              true, // bare root command (no subcommand)
}

// globalFlags lists flag names to exclude from params collection.
var globalFlags = map[string]bool{
	"json":   true,
	"quiet":  true,
	"lang":   true,
	"source": true,
}

// argKeyNames maps command paths to their positional argument key name.
var argKeyNames = map[string]string{
	"adp parse local":    "path",
	"adp parse url":      "url",
	"adp parse base64":   "base64",
	"adp parse query":    "task_id",
	"adp extract local":  "path",
	"adp extract url":    "url",
	"adp extract base64": "base64",
	"adp extract query":  "task_id",
}

var (
	mu        sync.Mutex
	current   *state
	lastError error
)

type state struct {
	command   string
	params    map[string]any
	startedAt time.Time
	version   string
	source    string
	skip      bool
}

// Begin records the start of a command execution.
// Call this from PersistentPreRun.
func Begin(cmd *cobra.Command, args []string, version, source string) {
	mu.Lock()
	defer mu.Unlock()

	commandPath := cmd.CommandPath()

	s := &state{
		command:   commandPath,
		startedAt: time.Now().UTC(),
		version:   version,
		source:    source,
		params:    make(map[string]any),
	}

	// Check if this command should be skipped
	if skipCommands[commandPath] {
		s.skip = true
		current = s
		return
	}

	// Collect user-provided flags (exclude global flags)
	cmd.Flags().Visit(func(f *pflag.Flag) {
		if globalFlags[f.Name] {
			return
		}
		s.params[f.Name] = f.Value.String()
	})

	// Collect positional arguments with a meaningful key name
	if len(args) > 0 {
		key := argKeyNames[commandPath]
		if key == "" {
			key = "args"
		}
		if len(args) == 1 {
			s.params[key] = args[0]
		} else {
			s.params[key] = args
		}
	}

	current = s
	lastError = nil
}

// SetError records an error that occurred during command execution.
// Call this from command Run functions when an error is encountered.
func SetError(err error) {
	if err == nil {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	lastError = err
}

// End finalises the telemetry event and returns it.
// Returns nil if the command should not be tracked.
func End() *Event {
	mu.Lock()
	defer mu.Unlock()

	if current == nil || current.skip {
		return nil
	}

	now := time.Now().UTC()
	duration := now.Sub(current.startedAt).Milliseconds()

	status := "success"
	errMsg := ""
	if lastError != nil {
		status = "failure"
		errMsg = lastError.Error()
	}

	event := &Event{
		Command: current.command,
		Metadata: Metadata{
			Command:      current.command,
			Params:       current.params,
			StartedAt:    current.startedAt.Format(time.RFC3339),
			EndedAt:      now.Format(time.RFC3339),
			DurationMs:   duration,
			Status:       status,
			ErrorMessage: errMsg,
			CLIVersion:   current.version,
			OS:           runtime.GOOS,
			Arch:         runtime.GOARCH,
			Terminal:     detectTerminal(),
			Source:       current.source,
		},
	}

	// Reset state
	current = nil
	lastError = nil

	return event
}

// detectTerminal attempts to identify the terminal/shell environment.
func detectTerminal() string {
	// 1. Check TERM_PROGRAM (set by many modern terminals)
	if tp := os.Getenv("TERM_PROGRAM"); tp != "" {
		return strings.ToLower(tp)
	}

	// 2. Platform-specific detection
	if runtime.GOOS == "windows" {
		// WT_SESSION is set inside Windows Terminal
		if os.Getenv("WT_SESSION") != "" {
			return "windows-terminal"
		}
		// PSModulePath is a reliable indicator of PowerShell
		if os.Getenv("PSModulePath") != "" {
			return "powershell"
		}
		return "cmd"
	}

	// 3. Unix: check SHELL
	if shell := os.Getenv("SHELL"); shell != "" {
		parts := strings.Split(shell, "/")
		return parts[len(parts)-1]
	}

	return "unknown"
}
