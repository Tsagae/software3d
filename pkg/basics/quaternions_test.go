package basics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIdentityQuaternion(t *testing.T) {
	q := NewQuaternionFromAngleAndAxis(30, NewVector3(2, 7, 9).Normalized())

	q1 := NewIdentityQuaternion()
	for i := 0; i < 12; i++ {
		q1 = q1.Mul(&q)
	}

	id := NewIdentityQuaternion()
	assert.Truef(t, id.Equals(&q1), "Accumulated quaternion to get identity does not equal identity quaternion: %v should be equal to %v", q1, id)

}

func TestQuaternionRotation(t *testing.T) {
	v0 := NewVector3(15.6, -7.1, 13)
	v1 := v0
	q := NewQuaternionFromAngleAndAxis(30, NewVector3(2, 7, 9).Normalized())

	for i := 0; i < 12; i++ {
		v1 = q.Rotated(v1)
	}

	assert.True(t, v0.Equals(&v1), "Error in rotation accumulation on a vector (should be equal to the starting vector)")

	v1 = q.Rotated(v1)
	assert.False(t, v0.Equals(&v1), "Error in rotation accumulation on a vector (should be different to the starting vector)")

	q = NewQuaternionFromAngleAndAxis(180, Up())

	v := NewVector3(7, 2, 6)
	v1 = NewVector3(-7, 2, -6)
	v = q.Rotated(v)
	assert.Truef(t, v.Equals(&v1), "Error in vector rotation")

	zeroV := Vector3{}
	newZero := q.Rotated(zeroV)
	assert.True(t, zeroV.Equals(&newZero), "Zero vector rotated is not zero")
}

func TestNewQuaternionFromAngleAndAxis(t *testing.T) {
	q := NewQuaternionFromAngleAndAxis(180, Up())

	q2 := NewQuaternionFromScalars(0, 1, 0, 0)
	assert.Truef(t, q.Equals(&q2), "Error in quaternion constructor")
}

func TestNewQuaternionFromEulerAngles(t *testing.T) {
	q := NewQuaternionFromEulerAngles(-90, 0, -90)
	p := NewVector3(1, 1, 1)
	p = q.Rotated(p)
	expected := NewVector3(-1, -1, 1)
	assert.Truef(t, expected.Equals(&p), "Error in quaternion from euler angles got: %v, expected: %v", p, expected)
}
