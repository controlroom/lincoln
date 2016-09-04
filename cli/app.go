package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/controlroom/lincoln/config"
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
	appCmd.AddCommand(appUpDevCmd)
	appCmd.AddCommand(appUpCmd)
	RootCmd.AddCommand(appCmd)
}

func sourcePath() string {
	return metadata.GetMeta("app:currentSource")
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
	Short: "Status of loaded apps",
	Run:   appStatus,
}

func appStatus(c *cobra.Command, args []string) {
}

// ===  UpDev  ==================================================================
//
var appUpDevCmd = &cobra.Command{
	Use:   "up-dev [name]",
	Short: "Deploy app locally in development mode (requires local copy)",
	Run:   appUpDev,
}

func appUpDev(c *cobra.Command, args []string) {
}

// ===  Up  ==================================================================
//
var appUpCmd = &cobra.Command{
	Use:   "up [name] [sha || branch]",
	Short: "Deploy app locally from built images",
	Run:   appUp,
}

func appUp(c *cobra.Command, args []string) {
}

// ===  List  ===================================================================
//
var appListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all applications",
	Run:   appList,
}

func appList(c *cobra.Command, args []string) {
	if sourcePath() != "" {
		apps := config.FindAllLocalApps(sourcePath())

		table := getTable("Name", "Local", "Branch")
		for _, app := range apps {
			table.AppendLine(app.Config.Name, yellowOut("*"), app.Branch)
		}
		table.Render()
	}
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
		if sourcePath() == "" {
			fmt.Println("Source not set, please pass source path")
		} else {
			fmt.Printf("Source: %v", sourcePath())
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
