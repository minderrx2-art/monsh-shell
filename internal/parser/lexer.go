package parser

import (
	"strings"
)

type Lexer struct {
	input  string
	pos    int
	state  parseState
	curr   strings.Builder
	tokens []Token
}

type parseState int

const (
	stateNormal parseState = iota
	stateSingleQuote
	stateDoubleQuote
	stateEscape
)

func Tokenize(input string) []Token {
	l := &Lexer{
		input:  input,
		state:  stateNormal,
		tokens: []Token{},
	}
	return l.Tokenize()
}

func (l *Lexer) Tokenize() []Token {
	for l.pos = 0; l.pos < len(l.input); l.pos++ {
		r := rune(l.input[l.pos])
		switch l.state {

		case stateNormal:
			switch r {
			case ' ', '\t':
				l.flush()
			case '\'':
				l.state = stateSingleQuote
			case '"':
				l.state = stateDoubleQuote
			case '\\':
				l.state = stateEscape
			default:
				l.curr.WriteRune(r)
			}

		case stateSingleQuote:
			switch r {
			case '\'':
				l.state = stateNormal
			default:
				l.curr.WriteRune(r)
			}

		case stateDoubleQuote:
			switch r {
			case '"':
				l.state = stateNormal
			case '\\':
				l.pos = l.handleEscape()
			default:
				l.curr.WriteRune(r)
			}

		case stateEscape:
			l.curr.WriteRune(r)
			l.state = stateNormal
		}
	}

	l.flush()
	l.tokens = append(l.tokens, Token{Type: TokenEOF})
	return l.tokens
}

func (l *Lexer) flush() {
	if l.curr.Len() <= 3 {
		chars := l.curr.String()
		switch chars {
		case ">", "1>":
			l.tokens = append(l.tokens, Token{Type: TokenRedirectOut, Value: chars})
			l.curr.Reset()
		case "2>":
			l.tokens = append(l.tokens, Token{Type: TokenRedirectOutError, Value: chars})
			l.curr.Reset()
		case ">>":
			l.tokens = append(l.tokens, Token{Type: TokenRedirectAppend, Value: chars})
			l.curr.Reset()
		case "1>>":
			l.tokens = append(l.tokens, Token{Type: TokenRedirectAppend, Value: chars})
			l.curr.Reset()
		case "2>>":
			l.tokens = append(l.tokens, Token{Type: TokenRedirectAppendError, Value: chars})
			l.curr.Reset()
		case "<", "0<":
			l.tokens = append(l.tokens, Token{Type: TokenRedirectIn, Value: chars})
			l.curr.Reset()
		case "|":
			l.tokens = append(l.tokens, Token{Type: TokenPipe, Value: chars})
			l.curr.Reset()

		case " ", "":
			l.curr.Reset()
		default:
			l.tokens = append(l.tokens, Token{Type: TokenWord, Value: chars})
			l.curr.Reset()
		}
	} else if l.curr.Len() > 0 {
		l.tokens = append(l.tokens, Token{Type: TokenWord, Value: l.curr.String()})
		l.curr.Reset()
	}
}

func (l *Lexer) handleEscape() int {
	if l.pos+1 < len(l.input) {
		next := l.input[l.pos+1]
		switch next {
		case '\\', '"', '$', '`', '\n':
			l.curr.WriteByte(next)
			return l.pos + 1
		default:
			l.curr.WriteRune('\\')
		}
	} else {
		l.curr.WriteByte('\\')
	}
	return l.pos
}
