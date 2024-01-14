package graphics

import (
	"bufio"
	"github.com/tsagae/software3d/pkg/basics"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type MeshIterator struct {
	index int
	mesh  *Mesh
}

type VertexAttributes struct {
	position basics.Vector3
	color    basics.Vector3 //Range 0-65535
	normal   basics.Vector3
}

type TriangleConnectivity [3]int

// The mesh winding order is assumed to be counterclockwise
type Mesh struct {
	geometry     []VertexAttributes
	connectivity []TriangleConnectivity
}

/* Constructors */

func NewMesh(geometry []VertexAttributes, connectivity []TriangleConnectivity) Mesh {
	return Mesh{geometry, connectivity} //should copy the slices
}

func NewEmpyMesh() Mesh {
	return Mesh{nil, nil}
}

func (m *Mesh) GetTriangles() []Triangle {
	return m.getTrianglesWithNormals()
}

// GetTrianglesWithFaceNormals Returns a slice of triangles but ignores the mesh normals and uses the normals of the faces
func (m *Mesh) GetTrianglesWithFaceNormals() []Triangle {
	return m.getTrianglesWithoutNormals()
}

func (m *Mesh) getTrianglesWithNormals() []Triangle {
	triangles := make([]Triangle, len(m.connectivity))
	for i, v := range m.connectivity {
		triangles[i] = NewTriangleWithNormals(
			[3]basics.Vector3{m.geometry[v[0]].position, m.geometry[v[1]].position, m.geometry[v[2]].position},
			[3]basics.Vector3{m.geometry[v[0]].color, m.geometry[v[1]].color, m.geometry[v[2]].color},
			[3]basics.Vector3{m.geometry[v[0]].normal, m.geometry[v[1]].normal, m.geometry[v[2]].normal},
		)
	}
	return triangles
}

func (m *Mesh) getTrianglesWithoutNormals() []Triangle {
	triangles := make([]Triangle, len(m.connectivity))
	for i, v := range m.connectivity {
		triangles[i] = NewTriangle(
			[3]basics.Vector3{m.geometry[v[0]].position, m.geometry[v[1]].position, m.geometry[v[2]].position},
			[3]basics.Vector3{m.geometry[v[0]].color, m.geometry[v[1]].color, m.geometry[v[2]].color},
		)
	}
	return triangles
}

/* Mesh Iterator */

func (m *Mesh) Iterator() MeshIterator {
	return MeshIterator{
		index: 0,
		mesh:  m,
	}
}

// Next Returns the next triangle in the geometry. Undefined behavior when called after HasNext has returned false
func (m *MeshIterator) Next() Triangle {
	mesh := m.mesh
	connectivityItem := mesh.connectivity[m.index]
	tri := NewTriangleWithNormals(
		[3]basics.Vector3{mesh.geometry[connectivityItem[0]].position, mesh.geometry[connectivityItem[1]].position, mesh.geometry[connectivityItem[2]].position},
		[3]basics.Vector3{mesh.geometry[connectivityItem[0]].color, mesh.geometry[connectivityItem[1]].color, mesh.geometry[connectivityItem[2]].color},
		[3]basics.Vector3{mesh.geometry[connectivityItem[0]].normal, mesh.geometry[connectivityItem[1]].normal, mesh.geometry[connectivityItem[2]].normal},
	)
	m.index++
	return tri
}

// NextWithFaceNormals Returns the next triangle in the geometry, ignores the mesh normals and uses the normals of the faces. Undefined behavior when called after HasNext has returned false.
func (m *MeshIterator) NextWithFaceNormals() Triangle {
	mesh := m.mesh
	connectivityItem := mesh.connectivity[m.index]
	tri := NewTriangle(
		[3]basics.Vector3{mesh.geometry[connectivityItem[0]].position, mesh.geometry[connectivityItem[1]].position, mesh.geometry[connectivityItem[2]].position},
		[3]basics.Vector3{mesh.geometry[connectivityItem[0]].color, mesh.geometry[connectivityItem[1]].color, mesh.geometry[connectivityItem[2]].color},
	)
	m.index++
	return tri
}

// HasNext Returns true if the iterator can return at least another triangle
func (m *MeshIterator) HasNext() bool {
	return m.index < len(m.mesh.connectivity)
}

/* Mesh reader */

// NewMeshFromReader reads a mesh in obj format
func NewMeshFromReader(reader io.Reader, color basics.Vector3) (Mesh, error) {
	/*TODO: should make it more reliable by throwing the appropriate error when scanning incorrect lines
	for example:
	#8 vertices, 12 faces
	v -1.00000000 -1.00000000 -1.00000000
	vn -0.57735027 -0.57735027 0.57735027
	should throw some kind of format error on the second line since he should be reading 7 more vertices and not a normal
	*/
	var line string
	var nVertices, nFaces int
	mesh := NewEmpyMesh()

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		match, err := regexp.MatchString("\\d.*(vertices).*\\d.*faces", scanner.Text())
		if err != nil {
			return mesh, err
		}
		if match {
			break
		}
	}

	// Get vertices and faces length
	line = scanner.Text()
	splittedLine := strings.Split(line, ",")

	verticesString := splittedLine[0][1:]
	facesString := splittedLine[1][1:]

	nVertices, err := getNumberOfElements(verticesString)
	if err != nil {
		return mesh, err
	}

	nFaces, err = getNumberOfElements(facesString)
	if err != nil {
		return mesh, err
	}

	mesh.geometry = make([]VertexAttributes, nVertices)
	mesh.connectivity = make([]TriangleConnectivity, nFaces)

	//fmt.Printf("nVertices: %d nFaces: %d\n", nVertices, nFaces)

	// Get vertices
	for i := 0; i < nVertices; i++ {
		if !scanner.Scan() {
			return mesh, scanner.Err()
		}

		vertex, err := getVectorFromLine(scanner.Text())
		if err != nil {
			return mesh, err
		}
		mesh.geometry[i].position = vertex
		mesh.geometry[i].color = color //hardcoded vertex color
	}

	// Get normals
	for i := 0; i < nVertices; i++ {
		if !scanner.Scan() {
			return mesh, scanner.Err()
		}

		vertex, err := getVectorFromLine(scanner.Text())
		if err != nil {
			return mesh, err
		}
		//vertex.ThisNormalize() // normalization for normals, shouldn't be needed
		mesh.geometry[i].normal = vertex
	}

	var currentLine string
	for scanner.Scan() {
		currentLine = scanner.Text()
		match, err := regexp.MatchString("f( (\\d+//\\d+))+", currentLine)
		if err != nil {
			return mesh, err
		}
		if match {
			break
		}
	}

	connectivity, err := getConnectivityFromLine(currentLine)
	if err != nil {
		return mesh, err
	}
	mesh.connectivity[0] = connectivity
	// Get connectivity
	for i := 1; i < nFaces; i++ {
		if !scanner.Scan() {
			return mesh, scanner.Err()
		}
		connectivity, err := getConnectivityFromLine(scanner.Text())
		if err != nil {
			return mesh, err
		}
		mesh.connectivity[i] = connectivity
	}

	return mesh, nil
}

