package cmd

import (
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	"github.com/laiye-ai/adp-cli/internal/api"
	"github.com/laiye-ai/adp-cli/internal/config"
	"github.com/laiye-ai/adp-cli/internal/errors"
	"github.com/laiye-ai/adp-cli/internal/i18n"
)

// appIDCmd represents the app-id command group
var appIDCmd = &cobra.Command{
	Use:   "app-id",
	Short: i18n.T("app_id_description"),
	Long:  i18n.T("app_id_description"),
}

// listAppsCmd represents the app-id list command
var listAppsCmd = &cobra.Command{
	Use:   "list",
	Short: i18n.T("app_id_list_title"),
	Long:  i18n.T("app_id_list_title"),
	Run: func(cmd *cobra.Command, args []string) {
		appLabel, _ := cmd.Flags().GetString("app-label")
		appType, _ := cmd.Flags().GetInt("app-type")
		limit, _ := cmd.Flags().GetInt("limit")

		cfg, err := config.Load()
		if err != nil {
			cliErr := errors.ClassifyException(err, "app-id list")
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
			cliErr := errors.ClassifyException(err, "app-id list")
			formatterOut.ExitWithError(cliErr)
		}

		var appTypePtr *int
		if appType > 0 {
			appTypePtr = &appType
		}

		apps, err := client.ListApps(appTypePtr, limit)
		if err != nil {
			cliErr := errors.ClassifyException(err, "app-id list")
			formatterOut.ExitWithError(cliErr)
		}

		// Cache all apps before filtering
		if len(apps) > 0 {
			cacheData := map[string]interface{}{"apps": apps}
			cacheBytes, _ := json.Marshal(cacheData)
			os.WriteFile(config.GetCachePath(), cacheBytes, 0644)
		}

		// Filter by app_label if provided
		if appLabel != "" && len(apps) > 0 {
			var filtered []map[string]interface{}
			for _, app := range apps {
				if labels, ok := app["app_label"].([]interface{}); ok {
					for _, label := range labels {
						if labelStr, ok := label.(string); ok {
							if labelStr == appLabel {
								filtered = append(filtered, app)
								break
							}
						}
					}
				}
			}
			apps = filtered
		}

		if len(apps) == 0 {
			formatterOut.PrintJSON(map[string]interface{}{
				"apps":    []interface{}{},
				"message": i18n.T("no_applications_found"),
			})
			return
		}

		formatterOut.PrintJSON(map[string]interface{}{"apps": apps})
	},
}

// listAppsCacheCmd represents the app-id cache command
var listAppsCacheCmd = &cobra.Command{
	Use:   "cache",
	Short: i18n.T("app_id_list_cache_title"),
	Long:  i18n.T("app_id_list_cache_title"),
	Run: func(cmd *cobra.Command, args []string) {
		// Read from cache file
		cachePath := config.GetCachePath()
		data, err := os.ReadFile(cachePath)
		if err != nil {
			if os.IsNotExist(err) {
				formatterOut.PrintJSON(map[string]interface{}{
					"apps":    []interface{}{},
					"message": i18n.T("no_applications_found"),
				})
				return
			}
			cliErr := errors.ClassifyException(err, "app-id cache")
			formatterOut.ExitWithError(cliErr)
		}

		var cache struct {
			Apps []map[string]interface{} `json:"apps"`
		}
		if err := json.Unmarshal(data, &cache); err != nil {
			cliErr := errors.ClassifyException(err, "app-id cache")
			formatterOut.ExitWithError(cliErr)
		}

		formatterOut.PrintJSON(map[string]interface{}{"apps": cache.Apps})
	},
}

func init() {
	rootCmd.AddCommand(appIDCmd)

	appIDCmd.AddCommand(listAppsCmd)
	appIDCmd.AddCommand(listAppsCacheCmd)

	// list flags
	listAppsCmd.Flags().String("app-label", "", i18n.T("app_id_list_app_label"))
	listAppsCmd.Flags().Int("app-type", 0, i18n.T("app_id_list_app_type"))
	listAppsCmd.Flags().Int("limit", 120, i18n.T("app_id_list_limit"))
}
