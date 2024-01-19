package renderer

import (
	"github.com/stretchr/testify/assert"
	"github.com/tsagae/software3d/pkg/basics"
	"testing"
)

func TestFindIntersectionPoint(t *testing.T) {
	p0 := basics.NewVector3(-1, 0, +1)
	p1 := basics.NewVector3(+1, 0, +1)
	planePoint := basics.NewVector3(0, 0, 0)
	planeNormal := basics.NewVector3(1, 0, 0)
	plane := basics.NewPlaneFromPointNormal(&planePoint, &planeNormal)

	//plane intersects segment
	clippingPoint, isClipped, _, _ := FindIntersectionPoint(&p0, &p1, &plane)

	assert.True(t, isClipped)

	expected := basics.NewVector3(0, 0, 1)
	assert.True(t, clippingPoint.Equals(expected))

	//plane intersects segment
	clippingPoint, isClipped, _, _ = FindIntersectionPoint(&p1, &p0, &plane)

	assert.True(t, isClipped)

	assert.True(t, clippingPoint.Equals(expected))

	//plane doesn't intersect segment
	planePoint2 := basics.NewVector3(0, 1, 0)
	planeNormal2 := basics.NewVector3(0, 1, 0)
	plane2 := basics.NewPlaneFromPointNormal(&planePoint2, &planeNormal2)

	_, isClipped, _, _ = FindIntersectionPoint(&p0, &p1, &plane2)

	assert.False(t, isClipped)

	//points are the same
	p0 = basics.NewVector3(-1, 0, +1)
	p1 = basics.NewVector3(-1, 0, +1)

	//no intersection if the segment is a point
	clippingPoint, isClipped, _, _ = FindIntersectionPoint(&p0, &p1, &plane)
	assert.False(t, isClipped)

	//line intersects the plane but segment doesn't
	p0 = basics.NewVector3(1, 0, +1)
	p1 = basics.NewVector3(2, 0, +1)

	clippingPoint, isClipped, _, _ = FindIntersectionPoint(&p0, &p1, &plane)
	assert.False(t, isClipped)
}

func TestClipSegment(t *testing.T) {
	/*
					view from the top

		  p2: -2,0,2				p3: 3,0,2
					| \ |         *
					|   |
					|	|\
					| <-| \
					|	|  \
					|	|   \
					----|----
		  p0: -2,0,0 	      p1: 2,0,0

					plane looks to the left
	*/

	p0 := basics.NewVector3(-2, 0, 0)
	p1 := basics.NewVector3(2, 0, 0)
	p2 := basics.NewVector3(-2, 0, 2)
	p3 := basics.NewVector3(3, 0, 2)

	expectedIntersectionA := basics.LerpVector3(&p1, &p2, 0.5)
	expectedIntersectionB := basics.NewVector3(0, 0, 0)

	planePoint := basics.NewVector3(0, 0, 0)
	planeNormal := basics.NewVector3(-1, 0, 0)
	plane := basics.NewPlaneFromPointNormal(&planePoint, &planeNormal)

	// base case, segment is clipped normally
	a, b, isKept := ClipSegment(&p0, &p1, &plane)
	assert.True(t, isKept)
	assert.True(t, a.Equals(p0) && b.Equals(expectedIntersectionB))

	a, b, isKept = ClipSegment(&p1, &p2, &plane)
	assert.True(t, isKept)
	assert.True(t, a.Equals(expectedIntersectionA) && b.Equals(p2))

	// segment is kept intact in front of the plane
	a, b, isKept = ClipSegment(&p0, &p2, &plane)
	assert.True(t, isKept)
	assert.True(t, a.Equals(p0) && b.Equals(p2))

	// segment is discarded for being fully behind the plane
	_, _, isKept = ClipSegment(&p1, &p3, &plane)
	assert.False(t, isKept)

	// segment is discarded for being coplanar
	_, _, isKept = ClipSegment(&expectedIntersectionA, &expectedIntersectionB, &plane)
	assert.False(t, isKept)
}
