package parser

type RedirectType int

const (
	In RedirectType = iota
	Out
	Append
)

type Command struct {
	Name      string
	Args      []string
	Redirects []Redirect
}

type Redirect struct {
	Type   RedirectType
	Target string
}

type Pipeline struct {
	Commands []Command
}
