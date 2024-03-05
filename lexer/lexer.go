package lexer

import "monkey/token"

type Lexer struct {
	input        string
	position     int  // Current position in input (points to current char)
	readPosition int  // Current reading position in input (after current char) - we'll need to be able to peek further into the input after the current character
	ch           byte // Current char under examination (ascii values are sufficiently encompassed by 8 bits) - would have to be a 'rune' if we were supporting all of Unicode
}

func New(input string) *Lexer { // Returns a pointer to a Lexer struct
	l := &Lexer{input: input} // The address of the lexer
	l.readChar() // So that the first character is read - when we call NextToken() it will not be "EOF" with value 0
	return l
}

// If we were supporting more characters, like all of Unicode (including emoji's), then we would need to change how this is done - read position may go up by more than a byte
func (l *Lexer) readChar() { // Takes in a pointer to a lexer
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII for "NUL" - either at the end of the file or we haven't read anything yet
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) NextToken() token.Token {

	var tok token.Token

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok

}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) { // Keep reading until we hit a non-letter in the name of this identifier
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	// THIS is the place to sneak in new character allowed for identifier names
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}