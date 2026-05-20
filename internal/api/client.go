package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/laiye-ai/adp-cli/internal/config"
	"github.com/rs/zerolog/log"
)

const (
	timeout = 300 * time.Second
)

// TaskStatus constants
const (
	TaskStatusPending   = 0
	TaskStatusRunning   = 2
	TaskStatusSuccess   = 4
	TaskStatusFailed    = 5
	TaskStatusCancelled = 6
)

// Client is the API client for ADP backend
type Client struct {
	baseURL    string
	apiKey     string
	tenantName string
	httpClient *http.Client
}

// NewClient creates a new API client
func NewClient(cfg *config.Config) (*Client, error) {
	apiKey, err := config.GetAPIKey(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt API key: %w", err)
	}

	tenantName := cfg.TenantName
	if tenantName == "" {
		tenantName = "laiye"
	}

	return &Client{
		baseURL:    cfg.APIBaseURL,
		apiKey:     apiKey,
		tenantName: tenantName,
		httpClient: &http.Client{Timeout: timeout},
	}, nil
}

func (c *Client) request(method, endpoint string, data interface{}) (map[string]interface{}, error) {
	url := c.baseURL + endpoint

	var body io.Reader
	if data != nil {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request: %w", err)
		}
		body = bytes.NewBuffer(jsonData)
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Source", "CLI")
	req.Header.Set("X-Api-Key", c.apiKey)

	log.Debug().Str("method", method).Str("url", url).Msg("API request")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		var errMsg string
		var result map[string]interface{}
		if err := json.Unmarshal(respBody, &result); err == nil {
			if msg, ok := result["message"].(string); ok {
				errMsg = msg
			}
		}
		if errMsg == "" {
			errMsg = fmt.Sprintf("HTTP %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("status code %d: %s", resp.StatusCode, errMsg)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}

func (c *Client) encodeFileToBase64(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

// ParseSync synchronously parses a document
func (c *Client) ParseSync(fileURL, appID string, filePath, fileBase64, fileName string) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"app_id":    appID,
		"file_name": fileName,
	}

	if fileBase64 != "" {
		data["file_base64"] = fileBase64
	} else if filePath != "" {
		b64, err := c.encodeFileToBase64(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to encode file: %w", err)
		}
		data["file_base64"] = b64
	} else {
		data["file_url"] = fileURL
	}

	return c.request("POST", fmt.Sprintf("/open/agentic_doc_processor/%s/v1/app/doc/recognize", c.tenantName), data)
}

// ParseAsync creates an async parse task
func (c *Client) ParseAsync(fileURL, appID string, filePath, fileBase64, fileName string) (string, error) {
	data := map[string]interface{}{
		"app_id":    appID,
		"file_name": fileName,
	}

	if fileBase64 != "" {
		data["file_base64"] = fileBase64
	} else if filePath != "" {
		b64, err := c.encodeFileToBase64(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to encode file: %w", err)
		}
		data["file_base64"] = b64
	} else {
		data["file_url"] = fileURL
	}

	resp, err := c.request("POST", fmt.Sprintf("/open/agentic_doc_processor/%s/v1/app/doc/recognize/create/task", c.tenantName), data)
	if err != nil {
		return "", err
	}

	if data, ok := resp["data"].(map[string]interface{}); ok {
		if taskID, ok := data["task_id"].(string); ok {
			return taskID, nil
		}
	}
	return "", fmt.Errorf("task_id not found in response")
}

// QueryParseTask queries parse task status
func (c *Client) QueryParseTask(taskID string) (map[string]interface{}, error) {
	return c.request("GET", fmt.Sprintf("/open/agentic_doc_processor/%s/v1/app/doc/recognize/query/task/%s", c.tenantName, taskID), nil)
}

// ExtractSync synchronously extracts from a document
func (c *Client) ExtractSync(fileURL, appID string, filePath, fileBase64, fileName string, extractConfig map[string]interface{}) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"app_id":          appID,
		"file_name":       fileName,
		"with_rec_result": false,
	}

	if fileBase64 != "" {
		data["file_base64"] = fileBase64
	} else if filePath != "" {
		b64, err := c.encodeFileToBase64(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to encode file: %w", err)
		}
		data["file_base64"] = b64
	} else {
		data["file_url"] = fileURL
	}

	if extractConfig != nil {
		data["extract_config"] = extractConfig
	}

	return c.request("POST", fmt.Sprintf("/open/agentic_doc_processor/%s/v1/app/doc/extract", c.tenantName), data)
}

