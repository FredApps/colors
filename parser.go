package colors

import (
	"errors"
	"strings"
	"unsafe"

	"github.com/thmalt/colors/internal/ascii"
	"github.com/thmalt/colors/internal/css/token"
	"github.com/thmalt/colors/internal/css/tokenizer"
	"github.com/thmalt/colors/space"
)

const (
	maxOriginDepth = 16 // Max nesting of relative color origins
	mathMaxDepth   = 32 // Max nesting of math functions and parentheses
)

var (
	errInvalidColor    = errors.New("invalid color")
	errInvalidHexColor = errors.New("invalid hex color")

	errInvalidColorComponent = errors.New("invalid color component")

	errNotImplemented = errors.New("not implemented")

	errInvalidAlpha      = errors.New("invalid alpha channel")
	errInvalidRed        = errors.New("invalid red channel")
	errInvalidGreen      = errors.New("invalid green channel")
	errInvalidBlue       = errors.New("invalid blue channel")
	errInvalidHue        = errors.New("invalid hue channel")
	errInvalidChroma     = errors.New("invalid chroma channel")
	errInvalidLightness  = errors.New("invalid lightness channel")
	errInvalidCoordinate = errors.New("invalid coordinate channel")

	errInvalidChannel1 = errors.New("invalid channel 1")
	errInvalidChannel2 = errors.New("invalid channel 2")
	errInvalidChannel3 = errors.New("invalid channel 3")

	errInvalidSyntax = errors.New("invalid syntax")

	errParserExpectedCloseParenthesis = errors.New("expected closing parenthesis")
)

type parseError struct {
	err   error
	text  string
	msg   string
	got   numKind
	allow numKindSet
	fn    mathFunc
	none  bool
}

func (e *parseError) Error() string {
	var b strings.Builder
	if e.err != nil {
		b.WriteString(e.err.Error())
	}

	if e.fn > fnUnknown {
		if b.Len() > 0 {
			b.WriteString(": ")
		}
		b.WriteString(e.fn.String())
	}

	if e.allow != 0 {
		if b.Len() > 0 {
			b.WriteString(": ")
		}
		b.WriteString("expected ")
		e.allow.writeNames(&b, e.none)
		b.WriteString(", got ")
		b.WriteString(e.got.String())
	} else if e.msg != "" {
		if b.Len() > 0 {
			b.WriteString(": ")
		}
		b.WriteString(e.msg)
	}

	if e.text != "" {
		if b.Len() > 0 {
			b.WriteString(": ")
		}
		b.WriteString(e.text)
	}

	return b.String()
}

func (e *parseError) Unwrap() error {
	return e.err
}

// Parse parses a CSS color value from s.
func Parse(s string) (Color, error) {
	var p parser
	p.init(s)
	return p.parse()
}

// ParseBytes parses a CSS color value from b.
func ParseBytes(b []byte) (Color, error) {
	var p parser
	p.initBytes(b)
	return p.parse()
}

// parse requires a single complete color value. Nested colors call
// parseColor directly so their enclosing function consumes what follows.
func (p *parser) parse() (Color, error) {
	p.next()
	c, err := p.parseColor()
	if err != nil {
		return Color{}, err
	}
	if !p.t.Is(token.EOF) {
		start := p.t.TokenStart()
		p.setErr(errInvalidSyntax)
		p.setErrMsg("unexpected trailing input")
		p.setErrTextAt(start, p.t.Len())
		return Color{}, p.fail()
	}
	return c, nil
}

type parser struct {
	t tokenizer.Tokenizer

	mathSpace   bool
	originDepth int

	prevBeg, prevEnd int
	curBeg, curEnd   int

	err parseError
}

func (p *parser) init(text string) {
	p.t.Init(text)
}

func (p *parser) initBytes(b []byte) {
	p.t.Init(unsafe.String(unsafe.SliceData(b), len(b)))
}

func (p *parser) fail() error { return &p.err }

func (p *parser) setErr(err error) {
	if p.err.err == nil {
		p.err.err = err
	}
}

func (p *parser) setErrMsg(msg string) {
	if p.err.msg == "" {
		p.err.msg = msg
	}
}

