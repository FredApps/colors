package tokenizer

import (
	"strconv"

	"github.com/thmalt/colors/internal/ascii"
	"github.com/thmalt/colors/internal/css/token"
)

type Tokenizer struct {
	source

	buf []byte
	tok token.Token
	st  state
}

func New(text string) Tokenizer {
	var t Tokenizer
	t.init(text)
	return t
}

func (t *Tokenizer) init(text string) {
	t.source.init(text)
	t.tok = token.Token{}
	t.st = state{}
}

func (t *Tokenizer) Init(text string) {
	t.init(text)
}

func (t *Tokenizer) State() State { return t.st }

func (t *Tokenizer) Reset(st State) {
	t.st = st
	t.reset(st)
	t.next()
}

// Token position

func (t *Tokenizer) TokenStart() uint { return t.st.pos }
func (t *Tokenizer) TokenEnd() uint   { return t.pos }
func (t *Tokenizer) Segment() string  { return t.SegmentAt(t.TokenStart(), t.TokenEnd()) }

// Token access

func (t *Tokenizer) Token() token.Token { return t.tok }
func (t *Tokenizer) Kind() token.Kind   { return t.tok.Kind }
func (t *Tokenizer) Text() string       { return t.tok.Text }
func (t *Tokenizer) Number() float64    { return t.tok.Number }
func (t *Tokenizer) Delim() byte        { return t.tok.Delim }
func (t *Tokenizer) Flags() token.Flags { return t.tok.Flags }

func (t *Tokenizer) IsBlockStart() bool { return t.tok.IsBlockStart() }

// CloneText makes t.Text independent of the input (idempotent).
func (t *Tokenizer) CloneText() string { return t.tok.CloneText() }

// Unit returns the unit of the current Dimension token.
// Precondition: Kind() == token.Dimension.
func (t *Tokenizer) Unit() string { return t.tok.Text }

func (t *Tokenizer) Is(kind token.Kind) bool { return t.tok.Kind == kind }

// IsIdent reports whether the current token is an [token.Ident] whose text equals name,
// ASCII case-insensitively. name must be lowercase.
func (t *Tokenizer) IsIdent(name string) bool {
	return t.tok.Kind == token.Ident && ascii.EqualFold(t.tok.Text, name)
}

// IsFunction reports whether the current token is an [token.Function] whose text equals name,
// ASCII case-insensitively. name must be lowercase.
func (t *Tokenizer) IsFunction(name string) bool {
	return t.tok.Kind == token.Function && ascii.EqualFold(t.tok.Text, name)
}

func (t *Tokenizer) IsDelim(delim byte) bool {
	return t.tok.Kind == token.Delim && t.tok.Delim == delim
}

// Advancing

func (t *Tokenizer) Next() {
	t.skipWhitespaceAndComments()
	t.next()
}

func (t *Tokenizer) NextIncludingWhitespace() {
	for t.ch == '/' && t.hasAtLeast(1) && t.byteAt(1) == '*' {
		t.advance(2)
		t.skipComment()
	}
	t.next()
}

func (t *Tokenizer) NextIncludingWhitespaceAndComments() {
	t.next()
}

// Matching

func (t *Tokenizer) NextIf(kind token.Kind) bool {
	if t.tok.Kind == kind {
		t.Next()
		return true
	}
	return false
}

func (t *Tokenizer) NextIfIdent(name string) bool {
	if t.tok.Kind == token.Ident && ascii.EqualFold(t.tok.Text, name) {
		t.Next()
		return true
	}
	return false
}

func (t *Tokenizer) NextIfFunction(name string) bool {
	if t.tok.Kind == token.Function && ascii.EqualFold(t.tok.Text, name) {
		t.Next()
		return true
	}
	return false
}

func (t *Tokenizer) NextIfDelim(delim byte) bool {
	if t.tok.Kind == token.Delim && t.tok.Delim == delim {
		t.Next()
		return true
	}
	return false
}

// Block handling

func (t *Tokenizer) BlockEnd() uint {
	switch t.tok.Kind {
	case token.CloseParenthesis, token.CloseSquareBracket, token.CloseCurlyBracket:
		return t.TokenEnd()
	case token.Function, token.OpenParenthesis,
		token.OpenSquareBracket, token.OpenCurlyBracket:
		return t.endBlock(t.TokenEnd(), true)
	default:
		return t.endBlock(t.TokenEnd(), false)
	}
}

