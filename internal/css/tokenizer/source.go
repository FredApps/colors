package tokenizer

import (
	"strings"
	"unicode/utf8"

	"github.com/thmalt/colors/internal/ascii"
	"github.com/thmalt/colors/internal/hexlut"
)

const replacementCharStr = "\uFFFD"

type state struct {
	pos uint
}

func (s state) Position() uint { return s.pos }

type State = state

type source struct {
	text string

	state

	ch byte
}

func (s *source) init(text string) {
	s.text = text

	s.seek(0)
}

func (s *source) reset(state state) {
	s.state = state
	s.seek(s.pos)
}

// func (s *source) error(msg string) {}

// Len returns the length of the source in bytes.
func (s *source) Len() uint { return uint(len(s.text)) }

// SegmentAt returns the source segment in [start, end).
// It panics if start > end or either position is out of range.
func (s *source) SegmentAt(start, end uint) string { return s.text[start:end] }

func (s *source) at(p uint) byte         { return s.text[p] }
func (s *source) byteAt(off uint) byte   { return s.text[s.pos+off] }
func (s *source) hasAtLeast(n uint) bool { return s.pos+n < uint(len(s.text)) }

func (s *source) seek(p uint) {
	n := uint(len(s.text))
	if p >= n {
		s.pos, s.ch = n, 0
		return
	}
	s.pos, s.ch = p, s.text[p]
}

func (s *source) nextch()        { s.seek(s.pos + 1) }
func (s *source) advance(n uint) { s.seek(s.pos + n) }
func (s *source) back(n uint)    { s.seek(s.pos - n) }

func (s *source) wsEnd() uint {
	text := s.text
	p, n := s.pos, uint(len(text))
	for p < n && wsLUT[text[p]] {
		p++
	}
	return p
}

func (s *source) skipWhitespace() { s.seek(s.wsEnd()) }

func (s *source) skipWhitespaceAndComments() {
	text := s.text
	p, n := s.pos, uint(len(text))

	for {
		// skip whitespace
		for p < n && wsLUT[text[p]] {
			p++
		}
		if p >= n {
			break
		}
		rest := text[p:]
		if len(rest) < 2 || rest[0] != '/' || rest[1] != '*' {
			break
		}
		// skip comment
		i := strings.Index(rest[2:], "*/")
		if i < 0 {
			p = n
			break
		}
		p += uint(i) + 4
	}

	s.seek(p)
}

func (s *source) skipComment() (end uint) {
	end = uint(len(s.text))
	p := end
	if i := strings.Index(s.text[s.pos:], "*/"); i >= 0 {
		end = s.pos + uint(i)
		p = end + 2
	}
	s.seek(p)
	return
}

func (s *source) skipDigits() {
	text := s.text
	p, n := s.pos, uint(len(text))
	for p < n && isDigit(text[p]) {
		p++
	}
	s.seek(p)
}

func (s *source) appendEscape(buf []byte, p uint) ([]byte, uint) {
	text := s.text
	n := uint(len(text))

	c := text[p]
	if d := hexlut.Dec[c]; d < hexlut.NibbleLimit {
		r := rune(d)
		p++

		for i := 1; i < 6 && p < n; i++ {
			d := hexlut.Dec[text[p]]
			if d >= hexlut.NibbleLimit {
				break
			}
			r = r<<4 | rune(d)
			p++
		}

		// hex digits delimiter
		if p < n && wsLUT[text[p]] {
			if text[p] == '\r' && p+1 < n && text[p+1] == '\n' {
				p++
			}
			p++
		}

		if r == 0 {
			return append(buf, replacementCharStr...), p
		}

		return utf8.AppendRune(buf, r), p
	}

	if c == 0 {
		return append(buf, replacementCharStr...), p + 1
	}

	return append(buf, c), p + 1
}

func (s *source) scanName(buf []byte, start uint) (string, []byte) {
	text := s.text
	p, n := s.pos, uint(len(text))

	for p < n && identContinueLUT[text[p]] {
		p++
	}

	// EOF
	if p >= n {
		s.seek(n)
		return text[start:n], buf
	}

	c := text[p]
	if c != 0 && c != '\\' {
		s.seek(p)
		return text[start:p], buf
	}

	if c == '\\' {
		// EOF
		if p+1 >= n {
			s.seek(n)
			return text[start:p] + replacementCharStr, buf
		}

		// line break
		if isNewline(text[p+1]) {
			s.seek(p)
			return text[start:p], buf
		}
	}

	buf = append(buf, text[start:p]...)
	return s.scanNameSlow(buf, p)
}

