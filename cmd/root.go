package cmd

import (
	"fmt"
	"os"

	"github.com/laiye-ai/adp-cli/internal/i18n"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var version = "dev"

var (
	jsonMode   bool
	quietMode  bool
	langMode   string
	versionCmd = &cobra.Command{
		Use:   "version",
		Short: i18n.T("version_description"),
		Long:  i18n.T("version_description"),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("ADP CLI version %s\n", version)
		},
	}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "adp",
	Short: "ADP CLI - AI Document Platform Command Line Tool",
	Long: `ADP CLI - AI Document Platform Command Line Tool
Command line tool for AI Document Platform, providing complete document processing capabilities.
Supports document parsing, content extraction, and async task querying.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Configure zerolog
		// Log level: debug, info, warn, error
		logLevel := os.Getenv("ADP_LOG_LEVEL")
		if logLevel == "" {
			logLevel = "info"
		}

		if quietMode || logLevel == "error" {
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		} else if logLevel == "warn" {
			zerolog.SetGlobalLevel(zerolog.WarnLevel)
		} else if logLevel == "debug" {
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		} else {
			zerolog.SetGlobalLevel(zerolog.InfoLevel)
		}

		if jsonMode {
			log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
		} else {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		}

		// Sync flags to formatter
		formatterOut.SetJSONMode(jsonMode)
		formatterOut.SetQuietMode(quietMode)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Command execution failed")
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().BoolVar(&jsonMode, "json", false, "Output in JSON format")
	rootCmd.PersistentFlags().BoolVar(&quietMode, "quiet", false, "Suppress all output except errors")
	rootCmd.PersistentFlags().StringVar(&langMode, "lang", "", "Set language (en or zh)")

	// Version command
	rootCmd.AddCommand(versionCmd)

	// Custom HelpFunc to ensure i18n is applied before rendering help
	// This is needed because cobra's --help bypasses OnInitialize
	defaultHelpFunc := rootCmd.HelpFunc()
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		// Apply language setting (mirrors initConfig logic)
		// langMode is already set from --lang flag by cobra
		lang := viper.GetString("lang")
		if langMode == "" {
			langMode = lang
		}
		if langMode != "" {
			i18n.SetLanguage(langMode)
		}
		reloadCommandTranslations()
		defaultHelpFunc(cmd, args)
	})
}

func initConfig() {
	// --lang flag always takes precedence over config file
	// langMode is set by the flag; only check viper if langMode is empty (no flag provided)
	if langMode == "" {
		langMode = viper.GetString("lang")
	}
	if langMode != "" {
		i18n.SetLanguage(langMode)
	}
	// Reload command translations
	reloadCommandTranslations()
}

// updateFlagHelp updates a flag's help text dynamically
func updateFlagHelp(cmd *cobra.Command, flagName string, help string) {
	f := cmd.Flags().Lookup(flagName)
	if f != nil {
		f.Usage = help
	}
}

// reloadCommandTranslations updates all command Short/Long descriptions with current language
func reloadCommandTranslations() {
	// creditCmd
	creditCmd.Short = i18n.T("credit_description")
	creditCmd.Long = i18n.T("credit_description")
	updateFlagHelp(creditCmd, "api-key", i18n.T("credit_api_key"))

	// configCmd
	configCmd.Short = i18n.T("config_description")
	configCmd.Long = i18n.T("config_description")
	configSetCmd.Short = i18n.T("config_set_title")
	configSetCmd.Long = i18n.T("config_set_title")
	updateFlagHelp(configSetCmd, "api-key", i18n.T("option_api_key"))
	updateFlagHelp(configSetCmd, "api-base-url", i18n.T("option_api_base_url"))
	configGetCmd.Short = i18n.T("config_get_title")
	configGetCmd.Long = i18n.T("config_get_title")
	configClearCmd.Short = i18n.T("config_clear_title")
	configClearCmd.Long = i18n.T("config_clear_title")
	updateFlagHelp(configClearCmd, "force", i18n.T("option_force_clear"))

	// appIDCmd
	appIDCmd.Short = i18n.T("app_id_description")
	appIDCmd.Long = i18n.T("app_id_description")
	listAppsCmd.Short = i18n.T("app_id_list_title")
	listAppsCmd.Long = i18n.T("app_id_list_title")
	listAppsCacheCmd.Short = i18n.T("app_id_list_cache_title")
	listAppsCacheCmd.Long = i18n.T("app_id_list_cache_title")
	updateFlagHelp(listAppsCmd, "app-label", i18n.T("app_id_list_app_label"))
	updateFlagHelp(listAppsCmd, "app-type", i18n.T("app_id_list_app_type"))
	updateFlagHelp(listAppsCmd, "limit", i18n.T("app_id_list_limit"))
	updateFlagHelp(listAppsCacheCmd, "app-label", i18n.T("app_id_list_app_label"))

	// customAppCmd
	customAppCmd.Short = i18n.T("custom_app_description")
	customAppCmd.Long = i18n.T("custom_app_description")
	createCustomAppCmd.Short = i18n.T("custom_app_create_title")
	createCustomAppCmd.Long = i18n.T("custom_app_create_title")
	updateCustomAppCmd.Short = i18n.T("custom_app_update_title")
	updateCustomAppCmd.Long = i18n.T("custom_app_update_title")
	getConfigCustomAppCmd.Short = i18n.T("custom_app_get_config_title")
	getConfigCustomAppCmd.Long = i18n.T("custom_app_get_config_title")
	deleteCustomAppCmd.Short = i18n.T("custom_app_delete_title")
	deleteCustomAppCmd.Long = i18n.T("custom_app_delete_title")
	deleteVersionCustomAppCmd.Short = i18n.T("custom_app_delete_version_title")
	deleteVersionCustomAppCmd.Long = i18n.T("custom_app_delete_version_title")
	aiGenerateCustomAppCmd.Short = i18n.T("custom_app_ai_generate_title")
	aiGenerateCustomAppCmd.Long = i18n.T("custom_app_ai_generate_title")

	// create flags
	updateFlagHelp(createCustomAppCmd, "api-key", i18n.T("custom_app_create_api_key"))
	updateFlagHelp(createCustomAppCmd, "app-name", i18n.T("custom_app_create_app_name"))
	updateFlagHelp(createCustomAppCmd, "app-label", i18n.T("custom_app_create_app_label"))
	updateFlagHelp(createCustomAppCmd, "extract-fields", i18n.T("custom_app_create_extract_fields"))
	updateFlagHelp(createCustomAppCmd, "parse-mode", i18n.T("custom_app_create_parse_mode"))
	updateFlagHelp(createCustomAppCmd, "enable-long-doc", i18n.T("custom_app_create_enable_long_doc"))
	updateFlagHelp(createCustomAppCmd, "long-doc-config", i18n.T("custom_app_create_long_doc_config"))

	// update flags
	updateFlagHelp(updateCustomAppCmd, "api-key", i18n.T("custom_app_update_api_key"))
	updateFlagHelp(updateCustomAppCmd, "app-id", i18n.T("custom_app_update_app_id"))
	updateFlagHelp(updateCustomAppCmd, "app-name", i18n.T("custom_app_update_app_name"))
	updateFlagHelp(updateCustomAppCmd, "app-label", i18n.T("custom_app_update_app_label"))
	updateFlagHelp(updateCustomAppCmd, "extract-fields", i18n.T("custom_app_update_extract_fields"))
	updateFlagHelp(updateCustomAppCmd, "parse-mode", i18n.T("custom_app_update_parse_mode"))
	updateFlagHelp(updateCustomAppCmd, "enable-long-doc", i18n.T("custom_app_update_enable_long_doc"))
	updateFlagHelp(updateCustomAppCmd, "long-doc-config", i18n.T("custom_app_update_long_doc_config"))

	// get-config flags
	updateFlagHelp(getConfigCustomAppCmd, "api-key", i18n.T("custom_app_get_config_api_key"))
	updateFlagHelp(getConfigCustomAppCmd, "app-id", i18n.T("custom_app_get_config_app_id"))
	updateFlagHelp(getConfigCustomAppCmd, "config-version", i18n.T("custom_app_get_config_config_version"))

	// delete flags
	updateFlagHelp(deleteCustomAppCmd, "api-key", i18n.T("custom_app_delete_api_key"))
	updateFlagHelp(deleteCustomAppCmd, "app-id", i18n.T("custom_app_delete_app_id"))

	// delete-version flags
	updateFlagHelp(deleteVersionCustomAppCmd, "api-key", i18n.T("custom_app_delete_version_api_key"))
	updateFlagHelp(deleteVersionCustomAppCmd, "app-id", i18n.T("custom_app_delete_version_app_id"))
	updateFlagHelp(deleteVersionCustomAppCmd, "config-version", i18n.T("custom_app_delete_version_config_version"))

	// ai-generate flags
	updateFlagHelp(aiGenerateCustomAppCmd, "api-key", i18n.T("custom_app_ai_generate_api_key"))
	updateFlagHelp(aiGenerateCustomAppCmd, "app-id", i18n.T("custom_app_ai_generate_app_id"))
	updateFlagHelp(aiGenerateCustomAppCmd, "file-url", i18n.T("custom_app_ai_generate_file_url"))
	updateFlagHelp(aiGenerateCustomAppCmd, "file-local", i18n.T("custom_app_ai_generate_file_local"))
	updateFlagHelp(aiGenerateCustomAppCmd, "base64", i18n.T("custom_app_ai_generate_file_base64"))

	// parseCmd
	parseCmd.Short = i18n.T("parse_description")
	parseCmd.Long = i18n.T("parse_description")
	parseLocalCmd.Short = i18n.T("parse_local_title")
	parseLocalCmd.Long = i18n.T("parse_local_title")
	parseURLCmd.Short = i18n.T("parse_url_title")
	parseURLCmd.Long = i18n.T("parse_url_title")
	parseBase64Cmd.Short = i18n.T("parse_base64_title")
	parseBase64Cmd.Long = i18n.T("parse_base64_title")
	parseQueryCmd.Short = i18n.T("parse_query_title")
	parseQueryCmd.Long = i18n.T("parse_query_title")

	// extractCmd
	extractCmd.Short = i18n.T("extract_description")
	extractCmd.Long = i18n.T("extract_description")
	extractLocalCmd.Short = i18n.T("extract_local_title")
	extractLocalCmd.Long = i18n.T("extract_local_title")
	extractURLCmd.Short = i18n.T("extract_url_title")
	extractURLCmd.Long = i18n.T("extract_url_title")
	extractBase64Cmd.Short = i18n.T("extract_base64_title")
	extractBase64Cmd.Long = i18n.T("extract_base64_title")
	extractQueryCmd.Short = i18n.T("extract_query_title")
	extractQueryCmd.Long = i18n.T("extract_query_title")

	// schemaCmd
	schemaCmd.Short = i18n.T("schema_description")
	schemaCmd.Long = i18n.T("schema_description")

	// versionCmd
	versionCmd.Short = i18n.T("version_description")
	versionCmd.Long = i18n.T("version_description")

	// helpCmd
	helpCmd.Short = i18n.T("help_description")
	helpCmd.Long = i18n.T("help_description")

	// Update all flag descriptions for parse commands
	updateFlagHelp(parseLocalCmd, "app-id", i18n.T("option_app_id_parse"))
	updateFlagHelp(parseLocalCmd, "async", i18n.T("option_async"))
	updateFlagHelp(parseLocalCmd, "export", i18n.T("option_export"))
	updateFlagHelp(parseLocalCmd, "timeout", i18n.T("option_timeout"))
	updateFlagHelp(parseLocalCmd, "concurrency", i18n.T("option_concurrency"))

	updateFlagHelp(parseURLCmd, "app-id", i18n.T("option_app_id_parse"))
	updateFlagHelp(parseURLCmd, "async", i18n.T("option_async"))
	updateFlagHelp(parseURLCmd, "export", i18n.T("option_export"))
	updateFlagHelp(parseURLCmd, "timeout", i18n.T("option_timeout"))
	updateFlagHelp(parseURLCmd, "concurrency", i18n.T("option_concurrency"))

	updateFlagHelp(parseBase64Cmd, "app-id", i18n.T("option_app_id_parse"))
	updateFlagHelp(parseBase64Cmd, "async", i18n.T("option_async"))
	updateFlagHelp(parseBase64Cmd, "export", i18n.T("option_export"))
	updateFlagHelp(parseBase64Cmd, "timeout", i18n.T("option_timeout"))
	updateFlagHelp(parseBase64Cmd, "file-name", i18n.T("option_file_name"))
	updateFlagHelp(parseBase64Cmd, "concurrency", i18n.T("option_concurrency"))

	// Update all flag descriptions for extract commands
	updateFlagHelp(extractLocalCmd, "app-id", i18n.T("option_app_id_extract"))
	updateFlagHelp(extractLocalCmd, "async", i18n.T("option_async"))
	updateFlagHelp(extractLocalCmd, "export", i18n.T("option_export"))
	updateFlagHelp(extractLocalCmd, "timeout", i18n.T("option_timeout"))
	updateFlagHelp(extractLocalCmd, "concurrency", i18n.T("option_concurrency"))

	updateFlagHelp(extractURLCmd, "app-id", i18n.T("option_app_id_extract"))
	updateFlagHelp(extractURLCmd, "async", i18n.T("option_async"))
	updateFlagHelp(extractURLCmd, "export", i18n.T("option_export"))
	updateFlagHelp(extractURLCmd, "timeout", i18n.T("option_timeout"))
	updateFlagHelp(extractURLCmd, "concurrency", i18n.T("option_concurrency"))

	updateFlagHelp(extractBase64Cmd, "app-id", i18n.T("option_app_id_extract"))
	updateFlagHelp(extractBase64Cmd, "async", i18n.T("option_async"))
	updateFlagHelp(extractBase64Cmd, "export", i18n.T("option_export"))
	updateFlagHelp(extractBase64Cmd, "timeout", i18n.T("option_timeout"))
	updateFlagHelp(extractBase64Cmd, "file-name", i18n.T("option_file_name"))
	updateFlagHelp(extractBase64Cmd, "concurrency", i18n.T("option_concurrency"))
}
