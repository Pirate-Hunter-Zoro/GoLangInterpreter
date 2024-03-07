package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

type Parser struct {
	l *lexer.Lexer

	curToken token.Token // Current token under examination
	peekToken token.Token // The next token to be read, which we may need to know if curToken does not give us enough information
	// For example, let x = 5; 
	// When curToken is 5, peekToken will help us know if we are at the end of a line or at the beginning of a more complex expression
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// Read two tokens, so curToken and peekToken are both set
	p.NextToken()
	p.NextToken()

	return p
}

func (p *Parser) NextToken() {
	p.curToken = p.peekToken // Still null the first time this function is called - that's why we call it twice above
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	return nil
}