package renderer

import (
	"fmt"
	"github.com/tsagae/software3d/pkg/basics"
	"github.com/tsagae/software3d/pkg/graphics"
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

	lineN := p1.Sub(*p0)
	planeP := plane.Point
	planeN := plane.Normal

	p1MinP2 := planeP.Sub(*p0)
	d := lineN.Dot(planeN)
	if d.IsZero() {
		return basics.Vector3{}, false, test0, test1
	}
	k := p1MinP2.Dot(planeN) / lineN.Dot(planeN)

	nK := lineN.Mul(k)
	intersection := p0.Add(nK)
	return intersection, true, test0, test1
}

// ClipSegment returns the start and end points of a segment after being clipped against a plane and false if the segment is completely on the back side the plane or coplanar to the plane
func ClipSegment(p0, p1 *basics.Vector3, plane *basics.Plane) (basics.Vector3, basics.Vector3, bool) {
	intersection, foundIntersection, test0, test1 := FindIntersectionPoint(p0, p1, plane)

	if !foundIntersection { // there is no intersection
		if test0 == 1 {
			return *p0, *p1, true // segment is completely in front of the plane
		}
		return basics.Vector3{}, basics.Vector3{}, false // segment is completely behind the plane
	}

	// there is an intersection
	if test0 != 1 && test1 != 1 { // the segment is coplanar or is completely clipped
		return basics.Vector3{}, basics.Vector3{}, false
	}

	// the segment is clipped normally
	if test0 == 1 {
		return *p0, intersection, true
	}

	return intersection, *p1, true
}

// ClipTriangle fills the buffer with the triangles created by clipping t and returns it. The triangles are not appended in the buffer but are inserted from the beginning of the slice
func ClipTriangle(t *graphics.Triangle, p *basics.Plane) []graphics.Triangle {
	//fmt.Println("in:", t[0].Position, t[1].Position, t[2].Position, "plane: ", p)

	buffer := make([]graphics.Triangle, 0)
	type vert struct {
		vertex      graphics.Vertex
		behindPlane bool
	}
	behindCount := 0
	vertices := make([]vert, 3)
	for i, vertex := range t {
		vertices[i].vertex = vertex
		if p.TestPoint(&vertex.Position) == 0 {
			vertices[i].behindPlane = true
			behindCount++
		}
	}

	switch behindCount {
	case 0:
		buffer = append(buffer, *t)
		return buffer
	case 3:
		//fmt.Printf("triangle: %v %v %v is discarded against plane %v\n", t[0].Position, t[1].Position, t[2].Position, *p)
		return buffer
	}

	newVertices := make([]graphics.Vertex, 0, 3)
	for i := 0; i < 3; i++ {
		v1 := vertices[i]
		v2 := vertices[(i+1)%3]
		if !(v1.behindPlane || v2.behindPlane) {
			newVertices = append(newVertices, v2.vertex)
		} else if v1.behindPlane && v2.behindPlane {
			//nothing
		} else {
			inters, foundOrCoplanar, _, _ := FindIntersectionPoint(&v1.vertex.Position, &v2.vertex.Position, p)
			if !foundOrCoplanar {
				return buffer
			}
			w0, w1, w2 := t.FindWeightsPosition(&inters)
			interp := t.InterpolateVertexProps(w0, w1, w2)

			if !v1.behindPlane {
				newVertices = append(newVertices, interp)
			} else {
				newVertices = append(newVertices, interp, v2.vertex)
			}
		}
	}
	if !(len(newVertices) == 3 || len(newVertices) == 4) {
		panic(fmt.Sprintf("Triangle clipping produced a polygon with %d vertices", len(newVertices)))
	}
	foundNan := false
	if len(newVertices) == 3 {
		for _, vertex := range newVertices {
			if vertex.Position.X.IsNaN() || vertex.Position.Y.IsNaN() || vertex.Position.Z.IsNaN() {
				foundNan = true
				break
			}
		}
		if !foundNan {
			buffer = append(buffer, [3]graphics.Vertex(newVertices))
		}
		return buffer
	}

	foundNan = false
	for _, vertex := range [3]graphics.Vertex(newVertices[:3]) {
		if vertex.Position.X.IsNaN() || vertex.Position.Y.IsNaN() || vertex.Position.Z.IsNaN() {
			foundNan = true
			break
		}
	}
	if !foundNan {
		buffer = append(buffer, [3]graphics.Vertex(newVertices[:3]))
	}

	foundNan = false
	v2 := [3]graphics.Vertex{newVertices[2], newVertices[3], newVertices[0]}
	for _, vertex := range v2 {
		if vertex.Position.X.IsNaN() || vertex.Position.Y.IsNaN() || vertex.Position.Z.IsNaN() {
			foundNan = true
			break
		}
	}
	if !foundNan {
		buffer = append(buffer, v2)
	}
	return buffer
}

func ClipTriangleAgainsPlanes(triangle *graphics.Triangle, planes []basics.Plane) []graphics.Triangle {
	if len(planes) == 0 {
		return []graphics.Triangle{*triangle}
	}
	outTriangles := make([]graphics.Triangle, 0)
	triangles := ClipTriangle(triangle, &planes[0])
	/*
		for i, t := range triangles {
			for j, vertex := range t {
				if math.IsNaN(float64(vertex.Position.X)) ||
					math.IsNaN(float64(vertex.Position.Y)) ||
					math.IsNaN(float64(vertex.Position.Z)) {
					fmt.Println("out triangle num:", i, "vertex num: ", j)
					fmt.Println("tri pos:", t[0].Position, t[1].Position, t[2].Position)
					panic("NaN after clipping")
				}
			}
		}
	*/

	for _, t := range triangles {
		outTriangles = append(outTriangles, ClipTriangleAgainsPlanes(&t, planes[1:])...)
	}
	return outTriangles
}