func (p *parser) setErrExpected(got numKind, allow numKindSet, none bool) {
	if p.err.msg == "" && p.err.allow == 0 {
		p.err.msg = "expected"
		p.err.got = got
		p.err.allow = allow
		p.err.none = none
	}
}

func (p *parser) setErrText(text string) {
	if p.err.text == "" {
		p.err.text = text
	}
}

func (p *parser) setErrTextAt(start, end uint) {
	p.setErrText(p.t.SegmentAt(start, end))
}

func (p *parser) next() { p.t.Next() }

func (p *parser) mathNext() {
	p.t.NextIncludingWhitespace()
	p.mathSpace = p.t.Is(token.Whitespace)
	if p.mathSpace {
		p.t.Next()
	}
}

func (p *parser) ParseColor() (Color, error) {
	c, err := p.parseColor()
	if err == nil {
		return c, nil
	}

	if err == &p.err {
		e := p.err
		p.err = parseError{}
		return c, &e
	}

	return c, err
}

func (p *parser) parseColor() (Color, error) {
	switch p.t.Kind() {
	case token.Hash:
		return p.parseHexColor()

	case token.Ident:
		return p.parseNamedColor()

	case token.Function:
		return p.parseColorFunc()

	default:
		return Color{}, errInvalidColor
	}
}

func (p *parser) parseHexColor() (Color, error) {
	c, ok := TryHex(p.t.Text())
	p.next()
	if !ok {
		return Color{}, errInvalidHexColor
	}
	return c, nil
}

func (p *parser) parseNamedColor() (Color, error) {
	if c, ok := lookupNamedColor(p.t.Text()); ok {
		p.next()
		return c, nil
	}
	p.setErr(errInvalidColor)
	p.setErrMsg("unknown named color")
	p.setErrTextAt(p.t.TokenStart(), p.t.TokenEnd())
	return Color{}, p.fail()
}

func (p *parser) parseColorFunc() (Color, error) {
	name := p.t.Text()

	var buf [16]byte
	n := copy(buf[:], name)
	ascii.ToLowerBytes(buf[:n])

	p.next()

	originColor, err := p.parseOriginColor()
	if err != nil {
		return Color{}, err
	}

	switch string(buf[:n]) {
	case "rgb", "rgba":
		return p.parseRgb(originColor)
	case "hsl", "hsla":
		return p.parseHxxLike(originColor, space.Hsl, space.ChIdentSetHSL, true)
	case "hsv":
		return p.parseHxxLike(originColor, space.Hsv, space.ChIdentSetHSV, false)
	case "hwb":
		return p.parseHxxLike(originColor, space.Hwb, space.ChIdentSetHWB, false)

	case "luv":
		return p.parseLabLike(originColor, space.Luv, 100, 215)
	case "lab":
		return p.parseLabLike(originColor, space.Lab, 100, 125)
	case "oklab":
		return p.parseLabLike(originColor, space.Oklab, 1, .4)

	case "lchuv":
		return p.parseLchLike(originColor, space.Lchuv, 100, 220)
	case "lch":
		return p.parseLchLike(originColor, space.Lch, 100, 150)
	case "oklch":
		return p.parseLchLike(originColor, space.Oklch, 1, .4)

	case "color":
		switch p.t.Kind() {
		case token.Ident:
			return p.parseColorSpace(originColor)
		case token.CloseParenthesis, token.EOF:
			p.setErr(errInvalidColor)
			p.setErrMsg("missing color space")
			return Color{}, p.fail()
		default:
			p.setErr(errInvalidColor)
			p.setErrMsg("expected color space")
			p.setErrTextAt(p.t.TokenStart(), p.t.TokenEnd())
			return Color{}, p.fail()
		}

	case "alpha":
		return p.parseRelativeAlpha(originColor)

	case "color-mix":
		return p.parseColorMix()
	}

	p.setErr(errInvalidColor)
	p.setErrMsg("unexpected function")
	p.setErrText(name)
	return Color{}, p.fail()
}

func (p *parser) parseOriginColor() (Color, error) {
	if !p.t.IsIdent("from") {
		return Color{}, nil
	}

	if p.originDepth >= maxOriginDepth {
		p.setErr(errInvalidColor)
		p.setErrMsg("origin color nesting too deep")
		return Color{}, p.fail()
	}

	p.originDepth++
	p.next()
	c, err := p.parseColor()
	p.originDepth--
	return c, err
}

