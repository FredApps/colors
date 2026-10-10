package token

import (
	"fmt"
	"strings"
)

type Token struct {
	Text   string
	Number float64
	Delim  byte

	Kind Kind
	Flags
}

// Unit returns the unit of the current Dimension token.
// Precondition: Kind() == token.Dimension.
func (t Token) Unit() string { return t.Text }

// CloneText makes t.Text independent of the input (idempotent).
func (t *Token) CloneText() string {
	if !t.OwnsText() {
		t.Text = strings.Clone(t.Text)
		t.Flags |= textOwned
	}
	return t.Text
}

func (t Token) Is(kind Kind) bool { return t.Kind == kind }
func (t Token) IsBlockStart() bool {
	switch t.Kind {
	case OpenSquareBracket, OpenParenthesis, OpenCurlyBracket, Function:
		return true
	default:
		return false
	}
}

func (t Token) String() string {
	switch t.Kind {
	case None, EOF, CDO, CDC, Colon, Semicolon, Comma,
		OpenParenthesis, OpenSquareBracket, OpenCurlyBracket,
		CloseParenthesis, CloseSquareBracket, CloseCurlyBracket:
		return t.Kind.String()
	case Delim:
		return fmt.Sprintf("%s(%q)", t.Kind, t.Delim)
	case Number, Percentage:
		return fmt.Sprintf("%s(%g%s)", t.Kind, t.Number, t.Flags.string(true, false))
	case Dimension:
		return fmt.Sprintf("%s(%g %s%s)", t.Kind, t.Number, t.Text, t.Flags.string(true, false))
	default:
		return fmt.Sprintf("%s(%q%s)", t.Kind, t.Text, t.Flags.string(true, false))
	}
}
