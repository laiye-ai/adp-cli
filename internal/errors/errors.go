package errors

import (
	"fmt"
	"strings"
)

// Exit codes
const (
	ExitSuccess           = 0
	ExitGeneralError      = 1
	ExitParameterError    = 2
	ExitResourceNotFound  = 3
	ExitPermissionDenied  = 4
	ExitConflict          = 5
	ExitPartialFailure    = 6
)

// Error types
const (
	ErrorTypeAPI        = "API_ERROR"
	ErrorTypeNetwork    = "NETWORK_ERROR"
	ErrorTypeAuth       = "AUTH_ERROR"
	ErrorTypeParam      = "PARAM_ERROR"
	ErrorTypeResource   = "RESOURCE_ERROR"
	ErrorTypeSystem     = "SYSTEM_ERROR"
	ErrorTypeConflict   = "CONFLICT_ERROR"
)

// CLIError represents a structured CLI error
type CLIError struct {
	Type      string                 // Error type
	Code      int                    // Exit code
	Message   string                 // Error message
	Fix       string                 // Fix suggestion
	Retryable bool                   // Whether the error is retryable
	Details   map[string]interface{} // Additional details
}

func (e *CLIError) Error() string {
	return e.Message
}

// NewCLIError creates a new CLIError
func NewCLIError(message, errorType string, exitCode int, retryable bool, fix string, details map[string]interface{}) *CLIError {
	return &CLIError{
		Type:      errorType,
		Code:      exitCode,
		Message:   message,
		Fix:       fix,
		Retryable: retryable,
		Details:   details,
	}
}

// classifyByHTTPStatus attempts to classify an error by HTTP status code pattern.
// Returns nil if no status code pattern is found.
func classifyByHTTPStatus(errMsg, context string, originalErr error) *CLIError {
	// Match "status code NNN" pattern from API client responses
	statusPatterns := []struct {
		code      string
		errType   string
		exitCode  int
		retryable bool
		prefix    string
		fix       string
	}{
		{"status code 400", ErrorTypeParam, ExitParameterError, false, "Bad request", "Check the input parameters are correct."},
		{"status code 401", ErrorTypeAuth, ExitPermissionDenied, false, "Authentication error", "Check your API key is correct and has not expired."},
		{"status code 403", ErrorTypeAuth, ExitPermissionDenied, false, "Permission denied", "You do not have permission to access this resource."},
		{"status code 404", ErrorTypeResource, ExitResourceNotFound, false, "Resource not found", "Check the resource ID or path is correct."},
		{"status code 409", ErrorTypeConflict, ExitConflict, false, "Conflict", "The resource already exists or conflicts with the current state."},
		{"status code 429", ErrorTypeAPI, ExitGeneralError, true, "Rate limited", "Too many requests. Try again later."},
		{"status code 5", ErrorTypeAPI, ExitGeneralError, true, "Server error", "Try again later or contact support."},
	}

	for _, p := range statusPatterns {
		if strings.Contains(errMsg, p.code) {
			return NewCLIError(
				fmt.Sprintf("%s: %s", p.prefix, originalErr.Error()),
				p.errType,
				p.exitCode,
				p.retryable,
				p.fix,
				map[string]interface{}{"context": context},
			)
		}
	}
	return nil
}

