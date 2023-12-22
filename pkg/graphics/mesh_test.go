package graphics

import (
	"github.com/stretchr/testify/assert"
	"github.com/tsagae/software3d/pkg/basics"
	"image/color"
	"strings"
	"testing"
)

func TestMeshIterator(t *testing.T) {
	mesh := Mesh{
		geometry: []VertexAttributes{
			{basics.NewVector3(-1.0, -1.0, -1.0), basics.NewVector3(0, 0, 0), basics.NewVector3(-0.57735027, -0.57735027, -0.57735027)},
			{basics.NewVector3(-1.0, -1.0, 1.0), basics.NewVector3(0, 0, 0), basics.NewVector3(-0.57735027, -0.57735027, 0.57735027)},
			{basics.NewVector3(-1.0, 1.0, -1.0), basics.NewVector3(0, 0, 0), basics.NewVector3(-0.57735027, 0.57735027, -0.57735027)},
			{basics.NewVector3(-1.0, 1.0, 1.0), basics.NewVector3(0, 0, 0), basics.NewVector3(-0.57735027, 0.57735027, 0.57735027)},
			{basics.NewVector3(1.0, -1.0, -1.0), basics.NewVector3(0, 0, 0), basics.NewVector3(0.57735027, -0.57735027, -0.57735027)},
		},
		connectivity: []TriangleConnectivity{
			{0, 4, 1},
			{1, 2, 0},
		},
	}
	triangles := mesh.GetTriangles()

	trianglesFromIter := make([]Triangle, len(triangles))
	i := 0
	iter := mesh.Iterator()
	for iter.HasNext() {
		trianglesFromIter[i] = iter.Next()
		i++
	}
	assert.Equal(t, trianglesFromIter, triangles)
}

func TestNewMeshFromReader(t *testing.T) {
	mesh := Mesh{
		geometry: []VertexAttributes{
			{basics.NewVector3(-1.0, -1.0, -1.0), basics.NewVector3(0, 0, 0), basics.NewVector3(-0.57735027, -0.57735027, -0.57735027)},
			{basics.NewVector3(-1.0, -1.0, 1.0), basics.NewVector3(0, 0, 0), basics.NewVector3(-0.57735027, -0.57735027, 0.57735027)},
			{basics.NewVector3(-1.0, 1.0, -1.0), basics.NewVector3(0, 0, 0), basics.NewVector3(-0.57735027, 0.57735027, -0.57735027)},
			{basics.NewVector3(-1.0, 1.0, 1.0), basics.NewVector3(0, 0, 0), basics.NewVector3(-0.57735027, 0.57735027, 0.57735027)},
			{basics.NewVector3(1.0, -1.0, -1.0), basics.NewVector3(0, 0, 0), basics.NewVector3(0.57735027, -0.57735027, -0.57735027)},
		},
		connectivity: []TriangleConnectivity{
			{0, 4, 1},
			{1, 2, 0},
		},
	}
	meshText := `
# Exported from Wings 3D 2.2.9
#mtllib cube.mtl
o Cube1
#5 vertices, 2 faces
v -1.00000000 -1.00000000 -1.00000000
v -1.00000000 -1.00000000 1.00000000
v -1.00000000 1.00000000 -1.00000000
v -1.00000000 1.00000000 1.00000000
v 1.00000000 -1.00000000 -1.00000000
vn -0.57735027 -0.57735027 -0.57735027
vn -0.57735027 -0.57735027 0.57735027
vn -0.57735027 0.57735027 -0.57735027
vn -0.57735027 0.57735027 0.57735027
vn 0.57735027 -0.57735027 -0.57735027
g Cube1_default
#usemtl default
s 1
f 1//1 5//5 2//2
f 2//2 3//3 1//1
`
	meshFromReader, err := NewMeshFromReader(strings.NewReader(meshText),
		basics.Vector3FromColor(color.RGBA{
			R: 0,
			G: 0,
			B: 0,
			A: 0,
		}))

	assert.Nil(t, err, "Error while reading mesh")
	assert.Equal(t, mesh, meshFromReader, "Mesh from reader is not corrent")
}
