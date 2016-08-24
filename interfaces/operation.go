package interfaces

type Stack struct {
	Name string
	ID   string
}

type Operation interface {
	CreateStack(name string)
	ListStacks() []Stack
}
