package docker

import (
	"fmt"
	"os"

	"github.com/controlroom/lincoln/interfaces"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/container"
	"github.com/docker/engine-api/types/filters"
	"github.com/docker/engine-api/types/network"
)

func (op DockerOperation) RemoveContainer(container interfaces.Container) {
	client.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{
		Force: true,
	})
}

func Containers(containers []types.Container) []interfaces.Container {
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

func (op DockerOperation) StartContainer(opts interfaces.ContainerStartOptions) {
	resp, err := client.ImagePull(ctx, opts.Image, types.ImagePullOptions{})
	defer resp.Close()
	jsonmessage.DisplayJSONMessagesStream(resp, os.Stderr, os.Stderr.Fd(), true, nil)

	if err != nil {
		panic(err)
	}

	containerConfig := &container.Config{
		AttachStdin:  true,
		AttachStdout: true,
		AttachStderr: true,
		Image:        opts.Image,
		OpenStdin:    true,
		Labels: map[string]string{
			"RFStack": opts.Stack.Name,
		},
	}

	hostConfig := &container.HostConfig{
		Binds: opts.Volumes,
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

	fmt.Printf("Booted container: %s\n", opts.Name)

	if err != nil {
		panic(err)
	}

	if opts.Cmd != nil {
		fmt.Printf("Running cmd on container: %s\n", opts.Name)
		exec, err := client.ContainerExecCreate(
			ctx,
			res.ID,
			types.ExecConfig{
				Detach: true,
				Cmd:    opts.Cmd,
			},
		)

		if err != nil {
			panic(err)
		}

		err = client.ContainerExecStart(
			ctx,
			exec.ID,
			types.ExecStartCheck{},
		)

		if err != nil {
			panic(err)
		}
	}
}