// ExtractAsync creates an async extract task
func (c *Client) ExtractAsync(fileURL, appID string, filePath, fileBase64, fileName string, extractConfig map[string]interface{}) (string, error) {
	data := map[string]interface{}{
		"app_id":          appID,
		"file_name":       fileName,
		"with_rec_result": false,
	}

	if fileBase64 != "" {
		data["file_base64"] = fileBase64
	} else if filePath != "" {
		b64, err := c.encodeFileToBase64(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to encode file: %w", err)
		}
		data["file_base64"] = b64
	} else {
		data["file_url"] = fileURL
	}

	if extractConfig != nil {
		data["extract_config"] = extractConfig
	}

	resp, err := c.request("POST", fmt.Sprintf("/open/agentic_doc_processor/%s/v1/app/doc/extract/create/task", c.tenantName), data)
	if err != nil {
		return "", err
	}

	if data, ok := resp["data"].(map[string]interface{}); ok {
		if taskID, ok := data["task_id"].(string); ok {
			return taskID, nil
		}
	}
	return "", fmt.Errorf("task_id not found in response")
}

// QueryExtractTask queries extract task status
func (c *Client) QueryExtractTask(taskID string) (map[string]interface{}, error) {
	return c.request("GET", fmt.Sprintf("/open/agentic_doc_processor/%s/v1/app/doc/extract/query/task/%s", c.tenantName, taskID), nil)
}

// WaitForTask waits for a task to complete
func (c *Client) WaitForTask(taskID string, queryFunc func(string) (map[string]interface{}, error), timeout, interval int) (map[string]interface{}, error) {
	deadline := time.Now().Add(time.Duration(timeout) * time.Second)

	for time.Now().Before(deadline) {
		result, err := queryFunc(taskID)
		if err != nil {
			return nil, err
		}

		data := result
		if d, ok := result["data"].(map[string]interface{}); ok {
			data = d
		}

		status := int(data["status"].(float64))
		switch status {
		case TaskStatusSuccess:
			return result, nil
		case TaskStatusFailed:
			return nil, fmt.Errorf("task failed: %v", data["message"])
		case TaskStatusCancelled:
			return nil, fmt.Errorf("task cancelled")
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}

	return nil, fmt.Errorf("task timeout after %d seconds", timeout)
}

// ListApps lists available applications
func (c *Client) ListApps(appType *int, limit int) ([]map[string]interface{}, error) {
	endpoint := fmt.Sprintf("/open/agentic_doc_processor/%s/v1/app-list?limit=%d", c.tenantName, limit)
	if appType != nil {
		endpoint = fmt.Sprintf("/open/agentic_doc_processor/%s/v1/app-list?app_type=%d&limit=%d", c.tenantName, *appType, limit)
	}

	resp, err := c.request("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}

	var apps []map[string]interface{}
	if data, ok := resp["data"].(map[string]interface{}); ok {
		if list, ok := data["list"].([]interface{}); ok {
			for _, item := range list {
				if m, ok := item.(map[string]interface{}); ok {
					app := map[string]interface{}{
						"app_id":    m["id"],
						"app_name":  m["app_name"],
						"app_label": m["app_label"],
						"app_type":  m["app_type"],
					}
					apps = append(apps, app)
				}
			}
		}
	}

	return apps, nil
}

// CreateCustomApp creates a custom extraction app
func (c *Client) CreateCustomApp(appName string, extractFields []map[string]interface{}, parseMode string, enableLongDoc *bool, longDocConfig []map[string]interface{}, appLabel []string) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"app_name":       appName,
		"extract_fields": extractFields,
		"parse_mode":     parseMode,
	}

	if appLabel != nil {
		data["app_label"] = appLabel
	}

	if enableLongDoc != nil {
		data["enable_long_doc"] = *enableLongDoc
		if *enableLongDoc && longDocConfig != nil {
			data["long_doc_config"] = longDocConfig
		}
	}

	return c.request("POST", fmt.Sprintf("/open/agentic_doc_processor/%s/v1/app-manage/create", c.tenantName), data)
}

