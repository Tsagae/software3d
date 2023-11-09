package basics

import (
	"image/color"
	"math"
)

type Vector3 struct {
	X, Y, Z Scalar
}

type Vector4 struct {
	X, Y, Z, W Scalar
}

// Constructors
func NewVector3(x Scalar, y Scalar, z Scalar) Vector3 {
	return Vector3{x, y, z}
}

func NewVector4(x Scalar, y Scalar, z Scalar, w Scalar) Vector4 {
	return Vector4{x, y, z, w}
}

// Range 0-65535 Loses alpha value
func Vector3FromColor(c color.Color) Vector3 {
	r, g, b, _ := c.RGBA()
	return NewVector3(Scalar(r), Scalar(g), Scalar(b))
}

// Constants
func ZeroVector() Vector3 {
	return Vector3{0, 0, 0}
}

func Right() Vector3 {
	return Vector3{1, 0, 0}
}

func Left() Vector3 {
	return Vector3{-1, 0, 0}
}
func Up() Vector3 {
	return Vector3{0, 1, 0}
}

func Down() Vector3 {
	return Vector3{0, -1, 0}
}

func Forward() Vector3 {
	return Vector3{0, 0, 1}
}

func Backward() Vector3 {
	return Vector3{0, 0, -1}
}

// Equality and zero
func (v *Vector3) IsZero() bool {
	return v.Length() < epsilon
}

func (v *Vector3) Equals(h *Vector3) bool {
	return v.X.Equals(h.X) && v.Y.Equals(h.Y) && v.Z.Equals(h.Z)
}

// Mutable operations on this
func (v *Vector3) ThisNormalize() {
	length := v.Length()
	v.X /= length
	v.Y /= length
	v.Z /= length
}

func (v *Vector3) ThisMul(a Scalar) {
	v.X *= a
	v.Y *= a
	v.Z *= a
}

func (v *Vector3) ThisDiv(a Scalar) {
	v.X /= a
	v.Y /= a
	v.Z /= a
}

func (v *Vector3) ThisAdd(h Vector3) {
	v.X += h.X
	v.Y += h.Y
	v.Z += h.Z
}

func (v *Vector3) ThisSub(h Vector3) {
	v.X -= h.X
	v.Y -= h.Y
	v.Z -= h.Z
}

func (v *Vector3) ThisInvert() {
	v.X = -v.X
	v.Y = -v.Y
	v.Z = -v.Z
}

// Operations that do not change this

func (v Vector3) Normalized() Vector3 {
	v.ThisNormalize()
	return v
}

func (v Vector3) Mul(a Scalar) Vector3 {
	v.ThisMul(a)
	return v
}

// Per component multiplication
func (v *Vector3) MulComponents(h *Vector3) Vector3 {
	return Vector3{v.X * h.X, v.Y * h.Y, v.Z * h.Z}
}

func (v Vector3) Div(a Scalar) Vector3 {
	v.ThisDiv(a)
	return v
}

func (v Vector3) Add(h *Vector3) Vector3 {
	return Vector3{v.X + h.X, v.Y + h.Y, v.Z + h.Z}
}

func (v Vector3) Sub(h *Vector3) Vector3 {
	return Vector3{v.X - h.X, v.Y - h.Y, v.Z - h.Z}
}

func (v Vector3) Inverse() Vector3 {
	v.ThisInvert()
	return v
}

func (v *Vector3) Dot(h *Vector3) Scalar {
	return v.X*h.X + v.Y*h.Y + v.Z*h.Z
}

func (v *Vector4) Dot(h *Vector4) Scalar {
	return v.X*h.X + v.Y*h.Y + v.Z*h.Z + v.W*h.W
}

func (v *Vector3) Cross(h *Vector3) Vector3 {
	return Vector3{
		X: v.Y*h.Z - v.Z*h.Y,
		Y: v.Z*h.X - v.X*h.Z,
		Z: v.X*h.Y - v.Y*h.X,
	}
}

func (v Vector3) Length() Scalar {
	return Scalar(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}

// Angle in radiants
func (v Vector3) AngleBetween(h Vector3) Scalar {
	v.ThisNormalize()
	h.ThisNormalize()
	return Scalar(math.Atan2(float64(v.Cross(&h).Length()), float64(v.Dot(&h))))
}

// Range 0-255 with 255 alpha
func (v *Vector3) ToColor() color.RGBA {
	return color.RGBA{uint8(v.X), uint8(v.Y), uint8(v.Z), 255}
}

//TODO commutative functions ?
