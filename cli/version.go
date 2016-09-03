package cli

import (
	"fmt"

	"github.com/controlroom/lincoln/version"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Current version",
	Run:   renderVersion,
}

func renderVersion(c *cobra.Command, args []string) {
	fmt.Println(version.GetHumanVersion())
}
