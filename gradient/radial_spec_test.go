package gradient_test

import (
	"math"
	"testing"

	"github.com/thmalt/colors/gradient"
)

func checkRadialRadii(t *testing.T, s *gradient.RadialSpec, wantX, wantY float64) {
	t.Helper()
	rx, ry := s.RadiusPixel()
	assertNear(t, "pixel radius x", rx, wantX)
	assertNear(t, "pixel radius y", ry, wantY)
	r, err := s.Build()
	if err != nil {
		t.Fatal(err)
	}
	w, h := s.Size()
	cx, cy := s.Center()
	cx, cy = cx*w, cy*h
	for _, p := range [][2]float64{{cx + wantX, cy}, {cx - wantX, cy}, {cx, cy + wantY}, {cx, cy - wantY}} {
		x, y := s.Transform.PositionAt(p[0], p[1])
		assertNear(t, "position at radius", r.PositionAt(x, y), 1)
	}
	x, y := s.Transform.PositionAt(cx, cy)
	assertNear(t, "position at center", r.PositionAt(x, y), 0)
}

func TestRadialExplicitRadiiFollowUnitsAndResize(t *testing.T) {
	for _, tc := range []struct {
		name             string
		set              func(*gradient.RadialSpec)
		initial, resized [2]float64
	}{
		{"normalized ellipse", func(s *gradient.RadialSpec) { s.SetRadius(0.5, 1.5) }, [2]float64{50, 75}, [2]float64{20, 180}},
		{"normalized circle", func(s *gradient.RadialSpec) { s.SetCircleRadius(0.75) }, [2]float64{37.5, 37.5}, [2]float64{30, 30}},
		{"pixel ellipse", func(s *gradient.RadialSpec) { s.SetRadiusPixel(30, 45) }, [2]float64{30, 45}, [2]float64{30, 45}},
		{"pixel circle", func(s *gradient.RadialSpec) { s.SetCircleRadiusPixel(25) }, [2]float64{25, 25}, [2]float64{25, 25}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, beforeSize := range []bool{true, false} {
				s := gradient.NewRadialSpec()
				if beforeSize {
					tc.set(&s)
				}
				s.SetSize(200, 100)
				if !beforeSize {
					tc.set(&s)
				}
				checkRadialRadii(t, &s, tc.initial[0], tc.initial[1])
				s.SetSize(80, 240)
				checkRadialRadii(t, &s, tc.resized[0], tc.resized[1])
				// Pixel radii describe local geometry, before the affine transform.
				s.Translate(13, -27)
				s.Rotate(0.4)
				s.Scale(2, 0.5)
				checkRadialRadii(t, &s, tc.resized[0], tc.resized[1])
			}
		})
	}
}

func TestRadialRadiusSettersReplaceUnits(t *testing.T) {
	s := gradient.NewRadialSpec()
	s.SetSize(200, 100)
	s.SetCircleRadiusPixel(25)
	s.SetCircleRadius(1)
	x, y := s.Radius()
	assertNear(t, "configured normalized radius x", x, 1)
	assertNear(t, "configured normalized radius y", y, 1)
	checkRadialRadii(t, &s, 50, 50)
	s.SetCircleRadiusPixel(30)
	s.SetSize(400, 200)
	checkRadialRadii(t, &s, 30, 30)
	s.SetRadius(0.5, 1.5)
	checkRadialRadii(t, &s, 100, 150)
	s.SetRadiusPixel(10, 20)
	checkRadialRadii(t, &s, 10, 20)
}

func TestRadialEllipseSizeModes(t *testing.T) {
	for _, tc := range []struct {
		name                      string
		center, closest, farthest [2]float64
	}{
		{"centered", [2]float64{0.5, 0.5}, [2]float64{100, 50}, [2]float64{100, 50}},
		{"off center", [2]float64{0.25, 0.3}, [2]float64{50, 30}, [2]float64{150, 70}},
		{"opposite quadrant", [2]float64{0.75, 0.8}, [2]float64{50, 20}, [2]float64{150, 80}},
		{"outside box", [2]float64{-0.25, 1.3}, [2]float64{50, 30}, [2]float64{250, 130}},
	} {
		for _, mode := range []struct {
			name            string
			value           gradient.RadialSizeMode
			closest, corner bool
		}{
			{"farthest corner", gradient.RadialFarthestCorner, false, true},
			{"closest corner", gradient.RadialClosestCorner, true, true},
			{"farthest side", gradient.RadialFarthestSide, false, false},
			{"closest side", gradient.RadialClosestSide, true, false},
		} {
			t.Run(tc.name+"/"+mode.name, func(t *testing.T) {
				s := gradient.NewRadialSpec()
				s.SetSize(200, 100)
				s.SetCenter(tc.center[0], tc.center[1])
				s.SetSizeMode(mode.value)
				radii := tc.farthest
				if mode.closest {
					radii = tc.closest
				}
				if mode.corner {
					radii[0] *= math.Sqrt2
					radii[1] *= math.Sqrt2
				}
				checkRadialRadii(t, &s, radii[0], radii[1])
				if mode.corner {
					r, err := s.Build()
					if err != nil {
						t.Fatal(err)
					}
					// Exactly the selected nearest/farthest corner lies on the ellipse.
					position := math.Inf(1)
					if !mode.closest {
						position = 0
					}
					for _, p := range [][2]float64{{0, 0}, {200, 0}, {0, 100}, {200, 100}} {
						v := r.PositionAt(p[0], p[1])
						if mode.closest {
							position = math.Min(position, v)
						} else {
							position = math.Max(position, v)
						}
					}
					assertNear(t, "position at selected corner", position, 1)
				}
			})
		}
	}
}

