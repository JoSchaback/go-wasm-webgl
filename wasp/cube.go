package wasp

func NewCube() *Mesh {
	mesh := NewMesh(POSITION, NORMAL, RGB, UV)
	// floor
	mesh.AddVertex( /* pos */ 0, 0, 0 /* normal */, 0, 0, -1 /* rgb */, 1, 0, 0 /* uv */, 0, 0)
	mesh.AddVertex(0, 1, 0, 0, 0, -1, 1, 0, 0, 1, 0)
	mesh.AddVertex(1, 1, 0, 0, 0, -1, 1, 0, 0, 1, 1)
	mesh.AddVertex(1, 0, 0, 0, 0, -1, 1, 0, 0, 0, 1)
	// top
	mesh.AddVertex(0, 0, 1, 0, 0, 1, 1, 0, 1, 0, 0)
	mesh.AddVertex(1, 0, 1, 0, 0, 1, 1, 0, 1, 0, 1)
	mesh.AddVertex(1, 1, 1, 0, 0, 1, 1, 0, 1, 1, 1)
	mesh.AddVertex(0, 1, 1, 0, 0, 1, 1, 0, 1, 1, 0)

	// left side
	mesh.AddVertex(0, 0, 0, -1, 0, 0, 0, 1, 1, 0, 0)
	mesh.AddVertex(0, 0, 1, -1, 0, 0, 0, 1, 1, 1, 0)
	mesh.AddVertex(0, 1, 1, -1, 0, 0, 0, 1, 1, 1, 1)
	mesh.AddVertex(0, 1, 0, -1, 0, 0, 0, 1, 1, 0, 1)

	// right side
	mesh.AddVertex(1, 0, 0, 1, 0, 0, 1, 1, 0, 0, 0)
	mesh.AddVertex(1, 1, 0, 1, 0, 0, 1, 1, 0, 0, 1)
	mesh.AddVertex(1, 1, 1, 1, 0, 0, 1, 1, 0, 1, 1)
	mesh.AddVertex(1, 0, 1, 1, 0, 0, 1, 1, 0, 1, 0)

	// front
	mesh.AddVertex(0, 0, 0, 0, -1, 0, 0, 0, 1, 0, 0)
	mesh.AddVertex(1, 0, 0, 0, -1, 0, 0, 0, 1, 0, 1)
	mesh.AddVertex(1, 0, 1, 0, -1, 0, 0, 0, 1, 1, 1)
	mesh.AddVertex(0, 0, 1, 0, -1, 0, 0, 0, 1, 1, 0)

	// back
	mesh.AddVertex(0, 1, 1, 0, 1, 0, 0, 1, 0, 1, 0)
	mesh.AddVertex(1, 1, 1, 0, 1, 0, 0, 1, 0, 1, 1)
	mesh.AddVertex(1, 1, 0, 0, 1, 0, 0, 1, 0, 0, 1)
	mesh.AddVertex(0, 1, 0, 0, 1, 0, 0, 1, 0, 0, 0)

	for i := 0; i < mesh.VertexCount(); i += 4 {
		mesh.AddFace(uint16(i+0), uint16(i+1), uint16(i+2))
		mesh.AddFace(uint16(i+2), uint16(i+3), uint16(i+0))
	}

	mesh.Translate(-0.5, -0.5, -0.5)
	return mesh
}
