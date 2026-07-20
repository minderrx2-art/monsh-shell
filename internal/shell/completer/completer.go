package completer

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/chzyer/readline"
)

type ShellCompleter struct {
	ReadlineCompleter readline.PrefixCompleter
	lastPrefix        string
	tabPressed        bool
}

func (c *ShellCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	prefix := strings.TrimSpace(string(line[:pos]))
	rawMatches, offset := c.ReadlineCompleter.Do(line, pos)
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
