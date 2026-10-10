package colors

import (
	"math"
	"math/rand"

	"github.com/thmalt/colors/internal/ascii"
	"github.com/thmalt/colors/internal/css/token"
	"github.com/thmalt/colors/space"
)

func (p *parser) lookupMathFunc(name string) (mathFunc, bool) {
	fn := lookupMathFunc(name)
	if fn == fnUnknown {
		if ascii.EqualFold(name, "var") {
			p.setErrMsg("unexpected function")
		} else {
			p.setErrMsg("unknown function")
		}
		p.setErrText(name)
		return fnUnknown, false
	}
	return fn, true
}

// parseMathFunc parses a math function call. depth starts at 0.
func (p *parser) parseMathFunc(ctx *colorCompCtx, depth int) (numValue, bool) {
	fn, ok := p.lookupMathFunc(p.t.Text())
	if !ok {
		return numValue{}, false
	}
	p.mathNext()

	prev := p.err.fn
	if p.err.msg == "" {
		p.err.fn = fn
	}
	v, ok := p.mathCall(fn, ctx, depth)
	if ok {
		p.err.fn = prev
	}
	return v, ok
}

func (p *parser) mathCall(fn mathFunc, ctx *colorCompCtx, depth int) (numValue, bool) {
	if depth > mathMaxDepth {
		return numValue{}, false
	}

	switch fn {
	case fnCalc: // calc
		return p.mathCalc(ctx, depth)
	case fnMin: // min
		return p.mathMin(ctx, depth)
	case fnMax: // max
		return p.mathMax(ctx, depth)
	case fnClamp: // clamp
		return p.mathClamp(ctx, depth)
	case fnRound: // round
		return p.mathRound(ctx, depth)
	case fnMod: // mod
		return p.mathMod(ctx, depth)
	case fnRem: // rem
		return p.mathRem(ctx, depth)
	case fnSin: // sin
		return p.mathSin(ctx, depth)
	case fnCos: // cos
		return p.mathCos(ctx, depth)
	case fnTan: // tan
		return p.mathTan(ctx, depth)
	case fnAsin: // asin
		return p.mathAsin(ctx, depth)
	case fnAcos: // acos
		return p.mathAcos(ctx, depth)
	case fnAtan: // atan
		return p.mathAtan(ctx, depth)
	case fnAtan2: // atan2
		return p.mathAtan2(ctx, depth)
	case fnPow: // pow
		return p.mathPow(ctx, depth)
	case fnSqrt: // sqrt
		return p.mathSqrt(ctx, depth)
	case fnHypot: // hypot
		return p.mathHypot(ctx, depth)
	case fnLog: // log
		return p.mathLog(ctx, depth)
	case fnExp: // exp
		return p.mathExp(ctx, depth)
	case fnAbs: // abs
		return p.mathAbs(ctx, depth)
	case fnSign: // sign
		return p.mathSign(ctx, depth)
	case fnRandom: // random
		return p.mathRandom()
	default:
		p.setErrMsg("unsupported function")
		p.setErrText(p.t.Text())
		return numValue{}, false
	}

}

func (p *parser) mathArg(ctx *colorCompCtx, depth int) (numValue, bool) {
	return p.mathExpr(ctx, 0, depth)
}

func (p *parser) mathArgs(ctx *colorCompCtx, depth int, allow numKindSet, same bool, args []numValue) bool {
	for i := range args {
		if i > 0 && !p.mathMore() {
			p.setErrMsg("too few arguments")
			return false
		}

		v, ok := p.mathExpr(ctx, 0, depth)
		if !ok {
			return false
		}

		if !allow.has(v.kind) {
			p.setErrExpected(v.kind, allow, false)
			return false
		}

		if same && i > 0 && v.kind != args[0].kind {
			p.setErrMsg("incompatible types")
			return false
		}

		args[i] = v
	}
	return p.mathEnd()
}

func (p *parser) mathMore() bool {
	if !p.t.Is(token.Comma) {
		return false
	}
	p.mathNext()
	return true
}