// rgb( <percentage>#{3} , <alpha-value>? ) | rgb( <number>#{3} , <alpha-value>? )

func (p *parser) parseRgb(originColor Color) (Color, error) {
	ctx := &colorCompCtx{}
	isLegacy := true

	if originColor.IsValid() {
		const sp = space.Srgb

		originColor.mutTo(sp)
		ctx.setSpace(sp)
		ctx.setChannel3(
			originColor.c1*255,
			originColor.c2*255,
			originColor.c3*255,
			originColor.alpha,
		)

		isLegacy = false
	}

	red, ok := p.parseNumberOrPercentage(ctx, true, 255)
	if !ok {
		p.setErr(errInvalidRed)
		return Color{}, p.fail()
	}

	var green, blue, alpha numValue

	isLegacy = isLegacy && !red.isNone && p.t.NextIf(token.Comma)
	if isLegacy {
		var beforeAlpha uint
		if red.kind == kindPercent {
			start := p.t.TokenStart()
			green, ok = p.parsePercentage(ctx, false, 255)
			if !ok {
				p.setErr(errInvalidGreen)
				return Color{}, p.fail()
			}

			if !p.t.NextIf(token.Comma) {
				p.setErr(errInvalidSyntax)
				p.setErrMsg("expected ',' separator")
				p.setErrTextAt(start, p.t.TokenEnd())
				return Color{}, p.fail()
			}

			beforeAlpha = p.t.TokenStart()
			blue, ok = p.parsePercentage(ctx, false, 255)
			if !ok {
				p.setErr(errInvalidBlue)
				return Color{}, p.fail()
			}
		} else {
			start := p.t.TokenStart()
			green, ok = p.parseNumber(ctx, false)
			if !ok {
				p.setErr(errInvalidGreen)
				return Color{}, p.fail()
			}

			if !p.t.NextIf(token.Comma) {
				p.setErr(errInvalidSyntax)
				p.setErrMsg("expected ',' separator")
				p.setErrTextAt(start, p.t.TokenEnd())
				return Color{}, p.fail()
			}

			beforeAlpha = p.t.TokenStart()
			blue, ok = p.parseNumber(ctx, false)
			if !ok {
				p.setErr(errInvalidBlue)
				return Color{}, p.fail()
			}
		}

		alpha, ok = p.parseLegacyAlpha(ctx, beforeAlpha)
		if !ok {
			p.setErr(errInvalidAlpha)
			return Color{}, p.fail()
		}
	} else {
		green, ok = p.parseNumberOrPercentage(ctx, true, 255)
		if !ok {
			p.setErr(errInvalidGreen)
			return Color{}, p.fail()
		}

		blue, ok = p.parseNumberOrPercentage(ctx, true, 255)
		if !ok {
			p.setErr(errInvalidBlue)
			return Color{}, p.fail()
		}

		alpha, ok = p.parseModernAlpha(ctx)
		if !ok {
			p.setErr(errInvalidAlpha)
			return Color{}, p.fail()
		}
	}

	if !p.closeParenthesis() {
		return Color{}, errParserExpectedCloseParenthesis
	}

	return SrgbAlpha(
		clamp01(red.value*inv255),
		clamp01(green.value*inv255),
		clamp01(blue.value*inv255),
		alpha.value,
	), nil
}

// legacy ( <hue>, <percentage>, <percentage>, <alpha-value>? )
// modern ( [<hue> | none] [<percentage> | <number> | none] [<percentage> | <number> | none] [ / [<alpha-value> | none] ]? )

