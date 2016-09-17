package operations

import (
	"testing"

	"github.com/controlroom/lincoln/interfaces"
)

type TestBackend struct{}

func (t TestBackend) StartContainer(c interfaces.ContainerStartOptions) {
}

type TestOperation struct {
}

func TestDevStart(t *testing.T) {
}
