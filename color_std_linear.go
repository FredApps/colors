package colors

import (
	"image/color"

	"github.com/thmalt/colors/convert"
	"github.com/thmalt/colors/space"
)

// FromStdToLinear converts a [color.Color] to a [Color] in [space.LinearSrgb].
// A nil input returns an invalid [Color].
func FromStdToLinear[T color.Color](c T) Color {
	switch c := any(c).(type) {
	case Color:
		c.mutTo(space.LinearSrgb)
		return c
	case *Color:
		if c == nil {
			return Color{}
		}
		return c.To(space.LinearSrgb)

	case color.Alpha:
		return stdAlphaToLinear(c.A)
	case *color.Alpha:
		if c == nil {
			return Color{}
		}
		return stdAlphaToLinear(c.A)

	case color.Alpha16:
		return stdAlpha16ToLinear(c.A)
	case *color.Alpha16:
		if c == nil {
			return Color{}
		}
		return stdAlpha16ToLinear(c.A)

	case color.CMYK:
		return stdCMYKToLinear(c.C, c.M, c.Y, c.K)
	case *color.CMYK:
		if c == nil {
			return Color{}
		}
		return stdCMYKToLinear(c.C, c.M, c.Y, c.K)

	case color.Gray:
		return stdGrayToLinear(c.Y)
	case *color.Gray:
		if c == nil {
			return Color{}
		}
		return stdGrayToLinear(c.Y)

	case color.Gray16:
		return stdGray16ToLinear(c.Y)
	case *color.Gray16:
		if c == nil {
			return Color{}
		}
		return stdGray16ToLinear(c.Y)

	case color.NRGBA:
		return stdNRGBAToLinear(c.R, c.G, c.B, c.A)
	case *color.NRGBA:
		if c == nil {
			return Color{}
		}
		return stdNRGBAToLinear(c.R, c.G, c.B, c.A)

	case color.NRGBA64:
		return stdNRGBA64ToLinear(c.R, c.G, c.B, c.A)
	case *color.NRGBA64:
		if c == nil {
			return Color{}
		}
		return stdNRGBA64ToLinear(c.R, c.G, c.B, c.A)

	case color.RGBA:
		return stdRGBAToLinear(c.R, c.G, c.B, c.A)
	case *color.RGBA:
		if c == nil {
			return Color{}
		}
		return stdRGBAToLinear(c.R, c.G, c.B, c.A)

	case color.RGBA64:
		return stdRGBA64ToLinear(c.R, c.G, c.B, c.A)
	case *color.RGBA64:
		if c == nil {
			return Color{}
		}
		return stdRGBA64ToLinear(c.R, c.G, c.B, c.A)

	case color.YCbCr:
		R, G, B, _ := c.RGBA()
		return stdNRGB64ToLinear(uint16(R), uint16(G), uint16(B))
	case *color.YCbCr:
		if c == nil {
			return Color{}
		}
		R, G, B, _ := c.RGBA()
		return stdNRGB64ToLinear(uint16(R), uint16(G), uint16(B))

	case color.NYCbCrA:
		R, G, B, _ := c.YCbCr.RGBA()
		return stdNRGBA64ToLinear(uint16(R), uint16(G), uint16(B), uint16(c.A)*0x101)
	case *color.NYCbCrA:
		if c == nil {
			return Color{}
		}
		R, G, B, _ := c.YCbCr.RGBA()
		return stdNRGBA64ToLinear(uint16(R), uint16(G), uint16(B), uint16(c.A)*0x101)
	}

	if isNil(c) {
		return Color{}
	}

	R, G, B, A := c.RGBA()
	if A == 0 {
		return LinearSrgbAlpha(0, 0, 0, 0)
	}

	if A == maxUint16 {
		return LinearSrgb(
			convert.Srgb16Decode(uint16(R)),
			convert.Srgb16Decode(uint16(G)),
			convert.Srgb16Decode(uint16(B)),
		)
	}

	return LinearSrgbAlpha(
		convert.Srgb16Decode(uint16(R*maxUint16/A)),
		convert.Srgb16Decode(uint16(G*maxUint16/A)),
		convert.Srgb16Decode(uint16(B*maxUint16/A)),
		float64(A)*invMaxUint16,
	)
}

// StdAlphaToLinear returns a linear sRGB [Color] from a [color.Alpha].
func StdAlphaToLinear(c color.Alpha) Color {
	return stdAlphaToLinear(c.A)
}

// StdAlpha16ToLinear returns a linear sRGB [Color] from a [color.Alpha16].
func StdAlpha16ToLinear(c color.Alpha16) Color {
	return stdAlpha16ToLinear(c.A)
}

// StdCMYKToLinear returns a linear sRGB [Color] from a [color.CMYK].
func StdCMYKToLinear(c color.CMYK) Color {
	return stdCMYKToLinear(c.C, c.M, c.Y, c.K)
}

// StdGrayToLinear returns a linear sRGB [Color] from a [color.Gray].
func StdGrayToLinear(c color.Gray) Color {
	return stdGrayToLinear(c.Y)
}

// StdGray16ToLinear returns a linear sRGB [Color] from a [color.Gray16].
func StdGray16ToLinear(c color.Gray16) Color {
	return stdGray16ToLinear(c.Y)
}

