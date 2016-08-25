package cli

import "github.com/spf13/cobra"

func init() {
	appCmd.AddCommand(appListCmd)
	RootCmd.AddCommand(appCmd)
}

// ===  Base Command  ===========================================================
//
var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Deploy and manage Apps",
}

// ===  List  ===================================================================
//
var appListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all running apps",
	Run:   appList,
}

func appList(c *cobra.Command, args []string) {
}
