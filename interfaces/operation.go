package interfaces

type Stack struct {
	Name string
	ID   string
}

type Operation interface {
	// ---  Stacks  --------------------------------------------------------------
	CreateStack(name string)
	ListStacks() []Stack
	GetDefaultStack() *Stack
	SetDefaultStack(name string) error
	FindStacks(flags []map[string]string) []Stack
	FindStack(flags []map[string]string) *Stack
}
