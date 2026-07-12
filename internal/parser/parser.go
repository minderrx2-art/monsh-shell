package parser

import (
	"errors"
)

var (
	ErrEmptyInput          = errors.New("empty input")
	ErrExpectedFilename    = errors.New("expected filename after redirect")
	ErrUnexpectedToken     = errors.New("unexpected token")
	ErrExpectedCommandName = errors.New("expected command name")
)

type Parser struct {
	tokens []Token
	pos    int
}

func Parse(tokens []Token) (*Pipeline, error) {
	p := &Parser{tokens: tokens}
	return p.parsePipeline()
}

func (p *Parser) parsePipeline() (*Pipeline, error) {
	if len(p.tokens) == 0 || p.tokens[0].Type == TokenEOF {
		return nil, ErrEmptyInput
	}

	var commands []Command
	for {
		cmd, err := p.parseSimpleCommand()
		if err != nil {
			return nil, err
		}
		commands = append(commands, cmd)

		if p.atEOF() {
			break
		}
		if p.peek().Type == TokenPipe {
			p.advance()
			continue
		}
		return nil, ErrUnexpectedToken
	}

	return &Pipeline{Commands: commands}, nil
}

func (p *Parser) parseSimpleCommand() (Command, error) {
	cmd := Command{}

	if p.peek().Type != TokenWord {
		return cmd, ErrExpectedCommandName
	}
	cmd.Name = p.advance().Value

	for p.peek().Type == TokenWord {
		cmd.Args = append(cmd.Args, p.advance().Value)
	}

	for isRedirect(p.peek()) {
		redirect := p.advance()

		if p.peek().Type != TokenWord {
			return cmd, ErrExpectedFilename
		}
		target := p.advance().Value

		switch redirect.Type {
		case TokenRedirectOut:
			cmd.Redirects = append(cmd.Redirects, Redirect{Type: Out, Target: target})
		case TokenRedirectOutError:
			cmd.Redirects = append(cmd.Redirects, Redirect{Type: OutErr, Target: target})
		case TokenRedirectAppend:
			cmd.Redirects = append(cmd.Redirects, Redirect{Type: Append, Target: target})
		case TokenRedirectAppendError:
			cmd.Redirects = append(cmd.Redirects, Redirect{Type: AppendErr, Target: target})
		case TokenRedirectIn:
			cmd.Redirects = append(cmd.Redirects, Redirect{Type: In, Target: target})
		}
	}

	return cmd, nil
}

func (p *Parser) peek() Token {
	if p.pos >= len(p.tokens) {
		return Token{Type: TokenEOF}
	}
	return p.tokens[p.pos]
}

func (p *Parser) advance() Token {
	token := p.peek()
	p.pos++
	return token
}

func (p *Parser) atEOF() bool {
	return p.peek().Type == TokenEOF
}

func isRedirect(token Token) bool {
	switch token.Type {
	case TokenRedirectOut, TokenRedirectOutError, TokenRedirectAppend, TokenRedirectAppendError, TokenRedirectIn:
		return true
	default:
		return false
	}
}
