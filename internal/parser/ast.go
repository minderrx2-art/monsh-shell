package parser

type RedirectType int

const (
	In        RedirectType = iota // <
	Out                           // 1> >
	OutErr                        // 2>
	Append                        // >>
	AppendErr                     // 2>>
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
