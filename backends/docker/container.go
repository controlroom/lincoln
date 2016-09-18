package docker

import (
	"fmt"
	"io"
	"os"

	"github.com/controlroom/lincoln/interfaces"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"

	"github.com/docker/go-connections/nat"
)

func (op DockerOperation) RemoveContainer(container *interfaces.Container) {
	if container == nil {
		return
	}

	client.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{
		Force: true,
	})
}

func hasContainer(name string) bool {
	containerList, _ := client.ContainerList(ctx, types.ContainerListOptions{})
	for _, v := range containerList {
		if v.Names[0] == name {
			return true
		}
	}
	return false
}

func Containers(containers []types.Container) []interfaces.Container {
	if len(containers) == 0 {
		return nil
	}

	toContainers := make([]interfaces.Container, len(containers))

	for i, container := range containers {
		toContainers[i] = interfaces.Container{
			ID:      container.ID,
			Name:    container.Names[0],
			Running: container.State == "running",
		}
	}

	return toContainers
}

func (op DockerOperation) ListContainers() []interfaces.Container {
	options := types.ContainerListOptions{All: false}
	containers, err := client.ContainerList(ctx, options)
	if err != nil {
		panic(err)
	}

	return Containers(containers)
}

func (op DockerOperation) FindContainers(flags []map[string]string) []interfaces.Container {
	args := filters.NewArgs()
	for _, flag := range flags {
		for k, v := range flag {
			args.Add(k, v)
		}
	}

	options := types.ContainerListOptions{
		All:    true,
		Filter: args,
	}

	containers, err := client.ContainerList(ctx, options)
	if err != nil {
		panic(err)
	}

	return Containers(containers)
}

func (op DockerOperation) FindContainer(flags []map[string]string) *interfaces.Container {
	curr := op.FindContainers(flags)

	if curr == nil {
		return nil
	} else {
		return &curr[0]
	}
}

func (op DockerOperation) FindContainerByName(name string) *interfaces.Container {
	return op.FindContainer([]map[string]string{
		map[string]string{"name": name},
	})
}

func (op DockerOperation) InspectContainer(container *interfaces.Container) interfaces.ContainerInfo {
	data, err := client.ContainerInspect(ctx, container.ID)
	if err != nil {
		panic(err)
	}

	fmt.Println(data.NetworkSettings.Networks)
	return interfaces.ContainerInfo{}
}

// ===  Running Container  ======================================================
//
//

func buildPorts(strPorts []string) map[nat.Port]struct{} {
	ports := make(map[nat.Port]struct{}, len(strPorts))
	for _, port := range strPorts {
		var v struct{}
		ports[nat.Port(port)] = v
	}
	return ports
}

func buildPortBindings(strPortBindings []string) nat.PortMap {
	portBindings := make(nat.PortMap)
	for _, port := range strPortBindings {
		mapping, _ := nat.ParsePortSpec(port)
		portBindings[mapping[0].Port] = []nat.PortBinding{mapping[0].Binding}
	}
	return portBindings
}

func buildNetworkConfig(opts interfaces.ContainerStartOptions) *network.NetworkingConfig {
	return &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			opts.Stack.Name: &network.EndpointSettings{},
		},
	}
}

func buildHostConfig(opts interfaces.ContainerStartOptions) *container.HostConfig {
	return &container.HostConfig{
		CapAdd:       strslice.StrSlice(opts.CapAdd),
		Binds:        opts.Volumes,
		VolumesFrom:  opts.VolumesFrom,
		AutoRemove:   true,
		PortBindings: buildPortBindings(opts.PortBindings),
	}
}

func buildContainerConfig(opts interfaces.ContainerStartOptions) *container.Config {
	return &container.Config{
		Image:        opts.Image,
		WorkingDir:   "/src",
		Env:          opts.Env,
		Cmd:          strslice.StrSlice(opts.Cmd),
		ExposedPorts: buildPorts(opts.Ports),
		Labels: map[string]string{
			"RFStack": opts.Stack.Name,
		},
	}
}

func (op DockerOperation) StartContainer(opts interfaces.ContainerStartOptions) *interfaces.Container {
	getImage(opts.Image)

	curr := op.FindContainerByName(opts.Name)

	if curr != nil {
		fmt.Println("Proc already running, shutting down")

		client.ContainerRemove(ctx, curr.ID, types.ContainerRemoveOptions{
			Force: true,
		})
	}

	createRes, err := client.ContainerCreate(
		ctx,
		buildContainerConfig(opts),
		buildHostConfig(opts),
		buildNetworkConfig(opts),
		opts.Name,
	)

	if err != nil {
		panic(err)
	}

	err = client.ContainerStart(
		ctx,
		createRes.ID,
		types.ContainerStartOptions{},
	)

	if err != nil {
		client.ContainerRemove(ctx, createRes.ID, types.ContainerRemoveOptions{
			Force: true,
		})

		panic(err)
	}

	fmt.Printf("Booted container: %s\n", opts.Name)

	return &interfaces.Container{
		ID:      createRes.ID,
		Name:    opts.Name,
		Running: true,
	}
}

func (op DockerOperation) RunContainer(opts interfaces.ContainerStartOptions) *interfaces.Container {
	getImage(opts.Image)

	containerConfig := buildContainerConfig(opts)

	containerConfig.AttachStdin = true
	containerConfig.AttachStdout = true
	containerConfig.AttachStderr = true
	containerConfig.OpenStdin = true
	containerConfig.StdinOnce = true
	containerConfig.Tty = true

	createRes, err := client.ContainerCreate(
		ctx,
		containerConfig,
		buildHostConfig(opts),
		buildNetworkConfig(opts),
		opts.Name,
	)

	if err != nil {
		panic(err)
	}

	options := types.ContainerAttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	}

	attachRes, err := client.ContainerAttach(
		ctx,
		createRes.ID,
		options,
	)

	defer attachRes.Close()

	finished := make(chan struct{})

	in := NewInStream(os.Stdin)

	if err = in.SetRawTerminal(); err != nil {
		panic(err)
	}

	out := NewOutStream(os.Stdout)
	if err = out.SetRawTerminal(); err != nil {
		panic(err)
	}

	go func() {
		io.Copy(out, attachRes.Reader)
		close(finished)
	}()

	go func() {
		io.Copy(attachRes.Conn, in)
		fmt.Println("finished stdin")
	}()

	if err != nil {
		panic(err)
	}

	err = client.ContainerStart(
		ctx,
		createRes.ID,
		types.ContainerStartOptions{},
	)

	if err != nil {
		panic(err)
	}

	<-finished

	in.RestoreTerminal()
	out.RestoreTerminal()

	client.ContainerRemove(ctx, createRes.ID, types.ContainerRemoveOptions{})

	return nil
}
