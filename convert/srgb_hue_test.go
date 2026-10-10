package convert

import (
	"fmt"
	"math"
	"testing"
)

func TestSrgbToHslNumericallyNeutral(t *testing.T) {
	// 32 binary64 machine epsilons, relative to the largest channel magnitude.
	const tolerance = 0x1p-47
	for _, scale := range []float64{0x1p-500, 0.125, 0.5, 1, 2, 0x1p500} {
		for _, sign := range []float64{-1, 1} {
			for _, fraction := range []float64{0, 0.5, 1} {
				t.Run(fmt.Sprintf("scale_%g/sign_%g/fraction_%g", scale, sign, fraction), func(t *testing.T) {
					x := sign * scale
					y := sign * (scale - scale*tolerance*fraction)
					for _, rgb := range [][3]float64{{x, y, x}, {y, x, x}, {x, x, y}} {
						h, s, l := SrgbToHsl(rgb[0], rgb[1], rgb[2])
						wantL := (x + y) * 0.5
						if h != 0 || s != 0 || l != wantL {
							t.Errorf("SrgbToHsl(%v) = (%g, %g, %g); want (0, 0, %g)", rgb, h, s, l, wantL)
						}
					}
				})
			}
		}
	}
	if h, s, l := SrgbToHsl(0, 0, 0); h != 0 || s != 0 || l != 0 {
		t.Errorf("black = (%g, %g, %g); want (0, 0, 0)", h, s, l)
	}
}

func TestSrgbToHslPreservesChromaticColors(t *testing.T) {
	for _, tc := range []struct {
		name string
		rgb  [3]float64
		hsl  [3]float64
	}{
		{"red", [3]float64{1, 0, 0}, [3]float64{0, 1, 0.5}},
		{"green", [3]float64{0, 1, 0}, [3]float64{120, 1, 0.5}},
		{"blue", [3]float64{0, 0, 1}, [3]float64{240, 1, 0.5}},
		{"ordinary", [3]float64{0.25, 0.5, 0.75}, [3]float64{210, 0.5, 0.5}},
		{"above_neutral_threshold", [3]float64{1, 1 - 0x1p-46, 1}, [3]float64{300, 1, 1 - 0x1p-47}},
		{"dark_saturated", [3]float64{1e-300, 0, 0}, [3]float64{0, 1, 0.5e-300}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			h, s, l := SrgbToHsl(tc.rgb[0], tc.rgb[1], tc.rgb[2])
			got := [3]float64{h, s, l}
			for i, want := range tc.hsl {
				if math.IsNaN(got[i]) || math.IsInf(got[i], 0) || got[i] != want {
					t.Errorf("HSL channel %d = %.17g; want %.17g", i, got[i], want)
				}
			}
		})
	}
}
