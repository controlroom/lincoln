package cli

import "github.com/spf13/cobra"

func init() {
	appCmd.AddCommand(appStatusCmd)
	RootCmd.AddCommand(appCmd)
}

// ===  Base Command  ===========================================================
//
var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Manage running apps",
}

// ===  Status  ==================================================================
//
var appStatusCmd = &cobra.Command{
	Use:   "status [name]",
	Short: "App status",
	Run:   appStatus,
}

func appStatus(c *cobra.Command, args []string) {
}
