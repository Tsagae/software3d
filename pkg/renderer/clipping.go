package renderer

import (
	"github.com/tsagae/software3d/pkg/basics"
)

// FindIntersectionPoint return the point where the plane intersects the segment and true if there is one or if one of them is coplanar, false if it's not found, they're the same point. Returns also the 2 tests against the plane for the 2 points
func FindIntersectionPoint(p0, p1 *basics.Vector3, plane *basics.Plane) (basics.Vector3, bool, uint8, uint8) {
	test0 := plane.TestPoint(p0)
	test1 := plane.TestPoint(p1)

	if test0 == 2 {
		return *p0, true, test0, test1
	}
	if test1 == 2 {
		return *p1, true, test0, test1
	}
	if test0 == test1 {
		return basics.Vector3{}, false, test0, test1
	}

	lineN := p1.Sub(p0)
	planeP := plane.Point
	planeN := plane.Normal

	p1MinP2 := planeP.Sub(p0)
	d := lineN.Dot(&planeN)
	if d.IsZero() {
		return basics.Vector3{}, false, test0, test1
	}
	k := p1MinP2.Dot(&planeN) / lineN.Dot(&planeN)

	nK := lineN.Mul(k)
	intersection := p0.Add(&nK)
	return intersection, true, test0, test1
}
