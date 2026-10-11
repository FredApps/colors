package colors

import (
	"image/color"
	"reflect"

	"github.com/thmalt/colors/convert"
	"github.com/thmalt/colors/space"
)

// RGBA implements the [color.Color] interface
func (c Color) RGBA() (r, g, b, a uint32) {
	alpha16 := clamp01(c.alpha) * maxUint16
	a = uint32(alpha16 + 0.5)

	if a == 0 {
		return 0, 0, 0, 0
	}

	switch c.space {
	case space.Srgb:
		r = uint32(clamp01(c.c1)*alpha16 + 0.5)
		g = uint32(clamp01(c.c2)*alpha16 + 0.5)
		b = uint32(clamp01(c.c3)*alpha16 + 0.5)
		return
	case space.Hsl, space.Hsv, space.Hwb:
		fr, fg, fb := c.srgb()
		r = uint32(clamp01(fr)*alpha16 + 0.5)
		g = uint32(clamp01(fg)*alpha16 + 0.5)
		b = uint32(clamp01(fb)*alpha16 + 0.5)
		return
	case space.LinearSrgb:
		r = uint32(convert.Srgb16Encode(c.c1)) * a / maxUint16
		g = uint32(convert.Srgb16Encode(c.c2)) * a / maxUint16
		b = uint32(convert.Srgb16Encode(c.c3)) * a / maxUint16
		return
	default:
		fr, fg, fb := c.linearSrgb()
		r = uint32(convert.Srgb16Encode(fr)) * a / maxUint16
		g = uint32(convert.Srgb16Encode(fg)) * a / maxUint16
		b = uint32(convert.Srgb16Encode(fb)) * a / maxUint16
		return
	}
}

// ToRGBA64 converts the color to an sRGB [color.RGBA64].
func (c Color) ToRGBA64() color.RGBA64 {
	alpha16 := clamp01(c.alpha) * maxUint16
	a := uint16(alpha16 + 0.5)

	if a == 0 {
		return color.RGBA64{}
	}

	switch c.space {
	case space.Srgb:
		r := uint16(clamp01(c.c1)*alpha16 + 0.5)
		g := uint16(clamp01(c.c2)*alpha16 + 0.5)
		b := uint16(clamp01(c.c3)*alpha16 + 0.5)
		return color.RGBA64{R: r, G: g, B: b, A: a}
	case space.Hsl, space.Hsv, space.Hwb:
		fr, fg, fb := c.srgb()
		r := uint16(clamp01(fr)*alpha16 + 0.5)
		g := uint16(clamp01(fg)*alpha16 + 0.5)
		b := uint16(clamp01(fb)*alpha16 + 0.5)
		return color.RGBA64{R: r, G: g, B: b, A: a}
	case space.LinearSrgb:
		r := uint16(uint32(convert.Srgb16Encode(c.c1)) * uint32(a) / maxUint16)
		g := uint16(uint32(convert.Srgb16Encode(c.c2)) * uint32(a) / maxUint16)
		b := uint16(uint32(convert.Srgb16Encode(c.c3)) * uint32(a) / maxUint16)
		return color.RGBA64{R: r, G: g, B: b, A: a}
	default:
		fr, fg, fb := c.linearSrgb()
		r := uint16(uint32(convert.Srgb16Encode(fr)) * uint32(a) / maxUint16)
		g := uint16(uint32(convert.Srgb16Encode(fg)) * uint32(a) / maxUint16)
		b := uint16(uint32(convert.Srgb16Encode(fb)) * uint32(a) / maxUint16)
		return color.RGBA64{R: r, G: g, B: b, A: a}
	}
}