func (s *source) scanNameSlow(buf []byte, p uint) (string, []byte) {
	text := s.text
	n := uint(len(text))
	seg := p

	for p < n {
		if identContinueLUT[text[p]] {
			p++
			continue
		}

		buf = append(buf, text[seg:p]...)

		switch text[p] {
		case 0: // NUL
			buf = append(buf, replacementCharStr...)
			p++
			seg = p

		case '\\': // escape
			// EOF
			if p+1 >= n {
				buf = append(buf, replacementCharStr...)
				p, seg = n, n
				goto done
			}

			// line break
			if isNewline(text[p+1]) {
				seg = p
				goto done
			}

			buf, p = s.appendEscape(buf, p+1)
			seg = p

		default:
			seg = p
			goto done
		}
	}

done:
	buf = append(buf, text[seg:p]...)
	s.seek(p)
	return string(buf), buf
}

func (s *source) scanString(buf []byte, quote byte) (string, []byte, bool) {
	text := s.text
	p, n := s.pos, uint(len(text))
	start := p

	stop := stopLUT(quote)

	for p < n && !stop[text[p]] {
		p++
	}

	// EOF
	if p >= n {
		s.seek(n)
		return text[start:p], buf, true
	}
	// quote
	if text[p] == quote {
		s.seek(p + 1)
		return text[start:p], buf, true
	}
	// newline
	if isNewline(text[p]) {
		return text[start:p], buf, false
	}

	buf = append(buf, text[start:p]...)
	return s.scanStringSlow(buf, p, quote)
}

func (s *source) scanStringSlow(buf []byte, p uint, quote byte) (string, []byte, bool) {
	text := s.text
	n := uint(len(text))
	seg := p

	stop := stopLUT(quote)

	for p < n {
		if !stop[text[p]] {
			p++
			continue
		}
		buf = append(buf, text[seg:p]...)

		switch c := text[p]; c {
		case quote:
			if c == quote {
				s.seek(p + 1)
				return string(buf), buf, true
			}
			seg = p
			p++

		case 0: // NUL
			buf = append(buf, replacementCharStr...)
			p++
			seg = p

		case '\\':
			// EOF
			if p+1 >= n {
				p, seg = n, n
				goto done
			}
			// line continuation
			if isNewline(text[p+1]) {
				p += 2
				if text[p-1] == '\r' && p < n && text[p] == '\n' {
					p++
				}
				seg = p
				continue
			}
			buf, p = s.appendEscape(buf, p+1)
			seg = p

		default:
			s.seek(p)
			return string(buf), buf, false
		}
	}

done:
	buf = append(buf, text[seg:p]...)
	s.seek(p)
	return string(buf), buf, true
}

func (s *source) scanUrl(buf []byte) (string, []byte, bool) {
	text := s.text
	p, n := s.pos, uint(len(text))
	start, seg := p, p

	for {
		for p < n && !urlSpecialLUT[text[p]] {
			p++
		}

		// EOF or escape
		if p >= n || text[p] != '\\' {
			break
		}

		// line break
		if p+1 < n && isNewline(text[p+1]) {
			break
		}

		buf = append(buf, text[seg:p]...)

		// EOF
		if p+1 >= n {
			buf = append(buf, replacementCharStr...)
			p, seg = n, n
			break
		}

		buf, p = s.appendEscape(buf, p+1)
		seg = p
	}

	end := p

	// skip whitespace
	for p < n && wsLUT[text[p]] {
		p++
	}

	switch {
	case p >= n:
	case text[p] == ')':
		p++
	default:
		p, end = skipBadUrl(text, p)
		s.seek(p)
		return text[start:end], buf, false
	}

	s.seek(p)
	if seg == start {
		return text[start:end], buf, true
	}

	buf = append(buf, text[seg:end]...)
	return string(buf), buf, true
}

// skipBadUrl consume the remnants of a bad url
//
// https://www.w3.org/TR/css-syntax-3/#consume-remnants-of-bad-url
func skipBadUrl(text string, p uint) (uint, uint) {
	n := uint(len(text))
	for p < n {
		switch text[p] {
		case ')':
			return p + 1, p
		case '\\':
			if p+1 < n && !isNewline(text[p+1]) {
				p += 2
			} else {
				p++
			}
		default:
			p++
		}
	}
	return n, n
}

