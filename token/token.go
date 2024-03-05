package token

type TokenType string // Making this a string makes it easier to debug

type Token struct {
	Type TokenType // To distinguish between "integers" and "right bracket" for example
	Literal string // The literal value of the token - so is the "integer" a 5 or a 10?
}

const (
	ILLEGAL 	= "ILLEGAL" // Signifies a token/character we don't know about
	EOF 		= "EOF"	// End of file - tells the parser that it can stop

	// Identifiers + literals
	IDENT		= "IDENT" // add, foobar, x, y, ...
	INT 		= "INT" // 123456

	// Operators
	ASSIGN		= "="
	PLUS		= "+"

	// Delimiters
	COMMA		= ","
	SEMICOLON 	= ";" 

	LPAREN 		= "("
	RPAREN		= ")"
	LBRACE		= "{"
	RBRACE		= "}"

	// Keywords
	FUNCTION	= "FUNCTION"
	LET 		= "LET"

)