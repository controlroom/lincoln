package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/controlroom/lincoln/metadata"
	"github.com/spf13/cobra"
)

var appGetAll bool

func init() {
	appGetCmd.Flags().BoolVar(&appGetAll, "all", false, "Clone all apps")

	appCmd.AddCommand(appStatusCmd)
	appCmd.AddCommand(appGetCmd)
	appCmd.AddCommand(appListCmd)
	appCmd.AddCommand(appSourceCmd)
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
var appSourceCmd = &cobra.Command{
	Use:   "source [path]",
	Short: "Get/set project source directory",
	Run:   appSource,
}

func exists(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		return err
	}
	return nil
}

func appSource(c *cobra.Command, args []string) {
	if len(args) == 0 {
		source := metadata.GetMeta("app:currentSource")
		if source == "" {
			fmt.Println("Source not set, please pass source path")
		} else {
			fmt.Printf("Source: %v", source)
		}
	} else {
		source, err := filepath.Abs(args[0])
		if err != nil {
			panic(err)
		}

		if err = exists(source); err == nil {
			metadata.PutMeta("app:currentSource", source)
			fmt.Printf("Source set: %v", source)
		} else {
			fmt.Printf("%v is not a directory", source)
		}
	}
}

// ===  Get  =================================================================
//
var appGetCmd = &cobra.Command{
	Use:   "get [apps]",
	Short: "Git clone application into projects folder",
	Long: `
Pull project source code into local destination. The path you pass into
destination should be the root folder for all your applications. Lincoln will
create a project folder and clone the project for you. Lincoln will also
remember where you cloned the project, so you can swap in a development version
to a stack with ease.
	`,
	RunE: appGet,
}

func appGet(c *cobra.Command, args []string) error {
	var apps []string

	if appGetAll == false && len(args) == 0 {
		return errors.New("Missing app name")
	}

	if appGetAll {
		apps = []string{"__all"}
	} else {
		apps = args
	}

	fmt.Println(apps)

	return nil
}
