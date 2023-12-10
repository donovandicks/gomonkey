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

type ErrNoPrefixParser struct {
	operator string
}

func (e ErrNoPrefixParser) Error() string {
	return fmt.Sprintf("no prefix parser found for %s", e.operator)
}

type ErrMissingOpener struct {
	expected string
}

func (e ErrMissingOpener) Error() string {
	return fmt.Sprintf("missing opening '%s'", e.expected)
}

type ErrMissingCloser struct {
	expected string
}

func (e ErrMissingCloser) Error() string {
	return fmt.Sprintf("missing closing '%s'", e.expected)
}
