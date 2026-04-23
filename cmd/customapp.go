package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/laiye-ai/adp-cli/internal/api"
	"github.com/laiye-ai/adp-cli/internal/config"
	"github.com/laiye-ai/adp-cli/internal/errors"
	"github.com/laiye-ai/adp-cli/internal/i18n"
)

// customAppCmd represents the custom-app command group
var customAppCmd = &cobra.Command{
	Use:   "custom-app",
	Short: i18n.T("custom_app_description"),
	Long:  i18n.T("custom_app_description"),
}

// createCustomAppCmd represents the custom-app create command
var createCustomAppCmd = &cobra.Command{
	Use:   "create",
	Short: i18n.T("custom_app_create_title"),
	Long:  i18n.T("custom_app_create_title"),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, _ := cmd.Flags().GetString("api-key")
		appName, _ := cmd.Flags().GetString("app-name")
		appLabelStr, _ := cmd.Flags().GetString("app-label")
		extractFieldsStr, _ := cmd.Flags().GetString("extract-fields")
		parseMode, _ := cmd.Flags().GetString("parse-mode")
		enableLongDocStr, _ := cmd.Flags().GetString("enable-long-doc")
		longDocConfigStr, _ := cmd.Flags().GetString("long-doc-config")

		// Validate parse-mode enum
		if err := errors.ValidateEnum(parseMode, []string{"advance", "standard", "agentic"}, "parse-mode"); err != nil {
			formatterOut.ExitWithError(err)
		}

		cfg, err := loadConfigWithOverride(apiKey)
		if err != nil {
			cliErr := errors.ClassifyException(err, "custom-app create")
			formatterOut.ExitWithError(cliErr)
		}

		client, err := api.NewClient(cfg)
		if err != nil {
			cliErr := errors.ClassifyException(err, "custom-app create")
			formatterOut.ExitWithError(cliErr)
		}

		// Parse extract fields
		extractFields, err := parseJSONParam(extractFieldsStr)
		if err != nil {
			cliErr := errors.NewCLIError(
				fmt.Sprintf("Invalid extract-fields: %v", err),
				errors.ErrorTypeParam,
				errors.ExitParameterError,
				false,
				"Check the extract-fields JSON format.",
				nil,
			)
			formatterOut.ExitWithError(cliErr)
		}

		// Parse app label
		var appLabel []string
		if appLabelStr != "" {
			appLabel, err = parseStringList(appLabelStr)
			if err != nil {
				cliErr := errors.NewCLIError(
					fmt.Sprintf("Invalid app-label: %v", err),
					errors.ErrorTypeParam,
					errors.ExitParameterError,
					false,
					"Check the app-label format.",
					nil,
				)
				formatterOut.ExitWithError(cliErr)
			}
		}

		// Parse enable long doc
		var enableLongDoc *bool
		if enableLongDocStr != "" {
			val := enableLongDocStr == "true" || enableLongDocStr == "1"
			enableLongDoc = &val
		}

		// Parse long doc config
		var longDocConfig []map[string]interface{}
		if longDocConfigStr != "" {
			longDocConfig, err = parseJSONParam(longDocConfigStr)
			if err != nil {
				cliErr := errors.NewCLIError(
					fmt.Sprintf("Invalid long-doc-config: %v", err),
					errors.ErrorTypeParam,
					errors.ExitParameterError,
					false,
					"Check the long-doc-config JSON format.",
					nil,
				)
				formatterOut.ExitWithError(cliErr)
			}
		}

		result, err := client.CreateCustomApp(appName, extractFields, parseMode, enableLongDoc, longDocConfig, appLabel)
		if err != nil {
			cliErr := errors.ClassifyException(err, "custom-app create")
			formatterOut.ExitWithError(cliErr)
		}

		formatterOut.PrintSuccess(i18n.T("app_created"))
		formatterOut.PrintJSON(result)
	},
}