// ToNRGBA64 converts the color to an sRGB [color.NRGBA64].
func (c Color) ToNRGBA64() color.NRGBA64 {
	a := uint16(clamp01(c.alpha)*maxUint16 + 0.5)

	switch c.space {
	case space.Srgb:
		r := uint16(clamp01(c.c1)*maxUint16 + 0.5)
		g := uint16(clamp01(c.c2)*maxUint16 + 0.5)
		b := uint16(clamp01(c.c3)*maxUint16 + 0.5)
		return color.NRGBA64{R: r, G: g, B: b, A: a}
	case space.Hsl, space.Hsv, space.Hwb:
		fr, fg, fb := c.srgb()
		r := uint16(clamp01(fr)*maxUint16 + 0.5)
		g := uint16(clamp01(fg)*maxUint16 + 0.5)
		b := uint16(clamp01(fb)*maxUint16 + 0.5)
		return color.NRGBA64{R: r, G: g, B: b, A: a}
	case space.LinearSrgb:
		r := convert.Srgb16Encode(c.c1)
		g := convert.Srgb16Encode(c.c2)
		b := convert.Srgb16Encode(c.c3)
		return color.NRGBA64{R: r, G: g, B: b, A: a}
	default:
		fr, fg, fb := c.linearSrgb()
		r := convert.Srgb16Encode(fr)
		g := convert.Srgb16Encode(fg)
		b := convert.Srgb16Encode(fb)
		return color.NRGBA64{R: r, G: g, B: b, A: a}
	}
}

// ToRGBA converts the color to an sRGB [color.RGBA].
func (c Color) ToRGBA() color.RGBA {
	alpha8 := clamp01(c.alpha) * maxUint8
	a := uint8(alpha8 + 0.5)

	if a == 0 {
		return color.RGBA{}
	}

	switch c.space {
	case space.Srgb:
		r := uint8(clamp01(c.c1)*alpha8 + 0.5)
		g := uint8(clamp01(c.c2)*alpha8 + 0.5)
		b := uint8(clamp01(c.c3)*alpha8 + 0.5)
		return color.RGBA{R: r, G: g, B: b, A: a}
	case space.Hsl, space.Hsv, space.Hwb:
		fr, fg, fb := c.srgb()
		r := uint8(clamp01(fr)*alpha8 + 0.5)
		g := uint8(clamp01(fg)*alpha8 + 0.5)
		b := uint8(clamp01(fb)*alpha8 + 0.5)
		return color.RGBA{R: r, G: g, B: b, A: a}
	case space.LinearSrgb:
		r := uint8(uint16(convert.Srgb8Encode(c.c1)) * uint16(a) / maxUint8)
		g := uint8(uint16(convert.Srgb8Encode(c.c2)) * uint16(a) / maxUint8)
		b := uint8(uint16(convert.Srgb8Encode(c.c3)) * uint16(a) / maxUint8)
		return color.RGBA{R: r, G: g, B: b, A: a}
	default:
		fr, fg, fb := c.linearSrgb()
		r := uint8(uint16(convert.Srgb8Encode(fr)) * uint16(a) / maxUint8)
		g := uint8(uint16(convert.Srgb8Encode(fg)) * uint16(a) / maxUint8)
		b := uint8(uint16(convert.Srgb8Encode(fb)) * uint16(a) / maxUint8)
		return color.RGBA{R: r, G: g, B: b, A: a}
	}
}

// ToNRGBA converts the color to an sRGB [color.NRGBA].
func (c Color) ToNRGBA() color.NRGBA {
	a := uint8(clamp01(c.alpha)*maxUint8 + 0.5)

	switch c.space {
	case space.Srgb:
		r := uint8(clamp01(c.c1)*maxUint8 + 0.5)
		g := uint8(clamp01(c.c2)*maxUint8 + 0.5)
		b := uint8(clamp01(c.c3)*maxUint8 + 0.5)
		return color.NRGBA{R: r, G: g, B: b, A: a}
	case space.Hsl, space.Hsv, space.Hwb:
		fr, fg, fb := c.srgb()
		r := uint8(clamp01(fr)*maxUint8 + 0.5)
		g := uint8(clamp01(fg)*maxUint8 + 0.5)
		b := uint8(clamp01(fb)*maxUint8 + 0.5)
		return color.NRGBA{R: r, G: g, B: b, A: a}
	case space.LinearSrgb:
		r := convert.Srgb8Encode(c.c1)
		g := convert.Srgb8Encode(c.c2)
		b := convert.Srgb8Encode(c.c3)
		return color.NRGBA{R: r, G: g, B: b, A: a}
	default:
		fr, fg, fb := c.linearSrgb()
		r := convert.Srgb8Encode(fr)
		g := convert.Srgb8Encode(fg)
		b := convert.Srgb8Encode(fb)
		return color.NRGBA{R: r, G: g, B: b, A: a}
	}
}

