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
		"option_concurrency":       "Concurrency for batch processing. Default: 1. Free users max=1, paid users max=2. Other values will not take effect.",
		"option_file_name":        "File name (default: document)",
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
		"custom_app_update_enable_long_doc":      "Enable long document support (true/false)",
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
		"error_invalid_concurrency": "Free users: 1, paid users: 2",
		"error_not_paid_user":       "You are a free user, maximum concurrency is 1",
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
		"option_concurrency":       "批量处理时的并发数。默认：1。免费用户最大为1，付费用户最大为2。输入其他数值，并发将不生效，请等待处理。",
		"option_file_name":        "文件名（默认: document）",
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
		"custom_app_update_enable_long_doc":         "启用长文档支持 (true/false)",
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
		"error_invalid_concurrency": "免费用户：1，付费用户：2",
		"error_not_paid_user":       "您是免费用户，最大并发数为1",
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
