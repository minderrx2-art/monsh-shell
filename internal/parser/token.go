package parser

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenWord
	TokenRedirectOut      // > 1>
	TokenRedirectOutError // 2>
	TokenRedirectAppend   // >>
	TokenRedirectIn       // < 0<
	TokenPipe
)

type Token struct {
	Type  TokenType
	Value string
}
