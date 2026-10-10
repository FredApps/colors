package colors_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/thmalt/colors"
	"github.com/thmalt/colors/space"
)

func TestGradientEmpty(t *testing.T) {
	opts := colors.InterpOptions{Space: space.Hsl, Hue: colors.HueLonger}
	for _, tc := range []struct {
		name  string
		build func() colors.Gradient
	}{
		{"zero value", func() colors.Gradient { return colors.Gradient{} }},
		{"constructor", func() colors.Gradient { return colors.NewGradient() }},
		{"constructor options", func() colors.Gradient { return colors.NewGradientWithOptions(opts) }},
		{"builder", func() colors.Gradient { return colors.NewGradientBuilder().Build() }},
		{"builder options", func() colors.Gradient { return colors.NewGradientBuilder().BuildWithOptions(opts) }},
		{"reset builder", func() colors.Gradient {
			return colors.NewGradientBuilder().AddStop(colors.Srgb(1, 0, 0)).Reset().Build()
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("empty gradient panicked: %v", r)
				}
			}()
			g := tc.build()
			if len(g.Stops()) != 0 {
				t.Fatal("empty gradient has stops")
			}
			for _, sample := range []float64{-2.5, 0, 0.5, 1, 2.5, math.NaN()} {
				for name, at := range gradientSamplers(g) {
					if got := at(sample); got != (colors.Color{}) {
						t.Errorf("%s(%g) = %v, want zero Color", name, sample, got)
					}
				}
			}
		})
	}
}

func TestGradientSingleStop(t *testing.T) {
	color := colors.SrgbAlpha(0.2, 0.4, 0.6, 0.3)
	opts := colors.InterpOptions{Space: space.Oklch, Premultiplied: true, Hue: colors.HueLonger}
	for _, offset := range []float64{math.NaN(), -0.5, 0, 0.3, 2} {
		t.Run(fmt.Sprint(offset), func(t *testing.T) {
			g := colors.NewGradientWithOptions(opts, colors.NewStopAt(color, offset))
			stops := g.Stops()
			if len(stops) != 1 {
				t.Fatalf("single-stop gradient has %d stops, want 1", len(stops))
			}
			if math.IsNaN(offset) {
				offset = 0
			}
			if stops[0].Offset != offset {
				t.Fatalf("stop offset = %g, want %g", stops[0].Offset, offset)
			}
			for _, sample := range []float64{-5, -0.5, 0, 0.3, 1, 5, math.NaN()} {
				for name, at := range gradientSamplers(g) {
					assertGradientColor(t, fmt.Sprintf("%s(%g)", name, sample), at(sample), color.To(opts.Space))
				}
			}
		})
	}
}

func TestGradientHintHuePaths(t *testing.T) {
	for _, hueMode := range []colors.HueInterpolation{colors.HueShorter, colors.HueLonger, colors.HueIncreasing, colors.HueDecreasing} {
		for _, hues := range [][2]float64{{0, 120}, {120, 0}, {350, 10}, {10, 350}, {40, 40}, {0, 180}, {180, 0}} {
			for _, hint := range []float64{0.001, 0.25, 0.5, 0.75, 0.999} {
				t.Run(fmt.Sprintf("%s/%g-%g/hint%g", hueMode, hues[0], hues[1], hint), func(t *testing.T) {
					opts := colors.InterpOptions{Space: space.Hsl, Hue: hueMode}
					left, right := colors.Hsl(hues[0], 1, 0.5), colors.Hsl(hues[1], 1, 0.5)
					checkHintedGradient(t, opts, left, right, hint)
				})
			}
		}
	}
}

func TestGradientHintLongerArcRegression(t *testing.T) {
	opts := colors.InterpOptions{Space: space.Hsl, Hue: colors.HueLonger}
	g := colors.NewGradientWithOptions(opts,
		colors.NewStopAt(colors.Hsl(0, 1, 0.5), 0), colors.NewHint(0.25),
		colors.NewStopAt(colors.Hsl(120, 1, 0.5), 1))
	// The sample is halfway between the synthetic stops at 1/12 and 1/6.
	// Their scalar weights follow sqrt(t); hue follows the original 240-degree
	// decreasing arc once, without adding another long arc between those stops.
	want := 360 - 240*(math.Sqrt(1.0/12)+math.Sqrt(1.0/6))/2
	hue, _, _ := g.At(0.125).Hsl()
	if math.Abs(hue-want) > 1e-10 {
		t.Fatalf("hinted longer-arc hue = %g, want %g", hue, want)
	}
}