// updateCustomAppCmd represents the custom-app update command
var updateCustomAppCmd = &cobra.Command{
	Use:   "update",
	Short: i18n.T("custom_app_update_title"),
	Long:  i18n.T("custom_app_update_title"),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, _ := cmd.Flags().GetString("api-key")
		appID, _ := cmd.Flags().GetString("app-id")
		appName, _ := cmd.Flags().GetString("app-name")
		appLabelStr, _ := cmd.Flags().GetString("app-label")
		extractFieldsStr, _ := cmd.Flags().GetString("extract-fields")
		parseMode, _ := cmd.Flags().GetString("parse-mode")
		enableLongDocStr, _ := cmd.Flags().GetString("enable-long-doc")
		longDocConfigStr, _ := cmd.Flags().GetString("long-doc-config")

		// Validate parse-mode enum
		if err := errors.ValidateEnum(parseMode, []string{"advance", "standard", "agentic"}, "parse-mode"); err != nil {
			formatterOut.ExitWithError(err)
		}

		cfg, err := loadConfigWithOverride(apiKey)
		if err != nil {
			cliErr := errors.ClassifyException(err, "custom-app update")
			formatterOut.ExitWithError(cliErr)
		}

		client, err := api.NewClient(cfg)
		if err != nil {
			cliErr := errors.ClassifyException(err, "custom-app update")
			formatterOut.ExitWithError(cliErr)
		}

		// Parse extract fields
		extractFields, err := parseJSONParam(extractFieldsStr)
		if err != nil {
			cliErr := errors.NewCLIError(
				fmt.Sprintf("Invalid extract-fields: %v", err),
				errors.ErrorTypeParam,
				errors.ExitParameterError,
				false,
				"Check the extract-fields JSON format.",
				nil,
			)
			formatterOut.ExitWithError(cliErr)
		}

		// Parse app label
		var appLabel []string
		if appLabelStr != "" {
			appLabel, err = parseStringList(appLabelStr)
			if err != nil {
				cliErr := errors.NewCLIError(
					fmt.Sprintf("Invalid app-label: %v", err),
					errors.ErrorTypeParam,
					errors.ExitParameterError,
					false,
					"Check the app-label format.",
					nil,
				)
				formatterOut.ExitWithError(cliErr)
			}
		}

		// Parse enable long doc (tri-state: nil = not provided, *bool otherwise)
		var enableLongDoc *bool
		if enableLongDocStr != "" {
			val := enableLongDocStr == "true" || enableLongDocStr == "1"
			enableLongDoc = &val
		}

		// Parse long doc config
		var longDocConfig []map[string]interface{}
		if longDocConfigStr != "" {
			longDocConfig, err = parseJSONParam(longDocConfigStr)
			if err != nil {
				cliErr := errors.NewCLIError(
					fmt.Sprintf("Invalid long-doc-config: %v", err),
					errors.ErrorTypeParam,
					errors.ExitParameterError,
					false,
					"Check the long-doc-config JSON format.",
					nil,
				)
				formatterOut.ExitWithError(cliErr)
			}
		}

		var appNamePtr *string
		if appName != "" {
			appNamePtr = &appName
		}

		result, err := client.UpdateCustomApp(appID, extractFields, parseMode, enableLongDoc, appNamePtr, appLabel, longDocConfig)
		if err != nil {
			cliErr := errors.ClassifyException(err, "custom-app update")
			formatterOut.ExitWithError(cliErr)
		}

		formatterOut.PrintSuccess(i18n.T("app_updated"))
		formatterOut.PrintJSON(result)
	},
}

// getConfigCustomAppCmd represents the custom-app get-config command
var getConfigCustomAppCmd = &cobra.Command{
	Use:   "get-config",
	Short: i18n.T("custom_app_get_config_title"),
	Long:  i18n.T("custom_app_get_config_title"),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, _ := cmd.Flags().GetString("api-key")
		appID, _ := cmd.Flags().GetString("app-id")
		configVersion, _ := cmd.Flags().GetString("config-version")

		cfg, err := loadConfigWithOverride(apiKey)
		if err != nil {
			cliErr := errors.ClassifyException(err, "custom-app get-config")
			formatterOut.ExitWithError(cliErr)
		}

		client, err := api.NewClient(cfg)
		if err != nil {
			cliErr := errors.ClassifyException(err, "custom-app get-config")
			formatterOut.ExitWithError(cliErr)
		}

		var configVersionPtr *string
		if configVersion != "" {
			configVersionPtr = &configVersion
		}

		result, err := client.GetCustomAppConfig(appID, configVersionPtr)
		if err != nil {
			cliErr := errors.ClassifyException(err, "custom-app get-config")
			formatterOut.ExitWithError(cliErr)
		}

		formatterOut.PrintJSON(result)
	},
}

