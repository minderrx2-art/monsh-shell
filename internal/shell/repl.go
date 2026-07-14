package shell

import (
	"errors"
	"fmt"
	"strings"

	"github.com/chzyer/readline"
	"github.com/minderrx2-art/monsh/internal/builtin"
	"github.com/minderrx2-art/monsh/internal/parser"
	"github.com/minderrx2-art/monsh/internal/path"
	"github.com/minderrx2-art/monsh/internal/runner"
)

type ShellCompleter struct {
	base readline.PrefixCompleter
}

func (c *ShellCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	newLine, offset := c.base.Do(line, pos)

	if len(newLine) == 0 {
		fmt.Print("\x07")
	}

	return newLine, offset
}

func listExecutables(prefix string) []string {
	list := []string{}
	matches, err := path.FindAll(prefix)
	if err != nil {
		return []string{}
	}
	if len(matches) == 0 {
		return []string{}
	}
	for _, match := range matches {
		list = append(list, match)
	}
	return list
}

func newReader() (*readline.Instance, error) {
	l, err := readline.NewEx(&readline.Config{
		Prompt:          "$ ",
		HistoryFile:     "/tmp/monsh.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
		AutoComplete: &ShellCompleter{
			base: *readline.NewPrefixCompleter(
				readline.PcItem("exit"),
				readline.PcItem("pwd"),
				readline.PcItem("cd"),
				readline.PcItem("type"),
				readline.PcItem("echo"),
				readline.PcItemDynamic(listExecutables),
			),
		},
	})
	if err != nil {
		return nil, err
	}
	return l, nil
}

func Start() error {
	for {
		reader, err := newReader()
		if err != nil {
			return err
		}

		defer reader.Close()

		line, err := reader.Readline()
		if err != nil {
			return err
		}
		tokens := parser.Tokenize(strings.TrimSpace(line))

		cmdPipeline, err := parser.Parse(tokens)

		if err != nil {
			if !errors.Is(err, parser.ErrEmptyInput) {
				fmt.Println(err)
			}
			continue
		}

		if len(cmdPipeline.Commands) == 1 && cmdPipeline.Commands[0].Redirects == nil {
			cmd := cmdPipeline.Commands[0]
			builtinFunc := builtinRouter(cmd.Name, cmd.Args)

			// Check for builtins
			if builtinFunc != nil {
				builtinFunc()
				continue
			}

			if exists, _, err := path.Find(cmd.Name); exists == true && err == nil {
				runner.Execute(cmd.Name, cmd.Args...)
			} else {
				fmt.Printf("%s: command not found\n", cmd.Name)
			}
		} else {
			if err := runner.ExecutePipeline(cmdPipeline); err != nil {
				// oops, something went wrong but its probably ok
			}
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
