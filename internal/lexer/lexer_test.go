package lexer_test

import (
	"testing"

	"github.com/donovandicks/gomonkey/internal/lexer"
	"github.com/donovandicks/gomonkey/internal/token"
	"github.com/stretchr/testify/assert"
)

func TestNextToken(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		expTokens []token.Token
	}{
		{
			name:  "special characters",
			input: "=+(){}[],;!-/*<>:",
			expTokens: []token.Token{
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.PLUS, Literal: "+"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.LBRACK, Literal: "["},
				{Type: token.RBRACK, Literal: "]"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.BANG, Literal: "!"},
				{Type: token.MINUS, Literal: "-"},
				{Type: token.FSLASH, Literal: "/"},
				{Type: token.STAR, Literal: "*"},
				{Type: token.LT, Literal: "<"},
				{Type: token.GT, Literal: ">"},
				{Type: token.COLON, Literal: ":"},
			},
		},
		{
			name:  "keywords",
			input: "fn let return if else true false while",
			expTokens: []token.Token{
				{Type: token.FUNCTION, Literal: "fn"},
				{Type: token.LET, Literal: "let"},
				{Type: token.RETURN, Literal: "return"},
				{Type: token.IF, Literal: "if"},
				{Type: token.ELSE, Literal: "else"},
				{Type: token.TRUE, Literal: "true"},
				{Type: token.FALSE, Literal: "false"},
				{Type: token.WHILE, Literal: "while"},
			},
		},
		{
			name:  "equals and not equals",
			input: "= == ! !=",
			expTokens: []token.Token{
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.EQ, Literal: "=="},
				{Type: token.BANG, Literal: "!"},
				{Type: token.NE, Literal: "!="},
			},
		},
		{
			name: "source code program",
			input: `let five = 5;
			let ten = 10;

			let add = fn(x, y) {
				x + y;
			};

			let result = add(five, ten);
			`,
			expTokens: []token.Token{
				{Type: token.LET, Literal: "let"},
				{Type: token.IDENT, Literal: "five"},
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.INT, Literal: "5"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.LET, Literal: "let"},
				{Type: token.IDENT, Literal: "ten"},
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.INT, Literal: "10"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.LET, Literal: "let"},
				{Type: token.IDENT, Literal: "add"},
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.FUNCTION, Literal: "fn"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.IDENT, Literal: "x"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.IDENT, Literal: "y"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.LBRACE, Literal: "{"},
				{Type: token.IDENT, Literal: "x"},
				{Type: token.PLUS, Literal: "+"},
				{Type: token.IDENT, Literal: "y"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.RBRACE, Literal: "}"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.LET, Literal: "let"},
				{Type: token.IDENT, Literal: "result"},
				{Type: token.ASSIGN, Literal: "="},
				{Type: token.IDENT, Literal: "add"},
				{Type: token.LPAREN, Literal: "("},
				{Type: token.IDENT, Literal: "five"},
				{Type: token.COMMA, Literal: ","},
				{Type: token.IDENT, Literal: "ten"},
				{Type: token.RPAREN, Literal: ")"},
				{Type: token.SEMICOLON, Literal: ";"},
				{Type: token.EOF, Literal: ""},
			},
		},
		{
			name:  "strings",
			input: `"hello, world!" "one"`,
			expTokens: []token.Token{
				{Type: token.STRING, Literal: "hello, world!"},
				{Type: token.STRING, Literal: "one"},
			},
		},
	}

	for _, testCase := range cases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l := lexer.NewLexer(tc.input)
			for _, exp := range tc.expTokens {
				tok := l.NextToken()
				assert.Equal(t, exp.Type, tok.Type)
				assert.Equal(t, exp.Literal, tok.Literal)
			}
		})
	}
}
