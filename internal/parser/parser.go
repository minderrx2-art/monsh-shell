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
		// Todo
		if token.Type == TokenOperation {
			words = append(words, token.Value)
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
		Name: "echo",
		Args: []string{"damn son"},
		Redirects: []Redirect{
			{Type: In, Target: "text.txt"},
		},
	}, nil
}

// [{1 echo} {1 damn son} {2 >} {1 text.txt} {0 }]
func parseSimpleCommand(tokens []Token, i *int) (Command, error) {
	cmd := Command{}

	// first word = command name
	cmd.Name = tokens[*i].Value
	*i++

	// read words until operator or EOF
	// [{1 damn son} |stop| {2 >} {1 text.txt} {0 }]
	for *i < len(tokens) && tokens[*i].Type == TokenWord {
		cmd.Args = append(cmd.Args, tokens[*i].Value)
		*i++
	}

	// read redirects attached to this command
	for *i < len(tokens) && tokens[*i].Type == TokenOperation {
		op := tokens[*i].Value
		*i++

		// if nothing is on the right, or the last value is not a TokenWord
		if *i >= len(tokens) || tokens[*i].Type != TokenWord {
			return cmd, errors.New("expected filename after redirect")
		}
		// set value after op as target
		target := tokens[*i].Value
		*i++

		switch op {
		case ">":
			cmd.Redirects = append(cmd.Redirects, Redirect{Type: Out, Target: target})
		case ">>":
			// cmd.Redirects = append(cmd.Redirects, Redirect{Type: RedirectAppend, Target: target})
		case "<":
			cmd.Redirects = append(cmd.Redirects, Redirect{Type: In, Target: target})
		}
	}

	return cmd, nil
}

func ParsePipeline(tokens []Token) (*Pipeline, error) {
	i := 0
	var commands []Command

	for {
		cmd, err := parseSimpleCommand(tokens, &i)
		if err != nil {
			return nil, err
		}
		commands = append(commands, cmd)

		if i >= len(tokens) || tokens[i].Type == TokenEOF {
			break
		}
		if tokens[i].Value == "|" {
			i++ // consume pipe, next loop parses next command
			continue
		}
		return nil, errors.New("unexpected token")
	}

	return &Pipeline{Commands: commands}, nil
}