func (s *source) endBlock(p uint, opener bool) uint {
	var stackBuf [16]byte
	stack := stackBuf[:0]

	var nameBuf [4]byte
	name := nameBuf[:0]

	text := s.text
	n := uint(len(text))

	if opener && p > 0 && p <= n {
		switch s.text[p-1] {
		case '(':
			stack = append(stack, ')')
		case '[':
			stack = append(stack, ']')
		case '{':
			stack = append(stack, '}')
		}
	}

	for p < n {
		c := text[p]

		// skip whitespace
		if wsLUT[c] {
			p++
			for p < n && wsLUT[text[p]] {
				p++
			}
			if p >= n {
				return n
			}

			name = name[:0]
			c = text[p]
		}

		if identContinueLUT[c] {
			for p < n && identContinueLUT[text[p]] {
				if len(name) < 4 {
					name = append(name, text[p])
				}
				p++
			}
			continue
		}

		switch c {
		case ')', ']', '}':
			p++
			if len(stack) == 0 {
				return p
			}
			if stack[len(stack)-1] == c {
				stack = stack[:len(stack)-1]
				if len(stack) == 0 && opener {
					return p
				}
			}

		case '(':
			p++
			if len(name) == 3 {
				ascii.ToLowerBytes(name)
				if string(name) == "url" {
					// skip whitespace
					for p < n && wsLUT[text[p]] {
						p++
					}
					if p >= n {
						return n
					}

					// function
					if q := text[p]; q == '"' || q == '\'' {
						stack = append(stack, ')')
						name = name[:0]
						continue
					}

					// scan to the first closing parenthesis
					for p < n && text[p] != ')' {
						if text[p] == '\\' {
							p++
						}
						p++
					}
					if p >= n {
						return n
					}
					p++

					name = name[:0]
					continue
				}
			}
			stack = append(stack, ')')

		case '[':
			p++
			stack = append(stack, ']')

		case '{':
			p++
			stack = append(stack, '}')

		case '\\': // escape
			if p+1 >= n {
				return n
			}
			if len(name) < 4 {
				var tb [4]byte
				tmp, tp := s.appendEscape(tb[:0], p+1)
				if len(tmp) == 1 && identContinueLUT[tmp[0]] {
					name = append(name, tmp[0])
				} else {
					name = nameBuf[:4]
				}
				p = tp
			} else {
				p += 2
			}
			continue

		case '"', '\'': // string
			q := c
			p++
			stop := stopLUT(q)
		strLoop:
			for p < n {
				if !stop[text[p]] {
					p++
					continue
				}
				switch text[p] {
				case q:
					p++
					break strLoop
				case '\\':
					p += 2
				case '\n', '\r', '\f': // bad string ends before the newline
					break strLoop
				default:
					p++
				}
			}

		case '/': // possible comment
			rest := text[p:]
			p++
			if len(rest) > 1 && rest[1] == '*' {
				i := strings.Index(rest[2:], "*/")
				if i < 0 {
					return n
				}
				p += uint(i) + 3
			}

		case 0: // NUL
			p++
			name = nameBuf[:4]
			continue

		default:
			p++
		}

		name = name[:0]
	}

	return n
}

func (s *source) isIdentStart() bool {
	switch s.ch {
	case 0: // NUL
		// EOF
		if s.pos >= s.Len() {
			return false
		}
		return true

	case '-':
		if !s.hasAtLeast(1) {
			return false
		}
		ch := s.byteAt(1)
		if isIdentStartLUT(ch) || ch == '-' {
			return true
		}
		return ch == '\\' && (!s.hasAtLeast(2) || !s.hasNewlineAt(2))

	case '\\': // escape
		if !s.hasAtLeast(1) {
			return false
		}
		return !s.hasNewlineAt(1)

	default:
		return isIdentStartLUT(s.ch)
	}
}

func (s *source) isDigit() bool {
	return '0' <= s.ch && s.ch <= '9'
}

func (s *source) isNewline() bool {
	return s.ch == '\n' || s.ch == '\r' || s.ch == '\f'
}

func (s *source) hasDigitsAt(offset uint) bool {
	ch := s.byteAt(offset)
	return '0' <= ch && ch <= '9'
}

func (s *source) hasNewlineAt(offset uint) bool {
	ch := s.byteAt(offset)
	return ch == '\n' || ch == '\r' || ch == '\f'
}

// 0-9
func isDigit(ch byte) bool { return '0' <= ch && ch <= '9' }

// \n \r \f
func isNewline(ch byte) bool { return ch == '\n' || ch == '\r' || ch == '\f' }
