package basics

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestCommutativeSum(t *testing.T) {
	//Vector v(1, 2, 3), w(4, 5, 6);
	v := NewVector3(1, 2, 3)
	w := NewVector3(4, 5, 6)

	assert.Equal(t, v.Add(&w), w.Add(&v), "Vector sum is not commutative")
}

func TestSub(t *testing.T) {
	left := Left()
	assert.Equal(t, NewVector3(2, 0, 0), Right().Sub(&left), "Incorrect vector subtraction")
}

func TestDiv(t *testing.T) {
	down := Down().Mul(2)
	assert.Equal(t, Down(), down.Div(2), "Incorrect vector division")
}

func TestCross(t *testing.T) {
	/*
		void unitTestCross() {
		    Vector v(1, 2, 3), w(4, 5, 6);
		    assert(isEqual(cross(v, w), -cross(w, v)));
		    // Resulting vector is orthogonal to both
		    assert(isZero(dot(cross(v,w), v)));
		    assert(isZero(dot(cross(v,w), w)));
		}
	*/
	v := NewVector3(1, 2, 3)
	w := NewVector3(4, 5, 6)

	c1 := v.Cross(&w)
	c2 := w.Cross(&v)

	assert.Equal(t, c1, c2.Inverse(), "Vector sum is not anticommutative")
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

	assert.True(t, t1.Equals(&v), "Normalization error")

	assert.True(t, d.Length().Equals(1), "Versor is not unitary")
}

func TestAngleBetween(t *testing.T) {
	up := Up()
	fmt.Printf("angle: %v expected: %v\n", up.AngleBetween(Forward()), math.Pi/2)
	assert.Equal(t, Scalar(math.Pi/2), up.AngleBetween(Forward()), "Error in angle between vectors")
}

func TestMulComponents(t *testing.T) {
	v := NewVector3(10, 20, 30)
	w := NewVector3(2, 3, 4)
	assert.Equal(t, NewVector3(20, 60, 120), v.MulComponents(&w))
}
