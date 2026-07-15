package shell

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/chzyer/readline"
	"github.com/minderrx2-art/monsh/internal/builtin"
	"github.com/minderrx2-art/monsh/internal/parser"
	"github.com/minderrx2-art/monsh/internal/path"
	"github.com/minderrx2-art/monsh/internal/runner"
)

type ShellCompleter struct {
	base       readline.PrefixCompleter
	lastPrefix string
	tabPressed bool
}

func (c *ShellCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	prefix := strings.TrimSpace(string(line[:pos]))
	rawMatches, offset := c.base.Do(line, pos)
	matches := make(map[string]struct{})

	for _, match := range rawMatches {
		matches[string(match)] = struct{}{}
	}
	uniqueMatches := slices.Sorted(maps.Keys(matches))

	if prefix != c.lastPrefix {
		c.lastPrefix = prefix
		c.tabPressed = false
	}

	if len(uniqueMatches) == 0 {
		fmt.Print("\x07")
		return nil, 0
	}

	if len(uniqueMatches) == 1 {
		c.tabPressed = false
		return [][]rune{[]rune(uniqueMatches[0])}, offset
	}

	names := make([]string, 0, len(uniqueMatches))

	// Mutli stage completions
	for _, suffix := range uniqueMatches {
		names = append(names, prefix+strings.TrimSpace(suffix))
	}
	slices.Sort(names)

	lcp := longestCommonPrefix(names)

	if len(lcp) > len(prefix) {
		c.tabPressed = false
		return [][]rune{[]rune(lcp[len(prefix):])}, offset
	}

	if !c.tabPressed {
		fmt.Print("\x07")
		c.tabPressed = true
		return nil, 0
	}

	fmt.Printf("\n%s\n", strings.Join(names, "  "))
	fmt.Printf("$ %s", prefix)
	c.tabPressed = false
	return nil, 0
}

func longestCommonPrefix(names []string) string {
	if len(names) == 0 {
		return ""
	}
	shortest := names[0]
	for _, s := range names[1:] {
		if len(s) < len(shortest) {
			shortest = s
		}
	}
	for i := 0; i < len(shortest); i++ {
		for _, s := range names {
			if s[i] != shortest[i] {
				return shortest[:i]
			}
		}
	}
	return shortest
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
	reader, err := newReader()
	for {
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
