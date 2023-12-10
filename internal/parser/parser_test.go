package parser_test

import (
	"testing"

	"github.com/donovandicks/gomonkey/internal/ast"
	"github.com/donovandicks/gomonkey/internal/lexer"
	"github.com/donovandicks/gomonkey/internal/parser"
	"github.com/donovandicks/gomonkey/internal/token"
	"github.com/stretchr/testify/assert"
)

func TestParser(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name         string
		input        string
		expected     []ast.Statement
		expectedErrs []string
	}{
		{
			name:  "valid identifier",
			input: "foo;",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.IDENT, Literal: "foo"},
					Expression: &ast.Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "foo"},
						Value: "foo",
					},
				},
			},
		},
		{
			name:  "valid integer literal",
			input: "5;",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Expression: &ast.IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Value: 5,
					},
				},
			},
		},
		{
			name:  "prefix bang",
			input: "!5;",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.BANG, Literal: "!"},
					Expression: &ast.PrefixExpression{
						Token:    token.Token{Type: token.BANG, Literal: "!"},
						Operator: "!",
						Right: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
		},
		{
			name:  "prefix minus",
			input: "-5;",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.MINUS, Literal: "-"},
					Expression: &ast.PrefixExpression{
						Token:    token.Token{Type: token.MINUS, Literal: "-"},
						Operator: "-",
						Right: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
		},
		{
			name:  "infix addition",
			input: "5 + 5",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Expression: &ast.InfixExpression{
						Token:    token.Token{Type: token.PLUS, Literal: "+"},
						Operator: "+",
						Left: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
						Right: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
		},
		{
			name:  "infix subtraction",
			input: "5 - 5",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Expression: &ast.InfixExpression{
						Token:    token.Token{Type: token.MINUS, Literal: "-"},
						Operator: "-",
						Left: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
						Right: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
		},
		{
			name:  "infix multiplication",
			input: "5 * 5",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Expression: &ast.InfixExpression{
						Token:    token.Token{Type: token.STAR, Literal: "*"},
						Operator: "*",
						Left: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
						Right: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
		},
		{
			name:  "infix division",
			input: "5 / 5",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Expression: &ast.InfixExpression{
						Token:    token.Token{Type: token.FSLASH, Literal: "/"},
						Operator: "/",
						Left: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
						Right: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
		},
		{
			name:  "infix equals",
			input: "5 == 5",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Expression: &ast.InfixExpression{
						Token:    token.Token{Type: token.EQ, Literal: "=="},
						Operator: "==",
						Left: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
						Right: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
		},
		{
			name:  "infix not equals",
			input: "5 != 5",
			expected: []ast.Statement{
				&ast.ExpressionStatement{
					Token: token.Token{Type: token.INT, Literal: "5"},
					Expression: &ast.InfixExpression{
						Token:    token.Token{Type: token.NE, Literal: "!="},
						Operator: "!=",
						Left: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
						Right: &ast.IntegerLiteral{
							Token: token.Token{Type: token.INT, Literal: "5"},
							Value: 5,
						},
					},
				},
			},
		},
		{
			name:  "valid return statement",
			input: "return 5;",
			expected: []ast.Statement{
				&ast.ReturnStatement{
					Token: token.Token{Type: token.RETURN, Literal: "return"},
					Value: &ast.IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Value: 5,
					},
				},
			},
		},
		{
			name:  "valid let statement",
			input: "let x = 5;",
			expected: []ast.Statement{
				&ast.LetStatement{
					Token: token.Token{Type: token.LET, Literal: "let"},
					Name: &ast.Identifier{
						Token: token.Token{Type: token.IDENT, Literal: "x"},
						Value: "x",
					},
					Value: &ast.IntegerLiteral{
						Token: token.Token{Type: token.INT, Literal: "5"},
						Value: 5,
					},
				},
			},
		},
		{
			name:  "invalid let statement: missing assignment operator",
			input: "let x 5;",
			expectedErrs: []string{
				"expected next token to be =, got INT instead",
			},
		},
		{
			name:         "invalid let statement: missing identifier",
			input:        "let = 10;",
			expectedErrs: []string{"expected next token to be IDENT, got = instead"},
		},
		{
			name:         "invalid let statement: missing identifier and assigner",
			input:        "let 10;",
			expectedErrs: []string{"expected next token to be IDENT, got INT instead"},
		},
	}

	for _, testCase := range cases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := lexer.NewLexer(tc.input)
			p := parser.NewParser(l)

			program := p.ParseProgram()
			assert.NotNil(t, program)
			if tc.expectedErrs == nil {
				assert.Nil(t, p.Errors())
			} else {
				assert.Equal(t, tc.expectedErrs, p.Errors())
			}

			if tc.expected != nil {
				assert.EqualValues(t, tc.expected, program.Statements)
			}
		})
	}
}

func TestParser_OperatorPrecedence(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "unary mixed with binary",
			input:    "-a * b",
			expected: "((-a) * b)",
		},
		{
			name:     "multiple unary",
			input:    "!-a",
			expected: "(!(-a))",
		},
		{
			name:     "multiple addition",
			input:    "a + b + c",
			expected: "((a + b) + c)",
		},
		{
			name:     "addition and subtraction",
			input:    "a + b - c",
			expected: "((a + b) - c)",
		},
		{
			name:     "multiple multiplication",
			input:    "a * b * c",
			expected: "((a * b) * c)",
		},
		{
			name:     "multiplication and division",
			input:    "a * b / c",
			expected: "((a * b) / c)",
		},
		{
			name:     "multiple in sequence",
			input:    "a - b * c / d * e + f",
			expected: "((a - (((b * c) / d) * e)) + f)",
		},
		{
			name:     "multiple statements",
			input:    "a + b; x * y",
			expected: "(a + b)(x * y)",
		},
		{
			name:     "comparison operations",
			input:    "a > b == c > d",
			expected: "((a > b) == (c > d))",
		},
		{
			name:     "comparison with unary and binary operations",
			input:    "a * -b != -c / d",
			expected: "((a * (-b)) != ((-c) / d))",
		},
	}

	for _, testCase := range cases {
		tc := testCase

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := lexer.NewLexer(tc.input)
			p := parser.NewParser(l)

			program := p.ParseProgram()
			assert.NotNil(t, program)

			assert.Equal(t, tc.expected, program.String())
		})
	}
}
