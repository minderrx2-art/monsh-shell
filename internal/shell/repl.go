package shell

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/minderrx2-art/monsh/internal/builtin"
	"github.com/minderrx2-art/monsh/internal/parser"
	"github.com/minderrx2-art/monsh/internal/runner"
)

func Start() error {
	for {
		fmt.Print("$ ")
		reader := bufio.NewReader(os.Stdin)
		line, _ := reader.ReadString('\n')

		tokens := parser.Tokenize(strings.TrimSpace(line))

		cmdPipeline, err := parser.ParsePipeline(tokens)

		if err != nil {
			fmt.Println(err)
		}
		runner.ExecutePipeline(cmdPipeline)
		if err != nil {
			fmt.Println(err)
			continue
		}
		// if err != nil {
		// 	if errors.Is(err, parser.ErrEmptyInput) {
		// 		continue
		// 	}
		// 	fmt.Println(err)
		// 	continue
		// }

		// builtinFunc := builtinRouter(cmd.Name, cmd.Args)

		// // Check for builtins
		// if builtinFunc != nil {
		// 	builtinFunc()
		// 	continue
		// }

		// if exists, _, err := path.Find(cmd.Name); exists == true && err == nil {
		// 	runner.Execute(cmd.Name, cmd.Args...)
		// } else {
		// 	fmt.Printf("%s: command not found\n", cmd.Name)
		// }
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
		return func() {
			err := builtin.Cd(rest[0])
			if err != nil {
				fmt.Println(err)
			}
		}
	default:
		return nil
	}
}