// deleteCustomAppCmd represents the custom-app delete command
var deleteCustomAppCmd = &cobra.Command{
	Use:   "delete",
	Short: i18n.T("custom_app_delete_title"),
	Long:  i18n.T("custom_app_delete_title"),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, _ := cmd.Flags().GetString("api-key")
		appID, _ := cmd.Flags().GetString("app-id")

		cfg, err := loadConfigWithOverride(apiKey)
		if err != nil {
			cliErr := errors.ClassifyException(err, "custom-app delete")
			formatterOut.ExitWithError(cliErr)
		}

		client, err := api.NewClient(cfg)
		if err != nil {
			cliErr := errors.ClassifyException(err, "custom-app delete")
			formatterOut.ExitWithError(cliErr)
		}

		result, err := client.DeleteCustomApp(appID)
		if err != nil {
			cliErr := errors.ClassifyException(err, "custom-app delete")
			// Idempotent: treat "not found" as success
			if errors.IsNotFoundError(err) {
				formatterOut.PrintWarning(fmt.Sprintf(i18n.T("not_found_may_deleted_app"), appID))
				formatterOut.PrintSuccess(i18n.T("app_delete_request_ok"))
				return
			}
			formatterOut.ExitWithError(cliErr)
		}

		// Check result for not-found in response body
		if isResultNotFound(result) {
			formatterOut.PrintWarning(fmt.Sprintf(i18n.T("not_found_may_deleted_app"), appID))
			formatterOut.PrintSuccess(i18n.T("app_delete_request_ok"))
			return
		}

		// Check result
		if data, ok := result["data"].(map[string]interface{}); ok {
			if success, ok := data["success"].(bool); ok && success {
				formatterOut.PrintSuccess(i18n.T("app_deleted"))
				return
			}
		}

		cliErr := errors.NewCLIError(
			"Failed to delete app",
			errors.ErrorTypeAPI,
			errors.ExitGeneralError,
			false,
			"Check the app ID and try again.",
			result,
		)
		formatterOut.ExitWithError(cliErr)
	},
}

// deleteVersionCustomAppCmd represents the custom-app delete-version command
var deleteVersionCustomAppCmd = &cobra.Command{
	Use:   "delete-version",
	Short: i18n.T("custom_app_delete_version_title"),
	Long:  i18n.T("custom_app_delete_version_title"),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, _ := cmd.Flags().GetString("api-key")
		appID, _ := cmd.Flags().GetString("app-id")
		configVersion, _ := cmd.Flags().GetString("config-version")

		cfg, err := loadConfigWithOverride(apiKey)
		if err != nil {
			cliErr := errors.ClassifyException(err, "custom-app delete-version")
			formatterOut.ExitWithError(cliErr)
		}

		client, err := api.NewClient(cfg)
		if err != nil {
			cliErr := errors.ClassifyException(err, "custom-app delete-version")
			formatterOut.ExitWithError(cliErr)
		}

		result, err := client.DeleteCustomAppVersion(appID, configVersion)
		if err != nil {
			cliErr := errors.ClassifyException(err, "custom-app delete-version")
			// Idempotent: treat "not found" as success
			if errors.IsNotFoundError(err) {
				formatterOut.PrintWarning(fmt.Sprintf(i18n.T("not_found_may_deleted_ver"), configVersion))
				formatterOut.PrintSuccess(i18n.T("version_delete_request_ok"))
				return
			}
			formatterOut.ExitWithError(cliErr)
		}

		// Check result for not-found in response body (API may return server_error with Chinese message)
		if isResultNotFound(result) {
			formatterOut.PrintWarning(fmt.Sprintf(i18n.T("not_found_may_deleted_ver"), configVersion))
			formatterOut.PrintSuccess(i18n.T("version_delete_request_ok"))
			return
		}

		if code, ok := result["code"].(string); ok && code == "success" {
			formatterOut.PrintSuccess(i18n.T("version_deleted"))
		} else {
			cliErr := errors.NewCLIError(
				"Failed to delete version",
				errors.ErrorTypeAPI,
				errors.ExitGeneralError,
				false,
				"Check the app ID and config version and try again.",
				result,
			)
			formatterOut.ExitWithError(cliErr)
		}
	},
}

