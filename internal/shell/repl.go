package shell

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/minderrx2-art/monsh/internal/builtin"
	"github.com/minderrx2-art/monsh/internal/parser"
	"github.com/minderrx2-art/monsh/internal/path"
	"github.com/minderrx2-art/monsh/internal/runner"
)

func Start() error {
	for {
		fmt.Print("$ ")
		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')

		tokens := parser.Tokenize(strings.TrimSpace(line))
		words := parser.Words(tokens)

		if len(words) == 0 {
			continue
		}
		command := words[0]
		var (
			rest []string
		)
		if len(words) > 1 {
			rest = words[1:]
		}

		routedFunc := router(command, rest)

		if routedFunc != nil {
			routedFunc()
			continue
		}

		if exists, _, err := path.Find(command); exists == true && err == nil {
			runner.Execute(command, rest...)
		} else {
			fmt.Printf("%s: command not found\n", command)
		}
	}
}

func router(command string, rest []string) func() {
	switch command {
	case "type":
		return func() { builtin.Type(rest[0]) }
	case "echo":
		return func() { builtin.Echo(rest) }
	case "exit":
		return func() { builtin.Exit() }
	default:
		return nil
	}
}
