package cmd

import (
	"fmt"

	"github.com/laiye-ai/adp-cli/internal/errors"
	"github.com/laiye-ai/adp-cli/internal/i18n"
	"github.com/spf13/cobra"
)

// humanReviewCmd represents the human-review command group
var humanReviewCmd = &cobra.Command{
	Use:   "human-review",
	Short: i18n.T("human_review_description"),
	Long:  i18n.T("human_review_description"),
}

// --- rule-create ---
var hrRuleCreateCmd = &cobra.Command{
	Use:   "rule-create",
	Short: i18n.T("human_review_rule_create_title"),
	Long:  i18n.T("human_review_rule_create_title"),
	Run: func(cmd *cobra.Command, args []string) {
		appID, _ := cmd.Flags().GetString("app-id")
		ruleName, _ := cmd.Flags().GetString("rule-name")
		ruleStatus, _ := cmd.Flags().GetString("rule-status")
		ruleStatusBool := ruleStatus == "true"
		ruleStr, _ := cmd.Flags().GetString("rule")
		ruleLogic, _ := cmd.Flags().GetInt("rule-logic")

		client, _ := initClientWithConfig("human-review rule-create")

		rules, err := parseJSONParam(ruleStr)
		if err != nil {
			formatterOut.ExitWithError(errors.NewCLIError(
				fmt.Sprintf("Invalid --rule JSON: %v", err),
				errors.ErrorTypeParam, errors.ExitParameterError, false,
				"Provide valid JSON array, e.g. '[{\"rule_dimension\":\"...\",\"rule_setting\":\"...\"}]'", nil,
			))
		}

		result, err := client.CreateCollaborationRule(appID, ruleName, ruleStatusBool, rules, ruleLogic)
		if err != nil {
			formatterOut.ExitWithError(errors.ClassifyException(err, "human-review rule-create"))
		}
		checkAPIResponse(result, "human-review rule-create")

		formatterOut.PrintSuccess(i18n.T("human_review_rule_created"))
		formatterOut.PrintJSON(result)
	},
}

// --- get-config ---
var hrGetConfigCmd = &cobra.Command{
	Use:   "get-config",
	Short: i18n.T("human_review_get_config_title"),
	Long:  i18n.T("human_review_get_config_title"),
	Run: func(cmd *cobra.Command, args []string) {
		appID, _ := cmd.Flags().GetString("app-id")
		client, _ := initClientWithConfig("human-review get-config")

		result, err := client.GetCollaborationRule(appID)
		if err != nil {
			formatterOut.ExitWithError(errors.ClassifyException(err, "human-review get-config"))
		}
		checkAPIResponse(result, "human-review get-config")

		formatterOut.PrintJSON(result)
	},
}

// --- rule-update ---
var hrRuleUpdateCmd = &cobra.Command{
	Use:   "rule-update",
	Short: i18n.T("human_review_rule_update_title"),
	Long:  i18n.T("human_review_rule_update_title"),
	Run: func(cmd *cobra.Command, args []string) {
		appID, _ := cmd.Flags().GetString("app-id")
		ruleName, _ := cmd.Flags().GetString("rule-name")
		ruleStatus, _ := cmd.Flags().GetString("rule-status")
		ruleStatusBool := ruleStatus == "true"
		ruleStr, _ := cmd.Flags().GetString("rule")
		ruleLogic, _ := cmd.Flags().GetInt("rule-logic")

		client, _ := initClientWithConfig("human-review rule-update")

		rules, err := parseJSONParam(ruleStr)
		if err != nil {
			formatterOut.ExitWithError(errors.NewCLIError(
				fmt.Sprintf("Invalid --rule JSON: %v", err),
				errors.ErrorTypeParam, errors.ExitParameterError, false,
				"Provide valid JSON array, e.g. '[{\"rule_dimension\":\"...\",\"rule_setting\":\"...\"}]'", nil,
			))
		}

		result, err := client.UpdateCollaborationRule(appID, ruleName, ruleStatusBool, rules, ruleLogic)
		if err != nil {
			formatterOut.ExitWithError(errors.ClassifyException(err, "human-review rule-update"))
		}
		checkAPIResponse(result, "human-review rule-update")

		formatterOut.PrintSuccess(i18n.T("human_review_rule_updated"))
		formatterOut.PrintJSON(result)
	},
}

// --- rule-delete ---
var hrRuleDeleteCmd = &cobra.Command{
	Use:   "rule-delete",
	Short: i18n.T("human_review_rule_delete_title"),
	Long:  i18n.T("human_review_rule_delete_title"),
	Run: func(cmd *cobra.Command, args []string) {
		appID, _ := cmd.Flags().GetString("app-id")
		client, _ := initClientWithConfig("human-review rule-delete")

		result, err := client.DeleteCollaborationRule(appID)
		if err != nil {
			formatterOut.ExitWithError(errors.ClassifyException(err, "human-review rule-delete"))
		}
		checkAPIResponse(result, "human-review rule-delete")

		formatterOut.PrintSuccess(i18n.T("human_review_rule_deleted"))
		formatterOut.PrintJSON(result)
	},
}

