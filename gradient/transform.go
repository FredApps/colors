package gradient

import "math"

// Transform represents a 2D affine transformation.
type Transform struct {
	m00, m01 float64
	m10, m11 float64
	tx, ty   float64
}

// Identity returns the identity transformation.
func Identity() Transform {
	return Transform{
		m00: 1,
		m11: 1,
	}
}

func (t *Transform) Reset() {
	t.m00 = 1
	t.m01 = 0
	t.m10 = 0
	t.m11 = 1
	t.tx = 0
	t.ty = 0
}

// TranslateX translates the transformation by x.
func (t *Transform) TranslateX(x float64) {
	t.tx += t.m00 * x
	t.ty += t.m10 * x
}

// TranslateY translates the transformation by y.
func (t *Transform) TranslateY(y float64) {
	t.tx += t.m01 * y
	t.ty += t.m11 * y
}

// Translate translates the transformation by x and y.
func (t *Transform) Translate(x, y float64) {
	t.tx += t.m00*x + t.m01*y
	t.ty += t.m10*x + t.m11*y
}

// ScaleX scales the transformation independently along the x-axes.
func (t *Transform) ScaleX(sx float64) {
	t.m00 *= sx
	t.m10 *= sx
}

// ScaleY scales the transformation independently along the y-axes.
func (t *Transform) ScaleY(sy float64) {
	t.m01 *= sy
	t.m11 *= sy
}

// Scale scales the transformation independently along the x-axes and y-axes.
func (t *Transform) Scale(sx, sy float64) {
	t.m00 *= sx
	t.m01 *= sy
	t.m10 *= sx
	t.m11 *= sy
}

// ScaleUniform scales the transformation uniformly by s.
func (t *Transform) ScaleUniform(s float64) {
	t.m00 *= s
	t.m01 *= s
	t.m10 *= s
	t.m11 *= s
}

// Rotate rotates the transformation by the specified angle in radians.
func (t *Transform) Rotate(angle float64) {
	sin, cos := math.Sincos(angle)

	m00 := t.m00*cos + t.m01*sin
	m01 := -t.m00*sin + t.m01*cos
	m10 := t.m10*cos + t.m11*sin
	m11 := -t.m10*sin + t.m11*cos

	t.m00 = m00
	t.m01 = m01
	t.m10 = m10
	t.m11 = m11
}

// SkewX skews the transformation along the x-axis by the specified angle in radians.
func (t *Transform) SkewX(angle float64) {
	k := math.Tan(angle)

	t.m01 += t.m00 * k
	t.m11 += t.m10 * k
}

// SkewY skews the transformation along the y-axis by the specified angle in radians.
func (t *Transform) SkewY(angle float64) {
	k := math.Tan(angle)

	t.m00 += t.m01 * k
	t.m10 += t.m11 * k
}

// Skew skews the transformation along the x-axis and y-axes by the specified angles in radians.
func (t *Transform) Skew(angleX, angleY float64) {
	kx := math.Tan(angleX)
	ky := math.Tan(angleY)

	m00 := t.m00
	m01 := t.m01
	m10 := t.m10
	m11 := t.m11

	t.m00 = m00 + m01*ky
	t.m01 = m00*kx + m01
	t.m10 = m10 + m11*ky
	t.m11 = m10*kx + m11
}

// Matrix sets the transformation matrix to the specified values.
func (t *Transform) Matrix(m Transform) {
	m00 := t.m00*m.m00 + t.m01*m.m10
	m01 := t.m00*m.m01 + t.m01*m.m11
	m10 := t.m10*m.m00 + t.m11*m.m10
	m11 := t.m10*m.m01 + t.m11*m.m11

	tx := t.m00*m.tx + t.m01*m.ty + t.tx
	ty := t.m10*m.tx + t.m11*m.ty + t.ty

	t.m00 = m00
	t.m01 = m01
	t.m10 = m10
	t.m11 = m11
	t.tx = tx
	t.ty = ty
}

// PositionAt transforms the specified position by the transformation.
func (t *Transform) PositionAt(x, y float64) (float64, float64) {
	return t.m00*x + t.m01*y + t.tx, t.m10*x + t.m11*y + t.ty
}

// VectorAt transforms the specified vector by the linear part of the transformation.
func (t *Transform) VectorAt(x, y float64) (float64, float64) {
	return t.m00*x + t.m01*y, t.m10*x + t.m11*y
}

// Determinant returns the determinant of the linear part of the transformation.
func (t *Transform) Determinant() float64 {
	return t.m00*t.m11 - t.m01*t.m10
}

// Inverse returns the inverse transformation and reports whether it exists.
func (t Transform) Inverse() (Transform, bool) {
	det := t.Determinant()
	if det == 0 {
		return Transform{}, false
	}

	invDet := 1 / det

	return Transform{
		m00: t.m11 * invDet,
		m01: -t.m01 * invDet,
		m10: -t.m10 * invDet,
		m11: t.m00 * invDet,

		tx: (t.m01*t.ty - t.m11*t.tx) * invDet,
		ty: (t.m10*t.tx - t.m00*t.ty) * invDet,
	}, true
}
