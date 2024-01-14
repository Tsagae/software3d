package graphics

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tsagae/software3d/pkg/basics"
	"testing"
)

func TestTriangle_GetSurfaceNormal(t *testing.T) {
	triangle := NewTriangle(
		[3]basics.Vector3{
			basics.NewVector3(1, 1, 0),
			basics.NewVector3(2, 3, 0),
			basics.NewVector3(4, 1, 0),
		},
		[3]basics.Vector3{
			basics.NewVector3(0, 0, 0),
			basics.NewVector3(0, 0, 0),
			basics.NewVector3(0, 0, 0),
		},
	)

	bw := basics.Backward()
	surfaceNormal := triangle.GetSurfaceNormal()

	fmt.Println(bw, surfaceNormal)
	assert.True(t, surfaceNormal.Equals(&bw), "Error in surface normal")
}

func TestTriangle_ThisApplyTransformation(t *testing.T) {
	triangle := NewTriangle(
		[3]basics.Vector3{
			basics.NewVector3(1, 1, 0),
			basics.NewVector3(2, 3, 0),
			basics.NewVector3(4, 1, 0),
		},
		[3]basics.Vector3{
			basics.NewVector3(0, 0, 0),
			basics.NewVector3(0, 0, 0),
			basics.NewVector3(0, 0, 0),
		},
	)

	transformation := basics.NewTransform(1, basics.NewQuaternionFromEulerAngles(180, 0, 0), basics.Vector3{})
	triangle.ThisApplyTransformation(&transformation)

	expected0 := basics.NewVector3(-1, 1, 0)
	expected1 := basics.NewVector3(-2, 3, 0)
	expected2 := basics.NewVector3(-4, 1, 0)

	assert.Truef(t, expected0.Equals(&triangle[0].Position), "Error in transforming Vertex 0, actual: %v, expected %v", triangle[0].Position, expected0)
	assert.Truef(t, expected1.Equals(&triangle[1].Position), "Error in transforming Vertex 1, actual: %v, expected %v", triangle[1].Position, expected1)
	assert.Truef(t, expected2.Equals(&triangle[2].Position), "Error in transforming Vertex 2, actual: %v, expected %v", triangle[2].Position, expected2)

}
