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
var vao1, vao2 *wasp.VertexArrayObject
var angle float32 = 0
var fsDraw js.Func

var viewMatrix = linalg.NewMatrix4()
var viewModelMatrix = linalg.NewMatrix4()
var projectionMatrix = linalg.NewMatrix4()

func main() {
	canvas := dom.FullPageCanvas()

	done := make(chan struct{}, 0)

	glDOM := canvas.JsValue().Call("getContext", "webgl2")
	gl = webgl.NewRenderingContext(glDOM)

	vsSource := `
	attribute    vec3 a_position;
	attribute    vec3 a_color;
	attribute    vec3 a_normal;

	varying      vec3 v_color;
	varying      vec3 v_normalInViewSpace;
	varying      vec3 v_positionInViewSpace;

	uniform      mat4 u_modelViewMatrix;
	uniform      mat4 u_projectionMatrix;
	uniform      mat3 u_normalMatrix;

	void main() {
	  // normalMatrix is computed in view spaces, meaning that it maps normals into
	  // view space as well. Lighting is usually computed in view space.
	  v_normalInViewSpace   = normalize( u_normalMatrix * a_normal);

	  // calculate vertex position in camera space (ie. after model and view transformations)
	  vec4 posInViewSpace4  = u_modelViewMatrix * vec4(a_position, 1.0);
	  v_positionInViewSpace = vec3( posInViewSpace4 );

	  // vertex colors are interpolated in fragement shader
	  v_color               = a_color;

	  // compute vertex position in screen (or "clip") space
	  gl_Position = u_projectionMatrix * posInViewSpace4;
	}`

	fsSource := `
	precision mediump float;

	varying lowp vec3 v_color;
	varying      vec3 v_normalInViewSpace;
	varying      vec3 v_positionInViewSpace;

	uniform      vec3 u_lightDirInViewSpace; // points towards light source

	void main() {

		// since the normals are interpolatde between pixels, we need
		// to normalize it every time
		vec3 normal = normalize(v_normalInViewSpace);

		// compute cosine (all vectors are unit vectors) between normal and
		// light direction in camera space. 
		float diffuse = max(dot(normal, -u_lightDirInViewSpace), 0.0);
		
		// specular term is 0 by default
		float specular = 0.0;

		// a specular reflection can only occure within the active area of diffuse
		// reflections
		if( diffuse > 0.0 ) {
			// Reflected light vector around normal in view space. The vector R will
			// be used to measure how much reflective light will be seen
			vec3 R = reflect(u_lightDirInViewSpace, normal);      
			
			// Vector to viewer	
			vec3 V = normalize(-v_positionInViewSpace); 
			
			// Compute the specular term as the cosine between reflection and 
			// viewer
    		float specAngle = max(dot(R, V), 0.0);
    		specular = pow(specAngle, 30.0);		
		}
		
		float ambient = 0.15;
		vec3 result   = (ambient + diffuse + specular) * v_color;
		gl_FragColor  = vec4(result, 1.0);		
	}`

	program = wasp.NewProgram(&gl, vsSource, fsSource)
	program.AttribTypes["a_position"] = wasp.POSITION
	program.AttribTypes["a_color"] = wasp.RGB
	program.AttribTypes["a_normal"] = wasp.NORMAL

	mesh := wasp.NewSphere(2)

	// set RGB value of all vertices to blue

	mesh.ForEachVertexComponent(wasp.RGB, func(rgb []float32) {
		rgb[0] = 0.2
		rgb[1] = 0.2
		rgb[2] = 0.8
	})
	vao1 = wasp.NewVAO(&gl, mesh, program)

	mesh = wasp.NewSphere(2)
	mesh.ForEachVertexComponent(wasp.RGB, func(rgb []float32) {
		rgb[0] = 0.2
		rgb[1] = 0.8
		rgb[2] = 0.2
	})
	mesh.ForEachVertexComponent(wasp.POSITION, func(pos []float32) {
		pos[0] -= 3
	})
	vao2 = wasp.NewVAO(&gl, mesh, program)

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

	gl.ClearColor(0, 0, 0, 1)

	gl.Enable(webgl.DEPTH_TEST)

	eye := linalg.Vector3{3, 3, 4}
	center := linalg.Vector3{0, 0, 0}
	up := linalg.Vector3{0, 0, 1}
	viewMatrix.LookAt(eye, center, up)

	// normal matrix is only inverted and transposed view matrix.
	// Thus we only need to update it when we update the view matrix.
	lightDirection := linalg.Vector3{0, 0, -100}
	//dom.Logf("1 lightDirection: %v", lightDirection)
	viewMatrix.MultVector3WriteBack(&lightDirection, 1)
	lightDirection.Normalize()
	//dom.Logf("2 lightDirection: %v", lightDirection)
	program.UniformVector3("u_lightDirInViewSpace", lightDirection) // in camera space!

	fsDraw = js.FuncOf(func(this js.Value, args []js.Value) interface{} {

		gl.Clear(webgl.COLOR_BUFFER_BIT | webgl.DEPTH_BUFFER_BIT)

		viewModelMatrix.Set(viewMatrix)
		program.UniformMatrix4fv("u_modelViewMatrix", &viewModelMatrix)
		normalMatrix := linalg.Matrix3{}
		normalMatrix.MakeNormalMatrix(viewModelMatrix)
		program.UniformMatrix3fv("u_normalMatrix", &normalMatrix)

		vao2.DrawElements()

		rot := linalg.NewMatrix4()
		rot.Rotation(angle, linalg.Vector3{0, 0, 1})
		angle += 0.01

		viewModelMatrix.Set(viewMatrix)
		viewModelMatrix.MultAssign(&rot)

		program.UniformMatrix4fv("u_modelViewMatrix", &viewModelMatrix)
		normalMatrix.MakeNormalMatrix(viewModelMatrix)
		program.UniformMatrix3fv("u_normalMatrix", &normalMatrix)

		vao1.DrawElements()
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
