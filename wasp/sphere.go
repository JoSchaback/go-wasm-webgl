package wasp

import (
	"math"
)

func NewSphere(subdivideCount int) *Mesh {

	mesh := NewMesh(POSITION, NORMAL, RGB, UV)

	// http://blog.andreaskahler.com/2009/06/creating-icosphere-mesh-in-code.html

	// create 12 vertices of a icosahedron
	t := float32((1 + math.Sqrt(5)) / 2)
	mesh.AddVertex(-1, t, 0, 0, 0, 0, 1, 1, 0, 0, 0)
	mesh.AddVertex(1, t, 0, 0, 0, 0, 1, 1, 0, 0, 0)
	mesh.AddVertex(-1, -t, 0, 0, 0, 0, 1, 1, 0, 0, 0)
	mesh.AddVertex(1, -t, 0, 0, 0, 0, 1, 1, 0, 0, 0)

	mesh.AddVertex(0, -1, t, 0, 0, 0, 1, 1, 0, 0, 0)
	mesh.AddVertex(0, 1, t, 0, 0, 0, 1, 1, 0, 0, 0)
	mesh.AddVertex(0, -1, -t, 0, 0, 0, 1, 1, 0, 0, 0)
	mesh.AddVertex(0, 1, -t, 0, 0, 0, 1, 1, 0, 0, 0)

	mesh.AddVertex(t, 0, -1, 0, 0, 0, 1, 1, 0, 0, 0)
	mesh.AddVertex(t, 0, 1, 0, 0, 0, 1, 1, 0, 0, 0)
	mesh.AddVertex(-t, 0, -1, 0, 0, 0, 1, 1, 0, 0, 0)
	mesh.AddVertex(-t, 0, 1, 0, 0, 0, 1, 1, 0, 0, 0)

	mesh.NormalizeAll(POSITION)

	mesh.AddFace(0, 11, 5)
	mesh.AddFace(0, 5, 1)
	mesh.AddFace(0, 1, 7)
	mesh.AddFace(0, 7, 10)
	mesh.AddFace(0, 10, 11)

	mesh.AddFace(1, 5, 9)
	mesh.AddFace(5, 11, 4)
	mesh.AddFace(11, 10, 2)
	mesh.AddFace(10, 7, 6)
	mesh.AddFace(7, 1, 8)

	mesh.AddFace(3, 9, 4)
	mesh.AddFace(3, 4, 2)
	mesh.AddFace(3, 2, 6)
	mesh.AddFace(3, 6, 8)
	mesh.AddFace(3, 8, 9)

	mesh.AddFace(4, 9, 5)
	mesh.AddFace(2, 4, 11)
	mesh.AddFace(6, 2, 10)
	mesh.AddFace(8, 6, 7)
	mesh.AddFace(9, 8, 1)

	for i := 0; i < subdivideCount; i++ {
		// we subdivide the face at the beginning of the
		// face slice. SubDivide will append (at the end) new Vertices
		// to the end of the mesh.vertices slice and three
		// additional faces to the mesh.faces slice.
		// Then, it will cut out the subvidided face, which
		// happens to be the first, such that the next
		// face to be subdivided is again at index 0.
		oldFacesCount := mesh.FaceCount()
		for j := 0; j < oldFacesCount; j++ {
			mesh.SubdivideFace(0)
		}
		// after each subdivision run, we need to adjust
		// the length of the (new) vertices because subdivision
		// only takes the mid points between triangles.
		mesh.NormalizeAll(POSITION)
	}

	normalOffset := mesh.Offset(NORMAL)
	posOffset := mesh.Offset(POSITION)

	mesh.ForEachVertex(func(vertex []float32) {
		vertex[normalOffset+0] = vertex[posOffset+0]
		vertex[normalOffset+1] = vertex[posOffset+1]
		vertex[normalOffset+2] = vertex[posOffset+2]
	})

	return mesh
}