// --- rule-ai-generate ---
var hrRuleAIGenerateCmd = &cobra.Command{
	Use:   "rule-ai-generate",
	Short: i18n.T("human_review_rule_ai_generate_title"),
	Long:  i18n.T("human_review_rule_ai_generate_title"),
	Run: func(cmd *cobra.Command, args []string) {
		appID, _ := cmd.Flags().GetString("app-id")
		fieldsStr, _ := cmd.Flags().GetString("fields")

		client, _ := initClientWithConfig("human-review ai-generate")

		var fields []map[string]interface{}
		if fieldsStr != "" {
			parsed, err := parseJSONParam(fieldsStr)
			if err != nil {
				formatterOut.ExitWithError(errors.NewCLIError(
					fmt.Sprintf("Invalid --fields JSON: %v", err),
					errors.ErrorTypeParam, errors.ExitParameterError, false,
					"Provide valid JSON array, e.g. '[{\"field_name\":\"...\",\"field_accuracy\":\"...\"}]'", nil,
				))
			}
			fields = parsed
		}

		result, err := client.RecommendCollaborationRule(appID, fields)
		if err != nil {
			formatterOut.ExitWithError(errors.ClassifyException(err, "human-review ai-generate"))
		}
		checkAPIResponse(result, "human-review ai-generate")

		formatterOut.PrintJSON(result)
	},
}

// --- task-create ---
var hrTaskCreateCmd = &cobra.Command{
	Use:   "task-create",
	Short: i18n.T("human_review_task_create_title"),
	Long:  i18n.T("human_review_task_create_title"),
	Run: func(cmd *cobra.Command, args []string) {
		appID, _ := cmd.Flags().GetString("app-id")
		localPath, _ := cmd.Flags().GetString("local")
		urlInput, _ := cmd.Flags().GetString("url")
		asyncMode, _ := cmd.Flags().GetBool("async")
		noWait, _ := cmd.Flags().GetBool("no-wait")
		export, _ := cmd.Flags().GetString("export")
		timeout, _ := cmd.Flags().GetInt("timeout")
		retry, _ := cmd.Flags().GetInt("retry")

		if localPath == "" && urlInput == "" {
			formatterOut.ExitWithError(errors.NewCLIError(
				"One of --local or --url must be provided",
				errors.ErrorTypeParam, errors.ExitParameterError, false,
				"Provide --local <path> or --url <url>.", nil,
			))
		}
		if localPath != "" && urlInput != "" {
			formatterOut.ExitWithError(errors.NewCLIError(
				"--local and --url are mutually exclusive",
				errors.ErrorTypeParam, errors.ExitParameterError, false,
				"Provide only one of --local or --url.", nil,
			))
		}

		if localPath != "" {
			processLocalFiles(localPath, appID, asyncMode, noWait, export, timeout, retry, "human-review")
		} else {
			urls := resolveURLInput(urlInput)
			processURLs(urls, appID, asyncMode, noWait, export, timeout, retry, "human-review")
		}
	},
}

// --- task-query ---
var hrTaskQueryCmd = &cobra.Command{
	Use:   "task-query",
	Short: i18n.T("human_review_task_query_title"),
	Long:  i18n.T("human_review_task_query_title"),
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		watch, _ := cmd.Flags().GetBool("watch")
		taskFile, _ := cmd.Flags().GetString("file")
		export, _ := cmd.Flags().GetString("export")
		timeout, _ := cmd.Flags().GetInt("timeout")

		if taskFile == "" && len(args) == 0 {
			formatterOut.ExitWithError(errors.NewCLIError(
				"Either task IDs or --file must be provided",
				errors.ErrorTypeParam, errors.ExitParameterError, false,
				"Provide task IDs as arguments or use --file to load from a JSON file.", nil,
			))
		}

		queryTasks(args, watch, taskFile, export, timeout, "human-review")
	},
}

// --- result-update ---
var hrResultUpdateCmd = &cobra.Command{
	Use:   "result-update",
	Short: i18n.T("human_review_result_update_title"),
	Long:  i18n.T("human_review_result_update_title"),
	Run: func(cmd *cobra.Command, args []string) {
		fileTaskID, _ := cmd.Flags().GetString("file-task-id")
		resultStr, _ := cmd.Flags().GetString("collaboration-result")

		client, _ := initClientWithConfig("human-review result-update")

		collaborationResult, err := parseJSONParam(resultStr)
		if err != nil {
			formatterOut.ExitWithError(errors.NewCLIError(
				fmt.Sprintf("Invalid --collaboration-result JSON: %v", err),
				errors.ErrorTypeParam, errors.ExitParameterError, false,
				"Provide valid JSON array, e.g. '[{\"field_name\":\"...\",\"field_type\":\"string\",\"field_values\":\"...\"}]'", nil,
			))
		}

		result, err := client.UpdateCollaborationResult(fileTaskID, collaborationResult)
		if err != nil {
			formatterOut.ExitWithError(errors.ClassifyException(err, "human-review result-update"))
		}
		checkAPIResponse(result, "human-review result-update")

		formatterOut.PrintSuccess(i18n.T("human_review_result_updated"))
		formatterOut.PrintJSON(result)
	},
}