func (t *Tokenizer) skipBlock() (end uint) {
	end = t.BlockEnd()
	t.seek(end)
	return
}

// SkipBlock scans to the end of the current block and returns its end position.
// The current token becomes stale; call Next* to continue tokenization.
func (t *Tokenizer) SkipBlock() (end uint) {
	end = t.skipBlock()
	t.tok = token.Token{}
	t.st = t.source.state
	return
}

// SkipBlockNext skips the current block and advances to the next significant token,
// skipping whitespace and comments. It returns the end position of the block,
// before any following whitespace or comments.
func (t *Tokenizer) SkipBlockNext() (end uint) {
	end = t.skipBlock()
	t.Next()
	return
}

// Token setters

func (t *Tokenizer) setKind(kind token.Kind) {
	t.tok.Kind = kind
}

func (t *Tokenizer) setFlags(flags token.Flags) {
	t.tok.Flags |= flags
}

func (t *Tokenizer) clearFlags(flags token.Flags) {
	t.tok.Flags &^= flags
}

func (t *Tokenizer) setText(text string) {
	t.tok.Text = text
}

func (t *Tokenizer) setDelim(delim byte) {
	t.tok.Delim = delim
	t.setKind(token.Delim)
}

// Token scanning

func (t *Tokenizer) next() {
	t.tok = token.Token{}
	t.st = t.source.state

	switch c := t.ch; c {
	case 0:
		// EOF
		if t.Position() >= t.Len() {
			t.setKind(token.EOF)
			return
		}
		t.consumeIdentLike()
		return

	case ' ', '\t', '\n', '\f', '\r':
		t.setText(t.consumeWhitespace())
		t.setKind(token.Whitespace)
		return

	case '[':
		t.nextch()
		t.setKind(token.OpenSquareBracket)
		return
	case ']':
		t.nextch()
		t.setKind(token.CloseSquareBracket)
		return

	case '(':
		t.nextch()
		t.setKind(token.OpenParenthesis)
		return
	case ')':
		t.nextch()
		t.setKind(token.CloseParenthesis)
		return

	case '{':
		t.nextch()
		t.setKind(token.OpenCurlyBracket)
		return
	case '}':
		t.nextch()
		t.setKind(token.CloseCurlyBracket)
		return

	case ',':
		t.nextch()
		t.setKind(token.Comma)
		return
	case ':':
		t.nextch()
		t.setKind(token.Colon)
		return
	case ';':
		t.nextch()
		t.setKind(token.Semicolon)
		return

	case '/':
		t.nextch()
		if t.ch == '*' {
			t.nextch()
			t.setText(t.consumeComment())
			t.setKind(token.Comment)
			return
		}
		t.setDelim(c)
		return

	case '#':
		t.nextch()
		switch {
		case t.isIdentStart():
			t.setText(t.consumeName(t.Position()))
			t.setKind(token.Hash)
			t.setFlags(token.IsHashID)
		case isIdentContinueLUT(t.ch):
			t.setText(t.consumeName(t.Position()))
			t.setKind(token.Hash)
		default:
			t.setDelim(c)
		}
		return

	case '@':
		t.nextch()
		if t.isIdentStart() {
			t.setText(t.consumeName(t.Position()))
			t.setKind(token.AtKeyword)
			return
		}
		t.setDelim(c)
		return

	case '\\':
		if t.hasAtLeast(1) && !t.hasNewlineAt(1) {
			t.consumeIdentLike()
			return
		}
		t.nextch()
		t.setDelim(c)
		return

	case '.':
		t.nextch()
		if isDigit(t.ch) {
			t.consumeNumeric(true)
			return
		}
		t.setDelim(c)
		return

	case '+':
		t.nextch()
		if isDigit(t.ch) {
			t.setFlags(token.HasSign)
			t.consumeNumeric(false)
			return
		}
		if t.ch == '.' && t.hasAtLeast(1) && t.hasDigitsAt(1) {
			t.nextch()

			t.setFlags(token.HasSign)
			t.consumeNumeric(true)
			return
		}
		t.setDelim(c)
		return

	case '-':
		t.nextch()
		switch t.ch {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			t.nextch()
			t.setFlags(token.HasSign)
			t.consumeNumeric(false)
			return

		case '.':
			if t.hasAtLeast(1) && t.hasDigitsAt(1) {
				t.nextch()
				t.setFlags(token.HasSign)
				t.consumeNumeric(true)
				return
			}

		case '-':
			t.nextch()
			if t.ch == '>' {
				t.nextch()
				t.setKind(token.CDC)
				return
			}
			t.setFlags(token.IsCustom)
			t.consumeIdentLike()
			return

		default:
			if isIdentStartLUT(t.ch) || t.ch == '\\' && (!t.hasAtLeast(1) || !t.hasNewlineAt(1)) {
				t.consumeIdentLike()
				return
			}
		}
		t.setDelim(c)
		return

	case '<':
		t.nextch()
		if t.hasAtLeast(2) && t.ch == '!' && t.byteAt(1) == '-' && t.byteAt(2) == '-' {
			t.advance(3)
			t.setKind(token.CDO)
			return
		}
		t.setDelim(c)
		return

	case '"', '\'':
		t.nextch()
		t.consumeString(c)
		return

	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		t.consumeNumeric(false)
		return

	case 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'_':
		t.consumeIdentLike()
		return

	default:
		if c >= 0x80 {
			t.consumeIdentLike()
			return
		}
		t.nextch()
		t.setDelim(c)
		return
	}
}