// FromStd converts a [color.Color] to a [Color] in [space.Srgb].
// A nil input returns an invalid [Color].
func FromStd[T color.Color](c T) Color {
	switch c := any(c).(type) {
	case Color:
		c.mutTo(space.Srgb)
		return c
	case *Color:
		if c == nil {
			return Color{}
		}
		return c.To(space.Srgb)

	case color.Alpha:
		return stdAlpha(c.A)
	case *color.Alpha:
		if c == nil {
			return Color{}
		}
		return stdAlpha(c.A)

	case color.Alpha16:
		return stdAlpha16(c.A)
	case *color.Alpha16:
		if c == nil {
			return Color{}
		}
		return stdAlpha16(c.A)

	case color.CMYK:
		return stdCMYK(c.C, c.M, c.Y, c.K)
	case *color.CMYK:
		if c == nil {
			return Color{}
		}
		return stdCMYK(c.C, c.M, c.Y, c.K)

	case color.Gray:
		return stdGray(c.Y)
	case *color.Gray:
		if c == nil {
			return Color{}
		}
		return stdGray(c.Y)

	case color.Gray16:
		return stdGray16(c.Y)
	case *color.Gray16:
		if c == nil {
			return Color{}
		}
		return stdGray16(c.Y)

	case color.NRGBA:
		return stdNRGBA(c.R, c.G, c.B, c.A)
	case *color.NRGBA:
		if c == nil {
			return Color{}
		}
		return stdNRGBA(c.R, c.G, c.B, c.A)

	case color.NRGBA64:
		return stdNRGBA64(c.R, c.G, c.B, c.A)
	case *color.NRGBA64:
		if c == nil {
			return Color{}
		}
		return stdNRGBA64(c.R, c.G, c.B, c.A)

	case color.RGBA:
		return stdRGBA(c.R, c.G, c.B, c.A)
	case *color.RGBA:
		if c == nil {
			return Color{}
		}
		return stdRGBA(c.R, c.G, c.B, c.A)

	case color.RGBA64:
		return stdRGBA64(c.R, c.G, c.B, c.A)
	case *color.RGBA64:
		if c == nil {
			return Color{}
		}
		return stdRGBA64(c.R, c.G, c.B, c.A)

	case color.YCbCr:
		R, G, B, _ := c.RGBA()
		return stdNRGB64(uint16(R), uint16(G), uint16(B))
	case *color.YCbCr:
		if c == nil {
			return Color{}
		}
		R, G, B, _ := c.RGBA()
		return stdNRGB64(uint16(R), uint16(G), uint16(B))

	case color.NYCbCrA:
		R, G, B, _ := c.YCbCr.RGBA()
		return stdNRGBA64(uint16(R), uint16(G), uint16(B), uint16(c.A)*0x101)
	case *color.NYCbCrA:
		if c == nil {
			return Color{}
		}
		R, G, B, _ := c.YCbCr.RGBA()
		return stdNRGBA64(uint16(R), uint16(G), uint16(B), uint16(c.A)*0x101)
	}

	if isNil(c) {
		return Color{}
	}

	R, G, B, A := c.RGBA()
	if A == 0 {
		return SrgbAlpha(0, 0, 0, 0)
	}

	a := float64(A)
	invA := 1 / a

	r := float64(R) * invA
	g := float64(G) * invA
	b := float64(B) * invA
	a *= invMaxUint16

	return SrgbAlpha(r, g, b, a)
}

// StdAlpha returns a sRGB [Color] from a [color.Alpha].
func StdAlpha(c color.Alpha) Color {
	return stdAlpha(c.A)
}

// StdAlpha16 returns a sRGB [Color] from a [color.Alpha16].
func StdAlpha16(c color.Alpha16) Color {
	return stdAlpha16(c.A)
}

// StdCMYK returns a sRGB [Color] from a [color.CMYK].
func StdCMYK(c color.CMYK) Color {
	return stdCMYK(c.C, c.M, c.Y, c.K)
}

// StdGray returns a sRGB [Color] from a [color.Gray].
func StdGray(c color.Gray) Color {
	return stdGray(c.Y)
}

// StdGray16 returns a sRGB [Color] from a [color.Gray16].
func StdGray16(c color.Gray16) Color {
	return stdGray16(c.Y)
}

