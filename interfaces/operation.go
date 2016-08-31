package interfaces

type Stack struct {
	Name string
	ID   string
}

type ContainerStartOptions struct {
	Image        string
	Name         string
	Cmd          []string
	Stack        Stack
	Volumes      []string
	Env          []string
	CapAdd       []string
	Ports        []string
	PortBindings []string
}

type ContainerInfo struct {
	Stack Stack
	IP    string
}

type Container struct {
	Name string
	ID   string
}

type Operation interface {
	// ---  General  --------------------------------------------------------------
	EnsureBootstrapped()

	// ---  Stacks  --------------------------------------------------------------
	CreateStack(name string)
	DestroyStack(name string)
	ListStacks() []Stack
	GetDefaultStack() *Stack
	SetDefaultStack(name string) error
	FindStacks(flags []map[string]string) []Stack
	FindStack(flags []map[string]string) *Stack

	// ---  Containers  -----------------------------------------------------------
	StartContainer(opts ContainerStartOptions) *Container
	RemoveContainer(container *Container)
	ListContainers() []Container
	FindContainers(flags []map[string]string) []Container
	FindContainer(flags []map[string]string) *Container
	FindContainerByName(name string) *Container
	InspectContainer(container *Container) ContainerInfo
}
