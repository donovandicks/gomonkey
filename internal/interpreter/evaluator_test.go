package interpreter_test

import (
	"testing"

	"github.com/donovandicks/gomonkey/internal/interpreter"
	"github.com/donovandicks/gomonkey/internal/lexer"
	"github.com/donovandicks/gomonkey/internal/object"
	"github.com/donovandicks/gomonkey/internal/parser"
	"github.com/stretchr/testify/assert"
)

func TestEvaluator(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		input  string
		output object.Object
	}{
		{
			name:   "integer literal",
			input:  "5",
			output: object.NewIntegerObject(5),
		},
	}

	for _, testCase := range cases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := lexer.NewLexer(tc.input)
			p := parser.NewParser(l)
			prog := p.ParseProgram()

			assert.Equal(t, tc.output, interpreter.Eval(prog))
		})
	}
}
