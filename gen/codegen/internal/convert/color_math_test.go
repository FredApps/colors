package convert

import (
	"testing"

	public "github.com/thmalt/colors/convert"
)

// The generator's handwritten conversion copies must agree with the runtime
// helpers; changes to one must not silently leave the other behind.
func TestHueConversionsMatchRuntime(t *testing.T) {
	for _, rgb := range [][3]float64{
		{0, 0, 0}, {0.5, 0.5, 0.5}, {1, 1, 1},
		{1, 1 - 0x1p-48, 1}, {1, 1 - 0x1p-47, 1}, {1, 1 - 0x1p-46, 1},
		{1e-300, 0, 0}, {0.25, 0.5, 0.75}, {-2, -2 + 0x1p-47, -2},
	} {
		h, s, l := SrgbToHsl(rgb[0], rgb[1], rgb[2])
		wh, ws, wl := public.SrgbToHsl(rgb[0], rgb[1], rgb[2])
		if h != wh || s != ws || l != wl {
			t.Errorf("SrgbToHsl(%v) = (%g, %g, %g); runtime = (%g, %g, %g)", rgb, h, s, l, wh, ws, wl)
		}
	}
	for _, hwb := range [][3]float64{
		{120, 0.8, 0.8}, {120, 0.6, 0.9}, {120, 0.8, 1},
		{120, 0.25, 0.75}, {120, 0, 1}, {120, 1, 0}, {120, 0.25, 0.5},
	} {
		h, s, v := HwbToHsv(hwb[0], hwb[1], hwb[2])
		wh, ws, wv := public.HwbToHsv(hwb[0], hwb[1], hwb[2])
		if h != wh || s != ws || v != wv {
			t.Errorf("HwbToHsv(%v) = (%g, %g, %g); runtime = (%g, %g, %g)", hwb, h, s, v, wh, ws, wv)
		}
	}
}
