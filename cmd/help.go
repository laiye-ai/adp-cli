package cmd

import (
	"github.com/spf13/cobra"
	"github.com/laiye-ai/adp-cli/internal/i18n"
)

// helpCmd represents the help command
var helpCmd = &cobra.Command{
	Use:   "help",
	Short: i18n.T("help_description"),
	Long:  i18n.T("help_description"),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			// Show help for specific subcommand
			foundCmd, _, err := cmd.Root().Find(args)
			if err == nil {
				foundCmd.Help()
			} else {
				cmd.Root().Help()
			}
		} else {
			cmd.Root().Help()
		}
	},
}

func init() {
	// Replace cobra's built-in help with our i18n-aware version
	rootCmd.SetHelpCommand(helpCmd)
}