func (p *parser) mathEnd() bool {
	switch p.t.Kind() {
	case token.CloseParenthesis:
		p.mathNext()
		return true
	case token.EOF:
		return true
	}
	return false
}

func (p *parser) mathCalc(ctx *colorCompCtx, depth int) (numValue, bool) {
	v, ok := p.mathArg(ctx, depth)
	if !ok || !p.mathEnd() {
		if p.t.Is(token.Comma) {
			p.setErrMsg("too many arguments")
		} else {
			p.setErrMsg("missing argument")
		}
		return numValue{}, false
	}
	return v, true
}

func (p *parser) mathMin(ctx *colorCompCtx, depth int) (numValue, bool) {
	acc, ok := p.mathArg(ctx, depth)
	if !ok {
		return numValue{}, false
	}

	for p.mathMore() {
		v, ok := p.mathArg(ctx, depth)
		if !ok || v.kind != acc.kind {
			return numValue{}, false
		}
		acc.value = min(acc.value, v.value)
	}

	if !p.mathEnd() {
		return numValue{}, false
	}

	return acc, true
}

func (p *parser) mathMax(ctx *colorCompCtx, depth int) (numValue, bool) {
	acc, ok := p.mathArg(ctx, depth)
	if !ok {
		return numValue{}, false
	}

	for p.mathMore() {
		v, ok := p.mathArg(ctx, depth)
		if !ok || v.kind != acc.kind {
			return numValue{}, false
		}
		acc.value = max(acc.value, v.value)
	}

	if !p.mathEnd() {
		return numValue{}, false
	}

	return acc, true
}

func (p *parser) mathClamp(ctx *colorCompCtx, depth int) (numValue, bool) {
	var args [3]numValue
	if !p.mathArgs(ctx, depth, allowAll, true, args[:]) {
		return numValue{}, false
	}

	return numValue{
		value: max(
			args[0].value,
			min(
				args[1].value,
				args[2].value,
			),
		),
		kind: args[0].kind,
	}, true
}

func (p *parser) mathRound(ctx *colorCompCtx, depth int) (numValue, bool) {
	strat := roundNearest
	if p.t.Is(token.Ident) {
		if s, ok := lookupRoundStrategy(p.t.Text()); ok {
			strat = s
			p.mathNext()
			if !p.mathMore() {
				return numValue{}, false
			}
		}
	}

	a, ok := p.mathArg(ctx, depth)
	if !ok {
		return numValue{}, false
	}

	step := numValue{value: 1, kind: kindNumber}
	if p.mathMore() {
		if step, ok = p.mathArg(ctx, depth); !ok || step.kind != a.kind {
			return numValue{}, false
		}
	} else if a.kind != kindNumber {
		return numValue{}, false
	}

	if !p.mathEnd() {
		return numValue{}, false
	}

	return numValue{value: roundTo(strat, a.value, step.value), kind: a.kind}, true
}

func (p *parser) mathMod(ctx *colorCompCtx, depth int) (numValue, bool) {
	var args [2]numValue
	if !p.mathArgs(ctx, depth, allowAll, true, args[:]) {
		return numValue{}, false
	}

	return numValue{
		value: cssMod(args[0].value, args[1].value),
		kind:  args[0].kind,
	}, true
}

func (p *parser) mathRem(ctx *colorCompCtx, depth int) (numValue, bool) {
	var args [2]numValue
	if !p.mathArgs(ctx, depth, allowAll, true, args[:]) {
		return numValue{}, false
	}

	return numValue{
		value: cssRem(args[0].value, args[1].value),
		kind:  args[0].kind,
	}, true
}

func (p *parser) mathSin(ctx *colorCompCtx, depth int) (numValue, bool) {
	v, ok := p.mathArg(ctx, depth)
	if !ok || !allowNumAngle.has(v.kind) || !p.mathEnd() {
		return numValue{}, false
	}

	return numValue{
		value: math.Sin(trigRad(v)),
		kind:  kindNumber,
	}, true
}

