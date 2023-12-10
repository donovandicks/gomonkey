package parser

import (
	"fmt"

	"github.com/donovandicks/gomonkey/internal/token"
)

type ErrNextTokenInvalid struct {
	expected token.TokenType
	actual   token.TokenType
}

func (e ErrNextTokenInvalid) Error() string {
	return fmt.Sprintf("expected next token to be %s, got %s instead", e.expected, e.actual)
}

type ErrParseError struct {
	expected string
	actual   string
}

func (e ErrParseError) Error() string {
	return fmt.Sprintf("error attempting to parse %s as a valid %s", e.actual, e.expected)
}
