package main

import (
	"fmt"
	"os"

	"github.com/donovandicks/gomonkey/internal/interpreter"
	"github.com/donovandicks/gomonkey/internal/lexer"
	"github.com/donovandicks/gomonkey/internal/object"
	"github.com/donovandicks/gomonkey/internal/parser"
)

func main() {
	fileName := os.Args[1]
	if fileName == "" {
		panic("must pass file name")
	}

	input, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	l := lexer.NewLexer(string(input))
	p := parser.NewParser(l)
	prog := p.ParseProgram()
	if errs := p.Errors(); len(errs) != 0 {
		for _, msg := range errs {
			fmt.Printf("ERROR: %s", msg)
		}
		return
	}

	env := object.NewEnv()
	evaled := interpreter.Eval(prog, env)
	if evaled != nil {
		fmt.Printf("%s\n", evaled.Inspect())
	}
}
