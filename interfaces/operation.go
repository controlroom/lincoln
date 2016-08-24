package interfaces

type Stack struct {
	Name string
	ID   string
}

type Operation interface {
	// ---  content  --------------------------------------------------------------
	CreateStack(name string)
	ListStacks() []Stack
}
