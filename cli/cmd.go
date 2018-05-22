package cli

type Cmd struct {
	ArgList
	Name string
	Desc string
	Help string
	Run  func([]string)
	Sub  []*Cmd
}
