package parser

import (
	"errors"
	"fmt"
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
		if token.Type == TokenPipe {
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

/*
[{1 echo} {1 damn son} {2 >} {1 test.txt} {0 }]
$ echo "damn son" 1> test.txt
[{1 echo} {1 damn son} {2 1>} {1 test.txt} {0 }]
$ echo "damn son" 2> test.txt
[{1 echo} {1 damn son} {3 2>} {1 test.txt} {0 }]
$ echo "damn son" 0< test.txt
[{1 echo} {1 damn son} {5 0<} {1 test.txt} {0 }]
$ echo "damn son" < test.txt
[{1 echo} {1 damn son} {5 <} {1 test.txt} {0 }]
*/
func parseSimpleCommand(tokens []Token, i *int) (Command, error) {
	cmd := Command{}

	// First word = command name
	cmd.Name = tokens[*i].Value
	*i++

	// Read words until operator or EOF
	// [{1 damn son} |stop| {2 >} {1 text.txt} {0 }]
	for *i < len(tokens) && tokens[*i].Type == TokenWord {
		cmd.Args = append(cmd.Args, tokens[*i].Value)
		*i++
	}

	// Read redirects attached to this command
	for *i < len(tokens) && tokens[*i].Type != TokenEOF {
		op := tokens[*i].Value
		*i++
		// If nothing is on the right, or the last value is not a TokenWord
		if *i >= len(tokens) || tokens[*i].Type != TokenWord {
			return cmd, errors.New("expected filename after redirect")
		}
		// Set value after op as target
		target := tokens[*i].Value
		*i++

		switch op {
		case ">", "1>":
			cmd.Redirects = append(cmd.Redirects, Redirect{Type: Out, Target: target})
		case "2>":
			cmd.Redirects = append(cmd.Redirects, Redirect{Type: OutErr, Target: target})
		case ">>":
			cmd.Redirects = append(cmd.Redirects, Redirect{Type: Append, Target: target})
		case "<", "0<":
			cmd.Redirects = append(cmd.Redirects, Redirect{Type: In, Target: target})
		}
		fmt.Println(cmd)
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
