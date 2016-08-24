package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var stackCmd = &cobra.Command{
	Use:   "stack",
	Short: "Operations for manipulating stacks",
}

var stackCreateCmd = &cobra.Command{
	Use:   "create network-name",
	Short: "Create stack",
	RunE:  createStack,
}

func createStack(c *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("Missing network name")
	}

	backend.CreateStack(args[0])
	return nil
}

var stackListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all stacks",
	Run:   listStacks,
}

func listStacks(c *cobra.Command, args []string) {
	fmt.Println(backend.ListStacks())
}

func init() {
	stackCmd.AddCommand(stackCreateCmd)
	stackCmd.AddCommand(stackListCmd)

	RootCmd.AddCommand(stackCmd)
}
