package colors

import (
	"math"
	"strconv"
	"strings"

	"github.com/thmalt/colors/internal/ascii"
	"github.com/thmalt/colors/space"
)

type colorCompCtx struct {
	channels [space.MaxChannelCount + 1]float64
	space    space.Space
	chSet    space.ChIdentSet
	err      error
}

func (ctx *colorCompCtx) setSpace(s space.Space) {
	ctx.space = s
	ctx.chSet = s.ChIdentSet()
}

func (ctx *colorCompCtx) setChannel3(c1, c2, c3, alpha float64) {
	ctx.channels[0] = c1
	ctx.channels[1] = c2
	ctx.channels[2] = c3
	ctx.channels[3] = alpha
}

func (ctx *colorCompCtx) channel(ch space.ChIdent) float64 {
	idx := uint8(ctx.space.ChIndex(ch))
	if idx >= uint8(len(ctx.channels)) {
		return 0
	}
	return ctx.channels[idx]
}

type numKind uint8

const (
	kindNumber numKind = iota
	kindPercent
	kindAngle
	kindTime
	kindLength
	kindFrequency
	kindResolution
	kindFlex
)

type numKindSet uint8

const (
	allowNumber     numKindSet = 1 << kindNumber
	allowPercent    numKindSet = 1 << kindPercent
	allowAngle      numKindSet = 1 << kindAngle
	allowTime       numKindSet = 1 << kindTime
	allowLength     numKindSet = 1 << kindLength
	allowFrequency  numKindSet = 1 << kindFrequency
	allowResolution numKindSet = 1 << kindResolution
	allowFlex       numKindSet = 1 << kindFlex

	allowAll = allowNumber | allowPercent | allowAngle

	allowNumPct   = allowNumber | allowPercent
	allowNumAngle = allowNumber | allowAngle

	allowHue = allowNumAngle

	allowNumPctAngle = allowNumber | allowPercent | allowAngle
)

func (k numKind) set() numKindSet { return 1 << (k & 7) }

func (s numKindSet) has(k numKind) bool { return s&(1<<(k&7)) != 0 }

var kindNames = [...]string{"number", "percentage", "angle", "time", "length", "frequency", "resolution", "flex"}

func (s numKindSet) writeNames(b *strings.Builder, none bool) {
	var names [6]string
	n := 0
	for i, name := range kindNames {
		if s&(1<<i) != 0 {
			names[n] = name
			n++
		}
	}
	if none {
		names[n] = "none"
		n++
	}
	switch n {
	case 1:
		b.WriteString(names[0])
	default:
		b.WriteString(strings.Join(names[:n-1], ", "))
		b.WriteString(" or ")
		b.WriteString(names[n-1])
	}
}

func (k numKind) String() string {
	if int(k) < len(kindNames) {
		return kindNames[k]
	}
	return "numKind(" + strconv.Itoa(int(k)) + ")"
}

type numValue struct {
	value   float64
	kind    numKind
	isNone  bool
	unitRel unitID
}

func (v numValue) sanitizeNumber() float64 {
	const channelLimit = 1e15

	x := v.value
	switch {
	case x != x:
		return 0
	case x > channelLimit:
		return channelLimit
	case x < -channelLimit:
		return -channelLimit
	}
	return x
}

func (v numValue) sanitizeHue() float64 {
	x := v.value
	if x != x || math.IsInf(x, 0) {
		return 0
	}

	if x >= 0 && x < 360 {
		return x
	}

	h := x - math.Floor(x/360)*360
	if h >= 360 || h < 0 {
		return 0
	}
	return h
}

func new3(sp space.Space, c1, c2, c3, alpha float64) Color {
	return Color{
		space: sp,
		c1:    c1,
		c2:    c2,
		c3:    c3,
		alpha: alpha,
	}
}

type mathFunc uint8

const (
	fnUnknown mathFunc = iota

	fnCalc   // calc
	fnMin    // min
	fnMax    // max
	fnClamp  // clamp
	fnRound  // round
	fnMod    // mod
	fnRem    // rem
	fnSin    // sin
	fnCos    // cos
	fnTan    // tan
	fnAsin   // asin
	fnAcos   // acos
	fnAtan   // atan
	fnAtan2  // atan2
	fnPow    // pow
	fnSqrt   // sqrt
	fnHypot  // hypot
	fnLog    // log
	fnExp    // exp
	fnAbs    // abs
	fnSign   // sign
	fnRandom // random
)