func (p *parser) mathCos(ctx *colorCompCtx, depth int) (numValue, bool) {
	v, ok := p.mathArg(ctx, depth)
	if !ok || !allowNumAngle.has(v.kind) || !p.mathEnd() {
		return numValue{}, false
	}

	return numValue{
		value: math.Cos(trigRad(v)),
		kind:  kindNumber,
	}, true
}

func (p *parser) mathTan(ctx *colorCompCtx, depth int) (numValue, bool) {
	v, ok := p.mathArg(ctx, depth)
	if !ok || !allowNumAngle.has(v.kind) || !p.mathEnd() {
		return numValue{}, false
	}

	return numValue{
		value: math.Tan(trigRad(v)),
		kind:  kindNumber,
	}, true
}

func (p *parser) mathAsin(ctx *colorCompCtx, depth int) (numValue, bool) {
	v, ok := p.mathArg(ctx, depth)
	if !ok || v.kind != kindNumber || !p.mathEnd() {
		return numValue{}, false
	}

	return numValue{
		value: math.Asin(v.value) * radToDeg,
		kind:  kindAngle,
	}, true
}

func (p *parser) mathAcos(ctx *colorCompCtx, depth int) (numValue, bool) {
	v, ok := p.mathArg(ctx, depth)
	if !ok || v.kind != kindNumber || !p.mathEnd() {
		return numValue{}, false
	}

	return numValue{
		value: math.Acos(v.value) * radToDeg,
		kind:  kindAngle,
	}, true
}

func (p *parser) mathAtan(ctx *colorCompCtx, depth int) (numValue, bool) {
	v, ok := p.mathArg(ctx, depth)
	if !ok || v.kind != kindNumber || !p.mathEnd() {
		return numValue{}, false
	}

	return numValue{
		value: math.Atan(v.value) * radToDeg,
		kind:  kindAngle,
	}, true
}

func (p *parser) mathAtan2(ctx *colorCompCtx, depth int) (numValue, bool) {
	var a [2]numValue
	if !p.mathArgs(ctx, depth, allowAll, true, a[:]) {
		return numValue{}, false
	}
	return numValue{
		value: math.Atan2(a[0].value, a[1].value) * radToDeg,
		kind:  kindAngle,
	}, true
}

func (p *parser) mathPow(ctx *colorCompCtx, depth int) (numValue, bool) {
	var args [2]numValue
	if !p.mathArgs(ctx, depth, allowNumber, true, args[:]) {
		return numValue{}, false
	}
	return numValue{
		value: math.Pow(args[0].value, args[1].value),
		kind:  args[0].kind,
	}, true
}

func (p *parser) mathSqrt(ctx *colorCompCtx, depth int) (numValue, bool) {
	v, ok := p.mathArg(ctx, depth)
	if !ok {
		return numValue{}, false
	}

	if v.kind != kindNumber {
		p.setErrExpected(v.kind, allowNumber, false)
		return numValue{}, false
	}

	if !p.mathEnd() {
		p.setErrMsg("too many arguments")
		return numValue{}, false
	}

	return numValue{value: math.Sqrt(v.value), kind: kindNumber}, true
}

func (p *parser) mathHypot(ctx *colorCompCtx, depth int) (numValue, bool) {
	acc, ok := p.mathArg(ctx, depth)
	if !ok {
		return numValue{}, false
	}

	if !allowNumPct.has(acc.kind) {
		p.setErrExpected(acc.kind, allowNumPct, false)
		return numValue{}, false
	}

	acc.value *= acc.value

	for p.mathMore() {
		v, ok := p.mathArg(ctx, depth)

		if !ok {
			return numValue{}, false
		}

		if v.kind != acc.kind {
			p.setErrMsg("expected arguments of the same type")
			return numValue{}, false
		}

		acc.value += v.value * v.value
	}

	if !p.mathEnd() {
		p.setErrMsg("unexpected token")
		return numValue{}, false
	}

	acc.value = math.Sqrt(acc.value)
	return acc, true
}