func getConnectivityFromLine(line string) (TriangleConnectivity, error) {
	// f 1//1 5//5 2//2
	splittedLine := strings.Split(line, " ")[1:]
	var connectivity TriangleConnectivity

	for i := 0; i < 3; i++ {
		temp := splittedLine[i]
		temp = strings.Split(temp, "//")[0]
		val, err := strconv.Atoi(temp)
		val--
		if err != nil {
			return connectivity, err
		}
		connectivity[i] = val

	}
	return connectivity, nil
}

func getVectorFromLine(line string) (basics.Vector3, error) {
	//v -1.00000000 -1.00000000 -1.00000000
	splittedLine := strings.Split(line, " ")[1:]
	vector := basics.ZeroVector()
	var coords [3]basics.Scalar

	for i := 0; i < 3; i++ {
		temp, err := strconv.ParseFloat(splittedLine[i], 64)
		if err != nil {
			return vector, err
		}
		coords[i] = basics.Scalar(temp)
	}
	vector.X = coords[0]
	vector.Y = coords[1]
	vector.Z = coords[2]
	return vector, nil
}

func getNumberOfElements(line string) (int, error) {
	//8 vertices
	//12 faces
	splittedLine := strings.Split(line, " ")
	return strconv.Atoi(splittedLine[0])
}
