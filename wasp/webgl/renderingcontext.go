package webgl

import (
	"fmt"
	"reflect"
	"runtime"
	"syscall/js"
	"unsafe"

	"github.com/schabby/linalg"
)

// RenderingContext is a wrapper around the js.Value handle that holds the instance of the javascript-bound WebGL2RenderingContext.
type RenderingContext struct {
	//Loaded  bool
	Js js.Value
	//version uint

	// Constant values
}

func NewRenderingContext(value js.Value) RenderingContext {
	return RenderingContext{Js: value}
}

func EmptyRenderingContext() RenderingContext {
	return RenderingContext{}
}

func sliceToByteSlice(s interface{}) []byte {
	switch s := s.(type) {
	case []int8:
		h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		return *(*[]byte)(unsafe.Pointer(h))
	case []int16:
		h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		h.Len *= 2
		h.Cap *= 2
		return *(*[]byte)(unsafe.Pointer(h))
	case []int32:
		h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		h.Len *= 4
		h.Cap *= 4
		return *(*[]byte)(unsafe.Pointer(h))
	case []int64:
		h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		h.Len *= 8
		h.Cap *= 8
		return *(*[]byte)(unsafe.Pointer(h))
	case []uint8:
		return s
	case []uint16:
		h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		h.Len *= 2
		h.Cap *= 2
		return *(*[]byte)(unsafe.Pointer(h))
	case []uint32:
		h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		h.Len *= 4
		h.Cap *= 4
		return *(*[]byte)(unsafe.Pointer(h))
	case []uint64:
		h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		h.Len *= 8
		h.Cap *= 8
		return *(*[]byte)(unsafe.Pointer(h))
	case []float32:
		h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		h.Len *= 4
		h.Cap *= 4
		return *(*[]byte)(unsafe.Pointer(h))
	case []float64:
		h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
		h.Len *= 8
		h.Cap *= 8
		return *(*[]byte)(unsafe.Pointer(h))
	default:
		panic(fmt.Sprintf("jsutil: unexpected value at sliceToBytesSlice: %T", s))
	}
}

func SliceToFloat32Array(s []float32) js.Value {

	a := js.Global().Get("Uint8Array").New(len(s) * 4)
	js.CopyBytesToJS(a, sliceToByteSlice(s))
	runtime.KeepAlive(s)
	buf := a.Get("buffer")
	return js.Global().Get("Float32Array").New(buf, a.Get("byteOffset"), a.Get("byteLength").Int()/4)
}

func SliceToTypedArray(s interface{}) js.Value {
	if s == nil {
		return js.Null()
	}

	switch s := s.(type) {
	case []int8:
		a := js.Global().Get("Uint8Array").New(len(s))
		js.CopyBytesToJS(a, sliceToByteSlice(s))
		runtime.KeepAlive(s)
		buf := a.Get("buffer")
		return js.Global().Get("Int8Array").New(buf, a.Get("byteOffset"), a.Get("byteLength"))
	case []int16:
		a := js.Global().Get("Uint8Array").New(len(s) * 2)
		js.CopyBytesToJS(a, sliceToByteSlice(s))
		runtime.KeepAlive(s)
		buf := a.Get("buffer")
		return js.Global().Get("Int16Array").New(buf, a.Get("byteOffset"), a.Get("byteLength").Int()/2)
	case []int32:
		a := js.Global().Get("Uint8Array").New(len(s) * 4)
		js.CopyBytesToJS(a, sliceToByteSlice(s))
		runtime.KeepAlive(s)
		buf := a.Get("buffer")
		return js.Global().Get("Int32Array").New(buf, a.Get("byteOffset"), a.Get("byteLength").Int()/4)
	case []uint8:
		a := js.Global().Get("Uint8Array").New(len(s))
		js.CopyBytesToJS(a, s)
		runtime.KeepAlive(s)
		return a
	case []uint16:
		a := js.Global().Get("Uint8Array").New(len(s) * 2)
		js.CopyBytesToJS(a, sliceToByteSlice(s))
		runtime.KeepAlive(s)
		buf := a.Get("buffer")
		return js.Global().Get("Uint16Array").New(buf, a.Get("byteOffset"), a.Get("byteLength").Int()/2)
	case []uint32:
		a := js.Global().Get("Uint8Array").New(len(s) * 4)
		js.CopyBytesToJS(a, sliceToByteSlice(s))
		runtime.KeepAlive(s)
		buf := a.Get("buffer")
		return js.Global().Get("Uint32Array").New(buf, a.Get("byteOffset"), a.Get("byteLength").Int()/4)
	case []float32:
		a := js.Global().Get("Uint8Array").New(len(s) * 4)
		js.CopyBytesToJS(a, sliceToByteSlice(s))
		runtime.KeepAlive(s)
		buf := a.Get("buffer")
		return js.Global().Get("Float32Array").New(buf, a.Get("byteOffset"), a.Get("byteLength").Int()/4)
	case []float64:
		a := js.Global().Get("Uint8Array").New(len(s) * 8)
		js.CopyBytesToJS(a, sliceToByteSlice(s))
		runtime.KeepAlive(s)
		buf := a.Get("buffer")
		return js.Global().Get("Float64Array").New(buf, a.Get("byteOffset"), a.Get("byteLength").Int()/8)
	default:
		panic(fmt.Sprintf("jsutil: unexpected value at SliceToTypedArray: %T", s))
	}
}

