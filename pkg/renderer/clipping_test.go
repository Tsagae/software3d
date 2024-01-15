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
	assert.True(t, clippingPoint.Equals(&expected))

	//plane intersects segment
	clippingPoint, isClipped, _, _ = FindIntersectionPoint(&p1, &p0, &plane)

	assert.True(t, isClipped)

	assert.True(t, clippingPoint.Equals(&expected))

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

	//plane does not segment
	clippingPoint, isClipped, _, _ = FindIntersectionPoint(&p0, &p1, &plane)
	assert.False(t, isClipped)
}