// ClassifyException classifies an error based on keywords
func ClassifyException(err error, context string) *CLIError {
	errMsg := strings.ToLower(err.Error())

	// 1. HTTP status code based classification (highest priority, most precise)
	if cliErr := classifyByHTTPStatus(errMsg, context, err); cliErr != nil {
		return cliErr
	}

	// 2. Auth errors — specific phrases only
	authKeywords := []string{"unauthorized", "invalid api key", "api key expired", "authentication failed"}
	for _, keyword := range authKeywords {
		if strings.Contains(errMsg, keyword) {
			return NewCLIError(
				fmt.Sprintf("Authentication error: %s", err.Error()),
				ErrorTypeAuth,
				ExitPermissionDenied,
				false,
				"Check your API key is correct and has not expired.",
				map[string]interface{}{"context": context},
			)
		}
	}

	// 3. Permission errors
	permKeywords := []string{"forbidden", "permission denied", "access denied"}
	for _, keyword := range permKeywords {
		if strings.Contains(errMsg, keyword) {
			return NewCLIError(
				fmt.Sprintf("Permission denied: %s", err.Error()),
				ErrorTypeAuth,
				ExitPermissionDenied,
				false,
				"You do not have permission to access this resource.",
				map[string]interface{}{"context": context},
			)
		}
	}

	// 4. Resource errors
	resourceKeywords := []string{"not found", "does not exist", "version_not_found", "app not found"}
	for _, keyword := range resourceKeywords {
		if strings.Contains(errMsg, keyword) {
			return NewCLIError(
				fmt.Sprintf("Resource not found: %s", err.Error()),
				ErrorTypeResource,
				ExitResourceNotFound,
				false,
				"Check the resource ID or path is correct.",
				map[string]interface{}{"context": context},
			)
		}
	}

	// 5. Conflict errors
	conflictKeywords := []string{"conflict", "already exists", "duplicate"}
	for _, keyword := range conflictKeywords {
		if strings.Contains(errMsg, keyword) {
			return NewCLIError(
				fmt.Sprintf("Conflict: %s", err.Error()),
				ErrorTypeConflict,
				ExitConflict,
				false,
				"The resource already exists or conflicts with the current state.",
				map[string]interface{}{"context": context},
			)
		}
	}

	// 6. File not found (local filesystem)
	fileKeywords := []string{"file not found", "no such file", "enoent", "path not found"}
	for _, keyword := range fileKeywords {
		if strings.Contains(errMsg, keyword) {
			return NewCLIError(
				fmt.Sprintf("File not found: %s", err.Error()),
				ErrorTypeResource,
				ExitResourceNotFound,
				false,
				"Check the file path is correct.",
				map[string]interface{}{"context": context},
			)
		}
	}

	// 7. Parameter errors — use specific phrases, avoid overly broad words like "parse"
	paramKeywords := []string{"failed to parse json", "failed to decode", "json decode", "json unmarshal",
		"path traversal", "invalid path", "unsupported file type", "invalid value", "missing required"}
	for _, keyword := range paramKeywords {
		if strings.Contains(errMsg, keyword) {
			return NewCLIError(
				fmt.Sprintf("Parameter error: %s", err.Error()),
				ErrorTypeParam,
				ExitParameterError,
				false,
				"Check the input parameters are correct.",
				map[string]interface{}{"context": context},
			)
		}
	}

	// 8. Network errors — transport-level failures only
	networkKeywords := []string{"dial tcp", "connection refused", "connection reset",
		"no such host", "dns lookup", "econnrefused", "econnreset", "i/o timeout",
		"tls handshake", "certificate", "network is unreachable"}
	for _, keyword := range networkKeywords {
		if strings.Contains(errMsg, keyword) {
			return NewCLIError(
				fmt.Sprintf("Network error: %s", err.Error()),
				ErrorTypeNetwork,
				ExitGeneralError,
				true,
				"Check your network connection and try again.",
				map[string]interface{}{"context": context},
			)
		}
	}

	// 9. Generic API errors
	if strings.Contains(errMsg, "status code") {
		return NewCLIError(
			fmt.Sprintf("API error: %s", err.Error()),
			ErrorTypeAPI,
			ExitGeneralError,
			true,
			"Try again later or contact support.",
			map[string]interface{}{"context": context},
		)
	}

	// 10. Default: System error
	return NewCLIError(
		fmt.Sprintf("System error: %s", err.Error()),
		ErrorTypeSystem,
		ExitGeneralError,
		false,
		"An unexpected error occurred.",
		map[string]interface{}{"context": context},
	)
}

// IsNotFoundError checks if an error is a "not found" error
func IsNotFoundError(err error) bool {
	errMsg := strings.ToLower(err.Error())
	notFoundKeywords := []string{"404", "not found", "does not exist", "version_not_found", "app not found"}
	for _, keyword := range notFoundKeywords {
		if strings.Contains(errMsg, keyword) {
			return true
		}
	}
	return false
}

// IsCLIErrorNotFound checks if a CLIError is a not found error
func IsCLIErrorNotFound(cliErr *CLIError) bool {
	return cliErr.Type == ErrorTypeResource && cliErr.Code == ExitResourceNotFound
}

// ValidateEnum checks if a value is in the allowed enum values
func ValidateEnum(value string, allowed []string, paramName string) *CLIError {
	for _, v := range allowed {
		if value == v {
			return nil
		}
	}
	return NewCLIError(
		fmt.Sprintf("Invalid value for --%s: %s (allowed: %v)", paramName, value, allowed),
		ErrorTypeParam,
		ExitParameterError,
		false,
		fmt.Sprintf("Use one of the allowed values: %v", allowed),
		nil,
	)
}
