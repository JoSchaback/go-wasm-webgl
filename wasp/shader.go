package wasp

import (
	"fmt"
	"syscall/js"

	"./webgl"
	"github.com/schabby/linalg"
)

// Program represents a shader program
type Program struct {
	handle      *js.Value
	gl          *webgl.RenderingContext
	attribs     map[string]uint32
	AttribTypes map[string]AttribType
	uniforms    map[string]js.Value
}

// Use binds this shader program.
func (c *Program) Use() {
	c.gl.UseProgram(*c.handle)
}

var tmpVec = make([]float32, 16)

func (c *Program) UniformMatrix4fv(name string, m *linalg.Matrix4) {
	loc := c.uniforms[name]
	//m.Copy(tmpVec)
	//c.gl.UniformMatrix4fv(loc, false, tmpVec, 0, 0)
	c.gl.UniformMatrix4fvCustom(loc, m)
}

// NewProgram creates a new Program instance by compiling the provided shaders
func NewProgram(gl *webgl.RenderingContext, vertexShaderSource string, fragmentShaderSource string) *Program {
	vertexShader := createShader(gl, webgl.VERTEX_SHADER, vertexShaderSource)
	fragmentShader := createShader(gl, webgl.FRAGMENT_SHADER, fragmentShaderSource)
	gl.CheckForError()
	prgmHandle := gl.CreateProgram()

	gl.AttachShader(prgmHandle, vertexShader)
	gl.AttachShader(prgmHandle, fragmentShader)
	gl.LinkProgram(prgmHandle)

	gl.CheckForError()

	if !gl.GetProgramParameter(prgmHandle, webgl.LINK_STATUS).Bool() {
		fmt.Println("unable to create shader program")
		fmt.Println(gl.GetProgramInfoLog(prgmHandle).String())
	}

	activeAttribCount := gl.GetProgramParameter(prgmHandle, webgl.ACTIVE_ATTRIBUTES).Int()

	program := Program{
		handle:      &prgmHandle,
		gl:          gl,
		attribs:     make(map[string]uint32),
		uniforms:    make(map[string]js.Value),
		AttribTypes: make(map[string]AttribType)}

	for index := 0; index < activeAttribCount; index++ {
		name, size, glType := gl.GetActiveAttrib(prgmHandle, uint32(index))
		fmt.Printf("found active attrib %s of size %d at index %d with type %d\n", name, size, index, glType)
		program.attribs[name] = uint32(index)
	}

	activeUniformCount := gl.GetProgramParameter(prgmHandle, webgl.ACTIVE_UNIFORMS).Int()

	for index := 0; index < activeUniformCount; index++ {
		unform := gl.GetActiveUniform(prgmHandle, uint32(index))
		name := unform.Get("name").String()
		size := unform.Get("size").Int()
		glType := unform.Get("type").Int()
		fmt.Printf("found active uniform %s of size %d at index %d with type %d\n", name, size, index, glType)
		program.uniforms[name] = gl.GetUniformLocation(prgmHandle, name)
	}

	return &program
}

func createShader(gl *webgl.RenderingContext, glType uint32, code string) js.Value {
	shader := gl.CreateShader(glType)
	gl.ShaderSource(shader, code)
	gl.CompileShader(shader)

	compileStatus := gl.GetShaderParameter(shader, webgl.COMPILE_STATUS)

	if !compileStatus.Bool() {
		compilerLog := gl.GetShaderInfoLog(shader)
		fmt.Printf("Error while compiling shader: %s", compilerLog)
		gl.DeleteShader(shader)
	}

	return shader
}
