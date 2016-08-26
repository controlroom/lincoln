package interfaces

type Stack struct {
	Name string
	ID   string
}

type ContainerStartOptions struct {
	Image   string
	Name    string
	Cmd     []string
	Stack   Stack
	Volumes []string
}

type Container struct {
	Name string
	ID   string
}

type Operation interface {
	// ---  Stacks  --------------------------------------------------------------
	CreateStack(name string)
	DestroyStack(name string)
	ListStacks() []Stack
	GetDefaultStack() *Stack
	SetDefaultStack(name string) error
	FindStacks(flags []map[string]string) []Stack
	FindStack(flags []map[string]string) *Stack

	// ---  Containers  -----------------------------------------------------------
	StartContainer(opts ContainerStartOptions)
	RemoveContainer(container Container)
	ListContainers() []Container
	FindContainers(flags []map[string]string) []Container
}
