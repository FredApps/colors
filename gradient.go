package colors

import (
	"math"
	"slices"
	"sort"
)

type Spread uint8

const (
	SpreadPad Spread = iota
	SpreadRepeat
	SpreadReflect
)

// Gradient stores a mixer and color stops for color interpolation.
// An empty gradient, including the zero value, returns an invalid zero Color
// when sampled. A single color stop produces a constant gradient in the
// interpolation space.
type Gradient struct {
	mixer Mixer

	stops []GradientStop
}

func newGradientWithOptions(opts InterpOptions, stops []GradientStop) Gradient {
	mixer := NewMixerWithOptions(opts)

	resolved := resolveStops(slices.Clone(stops), mixer)

	return Gradient{
		mixer: mixer,
		stops: resolved,
	}
}

func newGradient(stops []GradientStop) Gradient {
	return newGradientWithOptions(DefaultInterpOptions(), stops)
}

// NewGradientWithOptions creates a gradient with the
// specified interpolation options and color stops.
func NewGradientWithOptions(opts InterpOptions, stops ...GradientStop) Gradient {
	return newGradientWithOptions(opts, stops)
}

// NewGradient returns a new [Gradient] using [DefaultInterpOptions].
func NewGradient(stops ...GradientStop) Gradient {
	return newGradient(stops)
}

func NewHint(offset float64) GradientStop {
	return GradientStop{
		Color:  Color{},
		Offset: offset,
	}
}

func NewStop(color Color) GradientStop {
	return GradientStop{
		Color:  color,
		Offset: math.NaN(),
	}
}

func NewStopAt(color Color, offset float64) GradientStop {
	return GradientStop{
		Color:  color,
		Offset: offset,
	}
}

func NewStopsAt(color Color, offsets ...float64) []GradientStop {
	stops := make([]GradientStop, len(offsets))
	for i, offset := range offsets {
		stops[i] = GradientStop{Color: color, Offset: offset}
	}
	return stops
}

// Stops returns a copy of the gradient's color stops.
func (g Gradient) Stops() []GradientStop {
	return slices.Clone(g.stops)
}

// At returns the interpolated color at position t.
func (g Gradient) At(t float64) Color {
	return g.at(t)
}

// AtRepeat returns the interpolated color at position t using repeat spreading.
func (g Gradient) AtRepeat(t float64) Color {
	t = spreadRepeat(t)
	return g.at(t)
}

// AtReflect returns the interpolated color at position t using reflect spreading.
func (g Gradient) AtReflect(t float64) Color {
	t = spreadReflect(t)
	return g.at(t)
}

// AtSpread returns the interpolated color at position t using the specified spread method.
func (g Gradient) AtSpread(t float64, spread Spread) Color {
	switch spread {
	case SpreadRepeat:
		return g.AtRepeat(t)
	case SpreadReflect:
		return g.AtReflect(t)
	default:
		return g.At(t)
	}
}

func (g Gradient) at(t float64) Color {
	if len(g.stops) == 0 {
		return Color{}
	}

	i := g.findStop(t)

	if i == 0 {
		return g.stops[0].Color
	}

	if i == len(g.stops) {
		return g.stops[i-1].Color
	}

	a := &g.stops[i-1]
	b := &g.stops[i]

	seg := (t - a.Offset) * a.invRange

	if a.hint != nil {
		endWeight := 1.0
		if b.hint == a.hint {
			endWeight = b.weight
		}
		// A hint at an endpoint is a hard transition, leaving a constant
		// color on the remaining interval (even with longer hue interpolation).
		if a.weight == endWeight {
			if a.weight == 0 {
				return a.hint.left
			}
			if a.weight == 1 {
				return a.hint.right
			}
		}
		weight := a.weight + (endWeight-a.weight)*seg
		return g.mixer.UnsafeMix(a.hint.left, a.hint.right, weight)
	}

	return g.mixer.UnsafeMix(a.Color, b.Color, seg)
}

func (g Gradient) findStop(t float64) int {
	return sort.Search(len(g.stops), func(i int) bool {
		return g.stops[i].Offset > t
	})
}

func spreadRepeat(t float64) float64 {
	return t - math.Floor(t)
}

func spreadReflect(t float64) float64 {
	n := math.Floor(t)
	t -= n

	if int64(n)&1 != 0 {
		return 1 - t
	}

	return t
}
