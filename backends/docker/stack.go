package docker

import (
	"fmt"

	"github.com/controlroom/lincoln/interfaces"
	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
)

type DockerStack struct {
	Name string
	ID   string
}

func containsNetwork(name string) bool {
	netList, _ := client.NetworkList(ctx, types.NetworkListOptions{})
	for _, v := range netList {
		if v.Name == name {
			return true
		}
	}
	return false
}

func (op DockerOperation) CreateStack(name string) {
	if !containsNetwork(name) {
		stack, _ := client.NetworkCreate(ctx, name, types.NetworkCreate{
			CheckDuplicate: true,
			Labels:         map[string]string{"RFStack": "true"},
		})

		fmt.Printf("Created stack: %s\n", name)
		fmt.Println("Attaching default proxy container")

		op.StartContainer(interfaces.ContainerStartOptions{
			Name:  fmt.Sprintf("%s-root-proxy", name),
			Image: "jwilder/nginx-proxy",
			Stack: interfaces.Stack{
				Name: name,
				ID:   stack.ID,
			},
			Volumes: []string{"/var/run/docker.sock:/tmp/docker.sock:ro"},
		})

		fmt.Printf("Attached default proxy container: %s\n", name)
	} else {
		fmt.Printf("Stack already created: %s\n", name)
	}
}

func (op DockerOperation) DestroyStack(name string) {
	containers := op.FindContainers([]map[string]string{
		map[string]string{"label": fmt.Sprintf("RFStack=%s", name)},
	})

	for _, container := range containers {
		fmt.Printf("Removing container: %s\n", container.Name)
		op.RemoveContainer(container)
	}

	err := client.NetworkRemove(ctx, name)

	if err != nil {
		panic(err)
	}

	fmt.Printf("%s stack destroyed", name)
}

func (op DockerOperation) ListStacks() []interfaces.Stack {
	rfStacksFilter, _ := filters.ParseFlag("label=RFStack=true", filters.NewArgs())

	networks, _ := client.NetworkList(ctx, types.NetworkListOptions{
		Filters: rfStacksFilter,
	})

	stacks := make([]interfaces.Stack, len(networks))
	for i, network := range networks {
		stacks[i] = interfaces.Stack{
			Name: network.Name,
			ID:   network.ID,
		}
	}

	return stacks
}

func (op DockerOperation) FindStacks(flags []map[string]string) []interfaces.Stack {
	args := filters.NewArgs()
	for _, flag := range flags {
		for k, v := range flag {
			args.Add(k, v)
		}
	}

	curr, _ := client.NetworkList(
		ctx,
		types.NetworkListOptions{Filters: args},
	)

	if len(curr) == 0 {
		return nil
	} else {
		stacks := make([]interfaces.Stack, len(curr))

		for i, network := range curr {
			stacks[i] = interfaces.Stack{
				Name: network.Name,
				ID:   network.ID,
			}
		}

		return stacks
	}
}

func (op DockerOperation) FindStack(flags []map[string]string) *interfaces.Stack {
	curr := op.FindStacks(flags)

	if curr == nil {
		return nil
	} else {
		return &curr[0]
	}
}

func (op DockerOperation) GetDefaultStack() *interfaces.Stack {
	id := getMeta("stack:currentID")

	if id == "" {
		return nil
	}

	return op.FindStack([]map[string]string{
		map[string]string{"id": id},
	})
}

func (op DockerOperation) SetDefaultStack(name string) error {
	stack := op.FindStack([]map[string]string{
		map[string]string{"name": name},
	})

	if stack == nil {
		fmt.Println("Not a stack, create stack then set as default")
		return nil
	}

	putMeta("stack:currentID", stack.ID)
	return nil
}
