package wasp

import (
	"syscall/js"

	"./webgl"
)

type VertexArrayObject struct {
	handle     js.Value
	facesCount int
	gl         *webgl.RenderingContext
}

func NewVAO(gl *webgl.RenderingContext, mesh *Mesh, program *Program) *VertexArrayObject {
	vao := gl.CreateVertexArray()
	gl.BindVertexArray(vao)
	vertexBuffer := mesh.UploadVerticesToNewBuffer(gl)
	gl.BindBuffer(webgl.ARRAY_BUFFER, vertexBuffer)
	mesh.VertexAttribPointers(gl, program)
	faceBuffer := mesh.UploadFaceIndicesToNewBuffer(gl)
	gl.BindBuffer(webgl.ELEMENT_ARRAY_BUFFER, faceBuffer)
	return &VertexArrayObject{vao, mesh.FaceCount(), gl}
}

func (vao *VertexArrayObject) DrawElements() {
	vao.gl.BindVertexArray(vao.handle)
	vao.gl.DrawElements(webgl.TRIANGLES, uint32(vao.facesCount*3), webgl.UNSIGNED_SHORT, 0)
}
