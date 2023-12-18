package lexer

import (
	"github.com/donovandicks/gomonkey/internal/token"
)

type Lexer struct {
	input       string
	stringCache map[string]token.Token
	pos         int  // current position in input (current char)
	readPos     int  // current reading position in input (next char)
	ch          byte // current char
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:       input,
		stringCache: make(map[string]token.Token),
	}
	l.readChar()
	return l
}

func (l *Lexer) peek() byte {
	if l.readPos >= len(l.input) {
		return 0
	}

	return l.input[l.readPos]
}

func (l *Lexer) readChar() {
	l.ch = l.peek()
	l.pos = l.readPos
	l.readPos += 1
}

// readIdentifier reads an entire word at a time.
//
// The starting location of the word is the current lexer position at the time
// of call. The end is determined by advancing over the input until a non-letter
// byte is encountered. The characters between the start and end positions are
// returned as a single identifier.
func (l *Lexer) readIdentifier() string {
	pos := l.pos
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}

	return l.input[pos:l.pos]
}

func (l *Lexer) readNumber() string {
	pos := l.pos
	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[pos:l.pos]
}

func (l *Lexer) readString() token.Token {
	pos := l.pos + 1 // after the starting quote

	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}

	str := l.input[pos:l.pos]
	if t, ok := l.stringCache[str]; ok {
		return t
	}

	return token.NewStr(str)
}

// skipWhitespace advances the lexer over any whitespace characters
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
		l.readChar()
	}
}

func (l *Lexer) readSpecial(ch string) token.Token {
	var tok token.Token
	switch ch {
	case ";":
		tok = token.TokenSemi
	case "(":
		tok = token.TokenLParen
	case ")":
		tok = token.TokenRParen
	case ",":
		tok = token.TokenComma
	case "+":
		tok = token.TokenPlus
	case "-":
		tok = token.TokenMinus
	case "/":
		tok = token.TokenFSlash
	case "*":
		tok = token.TokenStar
	case "<":
		tok = token.TokenLT
	case ">":
		tok = token.TokenGT
	case ".":
		tok = token.TokenDot
	case "{":
		tok = token.TokenLBrace
	case "}":
		tok = token.TokenRBrace
	case "[":
		tok = token.TokenLBrack
	case "]":
		tok = token.TokenRBrack
	case ":":
		tok = token.TokenColon
	case "=":
		tok = token.TokenAssign
	case "==":
		tok = token.TokenEQ
	case "!":
		tok = token.TokenBang
	case "!=":
		tok = token.TokenNE
	}

	return tok
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peek() == '=' {
			l.readChar()
			tok = l.readSpecial("==")
		} else {
			tok = l.readSpecial("=")
		}
	case '!':
		if l.peek() == '=' {
			l.readChar()
			tok = l.readSpecial("!=")
		} else {
			tok = l.readSpecial("!")
		}
	case ';', '(', ')', ',', '+', '-', '/', '*', '<', '>', '.', '{', '}', '[', ']', ':':
		tok = l.readSpecial(string(l.ch))
	case '"':
		tok = l.readString()
	case 0:
		tok = token.TokenEOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			// exit early because the lexer has already been advanced in readIdentifier
			return tok
		} else if isDigit(l.ch) {
			tok = token.NewInt(l.readNumber())
			return tok
		} else {
			tok = token.New(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}