func (c *RenderingContext) AttachShader(program js.Value, shader js.Value) {
	c.Js.Call("attachShader", program, shader)
}

func (c *RenderingContext) LinkProgram(program js.Value) {
	c.Js.Call("linkProgram", program)
}

func (c *RenderingContext) Clear(mask uint32) {
	c.Js.Call("clear", mask)
}

func (c *RenderingContext) ClearColor(r, g, b, a float32) {
	c.Js.Call("clearColor", r, g, b, a)
}

func (c *RenderingContext) ClearDepth(depth float32) {
	c.Js.Call("clearDepth", depth)
}

func (c *RenderingContext) CreateProgram() js.Value {
	return c.Js.Call("createProgram")
}

func (c *RenderingContext) CreateShader(glEnum uint32) js.Value {
	return c.Js.Call("createShader", glEnum)
}

func (c *RenderingContext) ShaderSource(shaderHandle js.Value, code string) {
	c.Js.Call("shaderSource", shaderHandle, code)
}

func (c *RenderingContext) CompileShader(shader js.Value) {
	c.Js.Call("compileShader", shader)
}

func (c *RenderingContext) DrawElements(mode uint32, count uint32, itype uint32, offset uint32) {
	c.Js.Call("drawElements", mode, count, itype, offset)
}

func (c *RenderingContext) GetShaderParameter(shader js.Value, flag uint32) js.Value {
	return c.Js.Call("getShaderParameter", shader, flag)
}

func (c *RenderingContext) GetShaderInfoLog(shader js.Value) js.Value {
	return c.Js.Call("getShaderInfoLog", shader)
}

func (c *RenderingContext) GetProgramParameter(shader js.Value, flag uint32) js.Value {
	return c.Js.Call("getProgramParameter", shader, flag)
}

func (c *RenderingContext) GetProgramInfoLog(program js.Value) js.Value {
	return c.Js.Call("getProgramInfoLog", program)
}

func (c *RenderingContext) CreateVertexArray() js.Value {
	return c.Js.Call("createVertexArray")
}

func (c *RenderingContext) BindVertexArray(vao js.Value) js.Value {
	return c.Js.Call("bindVertexArray", vao)
}

func (c *RenderingContext) GetActiveAttrib(shader js.Value, index uint32) (string, int, int) {
	a := c.Js.Call("getActiveAttrib", shader, index)
	//fmt.Printf("we got something isNull? %t\n", a.IsNull())
	name := a.Get("name").String()
	glType := a.Get("type").Int()
	size := a.Get("size").Int()
	return name, size, glType
}

func (c *RenderingContext) GetAttribLocation(shader js.Value, name string) js.Value {
	return c.Js.Call("getAttribLocation", shader, name)
}

func (c *RenderingContext) GetActiveUniform(shader js.Value, enumT uint32) js.Value {
	return c.Js.Call("getActiveUniform", shader, enumT)
}

/*
func (c *RenderingContext) GetUniformLocation(shader js.Value, name string) js.Value {
	return c.Js.Call("getUniformLocation", shader, name)
}*/

func (c *RenderingContext) GetUniformLocation(shader js.Value, name string) js.Value {
	return c.Js.Call("getUniformLocation", shader, name)
}

