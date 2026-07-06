package parser

import (
	"strings"
)

type Lexer struct {
	input string
	pos   int
}

type parseState int

const (
	stateNormal parseState = iota
	stateSingleQuote
	stateDoubleQuote
	stateEscape
)

func Tokenize(input string) []Token {
	var tokens []Token
	var curr strings.Builder
	state := stateNormal

	// Write out the word currently in buffer
	flushWord := func() {
		if curr.Len() > 0 {
			tokens = append(tokens, Token{Type: TokenWord, Value: curr.String()})
			curr.Reset()
		}
	}

	// Walk through the input one rune at a time.
	for _, r := range input {
		switch state {
		case stateNormal:
			switch r {
			case ' ', '\t':
				// End of the current word. Emit it as a token and clear the buffer.
				flushWord()
			case '\'':
				// Enter single-quoted mode.
				state = stateSingleQuote
			case '"':
				// Enter double-quoted mode.
				state = stateDoubleQuote
			case '\\':
				// The next rune should be treated as escaped.
				state = stateEscape
			default:
				// Append the rune to the current word.
				curr.WriteRune(r)
			}

		case stateSingleQuote:
			if r == '\'' {
				// Closing single quote; return to normal parsing.
				state = stateNormal
			} else {
				curr.WriteRune(r)
			}

		case stateDoubleQuote:
			if r == '"' {
				// Closing double quote; return to normal parsing.
				state = stateNormal
			} else {
				curr.WriteRune(r)
			}

		case stateEscape:
			curr.WriteRune(r)
			state = stateNormal
		}

	}

	// Final flush of builder buffer
	flushWord()

	// Append EOF mark
	tokens = append(tokens, Token{Type: TokenEOF})
	return tokens
}

// Convert tokens to strings
// [{1 foo}, {1 bar}, {0 }] to ["foo", "bar", ""]
func Words(input []Token) []string {
	words := make([]string, len(input))
	for i, token := range input {
		words[i] = token.Value
	}
	return words
}
