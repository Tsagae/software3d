package basics

import (
	"testing"
)

func TestCumulation(t *testing.T) {
	/*
		yRotation := NewQuaternionFromAngleAndAxis(180, Up())
		t1 := NewTransform(1, NewIdentityQuaternion(), NewVector3(1, 1, 1))
		t2 := NewTransform(1, yRotation, ZeroVector())

		t3 := t1.Cumulate(&t2)
		for i := 0; i < 100; i++ {
			t3.ThisCumulate(&t2)
		}
		_ = t3*/
	/*
			 Transform T1(
		        3.5,
		        normalize(Quaternion(3, -3, 5, 9)),
		        Vector(5, 10, 15)
		    );
		    Transform T2(
		        0.2,
		        Quaternion::axisAngle(Versor::up(), 1.5),
		        Vector(-4, 7, 5.5)
		    );
		    Point p(0.2, 0.1, 0.6);

		    assert(isEqual(
		        T2(T1(p)),
		        (T2 * T1).apply(p)
		    ));

		    Transform T3 = T1.inverse();
		    assert(isEqual(
		        p,
		        T3(T1(p))
		    ));
	*/

	quat := NewQuaternionFromScalars(3, -3, 5, 9)
	quat.ThisNormalize()
	t1 := NewTransform(3.5, quat, NewVector3(5, 10, 15))

	p0 := NewVector3(0.2, 0.1, 0.6) //point
	p1 := NewVector3(0.2, 0.1, 0.6) //point

	t3 := t1.Inverse()

	t1.ApplyToPoint(&p0)
	t3.ApplyToPoint(&p0)

	if !p0.Equals(p1) {
		t.Errorf("Error in inverse transform")
	}
}