func (p *parser) parseHxxLike(originColor Color, sp space.Space, chSet space.ChIdentSet, allowLegacy bool) (Color, error) {
	ctx := &colorCompCtx{}
	isLegacy := allowLegacy

	if originColor.IsValid() {
		originColor.mutTo(sp)
		ctx.space = sp
		ctx.chSet = chSet
		ctx.setChannel3(
			originColor.c1,
			originColor.c2*100,
			originColor.c3*100,
			originColor.alpha,
		)

		isLegacy = false
	}

	hue, ok := p.parseHue(ctx, true)
	if !ok {
		p.setErr(errInvalidHue)
		return Color{}, p.fail()
	}

	var c2, c3, alpha numValue

	isLegacy = isLegacy && !hue.isNone && p.t.NextIf(token.Comma)
	if isLegacy {
		start := p.t.TokenStart()
		c2, ok = p.parsePercentage(ctx, false, 100)
		if !ok {
			p.setErr(errInvalidColorComponent)
			return Color{}, p.fail()
		}

		if !p.t.NextIf(token.Comma) {
			p.setErr(errInvalidSyntax)
			p.setErrMsg("expected ',' separator")
			p.setErrTextAt(start, p.t.TokenEnd())
			return Color{}, p.fail()
		}

		beforeAlpha := p.t.TokenStart()
		c3, ok = p.parsePercentage(ctx, false, 100)
		if !ok {
			p.setErr(errInvalidColorComponent)
			return Color{}, p.fail()
		}

		alpha, ok = p.parseLegacyAlpha(ctx, beforeAlpha)
		if !ok {
			p.setErr(errInvalidAlpha)
			return Color{}, p.fail()
		}
	} else {
		c2, ok = p.parseNumberOrPercentage(ctx, true, 100)
		if !ok {
			p.setErr(errInvalidColorComponent)
			return Color{}, p.fail()
		}

		c3, ok = p.parseNumberOrPercentage(ctx, true, 100)
		if !ok {
			p.setErr(errInvalidColorComponent)
			return Color{}, p.fail()
		}

		alpha, ok = p.parseModernAlpha(ctx)
		if !ok {
			p.setErr(errInvalidAlpha)
			return Color{}, p.fail()
		}
	}

	if !p.closeParenthesis() {
		return Color{}, errParserExpectedCloseParenthesis
	}

	const inv = 1 / 100.0
	return new3(
		sp,
		hue.value,
		c2.value*inv,
		c3.value*inv,
		alpha.value,
	), nil
}

/*
lab([from <color>]?
        [<percentage> | <number> | none]
        [<percentage> | <number> | none]
        [<percentage> | <number> | none]
        [ / [<alpha-value> | none] ]? )
*/

func (p *parser) parseLabLike(originColor Color, sp space.Space, lightnessRange, abRange float64) (Color, error) {
	ctx := &colorCompCtx{}

	if originColor.IsValid() {
		originColor.mutTo(sp)
		ctx.space = sp
		ctx.chSet = space.ChIdentSetLAB
		ctx.setChannel3(
			originColor.c1,
			originColor.c2,
			originColor.c3,
			originColor.alpha,
		)
	}

	lightness, ok := p.parseNumberOrPercentage(ctx, true, lightnessRange)
	if !ok {
		p.setErr(errInvalidLightness)
		return Color{}, p.fail()
	}

	a, ok := p.parseNumberOrPercentage(ctx, true, abRange)
	if !ok {
		p.setErr(errInvalidCoordinate)
		return Color{}, p.fail()
	}

	b, ok := p.parseNumberOrPercentage(ctx, true, abRange)
	if !ok {
		p.setErr(errInvalidCoordinate)
		return Color{}, p.fail()
	}

	alpha, ok := p.parseModernAlpha(ctx)
	if !ok {
		p.setErr(errInvalidAlpha)
		return Color{}, p.fail()
	}

	if !p.closeParenthesis() {
		return Color{}, errParserExpectedCloseParenthesis
	}

	return new3(
		sp,
		lightness.value,
		a.value,
		b.value,
		clamp01(alpha.value),
	), nil
}

/*
lch([from <color>]?
        [<percentage> | <number> | none]
        [<percentage> | <number> | none]
        [<hue> | none]
        [ / [<alpha-value> | none] ]? )
*/

func (p *parser) parseLchLike(originColor Color, sp space.Space, lightnessRange, chromaRange float64) (Color, error) {
	ctx := &colorCompCtx{}

	if originColor.IsValid() {
		originColor.mutTo(sp)
		ctx.space = sp
		ctx.chSet = space.ChIdentSetLCH
		ctx.setChannel3(
			originColor.c1,
			originColor.c2,
			originColor.c3,
			originColor.alpha,
		)
	}

	lightness, ok := p.parseNumberOrPercentage(ctx, true, lightnessRange)
	if !ok {
		p.setErr(errInvalidColorComponent)
		return Color{}, p.fail()
	}

	chroma, ok := p.parseNumberOrPercentage(ctx, true, chromaRange)
	if !ok {
		p.setErr(errInvalidChroma)
		return Color{}, p.fail()
	}

	hue, ok := p.parseHue(ctx, true)
	if !ok {
		p.setErr(errInvalidHue)
		return Color{}, p.fail()
	}

	alpha, ok := p.parseModernAlpha(ctx)
	if !ok {
		p.setErr(errInvalidAlpha)
		return Color{}, p.fail()
	}

	if !p.closeParenthesis() {
		return Color{}, errParserExpectedCloseParenthesis
	}

	return new3(
		sp,
		lightness.value,
		chroma.value,
		hue.value,
		alpha.value,
	), nil
}

