package cmd

import (
	"github.com/laiye-ai/adp-cli/internal/errors"
	"github.com/laiye-ai/adp-cli/internal/i18n"
	"github.com/spf13/cobra"
)

// parseCmd represents the parse command group
var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: i18n.T("parse_description"),
	Long:  i18n.T("parse_description"),
}

// parseLocalCmd represents the parse local command
var parseLocalCmd = &cobra.Command{
	Use:   "local",
	Short: i18n.T("parse_local_title"),
	Long:  i18n.T("parse_local_title"),
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		appID, _ := cmd.Flags().GetString("app-id")
		asyncMode, _ := cmd.Flags().GetBool("async")
		noWait, _ := cmd.Flags().GetBool("no-wait")
		export, _ := cmd.Flags().GetString("export")
		timeout, _ := cmd.Flags().GetInt("timeout")
		concurrency, _ := cmd.Flags().GetInt("concurrency")
		retry, _ := cmd.Flags().GetInt("retry")

		processLocalFiles(path, appID, asyncMode, noWait, export, timeout, concurrency, retry, "parse")
	},
}

// parseURLCmd represents the parse URL command
var parseURLCmd = &cobra.Command{
	Use:   "url",
	Short: i18n.T("parse_url_title"),
	Long:  i18n.T("parse_url_title"),
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		urlArg := args[0]
		appID, _ := cmd.Flags().GetString("app-id")
		asyncMode, _ := cmd.Flags().GetBool("async")
		noWait, _ := cmd.Flags().GetBool("no-wait")
		export, _ := cmd.Flags().GetString("export")
		timeout, _ := cmd.Flags().GetInt("timeout")
		concurrency, _ := cmd.Flags().GetInt("concurrency")
		retry, _ := cmd.Flags().GetInt("retry")

		urls := resolveURLInput(urlArg)
		processURLs(urls, appID, asyncMode, noWait, export, timeout, concurrency, retry, "parse")
	},
}

// parseBase64Cmd represents the parse base64 command
var parseBase64Cmd = &cobra.Command{
	Use:   "base64",
	Short: i18n.T("parse_base64_title"),
	Long:  i18n.T("parse_base64_title"),
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appID, _ := cmd.Flags().GetString("app-id")
		asyncMode, _ := cmd.Flags().GetBool("async")
		noWait, _ := cmd.Flags().GetBool("no-wait")
		export, _ := cmd.Flags().GetString("export")
		timeout, _ := cmd.Flags().GetInt("timeout")
		fileName, _ := cmd.Flags().GetString("file-name")
		concurrency, _ := cmd.Flags().GetInt("concurrency")
		retry, _ := cmd.Flags().GetInt("retry")

		processBase64(args, appID, fileName, asyncMode, noWait, export, timeout, concurrency, retry, "parse")
	},
}

// parseQueryCmd represents the parse query command
var parseQueryCmd = &cobra.Command{
	Use:   "query",
	Short: i18n.T("parse_query_title"),
	Long:  i18n.T("parse_query_title"),
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		watch, _ := cmd.Flags().GetBool("watch")
		taskFile, _ := cmd.Flags().GetString("file")
		export, _ := cmd.Flags().GetString("export")
		timeout, _ := cmd.Flags().GetInt("timeout")
		concurrency, _ := cmd.Flags().GetInt("concurrency")

		if taskFile == "" && len(args) == 0 {
			formatterOut.ExitWithError(errors.NewCLIError(
				"Either task IDs or --file must be provided",
				errors.ErrorTypeParam,
				errors.ExitParameterError,
				false,
				"Provide task IDs as arguments or use --file to load from a JSON file.",
				nil,
			))
		}

		queryTasks(args, watch, taskFile, export, timeout, concurrency, "parse")
	},
}

func init() {
	rootCmd.AddCommand(parseCmd)

	parseCmd.AddCommand(parseLocalCmd)
	parseCmd.AddCommand(parseURLCmd)
	parseCmd.AddCommand(parseBase64Cmd)
	parseCmd.AddCommand(parseQueryCmd)

	// parse local flags
	parseLocalCmd.Flags().String("app-id", "", i18n.T("option_app_id_parse"))
	parseLocalCmd.Flags().Bool("async", false, i18n.T("option_async"))
	parseLocalCmd.Flags().Bool("no-wait", false, i18n.T("option_no_wait"))
	parseLocalCmd.Flags().String("export", "", i18n.T("option_export"))
	parseLocalCmd.Flags().Int("timeout", 900, i18n.T("option_timeout"))
	parseLocalCmd.Flags().Int("concurrency", 1, i18n.T("option_concurrency"))
	parseLocalCmd.Flags().Int("retry", 0, i18n.T("option_retry"))
	parseLocalCmd.MarkFlagRequired("app-id")

	// parse url flags
	parseURLCmd.Flags().String("app-id", "", i18n.T("option_app_id_parse"))
	parseURLCmd.Flags().Bool("async", false, i18n.T("option_async"))
	parseURLCmd.Flags().Bool("no-wait", false, i18n.T("option_no_wait"))
	parseURLCmd.Flags().String("export", "", i18n.T("option_export"))
	parseURLCmd.Flags().Int("timeout", 900, i18n.T("option_timeout"))
	parseURLCmd.Flags().Int("concurrency", 1, i18n.T("option_concurrency"))
	parseURLCmd.Flags().Int("retry", 0, i18n.T("option_retry"))
	parseURLCmd.MarkFlagRequired("app-id")

	// parse base64 flags
	parseBase64Cmd.Flags().String("app-id", "", i18n.T("option_app_id_parse"))
	parseBase64Cmd.Flags().Bool("async", false, i18n.T("option_async"))
	parseBase64Cmd.Flags().Bool("no-wait", false, i18n.T("option_no_wait"))
	parseBase64Cmd.Flags().String("export", "", i18n.T("option_export"))
	parseBase64Cmd.Flags().Int("timeout", 900, i18n.T("option_timeout"))
	parseBase64Cmd.Flags().String("file-name", "document", i18n.T("option_file_name"))
	parseBase64Cmd.Flags().Int("concurrency", 1, i18n.T("option_concurrency"))
	parseBase64Cmd.Flags().Int("retry", 0, i18n.T("option_retry"))
	parseBase64Cmd.MarkFlagRequired("app-id")

	// parse query flags
	parseQueryCmd.Flags().Bool("watch", false, i18n.T("option_watch"))
	parseQueryCmd.Flags().String("file", "", i18n.T("option_task_file"))
	parseQueryCmd.Flags().String("export", "", i18n.T("option_export"))
	parseQueryCmd.Flags().Int("timeout", 900, i18n.T("option_watch_timeout"))
	parseQueryCmd.Flags().Int("concurrency", 1, i18n.T("option_concurrency"))
}