func TestRadialEllipseEdges(t *testing.T) {
	for _, center := range [][2]float64{{0, 0.5}, {1, 0.5}, {0.5, 0}, {0.5, 1}, {0, 0}, {1, 1}} {
		for _, mode := range []gradient.RadialSizeMode{gradient.RadialClosestSide, gradient.RadialClosestCorner} {
			s := gradient.NewRadialSpec()
			s.SetSize(200, 100)
			s.SetCenter(center[0], center[1])
			s.SetSizeMode(mode)
			rx, ry := s.RadiusPixel()
			if math.IsNaN(rx) || math.IsNaN(ry) || math.IsInf(rx, 0) || math.IsInf(ry, 0) {
				t.Errorf("center %v mode %v: nonfinite radii (%g, %g)", center, mode, rx, ry)
			}
			if _, err := s.Build(); err != gradient.ErrInvalidRadius {
				t.Errorf("center %v mode %v: Build error = %v, want %v", center, mode, err, gradient.ErrInvalidRadius)
			}
		}
		s := gradient.NewRadialSpec()
		s.SetSize(200, 100)
		s.SetCenter(center[0], center[1])
		rx := math.Max(center[0], 1-center[0]) * 200 * math.Sqrt2
		ry := math.Max(center[1], 1-center[1]) * 100 * math.Sqrt2
		checkRadialRadii(t, &s, rx, ry)
	}
}

func TestRadialCircleSizeModes(t *testing.T) {
	for _, tc := range []struct {
		mode   gradient.RadialSizeMode
		radius float64
	}{
		{gradient.RadialFarthestCorner, math.Hypot(150, 70)},
		{gradient.RadialClosestCorner, math.Hypot(50, 30)},
		{gradient.RadialFarthestSide, 150},
		{gradient.RadialClosestSide, 30},
	} {
		s := gradient.NewRadialSpec()
		s.SetSize(200, 100)
		s.SetCenter(0.25, 0.3)
		s.SetShape(gradient.RadialCircle)
		s.SetSizeMode(tc.mode)
		checkRadialRadii(t, &s, tc.radius, tc.radius)
	}
}

func TestRadialDegenerateSizeAndRadius(t *testing.T) {
	t.Run("unset size", func(t *testing.T) {
		s := gradient.NewRadialSpec()
		s.SetCircleRadius(1)
		if _, err := s.Build(); err != gradient.ErrInvalidSize {
			t.Fatalf("Build error = %v", err)
		}
	})
	t.Run("clamped size", func(t *testing.T) {
		s := gradient.NewRadialSpec()
		s.SetCircleRadius(1)
		s.SetSize(0, -1)
		checkRadialRadii(t, &s, 0.5, 0.5)
	})
	for _, radius := range []float64{0, -1, math.NaN(), math.Inf(1)} {
		s := gradient.NewRadialSpec()
		s.SetCircleRadius(radius)
		s.SetSize(200, 100)
		if _, err := s.Build(); err != gradient.ErrInvalidRadius {
			t.Errorf("radius %g: Build error = %v, want %v", radius, err, gradient.ErrInvalidRadius)
		}
	}
	for _, size := range []float64{math.NaN(), math.Inf(1)} {
		s := gradient.NewRadialSpec()
		s.SetCircleRadius(1)
		s.SetSize(size, 100)
		if _, err := s.Build(); err != gradient.ErrInvalidSize {
			t.Errorf("size %g: Build error = %v, want %v", size, err, gradient.ErrInvalidSize)
		}
	}
}
