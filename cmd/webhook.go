package cmd

import (
	"strconv"
	"strings"

	"github.com/laiye-ai/adp-cli/internal/errors"
	"github.com/laiye-ai/adp-cli/internal/i18n"
	"github.com/spf13/cobra"
)

// webhookCmd represents the webhook command group
var webhookCmd = &cobra.Command{
	Use:   "webhook",
	Short: i18n.T("webhook_description"),
	Long:  i18n.T("webhook_description"),
}

// parseEventTypes parses a comma-separated event-types string into ints.
// Event types: 1=start, 2=timeout, 3=completed, 4=failed.
func parseEventTypes(s string) ([]int, error) {
	if s == "" {
		return nil, nil
	}
	parts := strings.Split(s, ",")
	result := make([]int, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil, err
		}
		result = append(result, n)
	}
	return result, nil
}

func eventTypesOrExit(s, context string) []int {
	events, err := parseEventTypes(s)
	if err != nil || len(events) == 0 {
		formatterOut.ExitWithError(errors.NewCLIError(
			"--event-types must be a non-empty comma-separated list of integers",
			errors.ErrorTypeParam, errors.ExitParameterError, false,
			"Provide event types, e.g. --event-types 1,3,4 (1=start, 2=timeout, 3=completed, 4=failed)",
			nil,
		))
	}
	return events
}

// --- create ---
var webhookCreateCmd = &cobra.Command{
	Use:   "create",
	Short: i18n.T("webhook_create_title"),
	Long:  i18n.T("webhook_create_title"),
	Run: func(cmd *cobra.Command, args []string) {
		webhookURL, _ := cmd.Flags().GetString("webhook-url")
		eventTypesStr, _ := cmd.Flags().GetString("event-types")
		appIDStr, _ := cmd.Flags().GetString("app-id")

		client, _ := initClientWithConfig("webhook create")

		eventTypes := eventTypesOrExit(eventTypesStr, "webhook create")

		var appIDs []string
		if appIDStr != "" {
			appIDs, _ = parseStringList(appIDStr)
		}

		result, err := client.CreateWebhook(webhookURL, eventTypes, appIDs)
		if err != nil {
			formatterOut.ExitWithError(errors.ClassifyException(err, "webhook create"))
		}
		checkAPIResponse(result, "webhook create")

		formatterOut.PrintSuccess(i18n.T("webhook_created"))
		formatterOut.PrintJSON(result)
	},
}

// --- get-config ---
var webhookGetConfigCmd = &cobra.Command{
	Use:   "get-config",
	Short: i18n.T("webhook_get_config_title"),
	Long:  i18n.T("webhook_get_config_title"),
	Run: func(cmd *cobra.Command, args []string) {
		appIDStr, _ := cmd.Flags().GetString("app-id")

		client, _ := initClientWithConfig("webhook get-config")

		var appIDs []string
		if appIDStr != "" {
			appIDs, _ = parseStringList(appIDStr)
		}

		result, err := client.ListWebhooks(appIDs)
		if err != nil {
			formatterOut.ExitWithError(errors.ClassifyException(err, "webhook get-config"))
		}
		checkAPIResponse(result, "webhook get-config")

		formatterOut.PrintJSON(result)
	},
}

// --- update ---
var webhookUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: i18n.T("webhook_update_title"),
	Long:  i18n.T("webhook_update_title"),
	Run: func(cmd *cobra.Command, args []string) {
		webhookID, _ := cmd.Flags().GetString("webhook-id")
		webhookURL, _ := cmd.Flags().GetString("webhook-url")
		eventTypesStr, _ := cmd.Flags().GetString("event-types")
		appIDStr, _ := cmd.Flags().GetString("app-id")

		client, _ := initClientWithConfig("webhook update")

		eventTypes := eventTypesOrExit(eventTypesStr, "webhook update")

		var appIDs []string
		if appIDStr != "" {
			appIDs, _ = parseStringList(appIDStr)
		}

		result, err := client.UpdateWebhook(webhookID, webhookURL, eventTypes, appIDs)
		if err != nil {
			formatterOut.ExitWithError(errors.ClassifyException(err, "webhook update"))
		}
		checkAPIResponse(result, "webhook update")

		formatterOut.PrintSuccess(i18n.T("webhook_updated"))
		formatterOut.PrintJSON(result)
	},
}

