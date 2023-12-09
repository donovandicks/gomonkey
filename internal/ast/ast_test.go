package ast_test

import (
	"testing"

	"github.com/donovandicks/gomonkey/internal/ast"
	"github.com/donovandicks/gomonkey/internal/token"
	"github.com/stretchr/testify/assert"
)

func TestAST_ToString(t *testing.T) {
	cases := []struct {
		name     string
		input    *ast.Program
		expected string
	}{
		{
			name:     "valid program",
			expected: "let foo = bar;",
			input: &ast.Program{
				Statements: []ast.Statement{
					&ast.LetStatement{
						Token: token.Token{Type: token.LET, Literal: "let"},
						Name: &ast.Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "foo"},
							Value: "foo",
						},
						Value: &ast.Identifier{
							Token: token.Token{Type: token.IDENT, Literal: "bar"},
							Value: "bar",
						},
					},
				},
			},
		},
	}

	for _, testCase := range cases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.expected, tc.input.String())
		})
	}
}
