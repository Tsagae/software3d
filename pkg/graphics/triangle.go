package graphics

import (
	"github.com/tsagae/software3d/pkg/basics"
)

type Triangle [3]Vertex

/* Constructors */

// NewTriangle Orientation of vertices is clockwise
func NewTriangle(vertices [3]basics.Vector3, colors [3]basics.Vector3) Triangle {
	normal := computeNormalFromVertices(vertices[0], vertices[1], vertices[2])
	return [3]Vertex{
		{vertices[0], colors[0], normal},
		{vertices[1], colors[1], normal},
		{vertices[2], colors[2], normal},
	}
}

// NewTriangleWithNormals Orientation of vertices is clockwise
func NewTriangleWithNormals(vertices [3]basics.Vector3, colors [3]basics.Vector3, normals [3]basics.Vector3) Triangle {
	return [3]Vertex{
		{vertices[0], colors[0], normals[0]},
		{vertices[1], colors[1], normals[1]},
		{vertices[2], colors[2], normals[2]},
	}
}

/* Mutable operations on this */

func (t *Triangle) ThisApplyTransformation(transform *basics.Transform) {
	for i := 0; i < 3; i++ {
		transform.ApplyToPoint(&t[i].Position)
		transform.ApplyToVector(&t[i].Normal)
	}
}

/* Operations that do not change this */

func (t *Triangle) GetAverageZ() basics.Scalar {
	var sum basics.Scalar
	for i := 0; i < 3; i++ {
		sum += t[i].Position.Z
	}
	return sum / 3
}

func computeNormalFromVertices(v0 basics.Vector3, v1 basics.Vector3, v2 basics.Vector3) basics.Vector3 {
	u := v1.Sub(v0)
	v := v2.Sub(v0)
	normal := u.Cross(v).Normalized()
	return normal
}

// GetSurfaceNormal Normal of the triangle surface, ignores the mesh normals
func (t *Triangle) GetSurfaceNormal() basics.Vector3 {
	return computeNormalFromVertices(t[0].Position, t[1].Position, t[2].Position)
}

func (t *Triangle) InterpolateColor(w1, w2, w3 basics.Scalar) basics.Vector3 {
	return basics.Interpolate3(&t[0].Color, &t[1].Color, &t[2].Color, w1, w2, w3)
}

func (t *Triangle) InterpolatePosition(w1, w2, w3 basics.Scalar) basics.Vector3 {
	return basics.Interpolate3(&t[0].Position, &t[1].Position, &t[2].Position, w1, w2, w3)
}

func (t *Triangle) InterpolateVertexProps(w1, w2, w3 basics.Scalar) Vertex {
	return Vertex{
		basics.Interpolate3(&t[0].Position, &t[1].Position, &t[2].Position, w1, w2, w3),
		basics.Interpolate3(&t[0].Color, &t[1].Color, &t[2].Color, w1, w2, w3),
		basics.Interpolate3(&t[0].Normal, &t[1].Normal, &t[2].Normal, w1, w2, w3),
	}
}

func (t *Triangle) FindWeightsPosition(target *basics.Vector3) (basics.Scalar, basics.Scalar, basics.Scalar) {
	// most of this can be cached when finding weights inside the same triangle TODO
	return basics.FindWeights3(&t[0].Position, &t[1].Position, &t[2].Position, target)
}