func (p *parser) mathLog(ctx *colorCompCtx, depth int) (numValue, bool) {
	start := p.t.TokenStart()
	a, ok := p.mathArg(ctx, depth)
	if !ok {
		return numValue{}, false
	}

	if a.kind != kindNumber {
		p.setErrExpected(a.kind, allowNumber, false)
		p.setErrTextAt(start, p.t.TokenStart())
		return numValue{}, false
	}

	var v float64
	if p.mathMore() {
		start := p.t.TokenStart()
		b, ok := p.mathArg(ctx, depth)
		if !ok {
			return numValue{}, false
		}
		if b.kind != kindNumber {
			p.setErrExpected(b.kind, allowNumber, false)
			p.setErrTextAt(start, p.t.TokenStart())
			return numValue{}, false
		}

		switch b.value {
		case 2:
			v = math.Log2(a.value)
		case 10:
			v = math.Log10(a.value)
		default:
			v = math.Log(a.value) / math.Log(b.value)
		}
	} else {
		v = math.Log(a.value)
	}

	if !p.mathEnd() {
		p.setErrMsg("unexpected token")
		return numValue{}, false
	}

	return numValue{value: v, kind: kindNumber}, true
}

func (p *parser) mathExp(ctx *colorCompCtx, depth int) (numValue, bool) {
	v, ok := p.mathArg(ctx, depth)
	if !ok || v.kind != kindNumber || !p.mathEnd() {
		return numValue{}, false
	}
	return numValue{value: math.Exp(v.value), kind: kindNumber}, true
}

func (p *parser) mathAbs(ctx *colorCompCtx, depth int) (numValue, bool) {
	v, ok := p.mathArg(ctx, depth)
	if !ok || !p.mathEnd() {
		return numValue{}, false
	}
	return numValue{value: math.Abs(v.value), kind: v.kind}, true
}

func (p *parser) mathSign(ctx *colorCompCtx, depth int) (numValue, bool) {
	v, ok := p.mathArg(ctx, depth)
	if !ok || !p.mathEnd() {
		return numValue{}, false
	}
	return numValue{value: signf(v.value), kind: kindNumber}, true
}

func (p *parser) mathRandom() (numValue, bool) {
	return numValue{value: rand.Float64()}, false
}

func (p *parser) parseMathPrefix(ctx *colorCompCtx, depth int) (numValue, bool) {
	switch p.t.Kind() {
	case token.Ident:
		if ctx != nil && ctx.chSet != 0 {
			if chid := space.LookupChIdentFlags(p.t.Text()); ctx.chSet.Has(chid) {
				p.mathNext()
				return numValue{value: ctx.channel(chid), kind: kindNumber}, true
			}
		}

		var buf [16]byte
		n := copy(buf[:], p.t.Text())
		ascii.ToLowerBytes(buf[:n])

		var value float64
		switch string(buf[:n]) {
		case "e":
			value = math.E
		case "pi":
			value = math.Pi
		case "infinity":
			value = math.Inf(1)
		case "-infinity":
			value = math.Inf(-1)
		case "nan":
			value = math.NaN()
		default:
			p.setErrMsg("unknown constant")
			p.setErrTextAt(p.t.TokenStart(), p.t.TokenEnd())
			return numValue{}, false
		}

		p.mathNext()
		return numValue{value: value, kind: kindNumber}, true

	case token.Function:
		return p.parseMathFunc(ctx, depth+1)

	case token.OpenParenthesis:
		p.mathNext()
		if depth > mathMaxDepth {
			return numValue{}, false
		}
		v, ok := p.mathExpr(ctx, 0, depth+1)
		if !ok || !p.mathEnd() {
			return numValue{}, false
		}
		return v, true

	case token.Number:
		value := p.t.Number()
		p.mathNext()
		return numValue{value: value, kind: kindNumber}, true

	case token.Percentage:
		value := p.t.Number()
		p.mathNext()
		return numValue{value: value, kind: kindPercent}, true
	case token.Dimension:
		unit := p.t.Text()
		unitID := lookupUnit(unit)
		if unitID == 0 {
			if len(unit) > 0 && unit[0] == '-' {
				p.setErrMsg("whitespace required around operator")
			} else {
				p.setErrMsg("unknown unit")
			}
			p.setErrTextAt(p.t.TokenStart(), p.t.Position())
			return numValue{}, false
		}

		info := unitID.info()
		if info.rel != unitNone {
			p.setErrMsg("relative unit is not allowed")
			p.setErrTextAt(p.t.TokenStart(), p.t.TokenEnd())
			return numValue{}, false
		}

		value := p.t.Number() * info.scale
		p.mathNext()
		return numValue{value: value, kind: info.kind, unitRel: info.rel}, true
	}

	p.setErrMsg("unexpected token")
	p.setErrTextAt(p.t.TokenStart(), p.t.TokenEnd())
	return numValue{}, false
}

