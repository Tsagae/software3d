package basics

import (
	"math"
	"testing"
)

/*
	void unitTestSumIsCommutative() {
	    Vector v(1, 2, 3), w(4, 5, 6);
	    assert(isEqual(v+w, w+v));
	}
*/
func TestCommutativeSum(t *testing.T) {
	//Vector v(1, 2, 3), w(4, 5, 6);
	v := NewVector3(1, 2, 3)
	w := NewVector3(4, 5, 6)

	vW := v.Add(&w)
	wV := w.Add(&v)

	if !vW.Equals(&wV) {
		t.Errorf("Vector sum is not commutative")
	}

}

/*
	void unitTestCross() {
	    Vector v(1, 2, 3), w(4, 5, 6);
	    assert(isEqual(cross(v, w), -cross(w, v)));
	    // Resulting vector is orthogonal to both
	    assert(isZero(dot(cross(v,w), v)));
	    assert(isZero(dot(cross(v,w), w)));
	}
*/
func TestCross(t *testing.T) {
	v := NewVector3(1, 2, 3)
	w := NewVector3(4, 5, 6)

	c1 := v.Cross(&w)
	c2 := w.Cross(&v)

	if !(c1 == c2.Inverse()) {
		t.Errorf("Vector sum is not anticommutative")
	}
}

func TestVersor(t *testing.T) {
	/*
			Vector v(1, 2, 3);
		    Versor d(v);
		    assert(isEqual(d * length(v), v));
		    assert(isEqual(length(d.asVector()), 1));
		    assert(isEqual( angleBetween(Versor::up(), Versor::forward()), M_PI/2) );
	*/
	v := NewVector3(1, 2, 3)
	d := v.Normalized()

	t1 := d.Mul(v.Length())
	if !(t1.Equals(&v)) {
		t.Errorf("Normalization error")
	}

	if !(d.Length().Equals(1)) {
		t.Errorf("Versor is not unitary")
	}

}

func TestAngleBetween(t *testing.T) {
	if !(Up().AngleBetween(Forward()) == math.Pi/2) {
		t.Errorf("Error in angle between vectors")
	}
}
