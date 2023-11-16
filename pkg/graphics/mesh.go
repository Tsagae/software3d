package graphics

import (
	"bufio"
	"github.com/tsagae/software3d/pkg/basics"
	"os"
	"strconv"
	"strings"
)

type VertexAttributes struct {
	position basics.Vector3
	color    basics.Vector3 //Range 0-65535
	normal   basics.Vector3
}

type TriangleConnectivity [3]int

type Mesh struct {
	geometry     []VertexAttributes
	connectivity []TriangleConnectivity
}

// Constructors
func NewMesh(geometry []VertexAttributes, connectivity []TriangleConnectivity) Mesh {
	return Mesh{geometry, connectivity} //should copy the slices
}

func NewEmpyMesh() Mesh {
	return Mesh{nil, nil}
}

// if ignoreMeshNormas is true the normals are the ones of the surface of the triangles
func (m *Mesh) GetTriangles(ignoreMeshNormals bool) []Triangle {
	if ignoreMeshNormals {
		return m.getTrianglesWithoutNormals()
	}
	return m.getTrianglesWithNormals()
}

func (m *Mesh) getTrianglesWithNormals() []Triangle {
	triangles := make([]Triangle, len(m.connectivity))
	for i, v := range m.connectivity {
		triangles[i] = NewTriangleWithNormals([3]basics.Vector3{m.geometry[v[0]].position, m.geometry[v[1]].position, m.geometry[v[2]].position}, [3]basics.Vector3{m.geometry[v[0]].color, m.geometry[v[1]].color, m.geometry[v[2]].color}, [3]basics.Vector3{m.geometry[v[0]].normal, m.geometry[v[1]].normal, m.geometry[v[2]].normal})
	}
	return triangles
}

func (m *Mesh) getTrianglesWithoutNormals() []Triangle {
	triangles := make([]Triangle, len(m.connectivity))
	for i, v := range m.connectivity {
		triangleNormal := computeNormalFromVertices(m.geometry[v[0]].position, m.geometry[v[1]].position, m.geometry[v[2]].position)
		triangles[i] = NewTriangleWithNormals([3]basics.Vector3{m.geometry[v[0]].position, m.geometry[v[1]].position, m.geometry[v[2]].position}, [3]basics.Vector3{m.geometry[v[0]].color, m.geometry[v[1]].color, m.geometry[v[2]].color}, [3]basics.Vector3{triangleNormal, triangleNormal, triangleNormal})
	}
	return triangles
}

func ReadMeshFromFile(fileName string, color basics.Vector3) (Mesh, error) {
	var line string
	var nVertices, nFaces int
	mesh := NewEmpyMesh()

	readFile, err := os.Open(fileName)

	if err != nil {
		return mesh, err
	}

	fileScanner := bufio.NewScanner(readFile)

	// Skip first lines
	for i := 0; i < 4; i++ {
		if !fileScanner.Scan() {
			return mesh, fileScanner.Err()
		}
	}

	// Get vertices and faces length
	line = fileScanner.Text()
	splittedLine := strings.Split(line, ",")

	verticesString := splittedLine[0][1:]
	facesString := splittedLine[1][1:]

	nVertices, err = getNumberOfElements(verticesString)
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
		if !fileScanner.Scan() {
			return mesh, fileScanner.Err()
		}

		vertex, err := getVectorFromLine(fileScanner.Text())
		if err != nil {
			return mesh, err
		}
		mesh.geometry[i].position = vertex
		mesh.geometry[i].color = color //hardcoded vertex color
	}

	// Get normals
	for i := 0; i < nVertices; i++ {
		if !fileScanner.Scan() {
			return mesh, fileScanner.Err()
		}

		vertex, err := getVectorFromLine(fileScanner.Text())
		if err != nil {
			return mesh, err
		}
		//vertex.ThisNormalize() // normalization for normals, shouldn't be needed
		mesh.geometry[i].normal = vertex
	}

	// Skip 3 lines
	for i := 0; i < 3; i++ {
		if !fileScanner.Scan() {
			return mesh, fileScanner.Err()
		}
	}

	// Get connectivity
	for i := 0; i < nFaces; i++ {
		if !fileScanner.Scan() {
			return mesh, fileScanner.Err()
		}

		connectivity, err := getConnectivityFromLine(fileScanner.Text())
		if err != nil {
			return mesh, err
		}
		mesh.connectivity[i] = connectivity
	}

	readFile.Close()

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
