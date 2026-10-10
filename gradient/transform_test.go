package gradient_test

import (
	"math"
	"testing"

	"github.com/thmalt/colors/gradient"
)

func assertNear(t *testing.T, name string, got, want float64) {
	t.Helper()
	if math.IsNaN(got) || math.IsInf(got, 0) || math.Abs(got-want) > 1e-12*math.Max(1, math.Abs(want)) {
		t.Errorf("%s = %g, want %g", name, got, want)
	}
}

func TestSkewMatchesDedicatedAxes(t *testing.T) {
	for _, composed := range []bool{false, true} {
		base := gradient.Identity()
		if composed {
			base.Translate(7, -11)
			base.Rotate(0.3)
			base.Scale(2, 0.5)
		}
		for _, angle := range []float64{0, math.Pi / 4, -0.3} {
			for _, axis := range []string{"x", "y"} {
				combined, single := base, base
				if axis == "x" {
					combined.Skew(angle, 0)
					single.SkewX(angle)
				} else {
					combined.Skew(0, angle)
					single.SkewY(angle)
				}
				for _, point := range [][2]float64{{0, 0}, {2, 3}, {-5, 4}} {
					gotX, gotY := combined.PositionAt(point[0], point[1])
					wantX, wantY := single.PositionAt(point[0], point[1])
					assertNear(t, axis+" skew x", gotX, wantX)
					assertNear(t, axis+" skew y", gotY, wantY)
				}
			}
		}
	}
}

func TestSkewComposesBothAxes(t *testing.T) {
	base := gradient.Identity()
	base.Translate(4, -8)
	base.Rotate(0.2)
	base.Scale(2, 3)
	combined := base
	const angleX, angleY = 0.3, -0.4
	combined.Skew(angleX, angleY)
	inverse, ok := combined.Inverse()
	if !ok {
		t.Fatal("combined skew must be invertible")
	}
	for _, point := range [][2]float64{{0, 0}, {2, 3}, {-5, 4}} {
		x, y := point[0], point[1]
		// Skew applies both shears simultaneously, before the existing transform.
		wantX, wantY := base.PositionAt(x+math.Tan(angleX)*y, y+math.Tan(angleY)*x)
		gotX, gotY := combined.PositionAt(x, y)
		assertNear(t, "combined skew x", gotX, wantX)
		assertNear(t, "combined skew y", gotY, wantY)
		backX, backY := inverse.PositionAt(gotX, gotY)
		assertNear(t, "inverse x", backX, x)
		assertNear(t, "inverse y", backY, y)
	}
}