func (t *Tokenizer) consumeWhitespace() string {
	start := t.Position()
	t.skipWhitespace()
	return t.SegmentAt(start, t.Position())
}

func (t *Tokenizer) consumeComment() string {
	start := t.Position()
	end := t.skipComment()
	return t.SegmentAt(start, end)
}

func (t *Tokenizer) consumeNumeric(seenPoint bool) {
	if !seenPoint {
		t.skipDigits()
		if t.ch == '.' {
			t.nextch()
			seenPoint = true
		}
	}

	if seenPoint {
		t.skipDigits()
	}

	isFloat := seenPoint
	p := t.Position()

	if t.ch|0x20 == 'e' {
		t.nextch()
		if t.isDigit() {
			isFloat = true

			t.skipDigits()
		} else if (t.ch == '+' || t.ch == '-') && t.hasAtLeast(1) && t.hasDigitsAt(1) {
			isFloat = true

			t.nextch()
			t.skipDigits()
		} else {
			t.seek(p)
		}
	}

	t.tok.Number, _ = strconv.ParseFloat(t.Segment(), 64)

	if !isFloat {
		t.tok.Flags |= token.IsInteger
	}

	if t.ch == '%' {
		t.nextch()
		t.tok.Number *= 1 / 100.0
		t.setKind(token.Percentage)
		return
	}

	if t.isIdentStart() {
		t.setText(t.consumeName(t.Position()))
		t.setKind(token.Dimension)
		return
	}

	t.setKind(token.Number)
}

func (t *Tokenizer) consumeIdentLike() {
	t.setText(t.consumeName(t.TokenStart()))
	if t.ch == '(' {
		t.nextch()
		if ascii.EqualFold(t.tok.Text, "url") {
			if t.consumeUrl() {
				t.clearFlags(token.IsDecoded)
				return
			}
		}

		t.setKind(token.Function)
		return
	}

	t.setKind(token.Ident)
}

func (t *Tokenizer) consumeName(start uint) string {
	var text string
	text, t.buf = t.scanName(t.buf[:0], start)
	if len(t.buf) > 0 {
		t.setFlags(token.IsDecoded)
	}
	return text
}

func (t *Tokenizer) consumeString(quote byte) bool {
	var (
		text string
		ok   bool
	)

	text, t.buf, ok = t.scanString(t.buf[:0], quote)
	t.setText(text)

	if len(t.buf) > 0 {
		t.setFlags(token.IsDecoded)
	}

	if ok {
		t.setKind(token.String)
	} else {
		t.setKind(token.BadString)
	}

	return ok
}

// consumeUrl returns true if it emitted a Url or BadUrl token,
// false if the url has a quoted argument (caller emits Function).
func (t *Tokenizer) consumeUrl() (isFunc bool) {
	p := t.wsEnd()

	if p < t.Len() && (t.at(p) == '"' || t.at(p) == '\'') {
		return false
	}
	t.seek(p)

	var (
		text string
		ok   bool
	)

	text, t.buf, ok = t.scanUrl(t.buf[:0])
	t.setText(text)

	if !ok {
		t.setKind(token.BadUrl)
		return true
	}

	if len(t.buf) > 0 {
		t.setFlags(token.IsDecoded)
	}
	t.setKind(token.Url)
	return true
}
