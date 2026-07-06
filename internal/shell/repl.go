package shell

import (
	"bufio"
	"fmt"
	"os"
	"slices"
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
			rest = slices.DeleteFunc(words[1:], func(word string) bool {
				return word == ""
			})
		}

		builtinFunc := builtinRouter(command, rest)

		// Check for builtins
		if builtinFunc != nil {
			builtinFunc()
			continue
		}

		if exists, _, err := path.Find(command); exists == true && err == nil {
			runner.Execute(command, rest...)
		} else {
			fmt.Printf("%s: command not found\n", command)
		}
	}
}

func builtinRouter(command string, rest []string) func() {
	switch command {
	case "type":
		return func() { builtin.Type(rest[0]) }
	case "echo":
		return func() { builtin.Echo(rest) }
	case "exit":
		return func() { builtin.Exit() }
	case "pwd":
		return func() { builtin.Pwd() }
	case "cd":
		return func() { builtin.Cd(rest[0]) }
	default:
		return nil
	}
}
