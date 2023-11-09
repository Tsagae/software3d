package graphics

/*
import (
	"GoSDL/pkg/basics"
	"math"
)

type ProjectionMatrix basics.Matrix4

// Constructors
func NewProjectionMatrix(height basics.Scalar, width basics.Scalar, fov basics.Scalar, zFar basics.Scalar, zNear basics.Scalar) ProjectionMatrix {

	a := height / width
	f := basics.Scalar(1 / math.Tan((float64(fov)/2)*(math.Pi/180)))
	q := zFar / (zFar - zNear)
	mat := basics.NewMatrix4(
		//columns
		basics.NewVector4(a*f, 0, 0, 0),
		basics.NewVector4(0, f, 0, 0),
		basics.NewVector4(0, 0, q, -zNear*q),
		basics.NewVector4(0, 0, 1, 0),
	)
	return (ProjectionMatrix)(mat)

}

// Changes v
func ProjectionScaling(v *basics.Vector3, width basics.Scalar, height basics.Scalar) {
	v.X += 1
	v.Y += 1
	v.X *= 0.5 * width
	v.Y *= 0.5 * height
}

func (m *ProjectionMatrix) ProjectedVector(v *basics.Vector3) basics.Vector3 {
	o := basics.NewVector3(dot3(&m[0], v), dot3(&m[1], v), dot3(&m[2], v))
	w := dot3(&m[3], v)
	if w != 0 {
		o.X /= w
		o.Y /= w
		o.Z /= w
	}
	return o
}

func dot3(v *basics.Vector4, h *basics.Vector3) basics.Scalar {
	return v.X*h.X + v.Y*h.Y + v.Z*h.Z + v.W //h.w is implicitly 1
}

func (t *Triangle) ProjectedTriangle(projMat *ProjectionMatrix, width basics.Scalar, height basics.Scalar) Triangle {
	p0 := projMat.ProjectedVector(&t.Vertices[0])
	ProjectionScaling(&p0, width, height)

	p1 := projMat.ProjectedVector(&t.Vertices[1])
	ProjectionScaling(&p1, width, height)

	p2 := projMat.ProjectedVector(&t.Vertices[2])
	ProjectionScaling(&p2, width, height)

	return NewTriangle(p0, p1, p2, t.Color)
}

/*
func (m *Mesh) ProjectedMesh(projMat *ProjectionMatrix, width basics.Scalar, height basics.Scalar) Mesh {
	projectedTriangles := make([]Triangle, 0, len(m.Triangles))
	for _, tri := range m.Triangles {
		projectedTriangles = append(projectedTriangles, *tri.ProjectedTriangle(projMat, width, height))
	}
	return NewMesh(projectedTriangles)
}
*/
