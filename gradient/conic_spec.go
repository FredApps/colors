package gradient

import "math"

// ConicSpec specifies the parameters for a conic gradient.
type ConicSpec struct {
	Transform

	width  float64
	height float64

	centerX float64
	centerY float64

	startAngle float64

	conic *Conic
}

// NewConicSpec creates a new conic gradient specification.
func NewConicSpec() ConicSpec {
	return ConicSpec{
		Transform: Identity(),
		centerX:   half,
		centerY:   half,
	}
}

// SetSize sets the dimensions of the gradient in pixels.
// Width and height are clamped to a minimum of 1.
func (s *ConicSpec) SetSize(width, height float64) {
	s.width = width
	s.height = height
}

// SetCenter sets the normalized center point of the gradient.
func (s *ConicSpec) SetCenter(x, y float64) {
	s.centerX = x
	s.centerY = y
}

// SetStartAngle sets the starting angle of the gradient in turns.
func (s *ConicSpec) SetStartAngle(turn float64) {
	s.startAngle = turn
}

// Size returns the dimensions of the gradient.
func (s *ConicSpec) Size() (width, height float64) {
	return s.width, s.height
}

// Center returns the center point of the gradient.
func (s *ConicSpec) Center() (x, y float64) {
	return s.centerX, s.centerY
}

// StartAngle returns the starting angle of the gradient in turns.
func (s *ConicSpec) StartAngle() float64 {
	return s.startAngle
}

// Build creates a [Conic] gradient from the specification.
func (s *ConicSpec) Build() (Conic, error) {
	if s.conic == nil {
		s.conic = new(Conic)
	}

	c := s.conic

	if !isFinite(s.width) || !isFinite(s.height) ||
		s.width <= 0 || s.height <= 0 {
		return Conic{}, ErrInvalidSize
	}

	if !isFinite(s.centerX) || !isFinite(s.centerY) ||
		s.centerX < 0 || s.centerX > 1 ||
		s.centerY < 0 || s.centerY > 1 {
		return Conic{}, ErrInvalidCenter
	}

	if !isFinite(s.startAngle) {
		return Conic{}, ErrInvalidAngle
	}

	inv, ok := s.Transform.Inverse()
	if !ok {
		return Conic{}, ErrInvalidTransform
	}

	cx := s.centerX * s.width
	cy := s.centerY * s.height

	// Measure the angle relative to the center in untransformed coordinates.
	c.dx0 = inv.m00
	c.dx1 = inv.m01
	c.dx2 = inv.tx - cx

	c.dy0 = inv.m10
	c.dy1 = inv.m11
	c.dy2 = inv.ty - cy

	c.startAngle = s.startAngle - math.Floor(s.startAngle)

	return *c, nil
}
