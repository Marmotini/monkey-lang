// reads input, sends it to the interpreter for evaluation, prints the result and starts again
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/nerdysquirrel/monkey-lang/lexer"
	"github.com/nerdysquirrel/monkey-lang/token"
)

const PROMPT = ">> "

var (
	interactive bool
	fileName    string
)

func init() {
	flag.BoolVar(&interactive, "i", false, "interactive mode")
	flag.StringVar(&fileName, "f", "", "filename")
}

func main() {
	flag.Parse()

	if fileName != "" {
		lex := lexer.NewLexerFromFile(fileName)

		for tok := lex.NextToken(); tok.Type != token.EOF; tok = lex.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}

	if interactive && fileName == "" {
		StartInteractiveMode(os.Stdin, os.Stdout)
	}
}

func StartInteractiveMode(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)

	for {
		fmt.Printf(PROMPT)

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		lex := lexer.NewLexer(line)

		for tok := lex.NextToken(); tok.Type != token.EOF; tok = lex.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}