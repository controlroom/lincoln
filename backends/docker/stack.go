package docker

import (
	"fmt"

	"github.com/controlroom/lincoln/interfaces"
	"github.com/controlroom/lincoln/metadata"

	"github.com/docker/engine-api/types"
	"github.com/docker/engine-api/types/filters"
)

type DockerStack struct {
	Name string
	ID   string
}

func hasNetwork(name string) bool {
	netList, _ := client.NetworkList(ctx, types.NetworkListOptions{})
	for _, v := range netList {
		if v.Name == name {
			return true
		}
	}
	return false
}

func createNetwork(name string, isStack bool) types.NetworkCreateResponse {
	var labels map[string]string
	if isStack {
		labels = map[string]string{"RFStack": "true"}
	}

	network, err := client.NetworkCreate(ctx, name, types.NetworkCreate{
		CheckDuplicate: true,
		Labels:         labels,
	})

	if err != nil {
		panic(err)
	}

	return network
}

func (op DockerOperation) CreateStack(name string) {
	op.EnsureBootstrapped()

	if !hasNetwork(name) {
		createNetwork(name, true)

		fmt.Printf("Created stack: %s\n", name)
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
		op.RemoveContainer(&container)
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
	id := metadata.GetMeta("stack:currentID")

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

	metadata.PutMeta("stack:currentID", stack.ID)
	return nil
}