// alpha([from <color>] [ / [<alpha-value> | none] ]? )

func (p *parser) parseRelativeAlpha(originColor Color) (Color, error) {
	if !originColor.IsValid() {
		return Color{}, errInvalidAlpha
	}

	idx := uint8(originColor.space.ChIndex(space.ChIdentAlpha))
	if idx > space.MaxChannelCount {
		return Color{}, errInvalidAlpha
	}

	ctx := &colorCompCtx{
		space: originColor.space,
		chSet: space.ChIdentSetAlpha,
	}
	ctx.channels[idx] = originColor.alpha

	alpha, ok := p.parseModernAlpha(ctx)
	if !ok {
		p.setErr(errInvalidAlpha)
		return Color{}, p.fail()
	}

	if !p.closeParenthesis() {
		return Color{}, errParserExpectedCloseParenthesis
	}

	originColor.alpha = alpha.value
	return originColor, nil
}

func (p *parser) parseColorMix() (Color, error) {
	return Color{}, errNotImplemented
}

func (p *parser) parseColorSpace01(originColor Color, sp space.Space, chSet space.ChIdentSet) (Color, error) {
	ctx := &colorCompCtx{}

	if originColor.IsValid() {
		originColor.mutTo(sp)
		ctx.space = sp
		ctx.chSet = chSet
		ctx.setChannel3(
			originColor.c1,
			originColor.c2,
			originColor.c3,
			originColor.alpha,
		)
	}

	c1, ok := p.parseNumberOrPercentage(ctx, true, 1)
	if !ok {
		p.setErr(errInvalidChannel1)
		return Color{}, p.fail()
	}

	c2, ok := p.parseNumberOrPercentage(ctx, true, 1)
	if !ok {
		p.setErr(errInvalidChannel2)
		return Color{}, p.fail()
	}

	c3, ok := p.parseNumberOrPercentage(ctx, true, 1)
	if !ok {
		p.setErr(errInvalidChannel3)
		return Color{}, p.fail()
	}

	alpha, ok := p.parseModernAlpha(ctx)
	if !ok {
		p.setErr(errInvalidAlpha)
		return Color{}, p.fail()
	}

	if !p.closeParenthesis() {
		return Color{}, errParserExpectedCloseParenthesis
	}

	return new3(
		sp,
		c1.value,
		c2.value,
		c3.value,
		alpha.value,
	), nil
}

