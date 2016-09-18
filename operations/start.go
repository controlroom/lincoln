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

func buildVolumes(s *StartOperation) []string {
	s.Backend.CreateVolume(s.App.Config.Name)
	sv := make([]string, len(s.App.Config.SharedPaths))
	for i, path := range s.App.Config.SharedPaths {
		sv[i] = fmt.Sprintf("%s:%s", s.App.Config.Name, path)
	}
	return sv
}

func (s *StartOperation) StartDev(ops []string) error {
	config := s.App.Config
	sharedVolumes := buildVolumes(s)
	nodes := config.GetNodes()

	for _, op := range ops {
		node := nodes[op]

		fmt.Printf("%+v\n", config)
		s.Backend.StartContainer(interfaces.ContainerStartOptions{
			Image:       config.DevImage,
			Cmd:         strings.Split(node.Cmd, " "),
			Name:        fmt.Sprintf("%v-node-%v", config.Name, op),
			Volumes:     sharedVolumes,
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
		Volumes:     buildVolumes(s),
		VolumesFrom: []string{fmt.Sprintf("%v-dev-sync", s.App.Config.Name)},
		Stack:       s.Backend.GetDefaultStack(),
	})

	return nil
}
