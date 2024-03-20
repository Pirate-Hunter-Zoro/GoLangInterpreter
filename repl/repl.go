package repl

// REPL -> Read Eval Print Loop

import (
	"bufio"
	"fmt"
	"io"
	"monkey/lexer"
	"monkey/parser"
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
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}
		
		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}