package lexer

import (
	"github.com/donovandicks/gomonkey/internal/token"
)

type Lexer struct {
	input   string
	pos     int  // current position in input (current char)
	readPos int  // current reading position in input (next char)
	ch      byte // current char
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
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

func (l *Lexer) readString() string {
	pos := l.pos + 1 // after the starting quote

	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}

	return l.input[pos:l.pos]
}

// skipWhitespace advances the lexer over any whitespace characters
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
		l.readChar()
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peek() == '=' {
			l.readChar()
			tok = token.NewSpecial(token.EQ)
		} else {
			tok = token.NewSpecial(token.ASSIGN)
		}
	case ';':
		tok = token.NewSpecial(token.SEMICOLON)
	case '(':
		tok = token.NewSpecial(token.LPAREN)
	case ')':
		tok = token.NewSpecial(token.RPAREN)
	case ',':
		tok = token.NewSpecial(token.COMMA)
	case '+':
		tok = token.NewSpecial(token.PLUS)
	case '-':
		tok = token.NewSpecial(token.MINUS)
	case '/':
		tok = token.NewSpecial(token.FSLASH)
	case '*':
		tok = token.NewSpecial(token.STAR)
	case '<':
		tok = token.NewSpecial(token.LT)
	case '>':
		tok = token.NewSpecial(token.GT)
	case '.':
		tok = token.NewSpecial(token.DOT)
	case '!':
		if l.peek() == '=' {
			l.readChar()
			tok = token.NewSpecial(token.NE)
		} else {
			tok = token.NewSpecial(token.BANG)
		}
	case '{':
		tok = token.NewSpecial(token.LBRACE)
	case '}':
		tok = token.NewSpecial(token.RBRACE)
	case '[':
		tok = token.NewSpecial(token.LBRACK)
	case ']':
		tok = token.NewSpecial(token.RBRACK)
	case '"':
		tok = token.NewStr(l.readString())
	case ':':
		tok = token.NewSpecial(token.COLON)
	case 0:
		tok = token.NewEOF()
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
