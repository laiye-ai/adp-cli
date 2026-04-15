package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/laiye-ai/adp-cli/internal/api"
	"github.com/laiye-ai/adp-cli/internal/config"
	"github.com/laiye-ai/adp-cli/internal/errors"
	"github.com/laiye-ai/adp-cli/internal/i18n"
)

// creditCmd represents the credit command
var creditCmd = &cobra.Command{
	Use:   "credit",
	Short: i18n.T("credit_description"),
	Long:  i18n.T("credit_description"),
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, _ := cmd.Flags().GetString("api-key")

		cfg, err := loadConfigWithOverride(apiKey)
		if err != nil {
			cliErr := errors.ClassifyException(err, "credit")
			formatterOut.ExitWithError(cliErr)
		}

		if !config.IsConfigured(cfg) {
			cliErr := errors.NewCLIError(
				i18n.T("error_not_configured"),
				errors.ErrorTypeSystem,
				errors.ExitGeneralError,
				false,
				"Run 'adp config set' to configure your API key and base URL.",
				nil,
			)
			formatterOut.ExitWithError(cliErr)
		}

		client, err := api.NewClient(cfg)
		if err != nil {
			cliErr := errors.ClassifyException(err, "credit")
			formatterOut.ExitWithError(cliErr)
		}

		result, err := client.GetUserPaymentStatus()
		if err != nil {
			cliErr := errors.ClassifyException(err, "credit")
			formatterOut.ExitWithError(cliErr)
		}

		// Get credit from response (remaining_credits is at root level)
		credit := 0
		if c, ok := result["remaining_credits"].(float64); ok {
			credit = int(c)
		}

		portalURL := cfg.APIBaseURL

		if jsonMode {
			formatterOut.PrintJSON(map[string]interface{}{
				"credit":     credit,
				"portal_url": portalURL,
			})
		} else {
			formatterOut.PrintSection(i18n.T("credit_info"))
			formatterOut.PrintInfo(fmt.Sprintf("%s: %d", i18n.T("remaining_credits"), credit))
			formatterOut.PrintInfo(fmt.Sprintf("%s: %s", i18n.T("recharge_message"), portalURL))
		}
	},
}

func init() {
	rootCmd.AddCommand(creditCmd)
	creditCmd.Flags().String("api-key", "", i18n.T("credit_api_key"))
}