func TestGradientHintChannelsAndAlpha(t *testing.T) {
	for _, dst := range []space.Space{space.Srgb, space.Hsl, space.Hsv, space.Hwb, space.Oklch, space.Lch} {
		for _, premultiplied := range []bool{false, true} {
			for _, alpha := range [][2]float64{{1, 1}, {0.2, 0.8}, {0, 0.8}, {0.2, 0}, {0, 0}} {
				t.Run(fmt.Sprintf("%s/premult%t/alpha%g-%g", dst, premultiplied, alpha[0], alpha[1]), func(t *testing.T) {
					opts := colors.InterpOptions{Space: dst, Hue: colors.HueLonger, Premultiplied: premultiplied}
					left := colors.HslAlpha(20, 0.3, 0.2, alpha[0]).To(dst)
					right := colors.HslAlpha(100, 0.9, 0.7, alpha[1]).To(dst)
					checkHintedGradient(t, opts, left, right, 0.25)
					checkHintedGradient(t, opts, left, right, 0.75)
				})
			}
		}
	}
}

func TestGradientHintMissingComponentsMatchMixer(t *testing.T) {
	// Missing components currently propagate as NaN through the mixer. Hints
	// must preserve that behavior without contaminating otherwise finite channels.
	for _, dst := range []space.Space{space.Srgb, space.Hsl, space.Oklch} {
		for _, premultiplied := range []bool{false, true} {
			for channel := 0; channel < 4; channel++ {
				for _, side := range []string{"left", "right", "both"} {
					t.Run(fmt.Sprintf("%s/premult%t/channel%d/%s", dst, premultiplied, channel, side), func(t *testing.T) {
						a, b := []float64{0.2, 0.3, 0.4, 0.5}, []float64{0.8, 0.7, 0.6, 0.9}
						if side != "right" {
							a[channel] = math.NaN()
						}
						if side != "left" {
							b[channel] = math.NaN()
						}
						opts := colors.InterpOptions{Space: dst, Hue: colors.HueLonger, Premultiplied: premultiplied}
						checkHintedGradient(t, opts, colors.New(dst, a...), colors.New(dst, b...), 0.25)
					})
				}
			}
		}
	}
}

func TestGradientHintAtEndpoint(t *testing.T) {
	left, right := colors.HslAlpha(0, 0.3, 0.4, 0), colors.HslAlpha(120, 0.7, 0.6, 0.8)
	for _, hint := range []float64{0, 1} {
		for _, premultiplied := range []bool{false, true} {
			opts := colors.InterpOptions{Space: space.Hsl, Hue: colors.HueLonger, Premultiplied: premultiplied}
			g := colors.NewGradientWithOptions(opts, colors.NewStopAt(left, 0), colors.NewHint(hint), colors.NewStopAt(right, 1))
			want := right
			if hint == 1 {
				want = left
			}
			for _, sample := range []float64{0, 0.1, 0.5, 0.9} {
				assertGradientColor(t, fmt.Sprintf("hint %g, at %g, premult %t", hint, sample, premultiplied), g.At(sample), want)
			}
			assertGradientColor(t, "last stop", g.At(1), right)
			assertGradientColor(t, "before first stop", g.At(-1), left)
		}
	}
}

func TestGradientHintDoesNotChangeNeighboringSegments(t *testing.T) {
	opts := colors.InterpOptions{Space: space.Hsl, Hue: colors.HueLonger, Premultiplied: true}
	a, b, c, d := colors.Hsl(10, 1, 0.5), colors.Hsl(50, 1, 0.5), colors.Hsl(120, 1, 0.5), colors.Hsl(250, 1, 0.5)
	g := colors.NewGradientWithOptions(opts,
		colors.NewStopAt(a, 0), colors.NewHint(0.1), colors.NewStopAt(b, 0.4),
		colors.NewHint(0.7), colors.NewStopAt(c, 0.8), colors.NewStopAt(d, 1))
	first := colors.NewGradientWithOptions(opts, colors.NewStopAt(a, 0), colors.NewHint(0.25), colors.NewStopAt(b, 1))
	second := colors.NewGradientWithOptions(opts, colors.NewStopAt(b, 0), colors.NewHint(0.75), colors.NewStopAt(c, 1))
	for _, sample := range []float64{0.05, 0.125, 0.4, 0.6, 0.9} {
		assertGradientColor(t, "first hint", g.At(sample*0.4), first.At(sample))
		assertGradientColor(t, "second hint", g.At(0.4+sample*0.4), second.At(sample))
		assertGradientColor(t, "unhinted segment", g.At(0.8+sample*0.2), colors.MixWith(c, d, sample, opts))
	}
}

