package parser_test

import (
	"testing"

	"github.com/donovandicks/gomonkey/internal/ast"
	"github.com/donovandicks/gomonkey/internal/lexer"
	"github.com/donovandicks/gomonkey/internal/parser"
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
			let 1
			`,
			expectedLength: 0,
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
			assert.Equal(t, tc.expectedLength, len(program.Statements))
			assert.Equal(t, tc.expectedErrs, p.Errors())
			for i, s := range program.Statements {
				testLetStatement(t, s, tc.expectedIdents[i])
			}
		})
	}
}

func TestParser_ReturnStatement(t *testing.T) {
	cases := []struct {
		name           string
		input          string
		expectedLength int
		expectedErrs   []string
	}{
		{
			name: "valid return statements",
			input: `
			return 5;
			return 10;
			return 100;
			`,
			expectedLength: 3,
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
			assert.Equal(t, tc.expectedLength, len(program.Statements))
			assert.Equal(t, tc.expectedErrs, p.Errors())

			for _, s := range program.Statements {
				stmt, ok := s.(*ast.ReturnStatement)
				assert.True(t, ok)
				assert.Equal(t, "return", stmt.TokenLiteral())
			}
		})
	}
}
