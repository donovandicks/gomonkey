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
	INST                = "INSTANCE"
)

var (
	Keywords = map[string]TokenType{
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
		"inst":   INST,
	}

	TokenEOF    Token = Token{Type: EOF, Literal: ""}
	TokenSemi         = Token{Type: SEMICOLON, Literal: ";"}
	TokenLParen       = Token{Type: LPAREN, Literal: "("}
	TokenRParen       = Token{Type: RPAREN, Literal: ")"}
	TokenComma        = Token{Type: COMMA, Literal: ","}
	TokenPlus         = Token{Type: PLUS, Literal: "+"}
	TokenMinus        = Token{Type: MINUS, Literal: "-"}
	TokenFSlash       = Token{Type: FSLASH, Literal: "/"}
	TokenStar         = Token{Type: STAR, Literal: "*"}
	TokenLT           = Token{Type: LT, Literal: "<"}
	TokenGT           = Token{Type: GT, Literal: ">"}
	TokenDot          = Token{Type: DOT, Literal: "."}
	TokenLBrace       = Token{Type: LBRACE, Literal: "{"}
	TokenRBrace       = Token{Type: RBRACE, Literal: "}"}
	TokenLBrack       = Token{Type: LBRACK, Literal: "["}
	TokenRBrack       = Token{Type: RBRACK, Literal: "]"}
	TokenColon        = Token{Type: COLON, Literal: ":"}
	TokenAssign       = Token{Type: ASSIGN, Literal: "="}
	TokenEQ           = Token{Type: EQ, Literal: "=="}
	TokenBang         = Token{Type: BANG, Literal: "!"}
	TokenNE           = Token{Type: NE, Literal: "!="}
)

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
func NewIdent(val string) Token     { return Token{Type: IDENT, Literal: val} }
func NewInt(val string) Token       { return Token{Type: INT, Literal: val} }
func NewStr(val string) Token       { return Token{Type: STRING, Literal: val} }

func LookupIdent(ident string) TokenType {
	if tok, ok := Keywords[ident]; ok {
		return tok
	}

	return IDENT
}
