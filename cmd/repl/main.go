package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/user"

	"github.com/donovandicks/gomonkey/internal/lexer"
	"github.com/donovandicks/gomonkey/internal/parser"
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
		p := parser.NewParser(l)

		program := p.ParseProgram()
		if errs := p.Errors(); len(errs) != 0 {
			for _, msg := range errs {
				io.WriteString(out, "\t"+msg+"\n")
			}
			continue
		}

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
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
