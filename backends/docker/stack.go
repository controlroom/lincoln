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
		client.NetworkCreate(ctx, name, types.NetworkCreate{
			CheckDuplicate: true,
			Labels:         map[string]string{"RFStack": "true"},
		})
		fmt.Printf("Created stack: %s\n", name)
	} else {
		fmt.Printf("Stack already created: %s\n", name)
	}
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
