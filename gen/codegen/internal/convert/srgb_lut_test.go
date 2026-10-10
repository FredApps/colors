package convert

import (
	"math"
	"testing"
)

func TestSrgbLUT(t *testing.T) {
	for _, tc := range []struct {
		name   string
		max    int
		decode func(int) float64
		encode func(float64) int
	}{
		{
			name:   "8-bit",
			max:    1<<8 - 1,
			decode: func(x int) float64 { return Srgb8Decode(uint8(x)) },
			encode: func(x float64) int { return int(Srgb8Encode(x)) },
		},
		{
			name:   "16-bit",
			max:    1<<16 - 1,
			decode: func(x int) float64 { return Srgb16Decode(uint16(x)) },
			encode: func(x float64) int { return int(Srgb16Encode(x)) },
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for i := 0; i <= tc.max; i++ {
				got := tc.decode(i)
				want := SrgbDecode(float64(i) / float64(tc.max))
				if got != want {
					t.Fatalf("decode(%d) = %g, want %g", i, got, want)
				}
				if encoded := tc.encode(got); encoded != i {
					t.Fatalf("encode(decode(%d)) = %d", i, encoded)
				}
			}

			for _, x := range []float64{math.Inf(-1), -1, 0} {
				if got := tc.encode(x); got != 0 {
					t.Errorf("encode(%g) = %d, want 0", x, got)
				}
			}
			for _, x := range []float64{1, 2, math.Inf(1)} {
				if got := tc.encode(x); got != tc.max {
					t.Errorf("encode(%g) = %d, want %d", x, got, tc.max)
				}
			}
		})
	}
}