// --- delete ---
var webhookDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: i18n.T("webhook_delete_title"),
	Long:  i18n.T("webhook_delete_title"),
	Run: func(cmd *cobra.Command, args []string) {
		webhookID, _ := cmd.Flags().GetString("webhook-id")

		client, _ := initClientWithConfig("webhook delete")

		result, err := client.DeleteWebhook(webhookID)
		if err != nil {
			formatterOut.ExitWithError(errors.ClassifyException(err, "webhook delete"))
		}
		checkAPIResponse(result, "webhook delete")

		formatterOut.PrintSuccess(i18n.T("webhook_deleted"))
		formatterOut.PrintJSON(result)
	},
}

// --- log ---
var webhookLogCmd = &cobra.Command{
	Use:   "log",
	Short: i18n.T("webhook_log_title"),
	Long:  i18n.T("webhook_log_title"),
	Run: func(cmd *cobra.Command, args []string) {
		webhookID, _ := cmd.Flags().GetString("webhook-id")
		startTime, _ := cmd.Flags().GetString("start-time")
		endTime, _ := cmd.Flags().GetString("end-time")

		client, _ := initClientWithConfig("webhook log")

		result, err := client.QueryWebhookLog(webhookID, startTime, endTime)
		if err != nil {
			formatterOut.ExitWithError(errors.ClassifyException(err, "webhook log"))
		}
		checkAPIResponse(result, "webhook log")

		formatterOut.PrintJSON(result)
	},
}

func init() {
	rootCmd.AddCommand(webhookCmd)

	webhookCmd.AddCommand(webhookCreateCmd)
	webhookCmd.AddCommand(webhookGetConfigCmd)
	webhookCmd.AddCommand(webhookUpdateCmd)
	webhookCmd.AddCommand(webhookDeleteCmd)
	webhookCmd.AddCommand(webhookLogCmd)

	// create flags
	webhookCreateCmd.Flags().String("webhook-url", "", i18n.T("webhook_option_webhook_url"))
	webhookCreateCmd.Flags().String("event-types", "", i18n.T("webhook_option_event_types"))
	webhookCreateCmd.Flags().String("app-id", "", i18n.T("webhook_option_app_id"))
	webhookCreateCmd.MarkFlagRequired("webhook-url")
	webhookCreateCmd.MarkFlagRequired("event-types")

	// get-config flags
	webhookGetConfigCmd.Flags().String("app-id", "", i18n.T("webhook_option_app_id"))

	// update flags
	webhookUpdateCmd.Flags().String("webhook-id", "", i18n.T("webhook_option_webhook_id"))
	webhookUpdateCmd.Flags().String("webhook-url", "", i18n.T("webhook_option_webhook_url"))
	webhookUpdateCmd.Flags().String("event-types", "", i18n.T("webhook_option_event_types"))
	webhookUpdateCmd.Flags().String("app-id", "", i18n.T("webhook_option_app_id"))
	webhookUpdateCmd.MarkFlagRequired("webhook-id")
	webhookUpdateCmd.MarkFlagRequired("webhook-url")
	webhookUpdateCmd.MarkFlagRequired("event-types")

	// delete flags
	webhookDeleteCmd.Flags().String("webhook-id", "", i18n.T("webhook_option_webhook_id"))
	webhookDeleteCmd.MarkFlagRequired("webhook-id")

	// log flags
	webhookLogCmd.Flags().String("webhook-id", "", i18n.T("webhook_option_webhook_id_optional"))
	webhookLogCmd.Flags().String("start-time", "", i18n.T("webhook_option_start_time"))
	webhookLogCmd.Flags().String("end-time", "", i18n.T("webhook_option_end_time"))
}
