package wasp

import (
	"syscall/js"

	"./webgl"
	dom "github.com/schabby/go-wasm-dom"
)

var jsDrawCallback, jsResizeCallback js.Func

func CreateWebGLApp(
	init func(webgl.RenderingContext),
	resize func(webgl.RenderingContext),
	draw func(webgl.RenderingContext, int)) {

	canvas := dom.FullPageCanvas()
	glDOM := canvas.JsValue().Call("getContext", "webgl2")
	gl := webgl.NewRenderingContext(glDOM)

	// call once, used to set up application
	init(gl)

	// call once manually, but also register resize event
	jsResizeCallback = js.FuncOf(func(this js.Value, inputs []js.Value) interface{} {

		dpr := js.Global().Get("window").Get("devicePixelRatio").Float()
		rect := canvas.JsValue().Call("getBoundingClientRect")
		width := int(rect.Get("width").Float() * dpr)
		height := int(rect.Get("height").Float() * dpr)
		canvas.SetWidthI(width)
		canvas.SetHeightI(height)

		gl.Width = width
		gl.Height = height

		resize(gl)
		return nil
	})
	js.Global().Set("pageResize", jsResizeCallback)
	js.Global().Call("pageResize")
	js.Global().Get("window").Call("addEventListener", "resize", jsResizeCallback)

	jsDrawCallback = js.FuncOf(func(this js.Value, inputs []js.Value) interface{} {
		draw(gl, inputs[0].Int())
		js.Global().Call("requestAnimationFrame", jsDrawCallback)
		return nil
	})
	//defer jsDrawCallback.Release()

	// start rendering cycles
	js.Global().Call("requestAnimationFrame", jsDrawCallback)
}

func LoadImage(url string, callback func(image js.Value)) {
	imageHandle := js.Global().Get("Image").New()
	imageHandle.Set("src", url)
	imageHandle.Call("addEventListener", "load", js.FuncOf(func(image js.Value, inputs []js.Value) interface{} {
		callback(image)
		return nil
	}))
}
