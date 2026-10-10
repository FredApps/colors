package colors_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/thmalt/colors"
	"github.com/thmalt/colors/space"
)

func TestHslNeutralMixIndependentOfSourceSpace(t *testing.T) {
	for _, via := range []space.Space{space.Srgb, space.LabD50, space.LuvD50, space.Oklab} {
		for _, gray := range []float64{0, 0.5, 1} {
			for _, premultiplied := range []bool{false, true} {
				t.Run(fmt.Sprintf("%v/gray_%g/premultiplied_%t", via, gray, premultiplied), func(t *testing.T) {
					input := colors.SrgbAlpha(gray, gray, gray, 0.25).To(via)
					r, g, b := input.Srgb()
					for channel, value := range []float64{r, g, b} {
						assertHueChannelNear(t, fmt.Sprintf("source sRGB channel %d", channel), value, gray)
					}
					other := colors.SrgbAlpha(0.5, 0.5, 0.5, 0.25)
					opts := colors.InterpOptions{Space: space.Hsl, Premultiplied: premultiplied}
					for _, fraction := range []float64{0, 0.25, 0.5, 0.75, 1} {
						mixed := colors.MixWith(input, other, fraction, opts)
						r, g, b := mixed.Srgb()
						want := gray + (0.5-gray)*fraction
						for channel, value := range []float64{r, g, b} {
							assertHueChannelNear(t, fmt.Sprintf("mix %g sRGB channel %d", fraction, channel), value, want)
						}
						if mixed.Alpha() != 0.25 {
							t.Errorf("mix alpha = %g; want 0.25", mixed.Alpha())
						}
					}
				})
			}
		}
	}
}

func TestHwbToHsvNormalizesAchromaticColors(t *testing.T) {
	for _, tc := range []struct {
		name string
		w, b float64
	}{
		{"equal", 0.8, 0.8},
		{"blackness_one", 0.8, 1},
		{"unequal", 0.6, 0.9},
		{"boundary", 0.25, 0.75},
		{"black", 0, 1},
		{"white", 1, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			original := colors.HwbAlpha(120, tc.w, tc.b, 0.25)
			converted := original.To(space.Hsv)
			h, s, v := converted.Hsv()
			wantGray := tc.w / (tc.w + tc.b)
			assertHueChannelNear(t, "hue", h, 120)
			assertHueChannelNear(t, "saturation", s, 0)
			assertHueChannelNear(t, "value", v, wantGray)
			r, g, b := converted.Srgb()
			for channel, got := range []float64{r, g, b} {
				assertHueChannelNear(t, fmt.Sprintf("sRGB channel %d", channel), got, wantGray)
			}
			if converted.Alpha() != original.Alpha() {
				t.Errorf("alpha = %g; want %g", converted.Alpha(), original.Alpha())
			}
		})
	}
}

func assertHueChannelNear(t *testing.T, label string, got, want float64) {
	t.Helper()
	if math.IsNaN(got) || math.IsInf(got, 0) || math.Abs(got-want) > 1e-12 {
		t.Errorf("%s = %.17g; want %.17g", label, got, want)
	}
}
