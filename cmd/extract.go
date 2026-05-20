package cmd

import (
	"github.com/spf13/cobra"
	"github.com/laiye-ai/adp-cli/internal/errors"
	"github.com/laiye-ai/adp-cli/internal/i18n"
)

// extractCmd represents the extract command group
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: i18n.T("extract_description"),
	Long:  i18n.T("extract_description"),
}

// extractLocalCmd represents the extract local command
var extractLocalCmd = &cobra.Command{
	Use:   "local",
	Short: i18n.T("extract_local_title"),
	Long:  i18n.T("extract_local_title"),
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		appID, _ := cmd.Flags().GetString("app-id")
		asyncMode, _ := cmd.Flags().GetBool("async")
		noWait, _ := cmd.Flags().GetBool("no-wait")
		export, _ := cmd.Flags().GetString("export")
		timeout, _ := cmd.Flags().GetInt("timeout")
		retry, _ := cmd.Flags().GetInt("retry")

		processLocalFiles(path, appID, asyncMode, noWait, export, timeout, retry, "extract")
	},
}

// extractURLCmd represents the extract URL command
var extractURLCmd = &cobra.Command{
	Use:   "url",
	Short: i18n.T("extract_url_title"),
	Long:  i18n.T("extract_url_title"),
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		urlArg := args[0]
		appID, _ := cmd.Flags().GetString("app-id")
		asyncMode, _ := cmd.Flags().GetBool("async")
		noWait, _ := cmd.Flags().GetBool("no-wait")
		export, _ := cmd.Flags().GetString("export")
		timeout, _ := cmd.Flags().GetInt("timeout")
		retry, _ := cmd.Flags().GetInt("retry")

		urls := resolveURLInput(urlArg)
		processURLs(urls, appID, asyncMode, noWait, export, timeout, retry, "extract")
	},
}

// extractBase64Cmd represents the extract base64 command
var extractBase64Cmd = &cobra.Command{
	Use:   "base64",
	Short: i18n.T("extract_base64_title"),
	Long:  i18n.T("extract_base64_title"),
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		appID, _ := cmd.Flags().GetString("app-id")
		asyncMode, _ := cmd.Flags().GetBool("async")
		noWait, _ := cmd.Flags().GetBool("no-wait")
		export, _ := cmd.Flags().GetString("export")
		timeout, _ := cmd.Flags().GetInt("timeout")
		fileName, _ := cmd.Flags().GetString("file-name")
		retry, _ := cmd.Flags().GetInt("retry")

		processBase64(args, appID, fileName, asyncMode, noWait, export, timeout, retry, "extract")
	},
}

// extractQueryCmd represents the extract query command
var extractQueryCmd = &cobra.Command{
	Use:   "query",
	Short: i18n.T("extract_query_title"),
	Long:  i18n.T("extract_query_title"),
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		watch, _ := cmd.Flags().GetBool("watch")
		taskFile, _ := cmd.Flags().GetString("file")
		export, _ := cmd.Flags().GetString("export")
		timeout, _ := cmd.Flags().GetInt("timeout")

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

		queryTasks(args, watch, taskFile, export, timeout, "extract")
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)

	extractCmd.AddCommand(extractLocalCmd)
	extractCmd.AddCommand(extractURLCmd)
	extractCmd.AddCommand(extractBase64Cmd)
	extractCmd.AddCommand(extractQueryCmd)

	// extract local flags
	extractLocalCmd.Flags().String("app-id", "", i18n.T("option_app_id_extract"))
	extractLocalCmd.Flags().Bool("async", false, i18n.T("option_async"))
	extractLocalCmd.Flags().Bool("no-wait", false, i18n.T("option_no_wait"))
	extractLocalCmd.Flags().String("export", "", i18n.T("option_export"))
	extractLocalCmd.Flags().Int("timeout", 900, i18n.T("option_timeout"))
	extractLocalCmd.Flags().Int("retry", 0, i18n.T("option_retry"))
	extractLocalCmd.MarkFlagRequired("app-id")

	// extract url flags
	extractURLCmd.Flags().String("app-id", "", i18n.T("option_app_id_extract"))
	extractURLCmd.Flags().Bool("async", false, i18n.T("option_async"))
	extractURLCmd.Flags().Bool("no-wait", false, i18n.T("option_no_wait"))
	extractURLCmd.Flags().String("export", "", i18n.T("option_export"))
	extractURLCmd.Flags().Int("timeout", 900, i18n.T("option_timeout"))
	extractURLCmd.Flags().Int("retry", 0, i18n.T("option_retry"))
	extractURLCmd.MarkFlagRequired("app-id")

	// extract base64 flags
	extractBase64Cmd.Flags().String("app-id", "", i18n.T("option_app_id_extract"))
	extractBase64Cmd.Flags().Bool("async", false, i18n.T("option_async"))
	extractBase64Cmd.Flags().Bool("no-wait", false, i18n.T("option_no_wait"))
	extractBase64Cmd.Flags().String("export", "", i18n.T("option_export"))
	extractBase64Cmd.Flags().Int("timeout", 900, i18n.T("option_timeout"))
	extractBase64Cmd.Flags().String("file-name", "document", i18n.T("option_file_name"))
	extractBase64Cmd.Flags().Int("retry", 0, i18n.T("option_retry"))
	extractBase64Cmd.MarkFlagRequired("app-id")

	// extract query flags
	extractQueryCmd.Flags().Bool("watch", false, i18n.T("option_watch"))
	extractQueryCmd.Flags().String("file", "", i18n.T("option_task_file"))
	extractQueryCmd.Flags().String("export", "", i18n.T("option_export"))
	extractQueryCmd.Flags().Int("timeout", 900, i18n.T("option_watch_timeout"))
}