func (p *parser) parseColorComponent(ctx *colorCompCtx, allowNone bool, allowKind numKindSet) (numValue, bool) {
	switch p.t.Kind() {
	case token.Ident:
		text := p.t.Text()
		if ascii.EqualFold(text, "none") {
			if allowNone {
				p.next()
				return numValue{isNone: true, kind: kindNumber}, true
			}
			p.setErrMsg("keyword is not allowed")
		} else if ctx.chSet == 0 {
			p.setErrMsg("channel keyword requires an origin color")
		} else if chid := space.LookupChIdentFlags(text); ctx.chSet.Has(chid) {
			p.next()
			return numValue{value: ctx.channel(chid), kind: kindNumber}, true
		} else if chid != 0 {
			p.setErrMsg("channel keyword is not valid in this color space")
		} else {
			p.setErrMsg("unknown channel keyword")
		}

		p.setErrText(text)
		return numValue{}, false

	case token.Function:
		state := p.t.State()
		value, ok := p.parseMathFunc(ctx, 0)
		if !ok {
			p.t.Reset(state)
			p.setErrMsg("invalid math expression")
			p.setErrTextAt(state.Position(), p.t.SkipBlockNext())
			return numValue{}, false
		}

		if !allowKind.has(value.kind) {
			p.setErrExpected(value.kind, allowKind, allowNone)
			p.setErrTextAt(state.Position(), p.t.TokenStart())
			return numValue{}, false
		}

		return value, ok

	case token.Number:
		if allowKind.has(kindNumber) {
			value := p.t.Number()
			p.next()
			return numValue{value: value, kind: kindNumber}, true
		}
		p.setErrExpected(kindNumber, allowKind, allowNone)
		p.setErrTextAt(p.t.TokenStart(), p.t.Position())
		return numValue{}, false

	case token.Percentage:
		if allowKind.has(kindPercent) {
			value := p.t.Number()
			p.next()
			return numValue{value: value, kind: kindPercent}, true
		}
		p.setErrExpected(kindPercent, allowKind, allowNone)
		p.setErrTextAt(p.t.TokenStart(), p.t.Position())
		return numValue{}, false

	case token.Dimension:
		unitID := lookupUnit(p.t.Text())

		if unitID == unitNone {
			p.setErrMsg("unknown unit")
		} else {
			info := unitID.info()
			if allowKind.has(info.kind) {
				v := p.t.Number() * info.scale
				p.next()
				return numValue{value: v, kind: info.kind}, true
			}
			p.setErrExpected(info.kind, allowKind, allowNone)
		}
		p.setErrTextAt(p.t.TokenStart(), p.t.Position())
		return numValue{}, false
	}

	p.setErrMsg("unexpected token")
	p.setErrTextAt(p.t.TokenStart(), p.t.Position())
	return numValue{}, false
}

func (p *parser) parseHue(ctx *colorCompCtx, allowNone bool) (numValue, bool) {
	v, ok := p.parseColorComponent(ctx, allowNone, allowHue)
	if !ok || v.isNone {
		return v, ok
	}
	v.value = v.sanitizeHue()
	v.kind = kindNumber
	return v, true
}

func (p *parser) parseNumberOrPercentage(ctx *colorCompCtx, allowNone bool, percentScale float64) (numValue, bool) {
	v, ok := p.parseColorComponent(ctx, allowNone, allowNumPct)
	if !ok || v.isNone {
		return v, ok
	}
	if v.kind == kindPercent {
		v.value *= percentScale
		v.kind = kindNumber
	}
	v.value = v.sanitizeNumber()
	return v, true
}

func (p *parser) parseNumber(ctx *colorCompCtx, allowNone bool) (numValue, bool) {
	return p.parseColorComponent(ctx, allowNone, allowNumber)
}

func (p *parser) parsePercentage(ctx *colorCompCtx, allowNone bool, percentScale float64) (numValue, bool) {
	v, ok := p.parseColorComponent(ctx, allowNone, allowPercent)
	if !ok || v.isNone {
		return v, ok
	}
	v.value *= percentScale
	v.kind = kindNumber
	v.value = v.sanitizeNumber()
	return v, true
}

func (p *parser) parseLegacyAlpha(ctx *colorCompCtx, start uint) (numValue, bool) {
	if p.t.Is(token.CloseParenthesis) || p.t.Is(token.EOF) {
		return numValue{value: 1, kind: kindNumber}, true
	}

	if !p.t.NextIf(token.Comma) {
		p.setErr(errInvalidSyntax)
		p.setErrMsg("expected ',' or ')'")
		p.setErrTextAt(start, p.t.TokenEnd())
		return numValue{}, false
	}

	v, ok := p.parseNumberOrPercentage(ctx, false, 1)
	if ok {
		v.value = clamp01(v.value)
	}
	return v, ok
}

func (p *parser) parseModernAlpha(ctx *colorCompCtx) (numValue, bool) {
	if !p.t.NextIfDelim('/') {
		alpha := 1.0
		if ctx.chSet.Has(space.ChIdentAlpha) {
			alpha = ctx.channel(space.ChIdentAlpha)
		}
		return numValue{value: alpha, kind: kindNumber}, true
	}

	v, ok := p.parseNumberOrPercentage(ctx, true, 1)
	if ok && !v.isNone {
		v.value = clamp01(v.value)
	}
	return v, ok
}

func (p *parser) closeParenthesis() bool {
	return p.t.Is(token.EOF) || p.t.NextIf(token.CloseParenthesis)
}
