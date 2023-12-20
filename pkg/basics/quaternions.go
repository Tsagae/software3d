package basics

import (
	"math"
)

type Quaternion struct {
	Re Scalar
	Im Vector3
}

/* Constructors */

// NewRawQuaternion Returns a non normalized quaternion
func NewRawQuaternion(re Scalar, im Vector3) Quaternion {
	return Quaternion{re, im}
}

// NewQuaternionFromAngleAndAxis Returns a normalized quaternion. The angle parameter is in degrees. Axis is not normalized
func NewQuaternionFromAngleAndAxis(angle Scalar, axis Vector3) Quaternion {
	angleRad := angle * math.Pi / 180
	re := math.Cos(float64(angleRad / 2))
	axis.ThisMul(Scalar(math.Sin(float64(angleRad / 2))))
	quat := NewRawQuaternion(Scalar(re), axis)
	quat.ThisNormalize()
	return quat
}

// NewQuaternionFromEulerAngles Returns a normalized quaternion. The angle parameters are in degrees. Applied in this order: yaw -> pitch -> roll
func NewQuaternionFromEulerAngles(yaw Scalar, pitch Scalar, roll Scalar) Quaternion {
	yawQuaternion := NewQuaternionFromAngleAndAxis(yaw, Up())
	pitchQuaternion := NewQuaternionFromAngleAndAxis(pitch, Right())
	rollQuaternion := NewQuaternionFromAngleAndAxis(roll, Forward())

	finalQuat := yawQuaternion.Mul(&pitchQuaternion)
	finalQuat.ThisMul(&rollQuaternion)
	return finalQuat
}

// NewQuaternionFromScalars Returns a non normalized quaternion. x, y, z are the imaginary part and w the real part
func NewQuaternionFromScalars(x Scalar, y Scalar, z Scalar, w Scalar) Quaternion {
	return Quaternion{w, NewVector3(x, y, z)}
}

// NewQuaternionFromMatrix Returns a non normalized quaternion from a 3x3 matrix
func NewQuaternionFromMatrix(m *Matrix3) Quaternion {
	x := m[0]
	y := m[1]
	z := m[2]
	qw := Scalar(math.Sqrt(float64(1+x.X+y.Y+z.Z)) / 2.0)
	qw4 := qw * 4
	qx := (y.Z - z.Y) / qw4
	qy := (z.X - x.Z) / qw4
	qz := (x.Y - y.X) / qw4
	return NewQuaternionFromScalars(qx, qy, qz, qw)
}

func NewIdentityQuaternion() Quaternion {
	return Quaternion{1, NewVector3(0, 0, 0)}
}

func NewForwardQuaternion() Quaternion {
	return Quaternion{0, Forward()}
}

func NewBackWardQuaternion() Quaternion {
	return Quaternion{0, Backward()}
}

/* Equality and zero */

func (q *Quaternion) IsZero() bool {
	return q.Re.IsZero() && q.Im.IsZero()
}

func (q *Quaternion) Equals(p *Quaternion) bool {
	pImInverse := p.Im.Inverse()
	return q.Re.Equals(p.Re) && q.Im.Equals(&p.Im) || q.Re.Equals(-p.Re) && q.Im.Equals(&pImInverse)
}

/* Mutable operations on this */

func (q *Quaternion) ThisConjugate() {
	q.Im.ThisInvert()
}

func (q *Quaternion) ThisAdd(p *Quaternion) {
	q.Im.ThisAdd(p.Im)
	q.Re += p.Re
}

// ThisMul Multiplication of quaternions/Accumulation of rotations. The result is normalized
func (q *Quaternion) ThisMul(p *Quaternion) {
	//can be optimized
	newRe := q.Re*p.Re - q.Im.Dot(&p.Im)
	newIm := p.Im.Mul(q.Re)
	newIm.ThisAdd(q.Im.Mul(p.Re))
	newIm.ThisAdd(q.Im.Cross(&p.Im))
	q.Im = newIm
	q.Re = newRe
	q.ThisNormalize() //normalization at the end to avoid error stacking
}

// ThisMulScalar Multiples both the imaginary part and the real part with the scalar a
func (q *Quaternion) ThisMulScalar(a Scalar) {
	q.Im.ThisMul(a)
	q.Re *= a
}

func (q *Quaternion) ThisNormalize() {
	scaling := math.Sqrt(float64(q.Im.X*q.Im.X + q.Im.Y*q.Im.Y + q.Im.Z*q.Im.Z + q.Re*q.Re))
	q.Im.ThisDiv(Scalar(scaling))
	q.Re /= Scalar(scaling)
}

/* Operations that do not change this */

func (q Quaternion) Conjugate() Quaternion {
	q.ThisConjugate()
	return q
}

func (q Quaternion) Add(p *Quaternion) Quaternion {
	q.ThisAdd(p)
	return q
}

func (q Quaternion) Mul(p *Quaternion) Quaternion {
	q.ThisMul(p)
	return q
}

func (q Quaternion) MulScalar(a Scalar) Quaternion {
	q.ThisMulScalar(a)
	return q
}

// Rotated Does not modify v
func (q Quaternion) Rotated(v Vector3) Vector3 {
	if v.IsZero() {
		return v
	}
	beforeLen := v.Length()
	//Can be optimized
	conj := q.Conjugate()
	vQuat := NewRawQuaternion(0, v)
	q.ThisMul(&vQuat) //this is always zero
	q.ThisMul(&conj)
	q.Im.ThisMul(beforeLen) //preserve scaling
	return q.Im             //I only need the imaginary part
}

// LookAt target is a point
func (orientation *Matrix3) LookAt(direction *Vector3) Quaternion {
	up := Up()

	z := *direction
	z.Normalized()

	x := up.Cross(&z)
	x.ThisNormalize()

	y := z.Cross(&x)

	/*
		qw= √(1 + m00 + m11 + m22) /2
		qx = (m21 - m12)/( 4 *qw)
		qy = (m02 - m20)/( 4 *qw)
		qz = (m10 - m01)/( 4 *qw)
	*/
	m := NewMatrix3(x, y, z)
	return NewQuaternionFromMatrix(&m)
	//rotation.ThisNormalize() not needed
}
