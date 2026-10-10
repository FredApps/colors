package convert

import "math"

// Conversion path (1 steps):
//
//	HWB
//	-> sRGB
func HwbToSrgb(h, w, b float64) (red, green, blue float64) {
	if sum := w + b; sum >= 1 {
		gray := w / sum
		return gray, gray, gray
	}

	if h < 0 || h >= 360 {
		h = wrap360(h)
	}

	h *= 1 / 60.0
	sector := int(h)
	f := h - float64(sector)

	switch sector {
	case 0:
		red, green = 1, f
	case 1:
		red, green = 1-f, 1
	case 2:
		green, blue = 1, f
	case 3:
		green, blue = 1-f, 1
	case 4:
		red, blue = f, 1
	default:
		red, blue = 1, 1-f
	}

	scale := 1 - w - b

	return red*scale + w, green*scale + w, blue*scale + w
}

// Conversion path (1 steps):
//
//	HWB
//	-> HSL
func HwbToHsl(h, w, b float64) (float64, float64, float64) {
	sum := w + b
	if sum >= 1 {
		l := w / sum
		return h, 0, l
	}

	v := 1 - b
	c := v - w
	l := (v + w) * 0.5

	s := c / (1 - math.Abs(2*l-1))

	return h, s, l
}

// Conversion path (1 steps):
//
//	HWB
//	-> HSV
func HwbToHsv(h, w, b float64) (float64, float64, float64) {
	if sum := w + b; sum >= 1 {
		return h, 0, w / sum
	}

	v := 1 - b

	if v == 0 {
		return h, 0, v
	}

	return h, 1 - w/v, v
}
