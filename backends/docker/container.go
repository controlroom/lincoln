package docker

import (
	"fmt"

	"github.com/controlroom/lincoln/interfaces"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/filters"
	"github.com/docker/engine-api/types/network"
	"github.com/docker/engine-api/types/strslice"
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
			ID:   container.ID,
			Name: container.Names[0],
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
		All:    false,
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

func (op DockerOperation) StartContainer(opts interfaces.ContainerStartOptions) *interfaces.Container {
	var containerId string

	resp, err := client.ImagePull(ctx, opts.Image, types.ImagePullOptions{})
	defer resp.Close()
	jsonmessage.DisplayJSONMessagesStream(resp, os.Stderr, os.Stderr.Fd(), true, nil)

	if err != nil {
		panic(err)
	}

	ports := make(map[nat.Port]struct{}, len(opts.Ports))
	for _, port := range opts.Ports {
		var v struct{}
		ports[nat.Port(port)] = v
	}

	portBindings := make(nat.PortMap)
	for _, port := range opts.PortBindings {
		mapping, _ := nat.ParsePortSpec(port)
		portBindings[mapping[0].Port] = []nat.PortBinding{mapping[0].Binding}
	}

	containerConfig := &container.Config{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Image:        opts.Image,
		Env:          opts.Env,
		OpenStdin:    true,
		Cmd:          strslice.StrSlice(opts.Cmd),
		ExposedPorts: ports,
		Labels: map[string]string{
			"RFStack": opts.Stack.Name,
		},
	}

	hostConfig := &container.HostConfig{
		CapAdd:       strslice.StrSlice(opts.CapAdd),
		Binds:        opts.Volumes,
		PortBindings: portBindings,
	}

	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			opts.Stack.Name: &network.EndpointSettings{},
		},
	}

	res, err := client.ContainerCreate(ctx, containerConfig, hostConfig, networkConfig, opts.Name)
	if err != nil {
		panic(err)
	}

	err = client.ContainerStart(
		ctx,
		res.ID,
		types.ContainerStartOptions{},
	)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Booted container: %s\n", opts.Name)
	containerId = res.ID

	return &interfaces.Container{
		ID:   containerId,
		Name: opts.Name,
	}
}
