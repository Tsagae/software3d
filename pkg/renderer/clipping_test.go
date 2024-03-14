package renderer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tsagae/software3d/pkg/basics"
	"github.com/tsagae/software3d/pkg/graphics"
	"math"
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

func TestClipTriangle1Side(t *testing.T) {
	zero := basics.Vector3{}
	up := basics.Up()
	plane := basics.NewPlaneFromPointNormal(&zero, &up)
	buffer := make([]graphics.Triangle, 0)

	// 1 side in front of plane
	tri := graphics.Triangle{
		graphics.Vertex{
			Position: basics.Vector3{X: -1, Y: -1},
		},
		graphics.Vertex{
			Position: basics.Vector3{X: 1, Y: -1},
		},
		graphics.Vertex{
			Position: basics.Vector3{Y: 1},
		}}
	expected := graphics.Triangle{
		graphics.Vertex{
			Position: basics.Vector3{X: 0.5},
		},
		graphics.Vertex{
			Position: basics.Vector3{Y: 1},
		},
		graphics.Vertex{
			Position: basics.Vector3{X: -0.5},
		}}
	buffer = ClipTriangle(&tri, &plane)
	assert.Equal(t, 1, len(buffer))
	assert.Equal(t, expected, buffer[0])

	tri = graphics.Triangle{
		graphics.Vertex{
			Position: basics.Vector3{Y: 1},
		},
		graphics.Vertex{
			Position: basics.Vector3{X: -1, Y: -1},
		},
		graphics.Vertex{
			Position: basics.Vector3{X: 1, Y: -1},
		}}
	expected = graphics.Triangle{
		graphics.Vertex{
			Position: basics.Vector3{X: -0.5},
		},
		graphics.Vertex{
			Position: basics.Vector3{X: 0.5},
		},
		graphics.Vertex{
			Position: basics.Vector3{Y: 1},
		}}
	buffer = ClipTriangle(&tri, &plane)
	assert.Equal(t, 1, len(buffer))
	assert.Equal(t, expected, buffer[0])

}

func TestClipTriangle2Sides(t *testing.T) {
	// 2 sides in front of plane
	zero := basics.Vector3{}
	up := basics.Up()
	plane := basics.NewPlaneFromPointNormal(&zero, &up)
	buffer := make([]graphics.Triangle, 0)

	tri := graphics.Triangle{
		graphics.Vertex{
			Position: basics.Vector3{X: -1, Y: 1},
		},
		graphics.Vertex{
			Position: basics.Vector3{Y: -1},
		},
		graphics.Vertex{
			Position: basics.Vector3{X: 1, Y: 1},
		}}
	expected := graphics.Triangle{
		graphics.Vertex{
			Position: basics.Vector3{X: -0.5},
		},
		graphics.Vertex{
			Position: basics.Vector3{X: 0.5},
		},
		graphics.Vertex{
			Position: basics.Vector3{X: 1, Y: 1},
		}}
	expectedSecond := graphics.Triangle{
		graphics.Vertex{
			Position: basics.Vector3{X: 1, Y: 1},
		},
		graphics.Vertex{
			Position: basics.Vector3{X: -1, Y: 1},
		},
		graphics.Vertex{
			Position: basics.Vector3{X: -0.5},
		}}
	buffer = ClipTriangle(&tri, &plane)
	assert.Equal(t, 2, len(buffer))
	assert.Equal(t, expected, buffer[0])
	assert.Equal(t, expectedSecond, buffer[1])
}

// tPos : {8.000000000000002 -3.000000000000001 3} {8.000000000000002 -3.000000000000001 3} {8.000000000000002 -8.000000000000002 8}
// plane &{{1.3333333333333333 -1 1} {-0.6 0 0.7999999999999999}}
func TestClipTriangleNan(t *testing.T) {
	tri := graphics.Triangle{
		{Position: basics.Vector3{8, -3, 3}},
		{Position: basics.Vector3{8, -3, 3}},
		{Position: basics.Vector3{8, -8, 8}},
	}
	point := basics.Vector3{1.3333333333333333, -1, 1}
	normal := basics.Vector3{-0.6, 0, 0.8}
	plane := basics.NewPlaneFromPointNormal(&point, &normal)

	result := ClipTriangle(&tri, &plane)
	for _, triangle := range result {
		for _, vertex := range triangle {
			fmt.Println(vertex.Position)
			assert.False(t, math.IsNaN(float64(vertex.Position.X)) ||
				math.IsNaN(float64(vertex.Position.Y)) ||
				math.IsNaN(float64(vertex.Position.Z)),
				"vertex position component is NaN")
		}

	}
}
