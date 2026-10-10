package colors

import "math"

// GradientStop represents a color and its position within a gradient.
type GradientStop struct {
	Color  Color
	Offset float64

	invRange float64

	// Hinted segments interpolate weights between the original color pair,
	// rather than selecting a new hue path between the synthetic stop colors.
	hint   *gradientHint
	weight float64
}

type gradientHint struct {
	left, right Color
}

// HasOffset reports whether the stop has an explicit offset.
func (s GradientStop) HasOffset() bool {
	return !math.IsNaN(s.Offset)
}

// IsHint reports whether the stop is an interpolation hint.
func (s GradientStop) IsHint() bool {
	return !s.Color.space.IsValid()
}