// UpdateCustomApp updates a custom extraction app
func (c *Client) UpdateCustomApp(appID string, extractFields []map[string]interface{}, parseMode string, enableLongDoc *bool, appName *string, appLabel []string, longDocConfig []map[string]interface{}) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"app_id":         appID,
		"extract_fields": extractFields,
		"parse_mode":     parseMode,
	}

	if appName != nil {
		data["app_name"] = *appName
	}
	if appLabel != nil {
		data["app_label"] = appLabel
	}
	if enableLongDoc != nil {
		data["enable_long_doc"] = *enableLongDoc
		if *enableLongDoc && longDocConfig != nil {
			data["long_doc_config"] = longDocConfig
		}
	}

	return c.request("POST", fmt.Sprintf("/open/agentic_doc_processor/%s/v1/app-manage/update", c.tenantName), data)
}

// GetCustomAppConfig gets custom app configuration
func (c *Client) GetCustomAppConfig(appID string, configVersion *string) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"app_id": appID,
	}
	if configVersion != nil {
		data["config_version"] = *configVersion
	}

	return c.request("POST", fmt.Sprintf("/open/agentic_doc_processor/%s/v1/app-manage/config", c.tenantName), data)
}

// DeleteCustomApp deletes a custom extraction app
func (c *Client) DeleteCustomApp(appID string) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"app_id": appID,
	}
	return c.request("POST", fmt.Sprintf("/open/agentic_doc_processor/%s/v1/app-manage/delete", c.tenantName), data)
}

// DeleteCustomAppVersion deletes a specific config version
func (c *Client) DeleteCustomAppVersion(appID, configVersion string) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"app_id":         appID,
		"config_version": configVersion,
	}
	return c.request("POST", fmt.Sprintf("/open/agentic_doc_processor/%s/v1/app-manage/version/delete", c.tenantName), data)
}

// AIGenerateFields generates extraction field recommendations
func (c *Client) AIGenerateFields(appID string, fileURL, fileLocal, fileBase64 string) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"app_id": appID,
	}

	if fileBase64 != "" {
		data["file_base64"] = fileBase64
	} else if fileURL != "" {
		data["file_url"] = fileURL
	} else if fileLocal != "" {
		b64, err := c.encodeFileToBase64(fileLocal)
		if err != nil {
			return nil, fmt.Errorf("failed to encode file: %w", err)
		}
		data["file_base64"] = b64
	}

	return c.request("POST", fmt.Sprintf("/open/agentic_doc_processor/%s/v1/app-manage/ai-recommend", c.tenantName), data)
}

// GetAccountInfo gets user account info including credit and concurrency
func (c *Client) GetAccountInfo() (map[string]interface{}, error) {
	return c.request("GET", fmt.Sprintf("/open/agentic_doc_processor/%s/v1/user/payment", c.tenantName), nil)
}

// GetUserConcurrencyLimit returns the max concurrency allowed for the current user
func (c *Client) GetUserConcurrencyLimit() (int, error) {
	resp, err := c.GetAccountInfo()
	if err != nil {
		return 1, err
	}
	if concurrency, ok := resp["concurrency"].(float64); ok {
		return int(concurrency), nil
	}
	return 1, nil
}

// HealthCheck checks API health
func (c *Client) HealthCheck() bool {
	_, err := c.request("GET", "/health", nil)
	return err == nil
}