// aiGenerateCustomAppCmd represents the custom-app ai-generate command
var aiGenerateCustomAppCmd = &cobra.Command{
	Use:   "ai-generate",
	Short: i18n.T("custom_app_ai_generate_title"),
	Long:  i18n.T("custom_app_ai_generate_title"),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, _ := cmd.Flags().GetString("api-key")
		appID, _ := cmd.Flags().GetString("app-id")
		fileURL, _ := cmd.Flags().GetString("file-url")
		fileLocal, _ := cmd.Flags().GetString("file-local")
		fileBase64, _ := cmd.Flags().GetString("base64")

		if fileURL == "" && fileLocal == "" && fileBase64 == "" {
			cliErr := errors.NewCLIError(
				"One of --file-url, --file-local, or --base64 must be provided",
				errors.ErrorTypeParam,
				errors.ExitParameterError,
				false,
				"Provide at least one of --file-url, --file-local, or --base64.",
				nil,
			)
			formatterOut.ExitWithError(cliErr)
		}

		cfg, err := loadConfigWithOverride(apiKey)
		if err != nil {
			cliErr := errors.ClassifyException(err, "custom-app ai-generate")
			formatterOut.ExitWithError(cliErr)
		}

		client, err := api.NewClient(cfg)
		if err != nil {
			cliErr := errors.ClassifyException(err, "custom-app ai-generate")
			formatterOut.ExitWithError(cliErr)
		}

		result, err := client.AIGenerateFields(appID, fileURL, fileLocal, fileBase64)
		if err != nil {
			cliErr := errors.ClassifyException(err, "custom-app ai-generate")
			formatterOut.ExitWithError(cliErr)
		}

		formatterOut.PrintJSON(result)
	},
}

