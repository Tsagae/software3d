package graphics

import (
	"github.com/tsagae/software3d/pkg/basics"
)

type Triangle [3]Vertex

type Vertex struct {
	Position basics.Vector3
	Color    basics.Vector3
	Normal   basics.Vector3
}

// Constructors
func NewTriangle(vertices [3]basics.Vector3, colors [3]basics.Vector3) Triangle {
	normal := computeNormalFromVertices(vertices[0], vertices[1], vertices[2])
	return [3]Vertex{
		Vertex{vertices[0], colors[0], normal},
		Vertex{vertices[1], colors[1], normal},
		Vertex{vertices[2], colors[2], normal},
	}
}

func NewTriangleWithNormals(vertices [3]basics.Vector3, colors [3]basics.Vector3, normals [3]basics.Vector3) Triangle {
	return [3]Vertex{
		Vertex{vertices[0], colors[0], normals[0]},
		Vertex{vertices[1], colors[1], normals[1]},
		Vertex{vertices[2], colors[2], normals[2]},
	}
}

// Mutable operations on this
func (t *Triangle) ThisApplyTransformation(transform *basics.Transform) {
	for i := 0; i < 3; i++ {
		transform.ApplyToPoint(&t[i].Position)
		transform.ApplyToVector(&t[i].Normal)
	}
	/*
		if transform.IsInverse() {
			t.Vertices[0].ThisAdd(&transform.Position)
			transform.Rotation.Rotate(&t.Vertices[0])
			t.Vertices[0].ThisMul(transform.Scale)

			t.Vertices[1].ThisAdd(&transform.Position)
			transform.Rotation.Rotate(&t.Vertices[1])
			t.Vertices[1].ThisMul(transform.Scale)

			t.Vertices[2].ThisAdd(&transform.Position)
			transform.Rotation.Rotate(&t.Vertices[2])
			t.Vertices[2].ThisMul(transform.Scale)
			return
		}
		t.Vertices[0].ThisMul(transform.Scale)
		transform.Rotation.Rotate(&t.Vertices[0])
		t.Vertices[0].ThisAdd(&transform.Position)

		t.Vertices[1].ThisMul(transform.Scale)
		transform.Rotation.Rotate(&t.Vertices[1])
		t.Vertices[1].ThisAdd(&transform.Position)

		t.Vertices[2].ThisMul(transform.Scale)
		transform.Rotation.Rotate(&t.Vertices[2])
		t.Vertices[2].ThisAdd(&transform.Position)
	*/
}

// Operations that do not change this
func (t *Triangle) GetAverageZ() basics.Scalar {
	var sum basics.Scalar
	for i := 0; i < 3; i++ {
		sum += t[i].Position.Z
	}
	return sum / 3
}

func computeNormalFromVertices(v0 basics.Vector3, v1 basics.Vector3, v2 basics.Vector3) basics.Vector3 {
	u := v1.Sub(&v0)
	v := v2.Sub(&v0)
	normal := u.Cross(&v)
	normal.ThisNormalize()
	return normal
}

// Normal of the triangle surface, ignores the mesh normals
func (t *Triangle) GetSurfaceNormal() basics.Vector3 {
	return computeNormalFromVertices(t[0].Position, t[1].Position, t[2].Position)
}
