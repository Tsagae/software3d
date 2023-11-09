package basics

import (
	"testing"
)

func TestQuaternionRotation(t *testing.T) {
	/*
			Vector v0(15.6, -7.1 , 13);
		    Quaternion q = Quaternion::axisAngle(
		        normalize(Vector(2, 7, 9)),
		        toRad(30)
			);
		    Vector v1 = v0;
		    for (int i = 0; i < 12; ++i) {
		        v1 = rotate(q, v1);
		    }
		    assert(isEqual(v0, v1));

		    Quaternion q1 = Quaternion::identity();
		    for (int i = 0; i < 12; ++i) {
		        q1 = q1*q;
		    }
		    assert(isEquivalent(q1, Quaternion::identity()));
	*/

	//Vector v0(15.6, -7.1 , 13);
	v0 := NewVector3(15.6, -7.1, 13)
	q := NewQuaternionFromAngleAndAxis(30, NewVector3(2, 7, 9).Normalized())
	v1 := v0

	for i := 0; i < 12; i++ {
		v1 = q.Rotated(v1)
	}

	if !(v0.Equals(&v1)) {
		t.Errorf("Error in rotation cumulation on a vector")
	}

	q1 := NewIdentityQuaternion()
	for i := 0; i < 12; i++ {
		q1 = q1.Mul(&q)
	}

	id := NewIdentityQuaternion()
	if !(q1.Equals(&id)) {
		t.Errorf("Error in rotation cumulation on a quaternion")
	}

	/*
			 Quaternion q = Quaternion::axisAngle(
		                Versor::up(),
		                toRad(180)
		                );
		    assert(isEqual(q, Quaternion(0, 1, 0, 0)));

		    Vector v(7, 2, 6);
		    assert(isEqual(rotate(q, v), Vector(-7, 2, -6)));
	*/

	q = NewQuaternionFromAngleAndAxis(180, Up())

	q2 := NewQuaternionFromScalars(0, 1, 0, 0)
	if !(q.Equals(&q2)) {
		t.Errorf("Error in quaternion constructors")
	}

	v := NewVector3(7, 2, 6)
	v1 = NewVector3(-7, 2, -6)
	v = q.Rotated(v)
	if !(v.Equals(&v1)) {
		t.Errorf("Error in vector rotations")
	}

}
