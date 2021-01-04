package wasp

import (
	"fmt"
	"math"
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

// SubdivideFace replaces a face with three smaller faces.
func (c *Mesh) SubdivideFace(faceIndex int) {
	// take each face, subdivide it, register the new vertices and new faces and throw the old
	// face out

	//				   v0
	//				   /\
	//				  /  \
	//				 /    \
	//		 v0Tov1 /------\ v2tov0
	//			   /\      /\
	//			  /  \    /  \
	//			 /    \  /    \
	//		  v1 ------\/------ v2
	//               v1tov2

	// create new vertices at the intersection between the existing vertices
	v0tov1 := make([]float32, c.vertexSize)
	v1tov2 := make([]float32, c.vertexSize)
	v2tov0 := make([]float32, c.vertexSize)

	v0VerID := c.faces[faceIndex+0]
	v1VerID := c.faces[faceIndex+1]
	v2VerID := c.faces[faceIndex+2]

	// get existing vertices
	v0 := c.VertexByID(v0VerID)
	v1 := c.VertexByID(v1VerID)
	v2 := c.VertexByID(v2VerID)

	// the coordinates of the new vertices are in the middle of the existing vertices.
	// We interpolate ALL components intentionally since most components are interpolatable
	// such as RGB, UV, NORMAL etc.
	for i := 0; i < c.vertexSize; i++ {
		//fmt.Printf("i %d, %v ---- %v\n", i, v0tov1, v0)
		v0tov1[i] = (v0[i] + v1[i]) / 2
		v1tov2[i] = (v1[i] + v2[i]) / 2
		v2tov0[i] = (v2[i] + v0[i]) / 2
	}

	// we are setting some freaky colors, mainly for debugging
	// purposes.
	colorOffset := c.Offset(RGB)
	v0tov1[colorOffset] = 0
	v0tov1[colorOffset+1] = 1
	v0tov1[colorOffset+2] = 0

	v1tov2[colorOffset] = 1
	v1tov2[colorOffset+1] = 0
	v1tov2[colorOffset+2] = 0

	v2tov0[colorOffset] = 1
	v2tov0[colorOffset+1] = 0
	v2tov0[colorOffset+2] = 1

	v0tov1Id := uint16(c.VertexCount() + 0)
	v1tov2Id := uint16(c.VertexCount() + 1)
	v2tov0Id := uint16(c.VertexCount() + 2)

	c.vertices = append(c.vertices, v0tov1...)
	c.vertices = append(c.vertices, v1tov2...)
	c.vertices = append(c.vertices, v2tov0...)

	// lower left triangle
	c.AddFace(v1VerID, v1tov2Id, v0tov1Id)

	// upper triangle
	c.AddFace(v0VerID, v0tov1Id, v2tov0Id)

	// lower right triangle
	c.AddFace(v1tov2Id, v2VerID, v2tov0Id)

	// middle triangle
	c.AddFace(v0tov1Id, v1tov2Id, v2tov0Id)

	c.RemoveFace(faceIndex)
}

func (c *Mesh) VertexByIndex(vertexIndex int) []float32 {
	return c.vertices[vertexIndex : vertexIndex+c.vertexSize]
}

func (c *Mesh) VertexByID(vertexID uint16) []float32 {
	return c.vertices[int(vertexID)*c.vertexSize : (int(vertexID)*c.vertexSize)+c.vertexSize]
}

func (c *Mesh) VertexFromFaceId(vertexIndex uint16) []float32 {
	return c.vertices[int(vertexIndex) : int(vertexIndex)+c.vertexSize]
}

func (c *Mesh) RemoveFace(faceIndex int) {
	c.faces = append(c.faces[:faceIndex], c.faces[faceIndex+3:]...)
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
	maxVertexCount := uint16(len(c.vertices) / c.vertexSize)
	if vertexId1 >= maxVertexCount || vertexId1 < 0 {
		panic(fmt.Sprintf("vertexId1 %d is out of bounds (maxVertexCount %d)", vertexId1, c.vertexSize))
	}
	if vertexId2 >= maxVertexCount || vertexId2 < 0 {
		panic(fmt.Sprintf("vertexId2 %d is out of bounds (maxVertexCount %d)", vertexId2, c.vertexSize))
	}
	if vertexId3 >= maxVertexCount || vertexId3 < 0 {
		panic(fmt.Sprintf("vertexId3 %d is out of bounds (maxVertexCount %d)", vertexId3, c.vertexSize))
	}

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
	panic(fmt.Sprintf("mesh does not contain AttribType %v", attribType))
}

func (c *Mesh) ForEachVertex(vFunc func([]float32)) {
	for i := 0; i < len(c.vertices); i += c.vertexSize {
		vFunc(c.vertices[i:(i + c.vertexSize)])
	}
}

func (c *Mesh) ForEachVertexComponent(attribType AttribType, vFunc func([]float32)) {
	offset := c.Offset(attribType)
	for i := 0; i < len(c.vertices); i += c.vertexSize {
		vFunc(c.vertices[i+offset : (i + offset + attribType.size)])
	}
}

func (c *Mesh) ForEachFace(vFunc func([]uint16)) {
	for i := 0; i < len(c.faces); i += 3 {
		vFunc(c.faces[i:(i + 3)])
	}
}

func (c *Mesh) NormalizeAll(attrib AttribType) {
	offset := c.Offset(attrib)
	for i := offset; i < len(c.vertices); i += c.vertexSize {
		squaredSum := float64(0)
		for p := 0; p < attrib.size; p++ {
			squaredSum += float64(c.vertices[i+p] * c.vertices[i+p])
		}
		squared := float32(math.Sqrt(squaredSum))
		for p := 0; p < attrib.size; p++ {
			c.vertices[i+p] /= squared
		}
	}
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