func (c *RenderingContext) DeleteShader(shader js.Value) {
	c.Js.Call("deleteShader", shader)
}

func (c *RenderingContext) UseProgram(shader js.Value) {
	c.Js.Call("useProgram", shader)
}

// EnableVertexAttribArray turns on the generic vertex
// attribute array at the specified index into the list
// of attribute arrays.
func (c *RenderingContext) EnableVertexAttribArray(index uint32) {
	c.Js.Call("enableVertexAttribArray", index)
}

func (c *RenderingContext) Enable(flags uint32) {
	c.Js.Call("enable", flags)
}

func (c *RenderingContext) CreateBuffer() js.Value {
	return c.Js.Call("createBuffer")
}

func (c *RenderingContext) BindBuffer(target uint32, buffer js.Value) {
	c.Js.Call("bindBuffer", target, buffer)
}

func (c *RenderingContext) BufferDataUInt16(target uint32, srcData []uint16, usage uint32) {
	//c.Js.Call("bufferData", target /* js.TypedArrayOf(srcData)*/, srcData, usage)
	c.Js.Call("bufferData", target, SliceToTypedArray(srcData), usage)
}

func (c *RenderingContext) UniformMatrix4fv(location js.Value, transpose bool, data []float32, srcOffset int, srcLength int) {
	c.Js.Call("uniformMatrix4fv", location, transpose, SliceToFloat32Array(data), srcOffset, srcLength)
}

func (c *RenderingContext) UniformMatrix4fvCustom(location js.Value, m *linalg.Matrix4) {
	js.Global().Call("uniformMatrix4fv", c.Js, location, m.M_0_0, m.M_0_1, m.M_0_2, m.M_0_3, m.M_1_0, m.M_1_1, m.M_1_2, m.M_1_3, m.M_2_0, m.M_2_1, m.M_2_2, m.M_2_3, m.M_3_0, m.M_3_1, m.M_3_2, m.M_3_3)
	//c.Js.Call("uniformMatrix4fv", location, transpose, SliceToFloat32Array(data), srcOffset, srcLength)
}

func (c *RenderingContext) UniformMatrix3fvCustom(location js.Value, m *linalg.Matrix3) {
	js.Global().Call("uniformMatrix3fv", c.Js, location, m.M_0_0, m.M_0_1, m.M_0_2, m.M_1_0, m.M_1_1, m.M_1_2, m.M_2_0, m.M_2_1, m.M_2_2)
	//c.Js.Call("uniformMatrix4fv", location, transpose, SliceToFloat32Array(data), srcOffset, srcLength)
}

func (c *RenderingContext) Uniform3f(location js.Value, f1, f2, f3 float32) {
	//c.Js.Call("bufferData", target /* js.TypedArrayOf(srcData)*/, srcData, usage)
	c.Js.Call("uniform3f", location, f1, f2, f3)
}

func (c *RenderingContext) BufferData(target uint32, srcData []float32, usage uint32) {
	//c.Js.Call("bufferData", target /* js.TypedArrayOf(srcData)*/, srcData, usage)
	c.Js.Call("bufferData", target, SliceToFloat32Array(srcData), usage)
}

func (c *RenderingContext) VertexAttribPointer(index uint32, size int, glenumType uint32, normalized bool, stride uint32, offset uint32) {
	c.Js.Call("vertexAttribPointer", index, size, glenumType, normalized, stride, offset)
}

func (c *RenderingContext) Viewport(x int, y int, width int, height int) {
	c.Js.Call("viewport", x, y, width, height)
}

func (c *RenderingContext) CheckForError() {
	msg := c.GetError()
	if msg != "" {
		js.Global().Get("console").Call("error", msg)
	}
}

func (c *RenderingContext) GetError() string {
	errorJs := c.Js.Call("getError")

	switch uint32(errorJs.Int()) {
	case NO_ERROR:
		return ""
	case INVALID_ENUM:
		return "invalid enum"
	case INVALID_VALUE:
		return "invalid value"
	case INVALID_OPERATION:
		return "invalid operation"
	case INVALID_FRAMEBUFFER_OPERATION:
		return "invalid framebuffer operation"
	case OUT_OF_MEMORY:
		return "out of memory"
	case CONTEXT_LOST_WEBGL:
		return "context lost webgl"
	}

	return "unknown error"
}
