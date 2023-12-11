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
		{
			name:   "if expression: truthy condition",
			input:  "if (10) { 10 } else { -10 }",
			output: object.NewIntegerObject(10),
		},
		{
			name:   "if expression: literal true condition",
			input:  "if (true) { 10 }",
			output: object.NewIntegerObject(10),
		},
		{
			name:   "if expression: else branch",
			input:  "if (false) { 10 } else { -10 }",
			output: object.NewIntegerObject(-10),
		},
		{
			name:   "if expression: null return",
			input:  "if (false) { 10 }",
			output: object.NullObject,
		},
		{
			name:   "return: top-level",
			input:  "return 5;",
			output: object.NewIntegerObject(5),
		},
		{
			name:   "return: mid block",
			input:  "1 + 5; return 3; 7 * 7",
			output: object.NewIntegerObject(3),
		},
		{
			name: "return: nested blocks",
			input: `
			if (true) {
				if (true) {
					return 10;
				}

				return 5;
			}
			`,
			output: object.NewIntegerObject(10),
		},
		{
			name:   "let binding: simple assignment",
			input:  "let x = 5; x;",
			output: object.NewIntegerObject(5),
		},
		{
			name:   "let binding: expression assignment",
			input:  "let x = 5 * 5; x;",
			output: object.NewIntegerObject(25),
		},
		{
			name:   "let binding: transitive assignment",
			input:  "let x = 5; let y = x; y;",
			output: object.NewIntegerObject(5),
		},
		{
			name:   "let binding: identifier operations",
			input:  "let x = 5; let y = 10; x + y;",
			output: object.NewIntegerObject(15),
		},
	}

	for _, testCase := range cases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := lexer.NewLexer(tc.input)
			p := parser.NewParser(l)
			prog := p.ParseProgram()
			env := object.NewEnv()

			assert.Equal(t, tc.output, interpreter.Eval(prog, env))
		})
	}
}

func TestEvaluator_Errors(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		input string
		err   *object.Err
	}{
		{
			name:  "type error: binary operator",
			input: "5 + true; 5;",
			err:   &object.Err{Msg: "type error: cannot perform '+' on INTEGER, BOOLEAN"},
		},
		{
			name:  "unknown operator: unary expression",
			input: "-true",
			err:   &object.Err{Msg: "invalid operator '-' for type BOOLEAN"},
		},
		{
			name:  "unknown operator: infix expression",
			input: "true + false",
			err:   &object.Err{Msg: "unknown operator '+' for types BOOLEAN, BOOLEAN"},
		},
		{
			name:  "let binding: unbound identifier",
			input: "x;",
			err:   &object.Err{Msg: "undefined identifier 'x'"},
		},
	}

	for _, testCase := range cases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := lexer.NewLexer(tc.input)
			p := parser.NewParser(l)
			prog := p.ParseProgram()
			env := object.NewEnv()

			assert.Equal(t, tc.err, interpreter.Eval(prog, env))
		})
	}
}
