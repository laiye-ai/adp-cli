package i18n

import (
	"fmt"
	"os"
	"strings"
)

var (
	currentLang = detectLanguage()
)

// Translations holds all translations
var Translations = map[string]map[string]string{
	"en": {
		// Global options
		"option_json":              "Output in JSON format",
		"option_quiet":             "Suppress all output except errors",
		"option_lang":              "Set language (en or zh)",

		// Config command
		"config_description":       "Authentication configuration management.",
		"config_set_title":         "Set or update API Key and API Base URL.",
		"config_get_title":         "View current configuration.",
		"config_clear_title":       "Clear local configuration.",
		"config_set_examples":      "Examples:",
		"error_not_configured":     "Configuration incomplete. Please run 'adp config set' to configure API Key and API Base URL.",
		"error_api_key_or_url_required": "At least one of --api-key or --api-base-url must be provided",
		"api_key_configured":       "API Key configured successfully",
		"api_base_url_configured":  "API Base URL configured successfully",
		"config_cleared":           "Configuration cleared",
		"confirm_clear_config":     "Are you sure you want to clear all configuration?",
		"option_force_clear":       "Skip confirmation and force clear",
		"option_api_key":          "API Key for authentication (format: xxxxxxxxxxxx)",
		"option_api_base_url":     "API Base URL (default: https://adp.laiye.com)",

		// Parse command
		"parse_description":        "Document parsing.",
		"parse_local_title":        "Parse local files or folders.",
		"parse_url_title":          "Parse URL files or URL list files.",
		"parse_base64_title":       "Parse base64 encoded file content.",
		"parse_query_title":        "Query parse async task status.",
		"option_app_id_parse":      "Application ID for parsing",
		"option_async":             "Process asynchronously",
		"option_export":            "Export results to JSON file",
		"option_timeout":           "Timeout for sync mode (seconds)",
		"option_file_name":        "Display name with extension (used by server to detect file type; does not read from disk). Default: document",
		"option_retry":             "Number of retries for retryable errors (default: 0)",
		"option_watch":             "Watch task until completion",
		"option_watch_timeout":    "Timeout for watch mode (seconds)",
		"option_no_wait":          "Submit tasks only, do not wait for results (use with --async)",
		"option_task_file":        "Read task IDs from JSON file (output of --no-wait)",

		// Extract command
		"extract_description":      "Document extraction.",
		"extract_local_title":      "Extract local files or folders.",
		"extract_url_title":        "Extract URL files or URL list files.",
		"extract_base64_title":     "Extract base64 encoded file content.",
		"extract_query_title":      "Query extract async task status.",
		"option_app_id_extract":    "Application ID for extraction",

		// App ID command
		"app_id_description":       "Application management.",
		"app_id_list_title":        "List all available application IDs.",
		"app_id_list_cache_title":  "List cached application IDs (fast).",
		"app_id_list_app_label":    "Filter applications by label (optional)",
		"app_id_list_app_type":     "Filter applications by type: 0=system preset, 1=custom app; omit to list all",
		"app_id_list_limit":        "Limit the number of returned applications (default: 120)",
		"no_applications_found":   "No applications found",

		// Custom App command
		"custom_app_description":                 "Custom extraction application management.",
		"custom_app_create_title":                "Create a custom extraction application.",
		"custom_app_create_api_key":              "API Key for authentication (optional if already configured)",
		"custom_app_create_app_name":             "Application name (max 50 characters)",
		"custom_app_create_app_label":            "Application labels (max 5 labels, optional, format: JSON array or comma-separated string)",
		"custom_app_create_extract_fields":       "Field configuration in JSON format or JSON file path (required). Fields must include: field_name, field_type, field_prompt. field_type: string=text, date=date, table=table. When field_type is table, field_prompt can be empty and must have sub_fields. Sub_fields must include: field_name, field_type, field_prompt. Sub_fields field_type cannot be table, only string or date. See documentation for examples.",
		"custom_app_create_parse_mode":           "Parse mode: standard=standard parsing; advance=advanced parsing; agentic=agentic parsing",
		"custom_app_create_enable_long_doc":      "Enable long document support (true/false, optional)",
		"custom_app_create_long_doc_config":      "Long document configuration in JSON format or JSON file path (optional, only valid when enable-long-doc=true)",
		"custom_app_update_title":               "Update a custom extraction application.",
		"custom_app_update_api_key":             "API Key for authentication (optional if already configured)",
		"custom_app_update_app_id":               "Application ID (required)",
		"custom_app_update_app_name":             "Application name (max 50 characters)",
		"custom_app_update_app_label":            "Application labels (max 5 labels, optional, format: JSON array or comma-separated string)",
		"custom_app_update_extract_fields":       "Field configuration in JSON format or JSON file path (required). Fields must include: field_name, field_type, field_prompt. field_type: string=text, date=date, table=table. When field_type is table, field_prompt can be empty and must have sub_fields. Sub_fields must include: field_name, field_type, field_prompt. Sub_fields field_type cannot be table, only string or date. See documentation for examples.",
		"custom_app_update_parse_mode":           "Parse mode: standard=standard parsing; advance=advanced parsing; agentic=agentic parsing",
		"custom_app_update_enable_long_doc":      "Enable long document support (true/false, optional; if omitted, server default behavior is preserved)",
		"custom_app_update_long_doc_config":     "Long document configuration in JSON format or JSON file path (optional, only valid when enable-long-doc=true)",
		"custom_app_get_config_title":           "Query custom app configuration.",
		"custom_app_get_config_api_key":         "API Key for authentication",
		"custom_app_get_config_app_id":          "Application ID",
		"custom_app_get_config_config_version":  "Configuration version (optional, default: latest)",
		"custom_app_delete_title":              "Delete custom app.",
		"custom_app_delete_api_key":             "API Key for authentication",
		"custom_app_delete_app_id":              "Application ID",
		"custom_app_delete_version_title":       "Delete specified config version.",
		"custom_app_delete_version_api_key":     "API Key for authentication",
		"custom_app_delete_version_app_id":      "Application ID",
		"custom_app_delete_version_config_version": "Configuration version to delete",
		"custom_app_ai_generate_title":          "AI generate extraction field recommendations.",
		"custom_app_ai_generate_api_key":        "API Key for authentication",
		"custom_app_ai_generate_app_id":         "Application ID",
		"custom_app_ai_generate_file_url":       "URL of sample document",
		"custom_app_ai_generate_file_local":      "Local of sample document",
		"custom_app_ai_generate_file_base64":     "Base64 encoded sample document content",
		"app_created":               "Custom app created successfully",
		"app_updated":               "Custom app updated successfully",
		"app_deleted":               "App deleted successfully",
		"app_delete_request_ok":     "Delete request succeeded",
		"version_deleted":          "Version deleted successfully",
		"version_delete_request_ok": "Delete request succeeded",
		"not_found_may_deleted_app": "App %s not found, it may have already been deleted.",
		"not_found_may_deleted_ver": "Version %s not found, it may have already been deleted.",

		// Credit command
		"credit_description":        "Query remaining credits.",
		"credit_api_key":           "API Key for authentication (optional if already configured)",
		"credit_info":             "Credit Information",
		"remaining_credits":       "Remaining Credits",
		"recharge_message":        "To recharge, please visit",

		// Help command
		"help_description":         "Display help information and available commands.",

		// Schema command
		"schema_description":        "Display command schema for Agent introspection.",

		// Version command
		"version_description":       "Print the version number.",

		// Common
		"error":                      "Error:",
		"warning":                    "Warning:",
		"no_valid_files":            "No valid files to process",
		"processing_files":          "Processing %d file(s)",
		"failed_to_process":         "Failed to process %s: %s",
		"results_exported_to":       "Results exported to: %s",
		"processing_urls":           "Processing %d URL(s) from file",
		"processing_url":            "Processing URL",
		"invalid_url_format":        "Invalid URL format: %s",
		"failed_to_process_url":     "Failed to process URL: %s - %s",
		"task_completed":            "Task completed",
		"submitting_tasks":         "Submitting %d task(s)",

		// Human Review command
		"human_review_description":           "Human-review collaboration management.",
		"human_review_rule_create_title":     "Create a human-review collaboration rule.",
		"human_review_get_config_title":      "View human-review collaboration rule.",
		"human_review_rule_update_title":     "Update a human-review collaboration rule.",
		"human_review_rule_delete_title":     "Delete a human-review collaboration rule.",
		"human_review_ai_generate_title":     "AI-recommended review rules.",
		"human_review_task_create_title":     "Create a human-review task.",
		"human_review_task_query_title":      "Query human-review async task status.",
		"human_review_result_update_title":   "Update document human-review result.",
		"human_review_option_app_id":         "Application ID",
		"human_review_option_rule_name":      "Rule name (must be unique)",
		"human_review_option_rule_status":    "Enable rule, \"true\" or \"false\" (default: true)",
		"human_review_option_rule":           "Review rules in JSON format, e.g. '[{\"rule_dimension\":\"...\",\"rule_setting\":\"...\"}]'",
		"human_review_option_rule_logic":     "Rule logic: 1=any condition, 2=all conditions (default: 1)",
		"human_review_option_fields":         "Fields with accuracy expectations in JSON format (optional)",
		"human_review_option_local":          "Local file or folder path",
		"human_review_option_url":            "File URL or URL list file path",
		"human_review_option_file_task_id":   "Document task ID to update",
		"human_review_option_collaboration_result": "Collaboration result in JSON format. Required fields: field_name, field_type (string/date/table), field_values (array). When field_type=table, use table_values: [[{field_name, field_type (string/date only), field_values}]]. Example: '[{\"field_name\":\"amount\",\"field_type\":\"string\",\"field_values\":[\"100\"]}]'",
		"human_review_rule_created":          "Collaboration rule created successfully",
		"human_review_rule_updated":          "Collaboration rule updated successfully",
		"human_review_rule_deleted":          "Collaboration rule deleted successfully",
		"human_review_result_updated":        "Document result updated successfully",

		// Webhook command
		"webhook_description":                "Webhook callback configuration management.",
		"webhook_create_title":               "Create a webhook configuration.",
		"webhook_get_config_title":           "List webhook configurations.",
		"webhook_update_title":               "Update a webhook configuration.",
		"webhook_delete_title":               "Delete a webhook configuration.",
		"webhook_log_title":                  "Query webhook push logs.",
		"webhook_option_webhook_url":         "Webhook callback URL",
		"webhook_option_event_types":         "Event types (comma-separated): 1=start, 2=timeout, 3=completed, 4=failed",
		"webhook_option_app_id":              "Application IDs (comma-separated, optional)",
		"webhook_option_webhook_id":          "Webhook configuration ID",
		"webhook_option_webhook_id_optional": "Webhook configuration ID (optional)",
		"webhook_option_start_time":          "Start time (yyyy-MM-dd HH:mm:ss, optional, default: last 72h)",
		"webhook_option_end_time":            "End time (yyyy-MM-dd HH:mm:ss, optional, default: last 72h)",
		"webhook_created":                    "Webhook created successfully",
		"webhook_updated":                    "Webhook updated successfully",
		"webhook_deleted":                    "Webhook deleted successfully",
	},
	"zh": {
		// Global options
		"option_json":              "输出为 JSON 格式",
		"option_quiet":             "除错误外抑制所有输出",
		"option_lang":              "设置语言 (en 或 zh)",

		// Config command
		"config_description":        "认证配置管理。",
		"config_set_title":         "设置或更新 API Key 和 API Base URL。",
		"config_get_title":         "查看当前配置。",
		"config_clear_title":       "清除本地配置。",
		"config_set_examples":      "示例:",
		"error_not_configured":     "配置未完成，请先运行 'adp config set' 配置 API Key 和 API Base URL。",
		"error_api_key_or_url_required": "必须提供 --api-key 或 --api-base-url 中的至少一个",
		"api_key_configured":       "API Key 配置成功",
		"api_base_url_configured":  "API Base URL 配置成功",
		"config_cleared":           "配置已清除",
		"confirm_clear_config":     "确定要清除所有配置吗?",
		"option_force_clear":       "跳过确认，强制清除",
		"option_api_key":          "API 认证密钥（格式: xxxxxxxxxxxx）",
		"option_api_base_url":     "API Base URL (默认: https://adp.laiye.com)",

		// Parse command
		"parse_description":         "文档解析。",
		"parse_local_title":        "解析本地文件或文件夹。",
		"parse_url_title":          "解析 URL 文件或 URL 列表文件。",
		"parse_base64_title":       "解析 base64 编码的文件内容。",
		"parse_query_title":        "查询解析异步任务状态。",
		"option_app_id_parse":      "解析的应用 ID",
		"option_async":             "异步处理",
		"option_export":            "导出结果到 JSON 文件",
		"option_timeout":           "同步模式超时时间（秒）",
		"option_file_name":        "文件显示名（含扩展名，用于服务端识别文件类型，不会从磁盘读取）。默认: document",
		"option_retry":             "可重试错误的重试次数（默认：0）",
		"option_watch":             "监视任务直到完成",
		"option_watch_timeout":    "监视模式超时时间（秒）",
		"option_no_wait":          "仅提交任务，不等待结果（与 --async 配合使用）",
		"option_task_file":        "从 JSON 文件读取任务 ID（--no-wait 的输出文件）",

		// Extract command
		"extract_description":       "文档抽取。",
		"extract_local_title":      "抽取本地文件或文件夹。",
		"extract_url_title":        "抽取 URL 文件或 URL 列表文件。",
		"extract_base64_title":     "抽取 base64 编码的文件内容。",
		"extract_query_title":      "查询抽取异步任务状态。",
		"option_app_id_extract":    "抽取的应用 ID",

		// App ID command
		"app_id_description":        "应用管理。",
		"app_id_list_title":        "列出所有可用的应用 ID。",
		"app_id_list_cache_title":  "从缓存获取应用ID列表（快速）。",
		"app_id_list_app_label":    "按标签过滤应用（可选）",
		"app_id_list_app_type":     "按应用类型过滤：0=系统预设应用，1=自定义应用；不传则返回全部应用",
		"app_id_list_limit":        "限制返回应用数量（默认：120）",
		"no_applications_found":    "未找到应用",

		// Custom App command
		"custom_app_description":                    "自定义抽取应用管理。",
		"custom_app_create_title":                  "创建自定义抽取应用。",
		"custom_app_create_api_key":               "API 认证密钥（可选，如果已配置则不需要）",
		"custom_app_create_app_name":               "应用名称 (最多 50 个字符)",
		"custom_app_create_app_label":              "应用标签（最多5个，可选，格式：JSON 数组或逗号分隔字符串）",
		"custom_app_create_extract_fields":          "字段配置，JSON 格式或 JSON 文件路径（必需）。字段需要包含 field_name（字段名）、field_type（字段类型）、field_prompt（字段提示词）。field_type 可选值：string（文本）、date（日期）、table（表格）。table 类型时 field_prompt 可为空，需包含 sub_fields（表格子字段）。表格子字段的 field_type 只能为 string 或 date。参见文档查看示例。",
		"custom_app_create_parse_mode":              "解析模式：standard=标准解析；advance=增强解析；agentic=智能体解析",
		"custom_app_create_enable_long_doc":         "启用长文档支持 (true/false，可选)",
		"custom_app_create_long_doc_config":         "长文档配置，JSON 格式或 JSON 文件路径（可选，仅当 enable-long-doc=true 时生效）",
		"custom_app_update_title":                  "更新自定义抽取应用。",
		"custom_app_update_api_key":                "API 认证密钥（可选，如果已配置则不需要）",
		"custom_app_update_app_id":                 "应用 ID（必填）",
		"custom_app_update_app_name":                "应用名称 (最多 50 个字符)",
		"custom_app_update_app_label":               "应用标签（最多5个，可选，格式：JSON 数组或逗号分隔字符串）",
		"custom_app_update_extract_fields":          "字段配置，JSON 格式或 JSON 文件路径（必需）。字段需要包含 field_name（字段名）、field_type（字段类型）、field_prompt（字段提示词）。field_type 可选值：string（文本）、date（日期）、table（表格）。table 类型时 field_prompt 可为空，需包含 sub_fields（表格子字段）。表格子字段的 field_type 只能为 string 或 date。参见文档查看示例。",
		"custom_app_update_parse_mode":             "解析模式：standard=标准解析；advance=增强解析；agentic=智能体解析",
		"custom_app_update_enable_long_doc":         "启用长文档支持 (true/false，可选；不传则保持服务端默认行为)",
		"custom_app_update_long_doc_config":       "长文档配置，JSON 格式或 JSON 文件路径（可选，仅当 enable-long-doc=true 时生效）",
		"custom_app_get_config_title":              "查询自定义应用配置。",
		"custom_app_get_config_api_key":            "API 认证密钥",
		"custom_app_get_config_app_id":             "应用 ID",
		"custom_app_get_config_config_version":     "配置版本 (可选，默认: 最新版)",
		"custom_app_delete_title":                   "删除自定义应用。",
		"custom_app_delete_api_key":                "API 认证密钥",
		"custom_app_delete_app_id":                "应用 ID",
		"custom_app_delete_version_title":          "删除指定配置版本。",
		"custom_app_delete_version_api_key":        "API 认证密钥",
		"custom_app_delete_version_app_id":         "应用 ID",
		"custom_app_delete_version_config_version":  "要删除的配置版本",
		"custom_app_ai_generate_title":             "AI 生成抽取字段推荐。",
		"custom_app_ai_generate_api_key":           "API 认证密钥",
		"custom_app_ai_generate_app_id":            "应用 ID",
		"custom_app_ai_generate_file_url":          "示例文档 URL",
		"custom_app_ai_generate_file_local":         "本地示例文档",
		"custom_app_ai_generate_file_base64":        "Base64 编码的示例文档内容",
		"app_created":               "自定义应用创建成功",
		"app_updated":               "自定义应用更新成功",
		"app_deleted":               "应用删除成功",
		"app_delete_request_ok":     "删除请求成功",
		"version_deleted":          "版本删除成功",
		"version_delete_request_ok": "删除请求成功",
		"not_found_may_deleted_app": "应用 %s 不存在，可能已被删除。",
		"not_found_may_deleted_ver": "版本 %s 不存在，可能已被删除。",

		// Credit command
		"credit_description":        "查询剩余资产数。",
		"credit_api_key":           "API 认证密钥（可选，如果已配置则不需要）",
		"credit_info":             "资产信息",
		"remaining_credits":       "剩余资产",
		"recharge_message":        "充值请访问",

		// Help command
		"help_description":         "显示帮助信息和可用命令。",

		// Schema command
		"schema_description":        "显示命令 schema，用于 Agent 自省。",

		// Version command
		"version_description":       "显示版本号。",

		// Common
		"error":                      "错误:",
		"warning":                    "警告:",
		"no_valid_files":            "没有有效的文件需要处理",
		"processing_files":          "正在处理 %d 个文件",
		"failed_to_process":         "处理失败 %s: %s",
		"results_exported_to":       "结果已导出到: %s",
		"processing_urls":           "正在处理文件中的 %d 个 URL",
		"processing_url":            "正在处理 URL",
		"invalid_url_format":        "无效的 URL 格式: %s",
		"failed_to_process_url":     "处理 URL 失败: %s - %s",
		"task_completed":            "任务完成",
		"submitting_tasks":         "正在提交 %d 个任务",

		// Human Review command
		"human_review_description":           "人机协同管理。",
		"human_review_rule_create_title":     "创建人机协同审核规则。",
		"human_review_get_config_title":      "查看人机协同审核规则。",
		"human_review_rule_update_title":     "修改人机协同审核规则。",
		"human_review_rule_delete_title":     "删除人机协同审核规则。",
		"human_review_ai_generate_title":     "AI 推荐审核规则。",
		"human_review_task_create_title":     "创建人机协同任务。",
		"human_review_task_query_title":      "查询人机协同异步任务状态。",
		"human_review_result_update_title":   "修改文档人机协同处理结果。",
		"human_review_option_app_id":         "应用 ID",
		"human_review_option_rule_name":      "规则名称（不能重复）",
		"human_review_option_rule_status":    "是否启用规则，\"true\" 或 \"false\"（默认：true）",
		"human_review_option_rule":           "审核规则，JSON 格式，如 '[{\"rule_dimension\":\"...\",\"rule_setting\":\"...\"}]'",
		"human_review_option_rule_logic":     "规则逻辑：1=任一条件，2=所有条件（默认：1）",
		"human_review_option_fields":         "字段及准确度期望，JSON 格式（可选）",
		"human_review_option_local":          "本地文件或文件夹路径",
		"human_review_option_url":            "文件 URL 或 URL 列表文件路径",
		"human_review_option_file_task_id":   "待修改的文档任务 ID",
		"human_review_option_collaboration_result": "协同结果，JSON 格式。必填字段：field_name、field_type（string/date/table）、field_values（数组）。当 field_type=table 时，使用 table_values: [[{field_name, field_type (仅 string/date), field_values}]]。示例：'[{\"field_name\":\"amount\",\"field_type\":\"string\",\"field_values\":[\"100\"]}]'",
		"human_review_rule_created":          "人机协同审核规则创建成功",
		"human_review_rule_updated":          "人机协同审核规则更新成功",
		"human_review_rule_deleted":          "人机协同审核规则删除成功",
		"human_review_result_updated":        "文档处理结果更新成功",

		// Webhook command
		"webhook_description":                "Webhook 回调配置管理。",
		"webhook_create_title":               "创建回调配置。",
		"webhook_get_config_title":           "查看回调配置列表。",
		"webhook_update_title":               "修改回调配置。",
		"webhook_delete_title":               "删除回调配置。",
		"webhook_log_title":                  "查询 Webhook 推送日志。",
		"webhook_option_webhook_url":         "回调服务器地址",
		"webhook_option_event_types":         "事件类型（逗号分隔）：1=任务开始，2=任务超时，3=任务完成，4=任务失败",
		"webhook_option_app_id":              "应用 ID（逗号分隔，可选）",
		"webhook_option_webhook_id":          "Webhook 配置 ID",
		"webhook_option_webhook_id_optional": "Webhook 配置 ID（可选）",
		"webhook_option_start_time":          "开始时间（yyyy-MM-dd HH:mm:ss，可选，默认最近 72 小时）",
		"webhook_option_end_time":            "结束时间（yyyy-MM-dd HH:mm:ss，可选，默认最近 72 小时）",
		"webhook_created":                    "Webhook 创建成功",
		"webhook_updated":                    "Webhook 更新成功",
		"webhook_deleted":                    "Webhook 删除成功",
	},
}