func init() {
	rootCmd.AddCommand(humanReviewCmd)

	humanReviewCmd.AddCommand(hrRuleCreateCmd)
	humanReviewCmd.AddCommand(hrGetConfigCmd)
	humanReviewCmd.AddCommand(hrRuleUpdateCmd)
	humanReviewCmd.AddCommand(hrRuleDeleteCmd)
	humanReviewCmd.AddCommand(hrRuleAIGenerateCmd)
	humanReviewCmd.AddCommand(hrTaskCreateCmd)
	humanReviewCmd.AddCommand(hrTaskQueryCmd)
	humanReviewCmd.AddCommand(hrResultUpdateCmd)

	// rule-create flags
	hrRuleCreateCmd.Flags().String("app-id", "", i18n.T("human_review_option_app_id"))
	hrRuleCreateCmd.Flags().String("rule-name", "", i18n.T("human_review_option_rule_name"))
	hrRuleCreateCmd.Flags().String("rule-status", "true", i18n.T("human_review_option_rule_status"))
	hrRuleCreateCmd.Flags().String("rule", "", i18n.T("human_review_option_rule"))
	hrRuleCreateCmd.Flags().Int("rule-logic", 1, i18n.T("human_review_option_rule_logic"))
	hrRuleCreateCmd.MarkFlagRequired("app-id")
	hrRuleCreateCmd.MarkFlagRequired("rule-name")
	hrRuleCreateCmd.MarkFlagRequired("rule")

	// get-config flags
	hrGetConfigCmd.Flags().String("app-id", "", i18n.T("human_review_option_app_id"))
	hrGetConfigCmd.MarkFlagRequired("app-id")

	// rule-update flags
	hrRuleUpdateCmd.Flags().String("app-id", "", i18n.T("human_review_option_app_id"))
	hrRuleUpdateCmd.Flags().String("rule-name", "", i18n.T("human_review_option_rule_name"))
	hrRuleUpdateCmd.Flags().String("rule-status", "true", i18n.T("human_review_option_rule_status"))
	hrRuleUpdateCmd.Flags().String("rule", "", i18n.T("human_review_option_rule"))
	hrRuleUpdateCmd.Flags().Int("rule-logic", 1, i18n.T("human_review_option_rule_logic"))
	hrRuleUpdateCmd.MarkFlagRequired("app-id")
	hrRuleUpdateCmd.MarkFlagRequired("rule-name")
	hrRuleUpdateCmd.MarkFlagRequired("rule")

	// rule-delete flags
	hrRuleDeleteCmd.Flags().String("app-id", "", i18n.T("human_review_option_app_id"))
	hrRuleDeleteCmd.MarkFlagRequired("app-id")

	// ai-generate flags
	hrRuleAIGenerateCmd.Flags().String("app-id", "", i18n.T("human_review_option_app_id"))
	hrRuleAIGenerateCmd.Flags().String("fields", "", i18n.T("human_review_option_fields"))
	hrRuleAIGenerateCmd.Flags().Int("timeout", 900, i18n.T("option_timeout"))
	hrRuleAIGenerateCmd.MarkFlagRequired("app-id")

	// task-create flags
	hrTaskCreateCmd.Flags().String("app-id", "", i18n.T("human_review_option_app_id"))
	hrTaskCreateCmd.Flags().String("local", "", i18n.T("human_review_option_local"))
	hrTaskCreateCmd.Flags().String("url", "", i18n.T("human_review_option_url"))
	hrTaskCreateCmd.Flags().Bool("async", false, i18n.T("option_async"))
	hrTaskCreateCmd.Flags().Bool("no-wait", false, i18n.T("option_no_wait"))
	hrTaskCreateCmd.Flags().String("export", "", i18n.T("option_export"))
	hrTaskCreateCmd.Flags().Int("timeout", 900, i18n.T("option_timeout"))
	hrTaskCreateCmd.Flags().Int("retry", 0, i18n.T("option_retry"))
	hrTaskCreateCmd.MarkFlagRequired("app-id")

	// task-query flags
	hrTaskQueryCmd.Flags().Bool("watch", false, i18n.T("option_watch"))
	hrTaskQueryCmd.Flags().String("file", "", i18n.T("option_task_file"))
	hrTaskQueryCmd.Flags().String("export", "", i18n.T("option_export"))
	hrTaskQueryCmd.Flags().Int("timeout", 900, i18n.T("option_watch_timeout"))

	// result-update flags
	hrResultUpdateCmd.Flags().String("file-task-id", "", i18n.T("human_review_option_file_task_id"))
	hrResultUpdateCmd.Flags().String("collaboration-result", "", i18n.T("human_review_option_collaboration_result"))
	hrResultUpdateCmd.MarkFlagRequired("file-task-id")
	hrResultUpdateCmd.MarkFlagRequired("collaboration-result")
}