func (p *parser) mathBinary(op byte, a, b numValue) (numValue, bool) {
	switch op {
	case '+':
		if a.kind != b.kind || a.unitRel != b.unitRel {
			p.setErrMsg("incompatible types")
			return numValue{}, false
		}
		return numValue{value: a.value + b.value, kind: a.kind | b.kind}, true
	case '-':
		if a.kind != b.kind || a.unitRel != b.unitRel {
			p.setErrMsg("incompatible types")
			return numValue{}, false
		}
		return numValue{value: a.value - b.value, kind: a.kind | b.kind}, true

	case '*':
		switch {
		case a.kind == kindNumber:
			return numValue{value: a.value * b.value, kind: b.kind, unitRel: b.unitRel}, true
		case b.kind == kindNumber:
			return numValue{value: a.value * b.value, kind: a.kind, unitRel: a.unitRel}, true
		}
		p.setErrMsg("at least one operand must be a number")
		return numValue{}, false

	case '/':
		switch {
		case b.kind == kindNumber:
			return numValue{value: a.value / b.value, kind: a.kind, unitRel: a.unitRel}, true
		case a.kind == b.kind && a.unitRel == b.unitRel:
			return numValue{value: a.value / b.value, kind: kindNumber}, true
		}
		p.setErrMsg("divisor must be a number or the same type as the dividend")
		return numValue{}, false

	default:
		p.setErrMsg("unknown operator")
		return numValue{}, false
	}
}

func (p *parser) mathInfix() (byte, int) {
	if p.t.Kind() != token.Delim {
		return 0, 0
	}
	switch d := p.t.Delim(); d {
	case '+', '-':
		return d, 1
	case '*', '/':
		return d, 2
	}
	return 0, 0
}

func (p *parser) mathExpr(ctx *colorCompCtx, minBP int, depth int) (numValue, bool) {
	start := p.t.TokenStart()
	lhs, ok := p.parseMathPrefix(ctx, depth)
	if !ok {
		return lhs, false
	}

	for {
		op, bp := p.mathInfix()
		if op == 0 {
			switch p.t.Kind() {
			case token.CloseParenthesis, token.EOF:
				return lhs, true
			case token.Comma:
				return lhs, true
			case token.Number, token.Percentage, token.Dimension, token.Ident, token.Function, token.OpenParenthesis:
				p.setErrMsg("missing operator")
				return numValue{}, false
			default:
				p.setErrMsg("unexpected token")
				// p.setErrTextAt(p.t.TokenStart(), p.t.TokenEnd())
				return numValue{}, false
			}
		}

		if bp < minBP {
			return lhs, true
		}

		addSub := op == '+' || op == '-'
		if addSub && !p.mathSpace {
			p.setErrMsg("whitespace required before operator")
			p.setErrTextAt(start, p.t.TokenEnd())
			return numValue{}, false
		}

		start = p.t.TokenStart()
		p.mathNext()

		if addSub && !p.mathSpace {
			p.setErrMsg("whitespace required after operator")
			p.setErrTextAt(start, p.t.TokenEnd())
			return numValue{}, false
		}

		rhs, ok := p.mathExpr(ctx, bp+1, depth)
		if !ok {
			return rhs, false
		}

		if lhs, ok = p.mathBinary(op, lhs, rhs); !ok {
			return lhs, false
		}

	}
}
