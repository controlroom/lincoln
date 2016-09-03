package cli

import (
	"errors"

	"github.com/spf13/cobra"
)

func init() {
	appCmd.AddCommand(appStatusCmd)
	appCmd.AddCommand(appSourceCmd)
	appCmd.AddCommand(appListCmd)
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

// ===  List  ===================================================================
//
var appListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all applications",
	Run:   appList,
}

func appList(c *cobra.Command, args []string) {
}

// ===  Source  =================================================================
//
var appSourceCmd = &cobra.Command{
	Use:   "source [name] [destination]",
	Short: "Pull git source for application",
	Long: `
Pull project source code into local destination. The path you pass into
destination should be the root folder for all your applications. Lincoln will
create a project folder and clone the project for you. Lincoln will also
remember where you cloned the project, so you can swap in a development version
to a stack with ease.
	`,
	RunE: appSource,
}

func appSource(c *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("Missing app name & destination")
	} else if len(args) == 1 {
		return errors.New("Missing destination")
	}

	return nil
}
