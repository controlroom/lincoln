package cli

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	stackCmd.AddCommand(stackCreateCmd)
	stackCmd.AddCommand(stackListCmd)
	stackCmd.AddCommand(stackUseCmd)
	stackCmd.AddCommand(stackCurrentCmd)

	RootCmd.AddCommand(stackCmd)
}

// ===  Base Command  ===========================================================
//
var stackCmd = &cobra.Command{
	Use:   "stack",
	Short: "Operations for manipulating stacks",
}

// ===  Create  =================================================================
//
var stackCreateCmd = &cobra.Command{
	Use:   "create (name)",
	Short: "Create stack",
	RunE:  createStack,
}

func createStack(c *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("Missing stack name")
	}

	backend.CreateStack(args[0])
	return nil
}

// ===  List  ===================================================================
//
var stackListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all stacks",
	Run:   listStacks,
}

func listStacks(c *cobra.Command, args []string) {
	fmt.Println(backend.ListStacks())
}

// ===  Current  ================================================================
var stackCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Show default stack",
	Run:   currentStack,
}

func currentStack(c *cobra.Command, args []string) {
	stack := backend.GetDefaultStack()

	if stack == nil {
		fmt.Println("You have yet to select a default stack")
	} else {
		fmt.Println(stack.Name)
	}
}

// ===  Use  ====================================================================
//
var stackUseCmd = &cobra.Command{
	Use:   "use (name)",
	Short: "Set default stack",
	RunE:  useStack,
}

func useStack(c *cobra.Command, args []string) error {
	if len(args) == 0 {
		return errors.New("Missing stack name")
	}

	backend.SetDefaultStack(args[0])
	return nil
}