func lookupMathFunc(name string) mathFunc {
	var buf [16]byte
	n := copy(buf[:], name)
	ascii.ToLowerBytes(buf[:n])

	switch string(buf[:n]) {
	case "calc":
		return fnCalc
	case "min":
		return fnMin
	case "max":
		return fnMax
	case "clamp":
		return fnClamp
	case "round":
		return fnRound
	case "mod":
		return fnMod
	case "rem":
		return fnRem
	case "sin":
		return fnSin
	case "cos":
		return fnCos
	case "tan":
		return fnTan
	case "asin":
		return fnAsin
	case "acos":
		return fnAcos
	case "atan":
		return fnAtan
	case "atan2":
		return fnAtan2
	case "pow":
		return fnPow
	case "sqrt":
		return fnSqrt
	case "hypot":
		return fnHypot
	case "log":
		return fnLog
	case "exp":
		return fnExp
	case "abs":
		return fnAbs
	case "sign":
		return fnSign
	// case "random":
	// return fnRandom
	default:
		return fnUnknown
	}
}

func (f mathFunc) String() string {
	switch f {
	case fnUnknown:
		return "mathFunc(<unknown>)"
	case fnCalc:
		return "calc"
	case fnMin:
		return "min"
	case fnMax:
		return "max"
	case fnClamp:
		return "clamp"
	case fnRound:
		return "round"
	case fnMod:
		return "mod"
	case fnRem:
		return "rem"
	case fnSin:
		return "sin"
	case fnCos:
		return "cos"
	case fnTan:
		return "tan"
	case fnAsin:
		return "asin"
	case fnAcos:
		return "acos"
	case fnAtan:
		return "atan"
	case fnAtan2:
		return "atan2"
	case fnPow:
		return "pow"
	case fnSqrt:
		return "sqrt"
	case fnHypot:
		return "hypot"
	case fnLog:
		return "log"
	case fnExp:
		return "exp"
	case fnAbs:
		return "abs"
	case fnSign:
		return "sign"
	case fnRandom:
		return "random"
	default:
		return "mathFunc(" + strconv.Itoa(int(f)) + ")"
	}
}

type roundStrategy uint8

const (
	roundNearest roundStrategy = iota
	roundUp
	roundDown
	roundToZero
)

func lookupRoundStrategy(s string) (roundStrategy, bool) {
	var buf [8]byte
	n := copy(buf[:], s)
	ascii.ToLowerBytes(buf[:n])

	switch string(buf[:n]) {
	case "nearest":
		return roundNearest, true
	case "up":
		return roundUp, true
	case "down":
		return roundDown, true
	case "to-zero":
		return roundToZero, true
	}
	return 0, false
}

func roundTo(s roundStrategy, a, step float64) float64 {
	if step == 0 || math.IsNaN(step) {
		return math.NaN()
	}
	step = math.Abs(step)
	q := a / step
	switch s {
	case roundUp:
		q = math.Ceil(q)
	case roundDown:
		q = math.Floor(q)
	case roundToZero:
		q = math.Trunc(q)
	default:
		q = math.Floor(q + 0.5)
	}
	return q * step
}

func signf(x float64) float64 {
	switch {
	case x > 0:
		return 1
	case x < 0:
		return -1
	}
	return x // +0, -0, NaN
}

func cssRem(a, b float64) float64 { return math.Mod(a, b) }

func cssMod(a, b float64) float64 {
	r := math.Mod(a, b)
	if r != 0 && (r < 0) != (b < 0) {
		r += b
	}
	return r
}

const degToRad = math.Pi / 180
const radToDeg = 180 / math.Pi

func trigRad(v numValue) float64 {
	if v.kind == kindAngle {
		return v.value * degToRad
	}
	return v.value
}

type unitID uint8

type unitInfo struct {
	scale float64
	kind  numKind
	rel   unitID
}

