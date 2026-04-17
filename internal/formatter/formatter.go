package formatter

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/laiye-ai/adp-cli/internal/errors"
)

// Formatter handles output formatting
type Formatter struct {
	jsonMode  bool
	quietMode bool
	isTTY     bool
}

// New creates a new Formatter
func New(jsonMode, quietMode bool) *Formatter {
	return &Formatter{
		jsonMode:  jsonMode,
		quietMode: quietMode,
		isTTY:     isTTY(),
	}
}

func isTTY() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}

// SetJSONMode sets JSON output mode
func (f *Formatter) SetJSONMode(enabled bool) {
	f.jsonMode = enabled
}

// SetQuietMode sets quiet mode
func (f *Formatter) SetQuietMode(enabled bool) {
	f.quietMode = enabled
}

// IsTTY returns whether the output is a TTY
func (f *Formatter) IsTTY() bool {
	return f.isTTY
}

// PrintSuccess prints a success message
func (f *Formatter) PrintSuccess(msg string) {
	if f.quietMode {
		return
	}
	if f.isTTY {
		fmt.Fprintf(os.Stderr, "\033[32m%s\033[0m\n", msg)
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", msg)
	}
}

// PrintError prints an error message
func (f *Formatter) PrintError(msg string) {
	if f.isTTY {
		fmt.Fprintf(os.Stderr, "\033[31m%s\033[0m\n", msg)
	} else {
		fmt.Fprintf(os.Stderr, "error: %s\n", msg)
	}
}

// PrintErrorJSON prints structured error as JSON
func (f *Formatter) PrintErrorJSON(errType, message, fix string, retryable bool, details map[string]interface{}) {
	if f.isTTY && !f.jsonMode {
		fmt.Fprintf(os.Stderr, "\033[31mError: %s\033[0m\n", message)
		if fix != "" {
			fmt.Fprintf(os.Stderr, "  Fix: %s\n", fix)
		}
		fmt.Fprintf(os.Stderr, "  Retryable: %v\n", retryable)
		if details != nil {
			fmt.Fprintf(os.Stderr, "  Details: %v\n", details)
		}
	} else {
		err := map[string]interface{}{
			"type":      errType,
			"message":   message,
			"fix":       fix,
			"retryable": retryable,
			"details":   details,
		}
		jsonBytes, _ := json.MarshalIndent(err, "", "  ")
		fmt.Fprintln(os.Stderr, string(jsonBytes))
	}
}

// PrintCLIError prints a CLIError with structured output
func (f *Formatter) PrintCLIError(err *errors.CLIError) {
	if f.isTTY && !f.jsonMode {
		fmt.Fprintf(os.Stderr, "\033[31m%s\033[0m\n", err.Message)
		if err.Fix != "" {
			fmt.Fprintf(os.Stderr, "  Fix: %s\n", err.Fix)
		}
		fmt.Fprintf(os.Stderr, "  Type: %s\n", err.Type)
		fmt.Fprintf(os.Stderr, "  Code: %d\n", err.Code)
		fmt.Fprintf(os.Stderr, "  Retryable: %v\n", err.Retryable)
		if err.Details != nil {
			fmt.Fprintf(os.Stderr, "  Details: %v\n", err.Details)
		}
	} else {
		output := map[string]interface{}{
			"type":      err.Type,
			"code":      err.Code,
			"message":   err.Message,
			"fix":       err.Fix,
			"retryable": err.Retryable,
			"details":   err.Details,
		}
		jsonBytes, _ := json.MarshalIndent(output, "", "  ")
		fmt.Fprintln(os.Stderr, string(jsonBytes))
	}
}

// ExitWithError prints error and exits with the error code
func (f *Formatter) ExitWithError(err *errors.CLIError) {
	f.PrintCLIError(err)
	os.Exit(err.Code)
}

// PrintWarning prints a warning message
func (f *Formatter) PrintWarning(msg string) {
	if f.quietMode {
		return
	}
	if f.isTTY {
		fmt.Fprintf(os.Stderr, "\033[33m%s\033[0m\n", msg)
	} else {
		fmt.Fprintf(os.Stderr, "warning: %s\n", msg)
	}
}

// PrintInfo prints an info message
func (f *Formatter) PrintInfo(msg string) {
	if f.quietMode {
		return
	}
	if f.isTTY {
		fmt.Fprintf(os.Stderr, "\033[34m%s\033[0m\n", msg)
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", msg)
	}
}

// PrintJSON prints data as JSON
func (f *Formatter) PrintJSON(data interface{}) {
	var output []byte
	if f.jsonMode || !f.isTTY {
		output, _ = json.MarshalIndent(data, "", "  ")
	} else {
		output, _ = json.MarshalIndent(data, "", "  ")
	}
	fmt.Println(string(output))
}

// PrintSection prints a section header
func (f *Formatter) PrintSection(title string) {
	if f.quietMode {
		return
	}
	if f.isTTY {
		separator := strings.Repeat("-", len(title))
		fmt.Fprintf(os.Stderr, "\n\033[1m%s\033[0m\n", title)
		fmt.Fprintf(os.Stderr, "\033[2m%s\033[0m\n", separator)
	} else {
		fmt.Fprintf(os.Stderr, "\n%s\n", title)
	}
}

// PrintProgress prints progress information
func (f *Formatter) PrintProgress(current, total int, message string) {
	if f.quietMode {
		return
	}
	if !f.isTTY || f.jsonMode {
		// JSON Lines format for Agent consumption
		progress := map[string]interface{}{
			"type":    "progress",
			"current": current,
			"total":   total,
		}
		if message != "" {
			progress["file"] = message
		}
		jsonBytes, _ := json.Marshal(progress)
		fmt.Fprintln(os.Stderr, string(jsonBytes))
		return
	}
	percentage := float64(current) / float64(total) * 100
	progress := fmt.Sprintf("[%d/%d] %.1f%%", current, total, percentage)
	if message != "" {
		progress += fmt.Sprintf(" - %s", message)
	}
	fmt.Fprintf(os.Stderr, "\033[36m%s\033[0m\n", progress)
}

// PrintTaskResult prints task result
func (f *Formatter) PrintTaskResult(taskID string, status int, result map[string]interface{}) {
	f.PrintInfo(fmt.Sprintf("Task_ID: %s", taskID))
	f.PrintInfo(fmt.Sprintf("Status: %s", getStatusText(status)))
	if result != nil {
		f.PrintJSON(result)
	}
}

// PrintFileList prints a list of files
func (f *Formatter) PrintFileList(files []string, showSize bool) {
	if f.quietMode {
		return
	}
	f.PrintSection("Files")
	for i, file := range files {
		fmt.Fprintf(os.Stderr, "%d\t%s\n", i+1, file)
	}
}

// PrintResults prints processing results
func (f *Formatter) PrintResults(results []map[string]interface{}, mode string) {
	if f.quietMode {
		return
	}
	if len(results) == 0 {
		return
	}
	f.PrintJSON(results)
}

func getStatusText(status int) string {
	switch status {
	case 0:
		return "PENDING"
	case 2:
		return "RUNNING"
	case 4:
		return "SUCCESS"
	case 5:
		return "FAILED"
	case 6:
		return "CANCELLED"
	default:
		return "UNKNOWN"
	}
}
