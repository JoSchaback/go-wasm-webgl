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
var texture js.Value
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

	vsSource := `
	attribute vec3 a_position;
	attribute vec2 a_texcoord;

	varying   vec2 v_texcoord;

	uniform   mat4 u_modelViewMatrix;
	uniform   mat4 u_projectionMatrix;

	void main() {
	  gl_Position = u_projectionMatrix * u_modelViewMatrix * vec4(a_position, 1.0);
	  v_texcoord = a_texcoord;
	}`

	fsSource := `precision mediump float;

	varying vec2 v_texcoord;

	// The texture.
	uniform sampler2D u_texture;

	void main() {
		gl_FragColor  = texture2D(u_texture, v_texcoord);
	}`

	program = wasp.NewProgram(&gl, vsSource, fsSource)
	program.AttribTypes["a_position"] = wasp.POSITION
	program.AttribTypes["a_texcoord"] = wasp.UV
	mesh := wasp.NewCube()

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
	program.UniformMatrix4fv("u_projectionMatrix", &projectionMatrix)

	texture = gl.CreateTexture()
	gl.BindTexture(webgl.TEXTURE_2D, texture)
	gl.TexImage2Duint8(webgl.TEXTURE_2D, 0, webgl.RGBA, 1, 1, 0, webgl.RGBA, webgl.UNSIGNED_BYTE, []uint8{0, 0, 255, 255})

	gl.ClearColor(0, 0, 0, 1)

	vao = wasp.NewVAO(&gl, mesh, program)

	gl.Enable(webgl.DEPTH_TEST)

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

		program.UniformMatrix4fv("u_modelViewMatrix", &viewMatrix)

		vao.DrawElements()
		js.Global().Call("requestAnimationFrame", fsDraw)
		return nil
	})

	defer fsDraw.Release()
	funcOfImageLoaded := js.FuncOf(imageLoaded)
	js.Global().Set("imageLoaded", funcOfImageLoaded)

	imageHandle := js.Global().Get("Image").New()
	imageHandle.Set("src", "crate_diff_texture.png")
	imageHandle.Call("addEventListener", "load", funcOfImageLoaded)

	js.Global().Call("requestAnimationFrame", fsDraw)

	<-done
}

func printMessage(this js.Value, inputs []js.Value) interface{} {
	fmt.Printf("Hello from within WASM")
	return nil
}

func imageLoaded(image js.Value, inputs []js.Value) interface{} {
	gl.BindTexture(webgl.TEXTURE_2D, texture)
	gl.TexImage2DImage(webgl.TEXTURE_2D, 0, webgl.RGBA, webgl.RGBA, webgl.UNSIGNED_BYTE, image)
	gl.GenerateMipmap(webgl.TEXTURE_2D)

	return nil
}
