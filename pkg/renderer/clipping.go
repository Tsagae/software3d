package renderer

import (
	"github.com/tsagae/software3d/pkg/basics"
)

// ClipSegment return the point where the plane intersects the segment and true if there is one, false if it's not found, they're the same point or one of them is coplanar
func ClipSegment(p0, p1 *basics.Vector3, plane *basics.Plane) (basics.Vector3, bool) {
	test0 := plane.TestPoint(p0)
	test1 := plane.TestPoint(p1)

	if test0 == 2 {
		return basics.Vector3{}, false
	}
	if test1 == 2 {
		return basics.Vector3{}, false
	}
	if test0 == test1 {
		return basics.Vector3{}, false
	}

	lineN := p1.Sub(p0)
	planeP := plane.Point()
	planeN := plane.Normal()

	p1MinP2 := planeP.Sub(p0)
	d := lineN.Dot(&planeN)
	if d.IsZero() {
		return basics.Vector3{}, false
	}
	k := p1MinP2.Dot(&planeN) / lineN.Dot(&planeN)

	nK := lineN.Mul(k)
	intersection := p0.Add(&nK)
	return intersection, true
}
