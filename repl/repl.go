package repl

// REPL -> Read Eval Print Loop

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) { // Read until you encounter a new line, take the just read line and pass it to our lexer
	scanner := bufio.NewScanner(in) // Like a BufferedReader in Java

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		// Read until you encounter a new line
		line := scanner.Text()
		// Take the just read line and pass it to our lexer
		l := lexer.New(line)

		// Print all the tokens the lexer gives us until we encounter EOF
		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}