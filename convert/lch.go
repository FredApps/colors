package convert

import "math"

const (
	EpsLch   = 150 * 1e-5
	EpsLchuv = 215 * 1e-5
	EpsOklch = 0.4 * 1e-5

	radToDeg = 180 / math.Pi
	degToRad = math.Pi / 180
)

// LxyToLch converts Lxy Cartesian coordinates to LCh polar coordinates.
// This conversion only changes the coordinate representation; it does not
// involve chromatic adaptation or any other color-space transformation.
func LxyToLch(l, x, y, eps float64) (float64, float64, float64) {
	c := math.Sqrt(x*x + y*y)
	if c <= eps {
		return l, c, 0
	}

	h := math.Atan2(y, x) * radToDeg
	if h < 0 {
		h += 360
	}

	return l, c, h
}

// LchToLxy converts LCh polar coordinates to Lxy Cartesian coordinates.
// This conversion only changes the coordinate representation; it does not
// involve chromatic adaptation or any other color-space transformation.
func LchToLxy(l, c, h float64) (float64, float64, float64) {
	rad := h * degToRad
	sin, cos := math.Sincos(rad)

	x := c * cos
	y := c * sin

	return l, x, y
}

func labToLch(l, x, y float64) (float64, float64, float64) {
	c := math.Sqrt(x*x + y*y)
	if c <= EpsLch {
		return l, c, 0
	}

	h := math.Atan2(y, x) * radToDeg
	if h < 0 {
		h += 360
	}

	return l, c, h
}

func luvToLchuv(l, x, y float64) (float64, float64, float64) {
	c := math.Sqrt(x*x + y*y)
	if c <= EpsLchuv {
		return l, c, 0
	}

	h := math.Atan2(y, x) * radToDeg
	if h < 0 {
		h += 360
	}

	return l, c, h
}

func oklabToOklch(l, x, y float64) (float64, float64, float64) {
	c := math.Sqrt(x*x + y*y)
	if c <= EpsOklch {
		return l, c, 0
	}

	h := math.Atan2(y, x) * radToDeg
	if h < 0 {
		h += 360
	}

	return l, c, h
}