// StdNRGBA returns a sRGB [Color] from a [color.NRGBA].
func StdNRGBA(c color.NRGBA) Color {
	return stdNRGBA(c.R, c.G, c.B, c.A)
}

// StdNRGBA64 returns a sRGB [Color] from a [color.NRGBA64].
func StdNRGBA64(c color.NRGBA64) Color {
	return stdNRGBA64(c.R, c.G, c.B, c.A)
}

// StdRGBA returns a sRGB [Color] from a [color.RGBA].
func StdRGBA(c color.RGBA) Color {
	return stdRGBA(c.R, c.G, c.B, c.A)
}

// StdRGBA64 returns a sRGB [Color] from a [color.RGBA64].
func StdRGBA64(c color.RGBA64) Color {
	return stdRGBA64(c.R, c.G, c.B, c.A)
}

// StdYCbCr returns a sRGB [Color] from a [color.YCbCr].
func StdYCbCr(c color.YCbCr) Color {
	R, G, B, _ := c.RGBA()
	return stdNRGB64(uint16(R), uint16(G), uint16(B))
}

// StdNYCbCrA returns a sRGB [Color] from a [color.NYCbCrA].
func StdNYCbCrA(c color.NYCbCrA) Color {
	R, G, B, _ := c.YCbCr.RGBA()
	return stdNRGBA64(uint16(R), uint16(G), uint16(B), uint16(c.A)*0x101)
}

func stdColor(R, G, B, A uint32) Color {
	if A == 0 {
		return SrgbAlpha(0, 0, 0, 0)
	}

	a := float64(A)
	invA := 1 / a

	r := float64(R) * invA
	g := float64(G) * invA
	b := float64(B) * invA
	a *= invMaxUint16

	return SrgbAlpha(r, g, b, a)
}

func stdAlpha(A uint8) Color {
	a := float64(A) * invMaxUint8
	return SrgbAlpha(1, 1, 1, a)
}

func stdAlpha16(A uint16) Color {
	a := float64(A) * invMaxUint16
	return SrgbAlpha(1, 1, 1, a)
}

func stdCMYK(C, M, Y, K uint8) Color {
	w := 1 - float64(K)*invMaxUint8

	r := (1 - float64(C)*invMaxUint8) * w
	g := (1 - float64(M)*invMaxUint8) * w
	b := (1 - float64(Y)*invMaxUint8) * w

	return Srgb(r, g, b)
}

func stdGray(Y uint8) Color {
	y := float64(Y) * invMaxUint8
	return Srgb(y, y, y)
}

func stdGray16(Y uint16) Color {
	y := float64(Y) * invMaxUint16
	return Srgb(y, y, y)
}

func stdNRGBA(R, G, B, A uint8) Color {
	r := float64(R) * invMaxUint8
	g := float64(G) * invMaxUint8
	b := float64(B) * invMaxUint8
	a := float64(A) * invMaxUint8

	return SrgbAlpha(r, g, b, a)
}

func stdNRGBA64(R, G, B, A uint16) Color {
	r := float64(R) * invMaxUint16
	g := float64(G) * invMaxUint16
	b := float64(B) * invMaxUint16
	a := float64(A) * invMaxUint16

	return SrgbAlpha(r, g, b, a)
}

func stdNRGB64(R, G, B uint16) Color {
	return Srgb(
		float64(R)*invMaxUint16,
		float64(G)*invMaxUint16,
		float64(B)*invMaxUint16,
	)
}

func stdRGBA(R, G, B, A uint8) Color {
	if A == 0 {
		return SrgbAlpha(0, 0, 0, 0)
	}

	a := float64(A)
	invA := 1 / a

	r := float64(R) * invA
	g := float64(G) * invA
	b := float64(B) * invA
	a *= invMaxUint8

	return SrgbAlpha(r, g, b, a)
}

func stdRGBA64(R, G, B, A uint16) Color {
	if A == 0 {
		return SrgbAlpha(0, 0, 0, 0)
	}

	a := float64(A)
	invA := 1 / a

	r := float64(R) * invA
	g := float64(G) * invA
	b := float64(B) * invA
	a *= invMaxUint16

	return SrgbAlpha(r, g, b, a)
}

func isNil(v any) bool {
	if v == nil {
		return true
	}

	rv := reflect.ValueOf(v)
	return rv.Kind() == reflect.Pointer && rv.IsNil()
}
