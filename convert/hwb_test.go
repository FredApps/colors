package convert

import (
	"math"
	"testing"
)

func TestHwbToHsv(t *testing.T) {
	for _, tc := range []struct {
		name       string
		w, b, s, v float64
	}{
		{"equal_normalized", 0.8, 0.8, 0, 0.5},
		{"unequal_normalized", 0.6, 0.9, 0, 0.4},
		{"blackness_one_normalized", 0.8, 1, 0, 0.8 / 1.8},
		{"boundary", 0.25, 0.75, 0, 0.25},
		{"just_above_boundary", 0.5, 0.5 + 0x1p-51, 0, 0.5 / (1 + 0x1p-51)},
		{"just_below_boundary", 0.5, 0.5 - 0x1p-51, 1 - 0.5/(0.5+0x1p-51), 0.5 + 0x1p-51},
		{"black", 0, 1, 0, 0},
		{"white", 1, 0, 0, 1},
		{"saturated", 0, 0, 1, 1},
		{"chromatic", 0.25, 0.5, 0.5, 0.5},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, hue := range []float64{-30, 0, 120, 270, 390} {
				h, s, v := HwbToHsv(hue, tc.w, tc.b)
				if h != hue || math.IsNaN(s) || math.IsNaN(v) ||
					math.Abs(s-tc.s) > 2e-15*math.Abs(tc.s) || math.Abs(v-tc.v) > 2e-15*math.Abs(tc.v) {
					t.Errorf("HSV = (%g, %.17g, %.17g); want (%g, %.17g, %.17g)", h, s, v, hue, tc.s, tc.v)
				}
				r, g, b := HsvToSrgb(h, s, v)
				wr, wg, wb := HwbToSrgb(hue, tc.w, tc.b)
				got, want := [3]float64{r, g, b}, [3]float64{wr, wg, wb}
				for i := range got {
					if math.IsNaN(got[i]) || math.IsInf(got[i], 0) || math.Abs(got[i]-want[i]) > 1e-15 {
						t.Errorf("sRGB channel %d via HSV = %g; direct HWB conversion = %g", i, got[i], want[i])
					}
				}
			}
		})
	}
}
