package operations

import (
	"errors"
	"fmt"
	"strings"

	"github.com/controlroom/lincoln/config"
	"github.com/controlroom/lincoln/interfaces"
)

type Startable interface {
	StartDev(ops []string) error
	RunDev(cmd []string) error
	TestRun(cmd string, args []string) error
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

func (s *StartOperation) TestRun(cmd string, args []string) error {
	testCmd := s.App.Config.Tests[cmd]

	if testCmd == "" {
		return errors.New("Not a test command")
	}

	config := s.App.Config

	stack := s.Backend.CreateStack(fmt.Sprintf("test-%v-%v", config.Name, cmd))
	defer s.Backend.DestroyStack(stack.Name)

	for _, dep := range config.Deps.Resources {
		nameVer := strings.Split(dep, ":")
		s.Backend.StartContainer(interfaces.ContainerStartOptions{
			Image:      dep,
			Name:       fmt.Sprintf("%v-test-dep-%v", config.Name, nameVer[0]),
			Stack:      stack,
			NetAliases: []string{nameVer[0]},
		})
	}

	s.Backend.RunContainer(interfaces.ContainerStartOptions{
		Image: s.App.Config.DevImage,
		// Cmd:   strings.Split(testCmd, " "),
		Cmd: []string{"bash"},
		Env: []string{
			"POSTGRES_SERVICE_HOST=postgres",
			"POSTGRES_USERNAME=postgres",
			"REDIS_HOST=redis",
		},
		Volumes:     buildVolumes(s),
		VolumesFrom: []string{fmt.Sprintf("%v-dev-sync", s.App.Config.Name)},
		Stack:       stack,
	})

	return nil
}