// SetLanguage sets the current language
func SetLanguage(lang string) error {
	lang = strings.ToLower(lang)
	if lang != "en" && lang != "zh" {
		return fmt.Errorf("unsupported language: %s", lang)
	}
	currentLang = lang
	return nil
}

// GetLanguage returns the current language
func GetLanguage() string {
	return currentLang
}

// T translates a key
func T(key string, args ...interface{}) string {
	translations, ok := Translations[currentLang]
	if !ok {
		translations = Translations["en"]
	}

	text, ok := translations[key]
	if !ok {
		// Fallback to English
		if enText, ok := Translations["en"][key]; ok {
			text = enText
		} else {
			text = key
		}
	}

	if len(args) > 0 {
		return fmt.Sprintf(text, args...)
	}
	return text
}

func detectLanguage() string {
	// 1. ADP_LANG environment variable (highest priority, explicit override)
	if lang := os.Getenv("ADP_LANG"); lang != "" {
		lang = strings.ToLower(lang)
		if strings.HasPrefix(lang, "zh") {
			return "zh"
		}
		return "en"
	}

	// 2. LC_ALL environment variable
	if lang := os.Getenv("LC_ALL"); lang != "" {
		lang = strings.ToLower(lang)
		if strings.HasPrefix(lang, "zh") {
			return "zh"
		}
	}

	// 3. LANG environment variable
	if lang := os.Getenv("LANG"); lang != "" {
		lang = strings.ToLower(lang)
		if strings.HasPrefix(lang, "zh") {
			return "zh"
		}
	}

	// Default to English
	return "en"
}