func init() {
	rootCmd.AddCommand(customAppCmd)

	customAppCmd.AddCommand(createCustomAppCmd)
	customAppCmd.AddCommand(updateCustomAppCmd)
	customAppCmd.AddCommand(getConfigCustomAppCmd)
	customAppCmd.AddCommand(deleteCustomAppCmd)
	customAppCmd.AddCommand(deleteVersionCustomAppCmd)
	customAppCmd.AddCommand(aiGenerateCustomAppCmd)

	// create flags
	createCustomAppCmd.Flags().String("api-key", "", i18n.T("custom_app_create_api_key"))
	createCustomAppCmd.Flags().String("app-name", "", i18n.T("custom_app_create_app_name"))
	createCustomAppCmd.Flags().String("app-label", "", i18n.T("custom_app_create_app_label"))
	createCustomAppCmd.Flags().String("extract-fields", "", i18n.T("custom_app_create_extract_fields"))
	createCustomAppCmd.Flags().String("parse-mode", "", i18n.T("custom_app_create_parse_mode"))
	createCustomAppCmd.Flags().String("enable-long-doc", "", i18n.T("custom_app_create_enable_long_doc"))
	createCustomAppCmd.Flags().String("long-doc-config", "", i18n.T("custom_app_create_long_doc_config"))
	createCustomAppCmd.MarkFlagRequired("app-name")
	createCustomAppCmd.MarkFlagRequired("extract-fields")
	createCustomAppCmd.MarkFlagRequired("parse-mode")

	// update flags
	updateCustomAppCmd.Flags().String("api-key", "", i18n.T("custom_app_update_api_key"))
	updateCustomAppCmd.Flags().String("app-id", "", i18n.T("custom_app_update_app_id"))
	updateCustomAppCmd.Flags().String("app-name", "", i18n.T("custom_app_update_app_name"))
	updateCustomAppCmd.Flags().String("app-label", "", i18n.T("custom_app_update_app_label"))
	updateCustomAppCmd.Flags().String("extract-fields", "", i18n.T("custom_app_update_extract_fields"))
	updateCustomAppCmd.Flags().String("parse-mode", "", i18n.T("custom_app_update_parse_mode"))
	updateCustomAppCmd.Flags().String("enable-long-doc", "", i18n.T("custom_app_update_enable_long_doc"))
	updateCustomAppCmd.Flags().String("long-doc-config", "", i18n.T("custom_app_update_long_doc_config"))
	updateCustomAppCmd.MarkFlagRequired("app-id")
	updateCustomAppCmd.MarkFlagRequired("extract-fields")
	updateCustomAppCmd.MarkFlagRequired("parse-mode")

	// get-config flags
	getConfigCustomAppCmd.Flags().String("api-key", "", i18n.T("custom_app_get_config_api_key"))
	getConfigCustomAppCmd.Flags().String("app-id", "", i18n.T("custom_app_get_config_app_id"))
	getConfigCustomAppCmd.Flags().String("config-version", "", i18n.T("custom_app_get_config_config_version"))
	getConfigCustomAppCmd.MarkFlagRequired("app-id")

	// delete flags
	deleteCustomAppCmd.Flags().String("api-key", "", i18n.T("custom_app_delete_api_key"))
	deleteCustomAppCmd.Flags().String("app-id", "", i18n.T("custom_app_delete_app_id"))
	deleteCustomAppCmd.MarkFlagRequired("app-id")

	// delete-version flags
	deleteVersionCustomAppCmd.Flags().String("api-key", "", i18n.T("custom_app_delete_version_api_key"))
	deleteVersionCustomAppCmd.Flags().String("app-id", "", i18n.T("custom_app_delete_version_app_id"))
	deleteVersionCustomAppCmd.Flags().String("config-version", "", i18n.T("custom_app_delete_version_config_version"))
	deleteVersionCustomAppCmd.MarkFlagRequired("app-id")
	deleteVersionCustomAppCmd.MarkFlagRequired("config-version")

	// ai-generate flags
	aiGenerateCustomAppCmd.Flags().String("api-key", "", i18n.T("custom_app_ai_generate_api_key"))
	aiGenerateCustomAppCmd.Flags().String("app-id", "", i18n.T("custom_app_ai_generate_app_id"))
	aiGenerateCustomAppCmd.Flags().String("file-url", "", i18n.T("custom_app_ai_generate_file_url"))
	aiGenerateCustomAppCmd.Flags().String("file-local", "", i18n.T("custom_app_ai_generate_file_local"))
	aiGenerateCustomAppCmd.Flags().String("base64", "", i18n.T("custom_app_ai_generate_file_base64"))
	aiGenerateCustomAppCmd.MarkFlagRequired("app-id")
}

func loadConfigWithOverride(apiKey string) (*config.Config, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	if apiKey != "" {
		// Temporarily set API key
		if err := config.SetAPIKey(apiKey); err != nil {
			return nil, err
		}
		cfg, err = config.Load()
		if err != nil {
			return nil, err
		}
	}

	return cfg, nil
}

func parseJSONParam(value string) ([]map[string]interface{}, error) {
	// Strip surrounding single quotes (Windows CMD treats them as literal characters)
	if len(value) >= 2 && value[0] == '\'' && value[len(value)-1] == '\'' {
		value = value[1 : len(value)-1]
	}

	// Try to parse as JSON
	var result []map[string]interface{}
	if err := json.Unmarshal([]byte(value), &result); err != nil {
		// Try as JSON file path
		data, err := os.ReadFile(value)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON or file not found: %s", value)
		}
		if err := json.Unmarshal(data, &result); err != nil {
			return nil, fmt.Errorf("invalid JSON: %v", err)
		}
	}
	return result, nil
}

func parseStringList(value string) ([]string, error) {
	// Try JSON array first
	var result []string
	if err := json.Unmarshal([]byte(value), &result); err != nil {
		// Try comma-separated
		parts := strings.Split(value, ",")
		result = make([]string, 0, len(parts))
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				result = append(result, p)
			}
		}
	}
	return result, nil
}

// isResultNotFound checks if API response indicates a "not found" condition.
// Handles both English and Chinese messages from the API.
func isResultNotFound(result map[string]interface{}) bool {
	msg, _ := result["message"].(string)
	if msg == "" {
		return false
	}
	lower := strings.ToLower(msg)
	return strings.Contains(lower, "not found") ||
		strings.Contains(lower, "does not exist") ||
		strings.Contains(msg, "不存在")
}
