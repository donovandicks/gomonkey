package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"

	"github.com/donovandicks/gomonkey/internal/lexer"
	"github.com/donovandicks/gomonkey/internal/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, PROMPT)
		input := scanner.Scan()

		if !input {
			return
		}

		line := scanner.Text()
		l := lexer.NewLexer(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Fprintf(out, "%+v\n", tok)
		}
	}
}

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello, %s!\n", user.Username)

	Start(os.Stdin, os.Stderr)
}
