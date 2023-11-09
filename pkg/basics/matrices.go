package basics

type Matrix3 [3]Vector3 //columns
type Matrix4 [4]Vector4 //columns

// Constructors
func NewMatrix3(a Vector3, b Vector3, c Vector3) Matrix3 {
	return Matrix3{a, b, c}
}

func NewCanonicalMatrix3() Matrix3 {
	return NewMatrix3(Right(), Up(), Forward())
}

func NewMatrix4(a Vector4, b Vector4, c Vector4, d Vector4) Matrix4 {
	return Matrix4{a, b, c, d}
}

// Methods that do not change this
func (m *Matrix3) MulVec(v *Vector3) Vector3 {
	return NewVector3(m[0].Dot(v), m[1].Dot(v), m[2].Dot(v))
}

func (m *Matrix4) MulVec(v *Vector4) Vector4 {
	return NewVector4(m[0].Dot(v), m[1].Dot(v), m[2].Dot(v), m[3].Dot(v))
}