const (
	unitNone unitID = iota

	// length
	unitEm    // em
	unitRem   // rem
	unitEx    // ex
	unitRex   // rex
	unitCap   // cap
	unitRcap  // rcap
	unitCh    // ch
	unitRch   // rch
	unitIc    // ic
	unitRic   // ric
	unitLh    // lh
	unitRlh   // rlh
	unitVw    // vw
	unitVh    // vh
	unitVi    // vi
	unitVb    // vb
	unitVmin  // vmin
	unitVmax  // vmax
	unitSvw   // svw
	unitSvh   // svh
	unitSvi   // svi
	unitSvb   // svb
	unitSvmin // svmin
	unitSvmax // svmax
	unitLvw   // lvw
	unitLvh   // lvh
	unitLvi   // lvi
	unitLvb   // lvb
	unitLvmin // lvmin
	unitLvmax // lvmax
	unitDvw   // dvw
	unitDvh   // dvh
	unitDvi   // dvi
	unitDvb   // dvb
	unitDvmin // dvmin
	unitDvmax // dvmax
	unitCqw   // cqw
	unitCqh   // cqh
	unitCqi   // cqi
	unitCqb   // cqb
	unitCqmin // cqmin
	unitCqmax // cqmax
	unitCm    // cm
	unitMm    // mm
	unitQ     // q
	unitIn    // in
	unitPt    // pt
	unitPc    // pc
	unitPx    // px

	// angle
	unitDeg  // deg
	unitGrad // grad
	unitRad  // rad
	unitTurn // turn

	// time
	unitS   // s
	unitMs  // ms
	unitHz  // Hz
	unitKHz // kHz

	// resolution
	unitDpi  // dpi
	unitDpcm // dpcm
	unitDppx // dppx
	unitX    // x

	// flex
	unitFr // fr

	unitCount
)

var unitTable = [256]unitInfo{
	// relative length
	unitEm:    {kind: kindLength, scale: 1, rel: unitEm},
	unitRem:   {kind: kindLength, scale: 1, rel: unitRem},
	unitEx:    {kind: kindLength, scale: 1, rel: unitEx},
	unitRex:   {kind: kindLength, scale: 1, rel: unitRex},
	unitCap:   {kind: kindLength, scale: 1, rel: unitCap},
	unitRcap:  {kind: kindLength, scale: 1, rel: unitRcap},
	unitCh:    {kind: kindLength, scale: 1, rel: unitCh},
	unitRch:   {kind: kindLength, scale: 1, rel: unitRch},
	unitIc:    {kind: kindLength, scale: 1, rel: unitIc},
	unitRic:   {kind: kindLength, scale: 1, rel: unitRic},
	unitLh:    {kind: kindLength, scale: 1, rel: unitLh},
	unitRlh:   {kind: kindLength, scale: 1, rel: unitRlh},
	unitVw:    {kind: kindLength, scale: 1, rel: unitVw},
	unitVh:    {kind: kindLength, scale: 1, rel: unitVh},
	unitVi:    {kind: kindLength, scale: 1, rel: unitVi},
	unitVb:    {kind: kindLength, scale: 1, rel: unitVb},
	unitVmin:  {kind: kindLength, scale: 1, rel: unitVmin},
	unitVmax:  {kind: kindLength, scale: 1, rel: unitVmax},
	unitSvw:   {kind: kindLength, scale: 1, rel: unitSvw},
	unitSvh:   {kind: kindLength, scale: 1, rel: unitSvh},
	unitSvi:   {kind: kindLength, scale: 1, rel: unitSvi},
	unitSvb:   {kind: kindLength, scale: 1, rel: unitSvb},
	unitSvmin: {kind: kindLength, scale: 1, rel: unitSvmin},
	unitSvmax: {kind: kindLength, scale: 1, rel: unitSvmax},
	unitLvw:   {kind: kindLength, scale: 1, rel: unitLvw},
	unitLvh:   {kind: kindLength, scale: 1, rel: unitLvh},
	unitLvi:   {kind: kindLength, scale: 1, rel: unitLvi},
	unitLvb:   {kind: kindLength, scale: 1, rel: unitLvb},
	unitLvmin: {kind: kindLength, scale: 1, rel: unitLvmin},
	unitLvmax: {kind: kindLength, scale: 1, rel: unitLvmax},
	unitDvw:   {kind: kindLength, scale: 1, rel: unitDvw},
	unitDvh:   {kind: kindLength, scale: 1, rel: unitDvh},
	unitDvi:   {kind: kindLength, scale: 1, rel: unitDvi},
	unitDvb:   {kind: kindLength, scale: 1, rel: unitDvb},
	unitDvmin: {kind: kindLength, scale: 1, rel: unitDvmin},
	unitDvmax: {kind: kindLength, scale: 1, rel: unitDvmax},
	unitCqw:   {kind: kindLength, scale: 1, rel: unitCqw},
	unitCqh:   {kind: kindLength, scale: 1, rel: unitCqh},
	unitCqi:   {kind: kindLength, scale: 1, rel: unitCqi},
	unitCqb:   {kind: kindLength, scale: 1, rel: unitCqb},
	unitCqmin: {kind: kindLength, scale: 1, rel: unitCqmin},
	unitCqmax: {kind: kindLength, scale: 1, rel: unitCqmax},

	// length
	unitCm: {kind: kindLength, scale: 96 / 2.54},
	unitMm: {kind: kindLength, scale: 96 / 25.4},
	unitQ:  {kind: kindLength, scale: 96 / 101.6},
	unitIn: {kind: kindLength, scale: 96},
	unitPt: {kind: kindLength, scale: 96 / 72.},
	unitPc: {kind: kindLength, scale: 16},
	unitPx: {kind: kindLength, scale: 1},

	// angle
	unitDeg:  {kind: kindAngle, scale: 1},
	unitGrad: {kind: kindAngle, scale: 0.9},
	unitRad:  {kind: kindAngle, scale: 180 / math.Pi},
	unitTurn: {kind: kindAngle, scale: 360},

	// time
	unitS:   {kind: kindTime, scale: 1},
	unitMs:  {kind: kindTime, scale: 1 / 1000.},
	unitHz:  {kind: kindTime, scale: 1},
	unitKHz: {kind: kindTime, scale: 1000},

	// resolution
	unitDpi:  {kind: kindResolution, scale: 1 / 96.},
	unitDpcm: {kind: kindResolution, scale: 2.54 / 96},
	unitDppx: {kind: kindResolution, scale: 1},
	unitX:    {kind: kindResolution, scale: 1},

	// flex
	unitFr: {kind: kindFlex, scale: 1},
}

