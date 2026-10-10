package gradient_test

import (
	"math"
	"testing"

	"github.com/thmalt/colors/gradient"
)

func TestConicTransformPreservesLocalPosition(t *testing.T) {
	for _, tc := range []struct {
		name  string
		apply func(*gradient.Transform)
	}{
		{"identity", func(*gradient.Transform) {}},
		{"translation", func(m *gradient.Transform) { m.Translate(23, -17) }},
		{"uniform scale", func(m *gradient.Transform) { m.Scale(2, 2) }},
		{"nonuniform scale", func(m *gradient.Transform) { m.Scale(3, 0.5) }},
		{"rotation", func(m *gradient.Transform) { m.Rotate(0.7) }},
		{"reflection", func(m *gradient.Transform) { m.Scale(-1, 2) }},
		{"skew", func(m *gradient.Transform) { m.Skew(0.2, -0.3) }},
		{"composite", func(m *gradient.Transform) {
			m.Translate(-31, 19)
			m.Rotate(-0.8)
			m.Scale(2, 0.5)
			m.SkewX(0.4)
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			spec := gradient.NewConicSpec()
			spec.SetSize(200, 100)
			spec.SetCenter(0.25, 0.7)
			spec.SetStartAngle(-0.17)
			base, err := spec.Build()
			if err != nil {
				t.Fatal(err)
			}
			tc.apply(&spec.Transform)
			transformed, err := spec.Build()
			if err != nil {
				t.Fatal(err)
			}
			for _, p := range [][2]float64{{50, 0}, {100, 70}, {50, 120}, {0, 70}, {10, 10}, {135, 110}} {
				x, y := spec.Transform.PositionAt(p[0], p[1])
				got, want := transformed.PositionAt(x, y), base.PositionAt(p[0], p[1])
				// Turns wrap at one, so compare the shortest angular difference.
				assertNear(t, "transformed turn difference", math.Remainder(got-want, 1), 0)
			}
		})
	}
}

func TestConicRejectsSingularTransform(t *testing.T) {
	spec := gradient.NewConicSpec()
	spec.SetSize(100, 100)
	spec.Scale(0, 1)
	if _, err := spec.Build(); err != gradient.ErrInvalidTransform {
		t.Fatalf("Build error = %v, want %v", err, gradient.ErrInvalidTransform)
	}
}