// StdNRGBAToLinear returns a linear sRGB [Color] from a [color.NRGBA].
func StdNRGBAToLinear(c color.NRGBA) Color {
	return stdNRGBAToLinear(c.R, c.G, c.B, c.A)
}

// StdNRGBA64ToLinear returns a linear sRGB [Color] from a [color.NRGBA64].
func StdNRGBA64ToLinear(c color.NRGBA64) Color {
	return stdNRGBA64ToLinear(c.R, c.G, c.B, c.A)
}

// StdRGBAToLinear returns a linear sRGB [Color] from a [color.RGBA].
func StdRGBAToLinear(c color.RGBA) Color {
	return stdRGBAToLinear(c.R, c.G, c.B, c.A)
}

// StdRGBA64ToLinear returns a linear sRGB [Color] from a [color.RGBA64].
func StdRGBA64ToLinear(c color.RGBA64) Color {
	return stdRGBA64ToLinear(c.R, c.G, c.B, c.A)
}

// StdYCbCrToLinear returns a linear sRGB [Color] from a [color.YCbCr].
func StdYCbCrToLinear(c color.YCbCr) Color {
	R, G, B, _ := c.RGBA()
	return stdNRGB64ToLinear(uint16(R), uint16(G), uint16(B))
}

// StdNYCbCrAToLinear returns a linear sRGB [Color] from a [color.NYCbCrA].
func StdNYCbCrAToLinear(c color.NYCbCrA) Color {
	R, G, B, _ := c.YCbCr.RGBA()
	return stdNRGBA64ToLinear(uint16(R), uint16(G), uint16(B), uint16(c.A)*0x101)
}

func stdColorToLinear(R, G, B, A uint32) Color {
	if A == 0 {
		return LinearSrgbAlpha(0, 0, 0, 0)
	}

	if A == maxUint16 {
		return LinearSrgb(
			convert.Srgb16Decode(uint16(R)),
			convert.Srgb16Decode(uint16(G)),
			convert.Srgb16Decode(uint16(B)),
		)
	}

	return LinearSrgbAlpha(
		convert.Srgb16Decode(uint16(R*maxUint16/A)),
		convert.Srgb16Decode(uint16(G*maxUint16/A)),
		convert.Srgb16Decode(uint16(B*maxUint16/A)),
		float64(A)*invMaxUint16,
	)
}

func stdAlphaToLinear(A uint8) Color {
	a := float64(A) * invMaxUint8
	return LinearSrgbAlpha(1, 1, 1, a)
}

func stdAlpha16ToLinear(A uint16) Color {
	a := float64(A) * invMaxUint16
	return LinearSrgbAlpha(1, 1, 1, a)
}

func stdCMYKToLinear(C, M, Y, K uint8) Color {
	w := maxUint16 - uint32(K)*scale8To16

	r := (maxUint16 - uint32(C)*scale8To16) * w / maxUint16
	g := (maxUint16 - uint32(M)*scale8To16) * w / maxUint16
	b := (maxUint16 - uint32(Y)*scale8To16) * w / maxUint16

	return LinearSrgb(
		convert.Srgb8Decode(uint8(r>>8)),
		convert.Srgb8Decode(uint8(g>>8)),
		convert.Srgb8Decode(uint8(b>>8)),
	)
}

func stdGrayToLinear(Y uint8) Color {
	y := convert.Srgb8Decode(Y)
	return LinearSrgb(y, y, y)
}

func stdGray16ToLinear(Y uint16) Color {
	y := convert.Srgb16Decode(Y)
	return LinearSrgb(y, y, y)
}

func stdNRGBAToLinear(R, G, B, A uint8) Color {
	return LinearSrgbAlpha(
		convert.Srgb8Decode(R),
		convert.Srgb8Decode(G),
		convert.Srgb8Decode(B),
		float64(A)*invMaxUint8,
	)
}

func stdNRGBA64ToLinear(R, G, B, A uint16) Color {
	return LinearSrgbAlpha(
		convert.Srgb16Decode(R),
		convert.Srgb16Decode(G),
		convert.Srgb16Decode(B),
		float64(A)*invMaxUint16,
	)
}

func stdNRGB64ToLinear(R, G, B uint16) Color {
	return LinearSrgb(
		convert.Srgb16Decode(R),
		convert.Srgb16Decode(G),
		convert.Srgb16Decode(B),
	)
}

func stdRGBAToLinear(R, G, B, A uint8) Color {
	if A == 0 {
		return LinearSrgbAlpha(0, 0, 0, 0)
	}

	return LinearSrgbAlpha(
		convert.Srgb8Decode(uint8(uint32(R)*maxUint8/uint32(A))),
		convert.Srgb8Decode(uint8(uint32(G)*maxUint8/uint32(A))),
		convert.Srgb8Decode(uint8(uint32(B)*maxUint8/uint32(A))),
		float64(A)*invMaxUint8,
	)
}

func stdRGBA64ToLinear(R, G, B, A uint16) Color {
	if A == 0 {
		return LinearSrgbAlpha(0, 0, 0, 0)
	}

	return LinearSrgbAlpha(
		convert.Srgb16Decode(uint16(uint32(R)*maxUint16/uint32(A))),
		convert.Srgb16Decode(uint16(uint32(G)*maxUint16/uint32(A))),
		convert.Srgb16Decode(uint16(uint32(B)*maxUint16/uint32(A))),
		float64(A)*invMaxUint16,
	)
}
