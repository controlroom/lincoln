package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/controlroom/lincoln/backends/docker"
	"github.com/controlroom/lincoln/interfaces"
	"github.com/spf13/cobra"
)

var backendName string
var backend interfaces.Operation

var RootCmd = &cobra.Command{
	Use:               "lincoln",
	Short:             "Microservices orchestrator",
	PersistentPreRunE: determineBackend,
}

func determineBackend(cmd *cobra.Command, args []string) error {
	switch backendName {
	case "docker":
		backend = docker.DockerOperation{}
		return nil
	}

	return errors.New("Invalid backend")
}

func Execute() {
	RootCmd.PersistentFlags().StringVar(&backendName, "backend", "docker", "Backend (docker, kube)")

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
