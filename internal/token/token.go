package token

type TokenType string

const (
	ILLEGAL   TokenType = "ILLEGAL"
	EOF                 = "EOF"
	IDENT               = "IDENT"
	INT                 = "INT"
	ASSIGN              = "="
	PLUS                = "+"
	MINUS               = "-"
	STAR                = "*"
	FSLASH              = "/"
	BANG                = "!"
	LT                  = "<"
	GT                  = ">"
	DOT                 = "."
	EQ                  = "=="
	NE                  = "!="
	COMMA               = ","
	SEMICOLON           = ";"
	LPAREN              = "("
	RPAREN              = ")"
	LBRACE              = "{"
	RBRACE              = "}"
	LBRACK              = "["
	RBRACK              = "]"
	COLON               = ":"
	FUNCTION            = "FUNCTION"
	LET                 = "LET"
	RETURN              = "RETURN"
	IF                  = "IF"
	ELSE                = "ELSE"
	TRUE                = "TRUE"
	FALSE               = "FALSE"
	STRING              = "STRING"
	WHILE               = "WHILE"
	FOR                 = "FOR"
	CLASS               = "CLASS"
)

var Keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"return": RETURN,
	"if":     IF,
	"else":   ELSE,
	"true":   TRUE,
	"false":  FALSE,
	"while":  WHILE,
	"for":    FOR,
	"class":  CLASS,
}

type Token struct {
	Type    TokenType
	Literal string
}

func New(tokenType TokenType, literal byte) Token {
	return Token{Type: tokenType, Literal: string(literal)}
}

func NewKeyword(kw string) Token {
	tt, ok := Keywords[kw]
	if !ok {
		panic("invalid keyword")
	}

	return Token{Type: tt, Literal: kw}
}

func NewSpecial(tt TokenType) Token { return Token{Type: tt, Literal: string(tt)} }

func NewIdent(val string) Token { return Token{Type: IDENT, Literal: val} }

func LookupIdent(ident string) TokenType {
	if tok, ok := Keywords[ident]; ok {
		return tok
	}

	return IDENT
}
