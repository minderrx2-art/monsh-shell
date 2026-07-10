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
		if curr.Len() <= 3 {
			chars := curr.String()
			switch chars {
			case ">", "1>":
				tokens = append(tokens, Token{Type: TokenRedirectOut, Value: chars})
				curr.Reset()
			case "2>":
				tokens = append(tokens, Token{Type: TokenRedirectOutError, Value: chars})
				curr.Reset()
			case ">>":
				tokens = append(tokens, Token{Type: TokenRedirectAppend, Value: chars})
				curr.Reset()
			case "1>>":
				tokens = append(tokens, Token{Type: TokenRedirectAppend, Value: chars})
				curr.Reset()
			case "<", "0<":
				tokens = append(tokens, Token{Type: TokenRedirectIn, Value: chars})
				curr.Reset()
			case " ", "":
				curr.Reset()
			default:
				tokens = append(tokens, Token{Type: TokenWord, Value: chars})
				curr.Reset()
			}
		} else if curr.Len() > 0 {
			tokens = append(tokens, Token{Type: TokenWord, Value: curr.String()})
			curr.Reset()
		}
	}

	// Walk through the input one rune at a time.
	for i := 0; i < len(input); i++ {
		r := rune(input[i])
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
			switch r {
			// Closing single quote; return to normal parsing.
			case '\'':
				state = stateNormal
			default:
				curr.WriteRune(r)
			}

		case stateDoubleQuote:
			switch r {
			// Closing double quote; return to normal parsing.
			case '"':
				state = stateNormal
			case '\\':
				i = handleEscapeCharacter(i, input, &curr)
			default:
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

func handleEscapeCharacter(index int, input string, curr *strings.Builder) int {
	// When not out of bounds
	if index+1 < len(input) {
		next := input[index+1]
		switch next {
		// when next character can be escaped
		case '\\', '"', '$', '`', '\n':
			// write it and skip over a character (since its written)
			curr.WriteByte(next)
			return index + 1
		default:
			// when it can't be escaped just write a \
			curr.WriteRune('\\')
		}

	} else {
		curr.WriteByte('\\')
	}
	return index
}
