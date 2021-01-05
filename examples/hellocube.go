package main

import (
	"fmt"

	"../wasp"
	"../wasp/webgl"

	dom "github.com/schabby/go-wasm-dom"
	"github.com/schabby/linalg"
)

var program *wasp.Program
var vao *wasp.VertexArrayObject
var angle float32 = 0

var viewMatrix = linalg.NewMatrix4()
var projectionMatrix = linalg.NewMatrix4()

func myInit(gl webgl.RenderingContext) {
	vsSource := `attribute vec3 aVertexPosition;
	attribute vec3 aVertexColor;

	varying lowp vec3 v_rgb;

	uniform mat4 uModelViewMatrix;
	uniform mat4 uProjectionMatrix;

	void main() {
	  gl_Position = uProjectionMatrix * uModelViewMatrix * vec4(aVertexPosition, 1.0);
	  //gl_Position = vec4(aVertexPosition, 1);
	  v_rgb = aVertexColor;
	}`

	fsSource := `
		varying lowp vec3 v_rgb;

		void main() {
			gl_FragColor  = vec4(v_rgb, 1.0);
		}`

	program = wasp.NewProgram(&gl, vsSource, fsSource)
	program.Use()
	program.AttribTypes["aVertexPosition"] = wasp.POSITION
	program.AttribTypes["aVertexColor"] = wasp.RGB
	mesh := wasp.NewCube()
	gl.ClearColor(0, 0, 0, 1)

	vao = wasp.NewVAO(&gl, mesh, program)

	gl.Enable(webgl.DEPTH_TEST)

}

func resize(gl webgl.RenderingContext) {
	width := gl.Width
	height := gl.Height
	gl.Viewport(0, 0, width, height)

	dom.Log(fmt.Sprintf("canvas is %d x %d", width, height))

	projectionMatrix.Projection(45, float32(width), float32(height), 0.1, 100)
	program.UniformMatrix4fv("uProjectionMatrix", &projectionMatrix)

}

func draw(gl webgl.RenderingContext, timestamp int) {

	gl.Clear(webgl.COLOR_BUFFER_BIT | webgl.DEPTH_BUFFER_BIT)

	eye := linalg.Vector3{2, 2, 3}
	center := linalg.Vector3{0, 0, 0}
	up := linalg.Vector3{0, 0, 1}
	viewMatrix.LookAt(eye, center, up)

	rot := linalg.NewMatrix4()
	rot.Rotation(angle, linalg.Vector3{0, 0, 1})
	angle += 0.01

	viewMatrix.MultAssign(&rot)

	program.UniformMatrix4fv("uModelViewMatrix", &viewMatrix)

	vao.DrawElements()
}

func main() {
	done := make(chan struct{}, 0)

	wasp.CreateWebGLApp(myInit, resize, draw)

	<-done
}
