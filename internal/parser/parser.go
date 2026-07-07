package parser

import (
	"errors"
	"slices"
)

var ErrEmptyInput = errors.New("empty input")

func Parse(tokens []Token) (*Command, error) {
	var words []string

	for _, token := range tokens {
		if token.Type == TokenEOF {
			break
		}
		if token.Type == TokenWord {
			words = append(words, token.Value)
		}
	}

	if len(words) == 0 {
		return nil, ErrEmptyInput
	}

	args := words[1:]
	args = slices.DeleteFunc(args, func(arg string) bool {
		return arg == ""
	})

	return &Command{
		Name: words[0],
		Args: args,
	}, nil
}