func TestGradientStopsCanBeReused(t *testing.T) {
	opts := colors.InterpOptions{Space: space.Hsl, Hue: colors.HueLonger}
	g := colors.NewGradientWithOptions(opts,
		colors.NewStopAt(colors.Hsl(0, 1, 0.5), 0), colors.NewHint(0.25), colors.NewStopAt(colors.Hsl(120, 1, 0.5), 1))
	stops := g.Stops()
	publicStops := make([]colors.GradientStop, len(stops))
	for i, stop := range stops {
		publicStops[i] = colors.NewStopAt(stop.Color, stop.Offset)
	}
	otherOpts := colors.InterpOptions{Space: space.Oklab, Premultiplied: true}
	reused, fresh := colors.NewGradientWithOptions(otherOpts, stops...), colors.NewGradientWithOptions(otherOpts, publicStops...)
	for _, sample := range []float64{0.01, 0.125, 0.5, 0.9} {
		assertGradientColor(t, "reused resolved stops", reused.At(sample), fresh.At(sample))
	}
	before := g.At(0.125)
	stops[0].Color = colors.Srgb(1, 1, 1)
	assertGradientColor(t, "Stops returns a copy", g.At(0.125), before)
}

func checkHintedGradient(t *testing.T, opts colors.InterpOptions, left, right colors.Color, hint float64) {
	t.Helper()
	g := colors.NewGradientWithOptions(opts, colors.NewStopAt(left, 0), colors.NewHint(hint), colors.NewStopAt(right, 1))
	stops := g.Stops()
	wantStops := 11
	if hint == 0.5 {
		wantStops = 2
	}
	if len(stops) != wantStops {
		t.Fatalf("hint %g: got %d resolved stops, want %d", hint, len(stops), wantStops)
	}
	exponent := math.Log(0.5) / math.Log(hint)
	// Retain the existing piecewise-linear approximation of the hint's power
	// curve, but apply the chosen color interpolation once to the original pair.
	for i := 1; i < len(stops); i++ {
		a, b := stops[i-1], stops[i]
		wa, wb := math.Pow(a.Offset, exponent), math.Pow(b.Offset, exponent)
		for _, fraction := range []float64{0.1, 0.5, 0.9} {
			sample := a.Offset + (b.Offset-a.Offset)*fraction
			weight := wa + (wb-wa)*fraction
			want := colors.MixWith(left, right, weight, opts)
			assertGradientColor(t, fmt.Sprintf("hint %g at %g", hint, sample), g.At(sample), want)
		}
	}
}

func gradientSamplers(g colors.Gradient) map[string]func(float64) colors.Color {
	return map[string]func(float64) colors.Color{
		"At":               g.At,
		"AtRepeat":         g.AtRepeat,
		"AtReflect":        g.AtReflect,
		"AtSpread pad":     func(t float64) colors.Color { return g.AtSpread(t, colors.SpreadPad) },
		"AtSpread repeat":  func(t float64) colors.Color { return g.AtSpread(t, colors.SpreadRepeat) },
		"AtSpread reflect": func(t float64) colors.Color { return g.AtSpread(t, colors.SpreadReflect) },
	}
}

func assertGradientColor(t *testing.T, context string, got, want colors.Color) {
	t.Helper()
	if got.Space() != want.Space() {
		t.Fatalf("%s: space %v, want %v", context, got.Space(), want.Space())
	}
	gc, wc := append(got.Channels(), got.Alpha()), append(want.Channels(), want.Alpha())
	for i := range gc {
		if math.IsNaN(gc[i]) && math.IsNaN(wc[i]) {
			continue
		}
		if math.IsNaN(gc[i]) || math.IsNaN(wc[i]) || math.Abs(gc[i]-wc[i]) > 1e-9 {
			t.Fatalf("%s: channel %d = %.15g, want %.15g", context, i, gc[i], wc[i])
		}
	}
}
