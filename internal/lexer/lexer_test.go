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
				token.TokenAssign,
				token.TokenPlus,
				token.TokenLParen,
				token.TokenRParen,
				token.TokenLBrace,
				token.TokenRBrace,
				token.TokenLBrack,
				token.TokenRBrack,
				token.TokenComma,
				token.TokenSemi,
				token.TokenBang,
				token.TokenMinus,
				token.TokenFSlash,
				token.TokenStar,
				token.TokenLT,
				token.TokenGT,
				token.TokenColon,
				token.TokenDot,
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
				token.TokenAssign,
				token.TokenEQ,
				token.TokenBang,
				token.TokenNE,
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
				token.TokenAssign,
				token.NewInt("5"),
				token.TokenSemi,
				token.NewKeyword("let"),
				token.NewIdent("ten"),
				token.TokenAssign,
				token.NewInt("10"),
				token.TokenSemi,
				token.NewKeyword("let"),
				token.NewIdent("add"),
				token.TokenAssign,
				token.NewKeyword("fn"),
				token.TokenLParen,
				token.NewIdent("x"),
				token.TokenComma,
				token.NewIdent("y"),
				token.TokenRParen,
				token.TokenLBrace,
				token.NewIdent("x"),
				token.TokenPlus,
				token.NewIdent("y"),
				token.TokenSemi,
				token.TokenRBrace,
				token.TokenSemi,
				token.NewKeyword("let"),
				token.NewIdent("result"),
				token.TokenAssign,
				token.NewIdent("add"),
				token.TokenLParen,
				token.NewIdent("five"),
				token.TokenComma,
				token.NewIdent("ten"),
				token.TokenRParen,
				token.TokenSemi,
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
