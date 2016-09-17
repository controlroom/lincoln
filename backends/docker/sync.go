package docker

import (
	"fmt"
	"net"
	"os"
	"os/exec"

	"github.com/controlroom/lincoln/config"
	"github.com/controlroom/lincoln/interfaces"
	"github.com/controlroom/lincoln/metadata"
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
			Stack:        op.GetDefaultStack(),
			Env:          []string{"VOLUME=/src"},
			Volumes:      []string{fmt.Sprintf("%v/.ssh/:/root/.ssh/", homeDir)},
			PortBindings: []string{fmt.Sprintf("%v:873", port)},
		})
		metadata.AppNS(app.Config.Name).Put("syncPort", fmt.Sprintf("%v", port))

		fmt.Println("Syncing source")
		op.Sync(app.Config.Name, app.Path, true)
	}
}

var rsyncOpts []string = []string{
	"-aizP", "--delete", "--exclude=log", "--exclude=tmp", "--exclude=.git",
}

func (op DockerOperation) Sync(app string, path string, quiet bool) {
	var opts string

	if quiet {
		opts = "-q"
	}

	port := metadata.AppNS(app).Get("syncPort")
	path = fmt.Sprintf("%v/.", path)
	uri := fmt.Sprintf("rsync://localhost:%v/volume/.", port)
	rsyncArgs := append(rsyncOpts, opts, path, uri)
	cmd := exec.Command("rsync", rsyncArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
