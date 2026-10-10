package token

import (
	"strings"

	"github.com/thmalt/colors/internal/hexlut"
)

type Flags uint16

const (
	// Number, Percentage, Dimension: no fraction or exponent
	IsInteger Flags = 1 << iota

	// Number, Percentage, Dimension: leading '+' or '-'
	HasSign

	// Ident, Function, AtKeyword, Hash, String, URL: Text differs from raw source
	IsDecoded

	// Ident, Function: name starts with "--"
	IsCustom

	// Hash: type "id" instead of "unrestricted"
	IsHashID

	textOwned
)

func (f Flags) IsInteger() bool { return f&IsInteger != 0 }
func (f Flags) HasSign() bool   { return f&HasSign != 0 }
func (f Flags) IsDecoded() bool { return f&IsDecoded != 0 }
func (f Flags) IsCustom() bool  { return f&IsCustom != 0 }
func (f Flags) IsHashID() bool  { return f&IsHashID != 0 }

// OwnsText reports whether Text is a standalone string that does not
// alias the input (unescaped into a fresh string, or already cloned).
func (f Flags) OwnsText() bool { return f&(IsDecoded|textOwned) != 0 }

var flagNames = [...]struct {
	mask  Flags
	long  string // "IsInteger", ...
	short string // "int", ...
}{
	{IsInteger, "IsInteger", "int"},
	{HasSign, "HasSign", "sign"},
	{IsDecoded, "IsDecoded", "decoded"},
	{IsCustom, "IsCustom", "custom"},
	{IsHashID, "IsHashID", "id"},
	{textOwned, "textOwned", "owned"},
}

func (f Flags) string(leadingSpace, long bool) string {
	if f == 0 {
		return ""
	}

	var b strings.Builder
	if leadingSpace {
		b.WriteByte(' ')
	}

	sep := " "
	if long {
		sep = " | "
	} else {
		b.WriteByte('[')
	}

	first := true
	for _, e := range flagNames {
		if f&e.mask == 0 {
			continue
		}
		if !first {
			b.WriteString(sep)
		}
		first = false
		if long {
			b.WriteString(e.long)
		} else {
			b.WriteString(e.short)
		}
		f &^= e.mask
	}

	if f != 0 {
		if !first {
			b.WriteString(sep)
		}
		b.WriteString("0x")
		v := byte(f>>8) * 2
		if v != 0 {
			b.WriteByte(hexlut.Enc[v])
			b.WriteByte(hexlut.Enc[v+1])
		}
		v = byte(f) * 2
		b.WriteByte(hexlut.Enc[v])
		b.WriteByte(hexlut.Enc[v+1])
	}

	if !long {
		b.WriteByte(']')
	}
	return b.String()
}

func (f Flags) String() string {
	return f.string(false, true)
}
