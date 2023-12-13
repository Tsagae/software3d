package basics

type Vector4 struct {
	X, Y, Z, W Scalar
}

func NewVector4(x Scalar, y Scalar, z Scalar, w Scalar) Vector4 {
	return Vector4{x, y, z, w}
}

func (v *Vector4) Dot(h *Vector4) Scalar {
	return v.X*h.X + v.Y*h.Y + v.Z*h.Z + v.W*h.W
}
