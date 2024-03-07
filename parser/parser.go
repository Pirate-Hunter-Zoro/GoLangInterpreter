package parser

import (
	"fmt"
	"strconv"
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

const (
	_ int = iota // Give the following constants incrementing numbers as values
	LOWEST
	EQUALS		// ==
	LESSGREATER // > or <
	SUM			// +
	PRODUCT		// *
	PREFIX		// -X or !X
	CALL		// myFunction(X)
	// Establishes order of operations
)

type Parser struct {
	l *lexer.Lexer

	errors []string

	curToken token.Token // Current token under examination
	peekToken token.Token // The next token to be read, which we may need to know if curToken does not give us enough information
	// For example, let x = 5; 
	// When curToken is 5, peekToken will help us know if we are at the end of a line or at the beginning of a more complex expression

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns map[token.TokenType]infixParseFn
}

// Map the given tokenType to the given prefix function
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// Map the given tokenType to the given post-fix function
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// Create a parser from a lexer
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:	 	l,
		errors:	[]string{},
	}

	// Make the map of prefix functions and throw in the function to handle identifiers
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)

	// Read two tokens, so curToken and peekToken are both set
	p.NextToken()
	p.NextToken()

	return p
}

func (p *Parser) parseIdentifier() ast.Expression {
	// NO advancing of the tokens
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) Errors() []string {
	return p.errors
}

// Report an error where the next token should have been the given token type but was something else
func (p *Parser) PeekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) NextToken() {
	p.curToken = p.peekToken // Still null the first time this function is called - that's why we call it twice above
	p.peekToken = p.l.NextToken()
}

// Simply parse together a list of statements by progressing through each token
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.NextToken()
	}

	return program
}

// So how do we parse statements?
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// So how do we parse expression statements?
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.NextToken()
	}

	return stmt
}

// How do we parse a general expression
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// Grab the function from our map
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}
	// Apply the function to get its value
	leftExp := prefix()

	return leftExp
}

// So how do we parse return statements?
func (p *Parser) parseReturnStatement() ast.Statement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.NextToken()

	// TODO: We're skipping the expressions until we encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.NextToken()
	}

	return stmt

}

// So how do we parse let statements?
func (p *Parser) parseLetStatement() ast.Statement {
	stmt := &ast.LetStatement{Token: p.curToken}

	// If we have a let statement, the next token had better be an identifier - and if it is then progress the token
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// If we have a let statement, the next token had better be an assignment - and if it is then progress the token
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: We're skipping the expressions until we encounter a semicolon
	for !p.curTokenIs(token.SEMICOLON) {
		p.NextToken()
	}

	return stmt
}

// How do we parse and integer literal?
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

// Helper to see what type the parser's current token is
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}
// Same for the peek token
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// Helper to progress the progress the parser's token IF AND ONLY IF the next token is what we expect it to be
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.NextToken()
		return true
	} else {
		p.PeekError(t)
		return false
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////
// Implementing the Pratt Parser:

type (
	prefixParseFn func() ast.Expression // No left side
	infixParseFn func(ast.Expression) ast.Expression // YES left side
)