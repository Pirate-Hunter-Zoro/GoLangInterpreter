package ast

import "monkey/token"

// Every node in our AST has to implement the Node interface
type Node interface {
	// Return the literal value of the token this node is associated with - used for debugging and testing
	TokenLiteral() string
}

type Statement interface {
	Node
	// Dummy method
	statementNode()
}

type Expression interface {
	Node
	// Dummy method
	expressionNode()
}

// The Program node will be the root node of every AST our parser produces.
type Program struct {
	// A slice of the AST nodes that implement the Statement interface
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// One field for the identifier, one for the expression that produces the value, and one for the token
type LetStatement struct {
	Token token.Token // The token.LET token
	Name *Identifier
	Value Expression
}
func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type Identifier struct {
	Token token.Token // The token.IDENT token
	Value string
}
func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }