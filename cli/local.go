package cli

import "github.com/spf13/cobra"

func init() {
	localCmd.AddCommand(localBootCmd)
	RootCmd.AddCommand(localCmd)
}

// ===  Base Command  ===========================================================
//
var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Dev mode",
}

// ===  Boot  ===================================================================
//
var localBootCmd = &cobra.Command{
	Use:   "bootstrap (path)",
	Short: "Pull down all dependencies",
	Run:   localBoot,
}

func localBoot(c *cobra.Command, args []string) {
}
