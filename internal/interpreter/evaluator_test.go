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
		{
			name:   "boolean literal",
			input:  "true",
			output: object.TrueBool,
		},
		{
			name:   "bang operator: bool literal",
			input:  "!true",
			output: object.FalseBool,
		},
		{
			name:   "bang operator: integer literal",
			input:  "!5",
			output: object.FalseBool,
		},
		{
			name:   "bang operator: double operator",
			input:  "!!true",
			output: object.TrueBool,
		},
		{
			name:   "minus operator: integer literal",
			input:  "-5",
			output: object.NewIntegerObject(-5),
		},
		{
			name:   "infix expression: integer addition",
			input:  "5 + 5",
			output: object.NewIntegerObject(10),
		},
		{
			name:   "infix expression: integer subtraction",
			input:  "5 - 5",
			output: object.NewIntegerObject(0),
		},
		{
			name:   "infix expression: integer multiplication",
			input:  "5 * 5",
			output: object.NewIntegerObject(25),
		},
		{
			name:   "infix expression: integer division",
			input:  "5 / 5",
			output: object.NewIntegerObject(1),
		},
		{
			name:   "comparison: integer less than",
			input:  "5 < 6",
			output: object.TrueBool,
		},
		{
			name:   "comparison: integer greater than",
			input:  "2 > 1",
			output: object.TrueBool,
		},
		{
			name:   "comparison: integer equality",
			input:  "2 == 3",
			output: object.FalseBool,
		},
		{
			name:   "comparison: integer inequality",
			input:  "2 != 1",
			output: object.TrueBool,
		},
		{
			name:   "comparison: boolean equality",
			input:  "true == true",
			output: object.TrueBool,
		},
		{
			name:   "comparison: boolean inequality",
			input:  "true != true",
			output: object.FalseBool,
		},
		{
			name:   "comparison: integer with boolean equality",
			input:  "(1 < 2) == true",
			output: object.TrueBool,
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
