package docker

import (
	"fmt"
	"net"
	"os"

	"github.com/controlroom/lincoln/config"
	"github.com/controlroom/lincoln/interfaces"
)

func freePort() int {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

// Startup containers required for syncing local application source code.
func (op DockerOperation) SetupSync(app *config.App) {
	devSyncName := fmt.Sprintf("%v-dev-sync", app.Config.Name)
	sync := op.FindContainerByName(devSyncName)

	if sync == nil {
		port := freePort()
		homeDir := os.Getenv("HOME")
		sync = op.StartContainer(interfaces.ContainerStartOptions{
			Name:         devSyncName,
			Image:        "sync:0.5.0",
			Stack:        *op.GetDefaultStack(),
			Env:          []string{"VOLUME=/src"},
			Volumes:      []string{fmt.Sprintf("%v/.ssh/:/root/.ssh/", homeDir)},
			PortBindings: []string{fmt.Sprintf("%v:873", port)},
		})
	}
}
