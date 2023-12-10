package parser_test

import (
	"testing"

	"github.com/donovandicks/gomonkey/internal/ast"
	"github.com/donovandicks/gomonkey/internal/lexer"
	"github.com/donovandicks/gomonkey/internal/parser"
	"github.com/donovandicks/gomonkey/internal/token"
	"github.com/stretchr/testify/assert"
)

func testLetStatement(t *testing.T, s ast.Statement, name string) {
	assert.Equal(t, "let", s.TokenLiteral())
	stmt, ok := s.(*ast.LetStatement)
	assert.True(t, ok)
	assert.Equal(t, name, stmt.Name.Value)
	assert.Equal(t, name, stmt.Name.TokenLiteral())
}

func testParserErrors(t *testing.T, p *parser.Parser, expectedErrs []string) {
	errs := p.Errors()
	assert.Equal(t, expectedErrs, errs)
}

func TestParser_LetStatement(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		input          string
		expectedLength int
		expectedIdents []string
		expectedErrs   []string
	}{
		{
			name: "valid let statements",
			input: `
			let x = 5;
			let y = 10;
			let foobar = 1337;
			`,
			expectedLength: 3,
			expectedIdents: []string{"x", "y", "foobar"},
		},
		{
			name: "invalid let statements",
			input: `
			let x 5;
			let = 10;
			let 1;
			`,
			expectedErrs: []string{
				"expected next token to be =, got INT instead",
				"expected next token to be IDENT, got = instead",
				"expected next token to be IDENT, got INT instead",
			},
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
			assert.Equal(t, tc.expectedErrs, p.Errors(), "mismatched error count")
			if len(tc.expectedIdents) > 0 {
				for i, s := range program.Statements {
					testLetStatement(t, s, tc.expectedIdents[i])
				}
			}
		})
	}
}

func TestParser_ExpressionStatements(t *testing.T) {
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
	}

	for _, testCase := range cases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := lexer.NewLexer(tc.input)
			p := parser.NewParser(l)

			program := p.ParseProgram()
			assert.NotNil(t, program)
			assert.Nil(t, p.Errors())

			assert.EqualValues(t, tc.expected, program.Statements)
			assert.Equal(t, tc.expectedErrs, p.Errors())
		})
	}
}
