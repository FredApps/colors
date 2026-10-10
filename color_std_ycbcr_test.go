package colors_test

import (
	"fmt"
	"image/color"
	"math"
	"testing"

	"github.com/thmalt/colors"
)

type ycbcrRGBAOnly struct{ color.Color }

func TestYCbCrImportsHonorStandardClipping(t *testing.T) {
	values := []uint8{0, 1, 64, 127, 128, 192, 254, 255}
	for _, y := range values {
		for _, cb := range values {
			for _, cr := range values {
				src := color.YCbCr{Y: y, Cb: cb, Cr: cr}
				r, g, b, _ := src.RGBA()
				want := [3]float64{float64(r) / 65535, float64(g) / 65535, float64(b) / 65535}
				for _, tc := range []struct {
					name string
					got  colors.Color
				}{
					{"value", colors.FromStd(src)}, {"pointer", colors.FromStd(&src)}, {"direct", colors.StdYCbCr(src)},
				} {
					gotR, gotG, gotB := tc.got.Srgb()
					for channel, got := range [3]float64{gotR, gotG, gotB} {
						if math.Abs(got-want[channel]) > 1e-14 {
							t.Fatalf("%s YCbCr(%d,%d,%d) channel %d = %g; want %g", tc.name, y, cb, cr, channel, got, want[channel])
						}
					}
					if tc.got.Alpha() != 1 {
						t.Fatalf("%s alpha = %g; want 1", tc.name, tc.got.Alpha())
					}
				}
			}
		}
	}
}

func TestNYCbCrAImportsPreserveClippedStraightColorAndAlpha(t *testing.T) {
	for _, base := range []color.YCbCr{{Y: 0, Cb: 255, Cr: 255}, {Y: 255, Cb: 0, Cr: 0}, {Y: 127, Cb: 127, Cr: 127}, {Y: 64, Cb: 192, Cr: 128}} {
		r, g, b, _ := base.RGBA()
		want := [3]float64{float64(r) / 65535, float64(g) / 65535, float64(b) / 65535}
		for _, alpha := range []uint8{0, 1, 64, 128, 254, 255} {
			src := color.NYCbCrA{YCbCr: base, A: alpha}
			for _, tc := range []struct {
				name string
				got  colors.Color
			}{
				{"value", colors.FromStd(src)}, {"pointer", colors.FromStd(&src)}, {"direct", colors.StdNYCbCrA(src)},
			} {
				t.Run(fmt.Sprintf("%d_%d_%d_%d_%s", base.Y, base.Cb, base.Cr, alpha, tc.name), func(t *testing.T) {
					gr, gg, gb := tc.got.Srgb()
					for channel, got := range [3]float64{gr, gg, gb} {
						if math.Abs(got-want[channel]) > 1e-14 {
							t.Errorf("channel %d = %g; want %g", channel, got, want[channel])
						}
					}
					if math.Abs(tc.got.Alpha()-float64(alpha)/255) > 1e-15 {
						t.Errorf("alpha = %g; want %g", tc.got.Alpha(), float64(alpha)/255)
					}
					if alpha != 0 {
						wrapped := colors.FromStd(ycbcrRGBAOnly{src})
						wr, wg, wb := wrapped.Srgb()
						// The generic path unpremultiplies 16-bit components, which can lose
						// less than one premultiplied unit compared with the straight fast path.
						tolerance := 1/(float64(alpha)*257) + 1e-14
						for channel, got := range [3]float64{gr, gg, gb} {
							if math.Abs(got-([3]float64{wr, wg, wb})[channel]) > tolerance {
								t.Errorf("wrapped channel %d differs by more than premultiplication rounding", channel)
							}
						}
					}
				})
			}
		}
	}
}
