package operations

import (
	"fmt"
	"strings"

	"github.com/controlroom/lincoln/config"
	"github.com/controlroom/lincoln/interfaces"
)

type Startable interface {
	StartDev(ops []string) error
	RunDev(cmd []string) error
}

type StartOperation struct {
	Backend interfaces.Operation
	App     *config.App
}

func (s *StartOperation) StartDev(ops []string) error {
	config := s.App.Config
	nodes := config.GetNodes()

	for _, op := range ops {
		node := nodes[op]

		fmt.Printf("%+v\n", config)
		s.Backend.StartContainer(interfaces.ContainerStartOptions{
			Image:       config.DevImage,
			Cmd:         strings.Split(node.Cmd, " "),
			Name:        fmt.Sprintf("%v-node-%v", config.Name, op),
			VolumesFrom: []string{fmt.Sprintf("%v-dev-sync", config.Name)},
			Stack:       s.Backend.GetDefaultStack(),
		})
	}

	return nil
}

func (s *StartOperation) RunDev(cmd []string) error {
	s.Backend.RunContainer(interfaces.ContainerStartOptions{
		Image:       s.App.Config.DevImage,
		Cmd:         cmd,
		VolumesFrom: []string{fmt.Sprintf("%v-dev-sync", s.App.Config.Name)},
		Stack:       s.Backend.GetDefaultStack(),
	})

	return nil
}
