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
			input: "=+(){}[],;!-/*<>:.",
			expTokens: []token.Token{
				token.NewSpecial(token.ASSIGN),
				token.NewSpecial(token.PLUS),
				token.NewSpecial(token.LPAREN),
				token.NewSpecial(token.RPAREN),
				token.NewSpecial(token.LBRACE),
				token.NewSpecial(token.RBRACE),
				token.NewSpecial(token.LBRACK),
				token.NewSpecial(token.RBRACK),
				token.NewSpecial(token.COMMA),
				token.NewSpecial(token.SEMICOLON),
				token.NewSpecial(token.BANG),
				token.NewSpecial(token.MINUS),
				token.NewSpecial(token.FSLASH),
				token.NewSpecial(token.STAR),
				token.NewSpecial(token.LT),
				token.NewSpecial(token.GT),
				token.NewSpecial(token.COLON),
				token.NewSpecial(token.DOT),
			},
		},
		{
			name:  "keywords",
			input: "fn let return if else true false while for class",
			expTokens: []token.Token{
				token.NewKeyword("fn"),
				token.NewKeyword("let"),
				token.NewKeyword("return"),
				token.NewKeyword("if"),
				token.NewKeyword("else"),
				token.NewKeyword("true"),
				token.NewKeyword("false"),
				token.NewKeyword("while"),
				token.NewKeyword("for"),
				token.NewKeyword("class"),
			},
		},
		{
			name:  "equals and not equals",
			input: "= == ! !=",
			expTokens: []token.Token{
				token.NewSpecial(token.ASSIGN),
				token.NewSpecial(token.EQ),
				token.NewSpecial(token.BANG),
				token.NewSpecial(token.NE),
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
				token.NewKeyword("let"),
				token.NewIdent("five"),
				token.NewSpecial(token.ASSIGN),
				token.NewInt("5"),
				token.NewSpecial(token.SEMICOLON),
				token.NewKeyword("let"),
				token.NewIdent("ten"),
				token.NewSpecial(token.ASSIGN),
				token.NewInt("10"),
				token.NewSpecial(token.SEMICOLON),
				token.NewKeyword("let"),
				token.NewIdent("add"),
				token.NewSpecial(token.ASSIGN),
				token.NewKeyword("fn"),
				token.NewSpecial(token.LPAREN),
				token.NewIdent("x"),
				token.NewSpecial(token.COMMA),
				token.NewIdent("y"),
				token.NewSpecial(token.RPAREN),
				token.NewSpecial(token.LBRACE),
				token.NewIdent("x"),
				token.NewSpecial(token.PLUS),
				token.NewIdent("y"),
				token.NewSpecial(token.SEMICOLON),
				token.NewSpecial(token.RBRACE),
				token.NewSpecial(token.SEMICOLON),
				token.NewKeyword("let"),
				token.NewIdent("result"),
				token.NewSpecial(token.ASSIGN),
				token.NewIdent("add"),
				token.NewSpecial(token.LPAREN),
				token.NewIdent("five"),
				token.NewSpecial(token.COMMA),
				token.NewIdent("ten"),
				token.NewSpecial(token.RPAREN),
				token.NewSpecial(token.SEMICOLON),
				token.TokenEOF,
			},
		},
		{
			name:  "strings",
			input: `"hello, world!" "one"`,
			expTokens: []token.Token{
				token.NewStr("hello, world!"),
				token.NewStr("one"),
			},
		},
		{
			name:  "strings: multiple same",
			input: `"hello" "hello" "hello"`,
			expTokens: []token.Token{
				token.NewStr("hello"),
				token.NewStr("hello"),
				token.NewStr("hello"),
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
