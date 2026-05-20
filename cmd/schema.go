package cmd

import (
	"github.com/spf13/cobra"
	"github.com/laiye-ai/adp-cli/internal/i18n"
)

// schemaCmd represents the schema command
var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: i18n.T("schema_description"),
	Long:  i18n.T("schema_description"),
	Run: func(cmd *cobra.Command, args []string) {
		schema := GetFullSchema()
		formatterOut.PrintJSON(schema)
	},
}

func init() {
	rootCmd.AddCommand(schemaCmd)
}

// GetFullSchema returns the complete command schema for Agent introspection
func GetFullSchema() map[string]interface{} {
	schema := map[string]interface{}{
		"version": "1.0.0",
		"commands": map[string]interface{}{
			"config": map[string]interface{}{
				"description": i18n.T("config_description"),
				"subcommands": map[string]interface{}{
					"set": map[string]interface{}{
						"description": i18n.T("config_set_title"),
						"options": []map[string]interface{}{
							{"name": "api-key", "type": "string", "required": false, "description": i18n.T("option_api_key")},
							{"name": "api-base-url", "type": "string", "required": false, "description": i18n.T("option_api_base_url")},
						},
					},
					"get": map[string]interface{}{
						"description": i18n.T("config_get_title"),
						"options":      []map[string]interface{}{},
					},
					"clear": map[string]interface{}{
						"description": i18n.T("config_clear_title"),
						"options": []map[string]interface{}{
							{"name": "force", "type": "boolean", "short": "y", "description": i18n.T("option_force_clear")},
						},
					},
				},
			},
			"app-id": map[string]interface{}{
				"description": i18n.T("app_id_description"),
				"subcommands": map[string]interface{}{
					"list": map[string]interface{}{
						"description": i18n.T("app_id_list_title"),
						"options": []map[string]interface{}{
							{"name": "app-label", "type": "string", "required": false, "description": i18n.T("app_id_list_app_label")},
							{"name": "app-type", "type": "integer", "required": false, "description": i18n.T("app_id_list_app_type")},
							{"name": "limit", "type": "integer", "required": false, "default": 120, "description": i18n.T("app_id_list_limit")},
						},
					},
					"cache": map[string]interface{}{
						"description": i18n.T("app_id_list_cache_title"),
						"options":     []map[string]interface{}{},
					},
				},
			},
			"credit": map[string]interface{}{
				"description": i18n.T("credit_description"),
				"options": []map[string]interface{}{
					{"name": "api-key", "type": "string", "required": false, "description": i18n.T("credit_api_key")},
				},
			},
			"parse": map[string]interface{}{
				"description": i18n.T("parse_description"),
				"subcommands": map[string]interface{}{
					"local": map[string]interface{}{
						"description":  i18n.T("parse_local_title"),
						"arguments":    []string{"file-path"},
						"required_args": 1,
						"options": []map[string]interface{}{
							{"name": "app-id", "type": "string", "required": true, "description": i18n.T("option_app_id_parse")},
							{"name": "async", "type": "boolean", "default": false, "description": i18n.T("option_async")},
							{"name": "no-wait", "type": "boolean", "default": false, "description": i18n.T("option_no_wait")},
							{"name": "export", "type": "string", "required": false, "description": i18n.T("option_export")},
							{"name": "timeout", "type": "integer", "default": 900, "description": i18n.T("option_timeout")},
							{"name": "retry", "type": "integer", "default": 0, "description": i18n.T("option_retry")},
						},
					},
					"url": map[string]interface{}{
						"description":   i18n.T("parse_url_title"),
						"arguments":     []string{"url"},
						"required_args": 1,
						"options": []map[string]interface{}{
							{"name": "app-id", "type": "string", "required": true, "description": i18n.T("option_app_id_parse")},
							{"name": "async", "type": "boolean", "default": false, "description": i18n.T("option_async")},
							{"name": "no-wait", "type": "boolean", "default": false, "description": i18n.T("option_no_wait")},
							{"name": "export", "type": "string", "required": false, "description": i18n.T("option_export")},
							{"name": "timeout", "type": "integer", "default": 900, "description": i18n.T("option_timeout")},
							{"name": "retry", "type": "integer", "default": 0, "description": i18n.T("option_retry")},
						},
					},
					"base64": map[string]interface{}{
						"description":   i18n.T("parse_base64_title"),
						"arguments":     []string{"base64-strings"},
						"required_args": 1,
						"options": []map[string]interface{}{
							{"name": "app-id", "type": "string", "required": true, "description": i18n.T("option_app_id_parse")},
							{"name": "async", "type": "boolean", "default": false, "description": i18n.T("option_async")},
							{"name": "no-wait", "type": "boolean", "default": false, "description": i18n.T("option_no_wait")},
							{"name": "export", "type": "string", "required": false, "description": i18n.T("option_export")},
							{"name": "timeout", "type": "integer", "default": 900, "description": i18n.T("option_timeout")},
							{"name": "file-name", "type": "string", "default": "document", "description": i18n.T("option_file_name")},
							{"name": "retry", "type": "integer", "default": 0, "description": i18n.T("option_retry")},
						},
					},
					"query": map[string]interface{}{
						"description":   i18n.T("parse_query_title"),
						"arguments":     []string{"task-ids"},
						"required_args": 0,
						"options": []map[string]interface{}{
							{"name": "watch", "type": "boolean", "default": false, "description": i18n.T("option_watch")},
							{"name": "file", "type": "string", "required": false, "description": i18n.T("option_task_file")},
							{"name": "export", "type": "string", "required": false, "description": i18n.T("option_export")},
							{"name": "timeout", "type": "integer", "default": 900, "description": i18n.T("option_watch_timeout")},
						},
					},
				},
			},
			"extract": map[string]interface{}{
				"description": i18n.T("extract_description"),
				"subcommands": map[string]interface{}{
					"local": map[string]interface{}{
						"description":   i18n.T("extract_local_title"),
						"arguments":     []string{"file-path"},
						"required_args": 1,
						"options": []map[string]interface{}{
							{"name": "app-id", "type": "string", "required": true, "description": i18n.T("option_app_id_extract")},
							{"name": "async", "type": "boolean", "default": false, "description": i18n.T("option_async")},
							{"name": "no-wait", "type": "boolean", "default": false, "description": i18n.T("option_no_wait")},
							{"name": "export", "type": "string", "required": false, "description": i18n.T("option_export")},
							{"name": "timeout", "type": "integer", "default": 900, "description": i18n.T("option_timeout")},
							{"name": "retry", "type": "integer", "default": 0, "description": i18n.T("option_retry")},
						},
					},
					"url": map[string]interface{}{
						"description":   i18n.T("extract_url_title"),
						"arguments":     []string{"url"},
						"required_args": 1,
						"options": []map[string]interface{}{
							{"name": "app-id", "type": "string", "required": true, "description": i18n.T("option_app_id_extract")},
							{"name": "async", "type": "boolean", "default": false, "description": i18n.T("option_async")},
							{"name": "no-wait", "type": "boolean", "default": false, "description": i18n.T("option_no_wait")},
							{"name": "export", "type": "string", "required": false, "description": i18n.T("option_export")},
							{"name": "timeout", "type": "integer", "default": 900, "description": i18n.T("option_timeout")},
							{"name": "retry", "type": "integer", "default": 0, "description": i18n.T("option_retry")},
						},
					},
					"base64": map[string]interface{}{
						"description":   i18n.T("extract_base64_title"),
						"arguments":     []string{"base64-strings"},
						"required_args": 1,
						"options": []map[string]interface{}{
							{"name": "app-id", "type": "string", "required": true, "description": i18n.T("option_app_id_extract")},
							{"name": "async", "type": "boolean", "default": false, "description": i18n.T("option_async")},
							{"name": "no-wait", "type": "boolean", "default": false, "description": i18n.T("option_no_wait")},
							{"name": "export", "type": "string", "required": false, "description": i18n.T("option_export")},
							{"name": "timeout", "type": "integer", "default": 900, "description": i18n.T("option_timeout")},
							{"name": "file-name", "type": "string", "default": "document", "description": i18n.T("option_file_name")},
							{"name": "retry", "type": "integer", "default": 0, "description": i18n.T("option_retry")},
						},
					},
					"query": map[string]interface{}{
						"description":   i18n.T("extract_query_title"),
						"arguments":     []string{"task-ids"},
						"required_args": 0,
						"options": []map[string]interface{}{
							{"name": "watch", "type": "boolean", "default": false, "description": i18n.T("option_watch")},
							{"name": "file", "type": "string", "required": false, "description": i18n.T("option_task_file")},
							{"name": "export", "type": "string", "required": false, "description": i18n.T("option_export")},
							{"name": "timeout", "type": "integer", "default": 900, "description": i18n.T("option_watch_timeout")},
						},
					},
				},
			},
			"human-review": map[string]interface{}{
				"description": i18n.T("human_review_description"),
				"subcommands": map[string]interface{}{
					"rule-create": map[string]interface{}{
						"description": i18n.T("human_review_rule_create_title"),
						"options": []map[string]interface{}{
							{"name": "app-id", "type": "string", "required": true, "description": i18n.T("human_review_option_app_id")},
							{"name": "rule-name", "type": "string", "required": true, "description": i18n.T("human_review_option_rule_name")},
							{"name": "rule-status", "type": "string", "default": "true", "enum": []string{"true", "false"}, "description": i18n.T("human_review_option_rule_status")},
							{"name": "rule", "type": "string", "required": true, "description": i18n.T("human_review_option_rule")},
							{"name": "rule-logic", "type": "integer", "default": 1, "description": i18n.T("human_review_option_rule_logic")},
						},
					},
					"get-config": map[string]interface{}{
						"description": i18n.T("human_review_get_config_title"),
						"options": []map[string]interface{}{
							{"name": "app-id", "type": "string", "required": true, "description": i18n.T("human_review_option_app_id")},
						},
					},
					"rule-update": map[string]interface{}{
						"description": i18n.T("human_review_rule_update_title"),
						"options": []map[string]interface{}{
							{"name": "app-id", "type": "string", "required": true, "description": i18n.T("human_review_option_app_id")},
							{"name": "rule-name", "type": "string", "required": true, "description": i18n.T("human_review_option_rule_name")},
							{"name": "rule-status", "type": "string", "default": "true", "enum": []string{"true", "false"}, "description": i18n.T("human_review_option_rule_status")},
							{"name": "rule", "type": "string", "required": true, "description": i18n.T("human_review_option_rule")},
							{"name": "rule-logic", "type": "integer", "default": 1, "description": i18n.T("human_review_option_rule_logic")},
						},
					},
					"rule-delete": map[string]interface{}{
						"description": i18n.T("human_review_rule_delete_title"),
						"options": []map[string]interface{}{
							{"name": "app-id", "type": "string", "required": true, "description": i18n.T("human_review_option_app_id")},
						},
					},
					"rule-ai-generate": map[string]interface{}{
						"description": i18n.T("human_review_rule_ai_generate_title"),
						"options": []map[string]interface{}{
							{"name": "app-id", "type": "string", "required": true, "description": i18n.T("human_review_option_app_id")},
							{"name": "fields", "type": "string", "required": false, "description": i18n.T("human_review_option_fields")},
							{"name": "timeout", "type": "integer", "default": 900, "description": i18n.T("option_timeout")},
						},
					},
					"task-create": map[string]interface{}{
						"description": i18n.T("human_review_task_create_title"),
						"options": []map[string]interface{}{
							{"name": "app-id", "type": "string", "required": true, "description": i18n.T("human_review_option_app_id")},
							{"name": "local", "type": "string", "required": false, "description": i18n.T("human_review_option_local")},
							{"name": "url", "type": "string", "required": false, "description": i18n.T("human_review_option_url")},
							{"name": "async", "type": "boolean", "default": false, "description": i18n.T("option_async")},
							{"name": "no-wait", "type": "boolean", "default": false, "description": i18n.T("option_no_wait")},
							{"name": "export", "type": "string", "required": false, "description": i18n.T("option_export")},
							{"name": "timeout", "type": "integer", "default": 900, "description": i18n.T("option_timeout")},
							{"name": "concurrency", "type": "integer", "default": 1, "description": i18n.T("option_concurrency")},
							{"name": "retry", "type": "integer", "default": 0, "description": i18n.T("option_retry")},
						},
					},
					"task-query": map[string]interface{}{
						"description":   i18n.T("human_review_task_query_title"),
						"arguments":     []string{"task-ids"},
						"required_args": 0,
						"options": []map[string]interface{}{
							{"name": "watch", "type": "boolean", "default": false, "description": i18n.T("option_watch")},
							{"name": "file", "type": "string", "required": false, "description": i18n.T("option_task_file")},
							{"name": "export", "type": "string", "required": false, "description": i18n.T("option_export")},
							{"name": "timeout", "type": "integer", "default": 900, "description": i18n.T("option_watch_timeout")},
							{"name": "concurrency", "type": "integer", "default": 1, "description": i18n.T("option_concurrency")},
						},
					},
					"result-update": map[string]interface{}{
						"description": i18n.T("human_review_result_update_title"),
						"options": []map[string]interface{}{
							{"name": "file-task-id", "type": "string", "required": true, "description": i18n.T("human_review_option_file_task_id")},
							{"name": "collaboration-result", "type": "string", "required": true, "description": i18n.T("human_review_option_collaboration_result")},
						},
					},
				},
			},
			"webhook": map[string]interface{}{
				"description": i18n.T("webhook_description"),
				"subcommands": map[string]interface{}{
					"create": map[string]interface{}{
						"description": i18n.T("webhook_create_title"),
						"options": []map[string]interface{}{
							{"name": "webhook-url", "type": "string", "required": true, "description": i18n.T("webhook_option_webhook_url")},
							{"name": "event-types", "type": "string", "required": true, "description": i18n.T("webhook_option_event_types")},
							{"name": "app-id", "type": "string", "required": false, "description": i18n.T("webhook_option_app_id")},
						},
					},
					"get-config": map[string]interface{}{
						"description": i18n.T("webhook_get_config_title"),
						"options": []map[string]interface{}{
							{"name": "app-id", "type": "string", "required": false, "description": i18n.T("webhook_option_app_id")},
						},
					},
					"update": map[string]interface{}{
						"description": i18n.T("webhook_update_title"),
						"options": []map[string]interface{}{
							{"name": "webhook-id", "type": "string", "required": true, "description": i18n.T("webhook_option_webhook_id")},
							{"name": "webhook-url", "type": "string", "required": true, "description": i18n.T("webhook_option_webhook_url")},
							{"name": "event-types", "type": "string", "required": true, "description": i18n.T("webhook_option_event_types")},
							{"name": "app-id", "type": "string", "required": false, "description": i18n.T("webhook_option_app_id")},
						},
					},
					"delete": map[string]interface{}{
						"description": i18n.T("webhook_delete_title"),
						"options": []map[string]interface{}{
							{"name": "webhook-id", "type": "string", "required": true, "description": i18n.T("webhook_option_webhook_id")},
						},
					},
					"log": map[string]interface{}{
						"description": i18n.T("webhook_log_title"),
						"options": []map[string]interface{}{
							{"name": "webhook-id", "type": "string", "required": false, "description": i18n.T("webhook_option_webhook_id_optional")},
							{"name": "start-time", "type": "string", "required": false, "description": i18n.T("webhook_option_start_time")},
							{"name": "end-time", "type": "string", "required": false, "description": i18n.T("webhook_option_end_time")},
						},
					},
				},
			},
			"custom-app": map[string]interface{}{
				"description": i18n.T("custom_app_description"),
				"subcommands": map[string]interface{}{
					"create": map[string]interface{}{
						"description": i18n.T("custom_app_create_title"),
						"options": []map[string]interface{}{
							{"name": "api-key", "type": "string", "required": false, "description": i18n.T("custom_app_create_api_key")},
							{"name": "app-name", "type": "string", "required": true, "description": i18n.T("custom_app_create_app_name")},
							{"name": "app-label", "type": "string", "required": false, "description": i18n.T("custom_app_create_app_label")},
							{"name": "extract-fields", "type": "string", "required": true, "description": i18n.T("custom_app_create_extract_fields")},
							{"name": "parse-mode", "type": "string", "required": true, "enum": []string{"advance", "standard", "agentic"}, "description": i18n.T("custom_app_create_parse_mode")},
							{"name": "enable-long-doc", "type": "string", "required": false, "description": i18n.T("custom_app_create_enable_long_doc")},
							{"name": "long-doc-config", "type": "string", "required": false, "description": i18n.T("custom_app_create_long_doc_config")},
						},
					},
					"update": map[string]interface{}{
						"description": i18n.T("custom_app_update_title"),
						"options": []map[string]interface{}{
							{"name": "api-key", "type": "string", "required": false, "description": i18n.T("custom_app_create_api_key")},
							{"name": "app-id", "type": "string", "required": true, "description": i18n.T("custom_app_update_app_id")},
							{"name": "app-name", "type": "string", "required": false, "description": i18n.T("custom_app_update_app_name")},
							{"name": "app-label", "type": "string", "required": false, "description": i18n.T("custom_app_update_app_label")},
							{"name": "extract-fields", "type": "string", "required": true, "description": i18n.T("custom_app_update_extract_fields")},
							{"name": "parse-mode", "type": "string", "required": true, "enum": []string{"advance", "standard", "agentic"}, "description": i18n.T("custom_app_update_parse_mode")},
							{"name": "enable-long-doc", "type": "string", "required": false, "description": i18n.T("custom_app_update_enable_long_doc")},
							{"name": "long-doc-config", "type": "string", "required": false, "description": i18n.T("custom_app_update_long_doc_config")},
						},
					},
					"get-config": map[string]interface{}{
						"description": i18n.T("custom_app_get_config_title"),
						"options": []map[string]interface{}{
							{"name": "api-key", "type": "string", "required": false, "description": i18n.T("custom_app_create_api_key")},
							{"name": "app-id", "type": "string", "required": true, "description": i18n.T("custom_app_get_config_app_id")},
							{"name": "config-version", "type": "string", "required": false, "description": i18n.T("custom_app_get_config_config_version")},
						},
					},
					"delete": map[string]interface{}{
						"description": i18n.T("custom_app_delete_title"),
						"options": []map[string]interface{}{
							{"name": "api-key", "type": "string", "required": false, "description": i18n.T("custom_app_create_api_key")},
							{"name": "app-id", "type": "string", "required": true, "description": i18n.T("custom_app_delete_app_id")},
						},
					},
					"delete-version": map[string]interface{}{
						"description": i18n.T("custom_app_delete_version_title"),
						"options": []map[string]interface{}{
							{"name": "api-key", "type": "string", "required": false, "description": i18n.T("custom_app_create_api_key")},
							{"name": "app-id", "type": "string", "required": true, "description": i18n.T("custom_app_delete_version_app_id")},
							{"name": "config-version", "type": "string", "required": true, "description": i18n.T("custom_app_delete_version_config_version")},
						},
					},
					"ai-generate": map[string]interface{}{
						"description": i18n.T("custom_app_ai_generate_title"),
						"options": []map[string]interface{}{
							{"name": "api-key", "type": "string", "required": false, "description": i18n.T("custom_app_create_api_key")},
							{"name": "app-id", "type": "string", "required": true, "description": i18n.T("custom_app_ai_generate_app_id")},
							{"name": "file-url", "type": "string", "required": false, "description": i18n.T("custom_app_ai_generate_file_url")},
							{"name": "file-local", "type": "string", "required": false, "description": i18n.T("custom_app_ai_generate_file_local")},
							{"name": "base64", "type": "string", "required": false, "description": i18n.T("custom_app_ai_generate_file_base64")},
						},
					},
				},
			},
		},
		"global_options": []map[string]interface{}{
			{"name": "lang", "type": "string", "enum": []string{"en", "zh"}, "description": "Set language (en or zh)"},
			{"name": "json", "type": "boolean", "description": "Output in JSON format"},
			{"name": "quiet", "type": "boolean", "description": "Suppress all output except errors"},
			{"name": "source", "type": "string", "description": "Caller identity (e.g. claude, cursor, chatgpt)"},
		},
	}

	return schema
}
