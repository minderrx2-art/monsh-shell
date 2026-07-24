package builtin

import (
	"fmt"
)

type CompleteCommand struct {
	savedCompletions map[string]string
}

func (c *CompleteCommand) register(name string, completion string) {
	c.savedCompletions[name] = completion
}

func (c *CompleteCommand) get(name string) string {
	return c.savedCompletions[name]
}

func NewCompleteCommand() *CompleteCommand {
	return &CompleteCommand{
		savedCompletions: make(map[string]string),
	}
}

func (c *CompleteCommand) Complete(args []string) {
	switch args[0] {
	case "-p":
		fmt.Println(c.get(args[1]))
	case "-C":
		c.register(args[2], args[1])
	default:
		fmt.Printf("complete: %s: no completion specification\n", args[1])
	}
}
