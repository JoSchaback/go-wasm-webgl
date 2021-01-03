package main

import (
	"fmt"
	"syscall/js"

	"../wasp"
	"../wasp/webgl"

	dom "github.com/schabby/go-wasm-dom"
	"github.com/schabby/linalg"
)

var gl webgl.RenderingContext

var program *wasp.Program
var vao *wasp.VertexArrayObject
var angle float32 = 0
var fsDraw js.Func

var viewMatrix = linalg.NewMatrix4()
var projectionMatrix = linalg.NewMatrix4()

func main() {
	canvas := dom.FullPageCanvas()

	done := make(chan struct{}, 0)

	glDOM := canvas.JsValue().Call("getContext", "webgl2")
	gl = webgl.NewRenderingContext(glDOM)

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
	program.AttribTypes["aVertexPosition"] = wasp.POSITION
	program.AttribTypes["aVertexColor"] = wasp.RGB

	mesh := wasp.NewSphere(2)

	dpr := js.Global().Get("window").Get("devicePixelRatio").Float()
	rect := canvas.JsValue().Call("getBoundingClientRect")
	width := int(rect.Get("width").Float() * dpr)
	height := int(rect.Get("height").Float() * dpr)
	canvas.SetWidthI(width)
	canvas.SetHeightI(height)
	gl.Viewport(0, 0, width, height)

	dom.Log(fmt.Sprintf("canvas is %d x %d", width, height))

	projectionMatrix.Projection(45, float32(width), float32(height), 0.1, 100)
	program.Use()
	program.UniformMatrix4fv("uProjectionMatrix", &projectionMatrix)

	gl.ClearColor(0, 0, 0, 1)

	vao = wasp.NewVAO(&gl, mesh, program)

	gl.Enable(webgl.DEPTH_TEST)

	//js.Global().Set("printMessage", js.FuncOf(printMessage))

	fsDraw = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		//now := args[0].Float()
		gl.Clear(webgl.COLOR_BUFFER_BIT | webgl.DEPTH_BUFFER_BIT)
		//dom.Log("hinter Clear()...")

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
		js.Global().Call("requestAnimationFrame", fsDraw)
		return nil
	})

	defer fsDraw.Release()

	js.Global().Call("requestAnimationFrame", fsDraw)

	<-done
}

func printMessage(this js.Value, inputs []js.Value) interface{} {
	fmt.Printf("Hello from within WASM")
	return nil
}
