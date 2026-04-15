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
)

// Error types
const (
	ErrorTypeAPI        = "API_ERROR"
	ErrorTypeNetwork    = "NETWORK_ERROR"
	ErrorTypeAuth       = "AUTH_ERROR"
	ErrorTypeParam      = "PARAM_ERROR"
	ErrorTypeResource   = "RESOURCE_ERROR"
	ErrorTypeSystem     = "SYSTEM_ERROR"
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

// ClassifyException classifies an error based on keywords
func ClassifyException(err error, context string) *CLIError {
	errMsg := strings.ToLower(err.Error())

	// Network errors
	networkKeywords := []string{"timeout", "connection", "network", "dns", "econnrefused", "connection error", "connection refused"}
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

	// Auth errors - 401
	authKeywords401 := []string{"401", "unauthorized", "invalid api key", "api key"}
	for _, keyword := range authKeywords401 {
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

	// Auth errors - 403
	authKeywords403 := []string{"403", "forbidden", "permission denied"}
	for _, keyword := range authKeywords403 {
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

	// Resource errors - 404
	resourceKeywords := []string{"404", "not found", "does not exist", "version_not_found", "app not found"}
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

	// File not found
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

	// Parameter errors
	paramKeywords := []string{"json", "decode", "parse", "path traversal", "invalid path", "unsupported"}
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

	// API errors
	if strings.Contains(errMsg, "api") || strings.Contains(errMsg, "status code") {
		return NewCLIError(
			fmt.Sprintf("API error: %s", err.Error()),
			ErrorTypeAPI,
			ExitGeneralError,
			true,
			"Try again later or contact support.",
			map[string]interface{}{"context": context},
		)
	}

	// Default: System error
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
