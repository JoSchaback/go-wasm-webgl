package wasp

import (
	"fmt"
	"syscall/js"

	"./webgl"
)

type Mesh struct {
	vertices    []float32
	faces       []uint16
	vertexSize  int
	attribTypes []AttribType
}

func NewMesh(attribTypes ...AttribType) *Mesh {

	vertexSize := 0

	for _, t := range attribTypes {
		vertexSize += int(t.size)
	}

	mesh := Mesh{
		attribTypes: attribTypes,
		vertexSize:  vertexSize,
		vertices:    make([]float32, 0, int(vertexSize*3*20)),
		faces:       make([]uint16, 0, 20),
	}

	return &mesh
}

func (c *Mesh) AddVertex(vertex ...float32) uint16 {
	if len(vertex) != c.vertexSize {
		panic(fmt.Sprintf("AddVertex: %d parameters provided, but vertexSize is %d", len(vertex), c.vertexSize))
	}
	// append vertex data to mesh vertices
	c.vertices = append(c.vertices, vertex...)

	return uint16(c.VertexCount() - 1)
}

func (c *Mesh) VertexCount() int {
	return len(c.vertices) / c.vertexSize
}

func (c *Mesh) FaceCount() int {
	return len(c.faces) / 3
}

func (c *Mesh) AddFace(vertexId1 uint16, vertexId2 uint16, vertexId3 uint16) {
	c.faces = append(c.faces, vertexId1, vertexId2, vertexId3)
}

func (m *Mesh) Translate(x float32, y float32, z float32) {
	offset := 0
	for _, t := range m.attribTypes {
		if t == POSITION {
			break
		} else {
			offset += int(t.size)
		}
	}
	for i := offset; i < len(m.vertices); i += m.vertexSize {
		m.vertices[i+0] += x
		m.vertices[i+1] += y
		m.vertices[i+2] += z
	}
}

func (c *Mesh) Offset(attribType AttribType) int {
	offset := 0
	for _, t := range c.attribTypes {
		if t == attribType {
			return offset
		}
		offset += int(t.size)
	}
	return 0
}

func (c *Mesh) UploadVerticesToNewBuffer(gl *webgl.RenderingContext) js.Value {
	buffer := gl.CreateBuffer()
	if buffer.IsNull() {
		panic("cound not create buffer")
	}
	gl.BindBuffer(webgl.ARRAY_BUFFER, buffer)
	gl.BufferData(webgl.ARRAY_BUFFER, c.vertices, webgl.STATIC_DRAW)
	gl.CheckForError()
	return buffer
}

func (c *Mesh) UploadFaceIndicesToNewBuffer(gl *webgl.RenderingContext) js.Value {
	buffer := gl.CreateBuffer()
	if buffer.IsNull() {
		panic("cound not create buffer")
	}
	gl.BindBuffer(webgl.ELEMENT_ARRAY_BUFFER, buffer)
	gl.BufferDataUInt16(webgl.ELEMENT_ARRAY_BUFFER, c.faces, webgl.STATIC_DRAW)
	gl.CheckForError()
	return buffer
}

func (c *Mesh) VertexAttribPointers(gl *webgl.RenderingContext, program *Program) {
	//fmt.Printf("vertexSize: %d, stride %d\n", c.vertexSize, c.vertexSize*4)
	for name, attribType := range program.AttribTypes {
		loc := program.attribs[name]
		gl.EnableVertexAttribArray(loc)
		//fmt.Printf("loc %d, size %d, offset %d\n", loc, attribType.size, c.Offset(attribType))
		gl.VertexAttribPointer(
			loc,
			int(attribType.size),
			webgl.FLOAT,
			false,
			uint32(c.vertexSize*4),
			uint32(c.Offset(attribType)*4))
	}

}
