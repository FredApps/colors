package colors_test

import (
	"testing"

	"github.com/thmalt/colors"
)

func TestParseRejectsURLs(t *testing.T) {
	for _, input := range []string{
		"url(x)", "URL(x)", "url(;)", "url(x;y)", "url()", "url(",
		"url(x", "url(  x  )", `url("x")`, `url('x')`,
		`url(\)x)`, "url(\\", "url(\x00)", "url(é)",
		"red url(x)", "#f00 url(x)", "rgb(255 0 0) url(x)",
		"rgb(from url(x) r g b)",
	} {
		t.Run(input, func(t *testing.T) {
			assertParseError(t, input)
		})
	}
}

func TestParseRejectsTrailingTokens(t *testing.T) {
	for _, input := range []string{
		"red blue", "#f00 #00f", "rgb(255 0 0))", "rgb(255 0 0);",
		"red,", "red 1", "red /**/ blue", "red /**/ ;",
		"rgb(from rgb(255 0 0) r g b) blue",
	} {
		t.Run(input, func(t *testing.T) {
			assertParseError(t, input)
		})
	}
}

func TestParseRejectsUnknownRelativeChannels(t *testing.T) {
	for _, input := range []string{
		"rgb(from red typo g b)",
		"rgb(from red r typo b)",
		"rgb(from red r g typo)",
		"rgb(from red r g b / typo)",
		"hsl(from red typo s l)",
		"hwb(from red h typo b)",
		"lab(from red l typo b)",
		"lch(from red l c typo)",
		"color(from red srgb typo g b)",
		"color(from red xyz x y typo)",
		"alpha(from red / typo)",
		// Known identifiers from a different color space must also be rejected.
		"rgb(from red h g b)",
		"hsl(from red r s l)",
	} {
		t.Run(input, func(t *testing.T) {
			assertParseError(t, input)
		})
	}
}

func TestParseWhitespaceAndComments(t *testing.T) {
	for _, input := range []string{
		"/* comment */ red",
		"/* first */ \t /* second */\nred",
		"rgb(255/**/ 0 0)",
		"rgb(255 /**/\t/**/\n0 0)",
		"rgb(255 0 0 / /**/ 100%)",
		"red /* trailing */ \t\r\n\f",
		"rgb(255 0 0) /**/ /**/ ",
		"rgb(from /**/ red /**/ r /**/ g /**/ b)",
		"rgb(from rgb(255/**/ 0 0) /**/ r g b)",
	} {
		t.Run(input, func(t *testing.T) {
			assertParseHex(t, input, "#ff0000")
		})
	}
}

func TestParseValidControls(t *testing.T) {
	for _, tc := range []struct{ input, hex string }{
		{" red ", "#ff0000"},
		{"rgb(255 0 0)", "#ff0000"},
		{"rgb(from red r g b)", "#ff0000"},
		{"rgb(from rgb(from red r g b) r g b)", "#ff0000"},
		{"color(from red srgb r g b)", "#ff0000"},
		{"rgb(from rgba(255, 0, 0, 0.5) r g b / alpha)", "#ff000080"},
		{"alpha(from red / alpha)", "#ff0000"},
		{"rgb(from red none g b)", "#000000"},
		{"rgb(255/**/0 0)", "#ff0000"},
		{"transparent", "#00000000"},
		{"rgba(255, 0, 0, 0.5)", "#ff000080"},
		{"#f008", "#ff000088"},
	} {
		t.Run(tc.input, func(t *testing.T) {
			assertParseHex(t, tc.input, tc.hex)
		})
	}
}

func assertParseError(t *testing.T, input string) {
	t.Helper()
	if c, err := colors.Parse(input); err == nil {
		t.Errorf("Parse(%q) = %v, nil; want an error", input, c)
	}
	if c, err := colors.ParseBytes([]byte(input)); err == nil {
		t.Errorf("ParseBytes(%q) = %v, nil; want an error", input, c)
	}
	original := colors.Hex("#00ff00")
	c := original
	if err := c.UnmarshalText([]byte(input)); err == nil {
		t.Errorf("UnmarshalText(%q) succeeded; want an error", input)
	}
	if c != original {
		t.Errorf("UnmarshalText(%q) changed the receiver after an error", input)
	}
}

func assertParseHex(t *testing.T, input, hex string) {
	t.Helper()
	parsers := []struct {
		name  string
		parse func() (colors.Color, error)
	}{
		{"Parse", func() (colors.Color, error) { return colors.Parse(input) }},
		{"ParseBytes", func() (colors.Color, error) { return colors.ParseBytes([]byte(input)) }},
		{"UnmarshalText", func() (colors.Color, error) {
			var c colors.Color
			err := c.UnmarshalText([]byte(input))
			return c, err
		}},
	}
	for _, parser := range parsers {
		c, err := parser.parse()
		if err != nil {
			t.Errorf("%s(%q): %v", parser.name, input, err)
		} else if got := c.Hex(); got != hex {
			t.Errorf("%s(%q).Hex() = %s; want %s", parser.name, input, got, hex)
		}
	}
}
