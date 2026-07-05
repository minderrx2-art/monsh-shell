package parser

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenWord
)

type Token struct {
	Type  TokenType
	Value string
}
