package cli

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/controlroom/lincoln/config"
	"github.com/controlroom/lincoln/metadata"
	"github.com/controlroom/lincoln/sync"
	"github.com/spf13/cobra"
)

// Bootstrap lincoln app command
func init() {
	appCmd := &cobra.Command{
		Use:   "app",
		Short: "Manage running apps",
	}

	appCmd.AddCommand(appStatusCmd)
	appCmd.AddCommand(appListCmd)
	appCmd.AddCommand(appSourceCmd)
	appCmd.AddCommand(appUpCmd)
	appCmd.AddCommand(appWatchCmd)
	appCmd.AddCommand(appDownCmd)
	attachGet(appCmd)
	attachUpDev(appCmd)

	RootCmd.AddCommand(appCmd)
}

func sourcePath() string {
	return metadata.GetMeta("app:currentSource")
}

func rpcClient() *sync.Client {
	client, err := sync.NewClient("localhost:9876", time.Millisecond*500)

	if err != nil {
		panic(err)
	}

	return client
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
	Use:   "up-dev [appName] (nodeName)",
	Short: "Deploy app locally in development mode (requires local copy)",
	RunE:  appUpDev,
}

func attachUpDev(appCmd *cobra.Command) {
	appCmd.AddCommand(appUpDevCmd)
}

func appUpDev(c *cobra.Command, args []string) error {
	if sourcePath() == "" {
		return errors.New("Missing source path")
	} else if len(args) == 0 {
		return errors.New("Missing app name")
	}

	app, err := config.FindLocalApp(sourcePath(), args[0])
	if err != nil {
		return err
	}

	backend.SetupSync(app)
	rpcClient().Watch(app.Config.Name, app.Path)
	return nil
}

// ===  Watch  ==================================================================
//
var appWatchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch for changes in projects that are in development mode",
	RunE:  appWatch,
}

func appWatch(c *cobra.Command, args []string) error {
	sync.StartServer(9876)
	return nil
}

// ===  Up  ==================================================================
//
var appUpCmd = &cobra.Command{
	Use:   "up [name] [sha || branch]",
	Short: "Deploy app locally from built images",
	RunE:  appUp,
}

func appUp(c *cobra.Command, args []string) error {
	return nil
}

// ===  Down  ===================================================================
//
var appDownCmd = &cobra.Command{
	Use:   "down [name]",
	Short: "Remove appliction from stack",
	RunE:  appDown,
}

func appDown(c *cobra.Command, args []string) error {
	if sourcePath() == "" {
		return errors.New("Missing source path")
	} else if len(args) == 0 {
		return errors.New("Missing app name")
	}

	app, err := config.FindLocalApp(sourcePath(), args[0])
	if err != nil {
		return err
	}

	rpcClient().UnWatch(app.Config.Name)
	return nil
}

// ===  List  ===================================================================
//
var appListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all applications",
	Run:   appList,
}

func appList(c *cobra.Command, args []string) {
	fmt.Println("-- App list")
	if sourcePath() != "" {
		apps := config.FindAllLocalApps(sourcePath())

		table := getTable("Name", "Local", "Branch", "Description")
		for _, app := range apps {
			table.AppendLine(
				app.Config.Name,
				yellowOut("*"),
				app.Branch,
				app.Config.Description,
			)
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
	Short: "Git clone application into source directory",
	Long: `
Pull project source code into local destination. The path you pass into
destination should be the root folder for all your applications. Lincoln will
create a project folder and clone the project for you. Lincoln will also
remember where you cloned the project, so you can swap in a development version
to a stack with ease.
	`,
	RunE: appGet,
}

var appGetAll bool

func attachGet(appCmd *cobra.Command) {
	appGetCmd.Flags().BoolVar(&appGetAll, "all", false, "Clone all apps")
	appCmd.AddCommand(appGetCmd)
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
