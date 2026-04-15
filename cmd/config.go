package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/laiye-ai/adp-cli/internal/config"
	"github.com/laiye-ai/adp-cli/internal/errors"
	"github.com/laiye-ai/adp-cli/internal/formatter"
	"github.com/laiye-ai/adp-cli/internal/i18n"
)

var formatterOut = formatter.New(false, false)

// configCmd represents the config command group
var configCmd = &cobra.Command{
	Use:   "config",
	Short: i18n.T("config_description"),
	Long:  i18n.T("config_description"),
}

// configSetCmd represents the config set command
var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: i18n.T("config_set_title"),
	Long:  i18n.T("config_set_title"),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, _ := cmd.Flags().GetString("api-key")
		apiBaseURL, _ := cmd.Flags().GetString("api-base-url")

		if apiKey == "" && apiBaseURL == "" {
			cliErr := errors.NewCLIError(
				i18n.T("error_api_key_or_url_required"),
				errors.ErrorTypeParam,
				errors.ExitParameterError,
				false,
				"Provide at least --api-key or --api-base-url.",
				nil,
			)
			formatterOut.ExitWithError(cliErr)
		}

		if apiKey != "" {
			if err := config.SetAPIKey(apiKey); err != nil {
				cliErr := errors.ClassifyException(err, "config set")
				formatterOut.ExitWithError(cliErr)
			}
			formatterOut.PrintSuccess(i18n.T("api_key_configured"))
		}

		if apiBaseURL != "" {
			if err := config.SetAPIBaseURL(apiBaseURL); err != nil {
				cliErr := errors.ClassifyException(err, "config set")
				formatterOut.ExitWithError(cliErr)
			}
			formatterOut.PrintSuccess(i18n.T("api_base_url_configured"))
		}
	},
}

// configGetCmd represents the config get command
var configGetCmd = &cobra.Command{
	Use:   "get",
	Short: i18n.T("config_get_title"),
	Long:  i18n.T("config_get_title"),
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.Load()
		if err != nil {
			cliErr := errors.ClassifyException(err, "config get")
			formatterOut.ExitWithError(cliErr)
		}

		summary := config.GetConfigSummary(cfg)
		formatterOut.PrintJSON(summary)
	},
}

// configClearCmd represents the config clear command
var configClearCmd = &cobra.Command{
	Use:   "clear",
	Short: i18n.T("config_clear_title"),
	Long:  i18n.T("config_clear_title"),
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")

		// Skip confirmation in non-TTY environment
		if !force && formatterOut.IsTTY() {
			fmt.Print(i18n.T("confirm_clear_config") + " [y/N] ")
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				return
			}
		}

		if err := config.Clear(); err != nil {
			cliErr := errors.ClassifyException(err, "config clear")
			formatterOut.ExitWithError(cliErr)
		}
		formatterOut.PrintSuccess(i18n.T("config_cleared"))
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configGetCmd)
	configCmd.AddCommand(configClearCmd)

	// config set flags
	configSetCmd.Flags().String("api-key", "", i18n.T("option_api_key"))
	configSetCmd.Flags().String("api-base-url", "", i18n.T("option_api_base_url"))

	// config clear flags
	configClearCmd.Flags().BoolP("force", "y", false, i18n.T("option_force_clear"))
}
