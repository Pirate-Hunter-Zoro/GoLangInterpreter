package ast

import (
	"monkey/token"
	"bytes"
)

// Every node in our AST has to implement the Node interface
type Node interface {
	// Return the literal value of the token this node is associated with - used for debugging and testing
	TokenLiteral() string
	// String represntation for this Node
	String() string
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
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// One field for the identifier, one for the expression that produces the value, and one for the token
type LetStatement struct {
	Token token.Token // The token.LET token
	Name *Identifier
	Value Expression
}
func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral()+" ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type Identifier struct {
	Token token.Token // The token.IDENT token
	Value string
}
func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string { return i.Value }

type ReturnStatement struct {
	Token		token.Token // The 'return' token
	ReturnValue Expression
}
func (rs *ReturnStatement) statementNode()		{}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token 		token.Token // the firtst token of the expression
	Expression 	Expression
}
func (es *ExpressionStatement) statementNode() 		 {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// AST node representation of an integer
type IntegerLiteral struct {
	Token token.Token
	Value int64 // We're going to have to convert from a string in the parsing function associated with token.INT in the parser, called parseIntegerLiteral
}
func (il *IntegerLiteral) expressionNode()		{}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string 		{ return il.Token.Literal }

// AST representation of a boolean
type Boolean struct {
	Token token.Token
	Value bool
}
func (b *Boolean) expressionNode()		{}
func (b *Boolean) TokenLiteral() string { return b.Token.Literal }
func (b *Boolean) String() string 		{ return b.Token.Literal }

// AST representation of an if statement
type IfExpression struct {
	Token token.Token // The 'if' token
	Condition Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}
func (ie *IfExpression) expressionNode()	  {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

// AST representation of a block statement
type BlockStatement struct {
	Token token.Token // the { token
	Statements []Statement
}
func (bs *BlockStatement) statementNode()		{}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// AST representation of a prefix expression
type PrefixExpression struct {
	Token 		token.Token // The prefix token, e.g. !
	Operator 	string
	Right 		Expression // Who KNOWS what that will be - hence a recursive parsing call eventually
}
func (pe *PrefixExpression) expressionNode()		{}
func (pe *PrefixExpression) TokenLiteral() string 	{ return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String()) // The expression to the right could be any kind of expression, but it will know how to represent itself as a string
	out.WriteString(")")

	return out.String()
}

// AST representation of an infix expression
type InfixExpression struct {
	Token 		token.Token // The operator token, e.g. +
	Left 		Expression
	Operator	string
	Right 		Expression
}
func (oe *InfixExpression) expressionNode()		 {}
func (oe *InfixExpression) TokenLiteral() string { return oe.Token.Literal }
func (oe *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(oe.Left.String())
	out.WriteString(" " + oe.Operator + " ")
	out.WriteString(oe.Right.String())
	out.WriteString(")")

	return out.String()
}