func (u unitID) info() unitInfo { return unitTable[u] }

func lookupUnit(unit string) unitID {
	var buf [8]byte
	n := copy(buf[:], unit)
	ascii.ToLowerBytes(buf[:n])

	switch string(buf[:n]) {
	case "em":
		return unitEm
	case "rem":
		return unitRem
	case "ex":
		return unitEx
	case "rex":
		return unitRex
	case "cap":
		return unitCap
	case "rcap":
		return unitRcap
	case "ch":
		return unitCh
	case "rch":
		return unitRch
	case "ic":
		return unitIc
	case "ric":
		return unitRic
	case "lh":
		return unitLh
	case "rlh":
		return unitRlh
	case "vw":
		return unitVw
	case "vh":
		return unitVh
	case "vi":
		return unitVi
	case "vb":
		return unitVb
	case "vmin":
		return unitVmin
	case "vmax":
		return unitVmax
	case "svw":
		return unitSvw
	case "svh":
		return unitSvh
	case "svi":
		return unitSvi
	case "svb":
		return unitSvb
	case "svmin":
		return unitSvmin
	case "svmax":
		return unitSvmax
	case "lvw":
		return unitLvw
	case "lvh":
		return unitLvh
	case "lvi":
		return unitLvi
	case "lvb":
		return unitLvb
	case "lvmin":
		return unitLvmin
	case "lvmax":
		return unitLvmax
	case "dvw":
		return unitDvw
	case "dvh":
		return unitDvh
	case "dvi":
		return unitDvi
	case "dvb":
		return unitDvb
	case "dvmin":
		return unitDvmin
	case "dvmax":
		return unitDvmax
	case "cqw":
		return unitCqw
	case "cqh":
		return unitCqh
	case "cqi":
		return unitCqi
	case "cqb":
		return unitCqb
	case "cqmin":
		return unitCqmin
	case "cqmax":
		return unitCqmax
	case "cm":
		return unitCm
	case "mm":
		return unitMm
	case "q":
		return unitQ
	case "in":
		return unitIn
	case "pt":
		return unitPt
	case "pc":
		return unitPc
	case "px":
		return unitPx
	case "deg":
		return unitDeg
	case "grad":
		return unitGrad
	case "rad":
		return unitRad
	case "turn":
		return unitTurn
	case "s":
		return unitS
	case "ms":
		return unitMs
	case "hz":
		return unitHz
	case "khz":
		return unitKHz
	case "dpi":
		return unitDpi
	case "dpcm":
		return unitDpcm
	case "dppx":
		return unitDppx
	case "x":
		return unitX
	case "fr":
		return unitFr
	default:
		return unitNone
	}
}
