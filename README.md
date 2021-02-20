# wasp: WebGL2 utility lib in Go on WASM
A small, stingy and light-weight (= little functionality) webgl2 layer in Go for WASM.

Wasp comes with a small set of convenience structs and functions that easy the work
with WebGL2. The main concepts are the following.
- ``shader.go``: Common shader functionality such as loading, compiling shader sources as well as caching attributes and uniforms.
- ``webglapp.go``: A setup for a fullpage WebGL2 application with ``requestAnimationFrame``, mouse events and resize callbacks. Reduces common boilerplate code on the application level.
- ``mesh.go``: a simple data structure for manipulating triangle meshes and uploading to WebGL2 as interleaved VertexArrayObjects.
- ``webgl/renderingcontext.go``: Wrapper for ``WebGl2RenderingContext`` in JavaScript for typesafety and convenience.

Check out the examples in the ``examples`` folder for further research.

## Examples
There are currently three examples that demonstrate the usage of Wasp.

### Simple Cube Example
The simple cube example in ``examples/simplecube.go`` shows a colored, rotating cube. It uses a very basic shader without lighting. 

![Wasp example showing simple rotating cube](https://github.com/schabby/go-wasm-webgl/raw/main/examples/simplecube.png "Wasp example showing simple rotating cube")

### Simple Sphere Example
The simple sphere example in ``examples/simplesphere.go`` shows a colored sphere without lighting.

![Wasp example showing simple colored sphere](https://github.com/schabby/go-wasm-webgl/raw/main/examples/simplesphere.png "Wasp example showing simple colored sphere")

### Simple Lighting Example
The example in ``examples/simplelighting.go`` shows spheres with classical lighting shading. 

![Wasp simple light shading example](https://github.com/schabby/go-wasm-webgl/raw/main/examples/simplelighting.png "Wasp example showing simple lighting shading")

### Simple Texture Example
The example in ``examples/simpletexture.go`` shows a crate with a textured mapped on its faces.

![Wasp simple texturing example](https://github.com/schabby/go-wasm-webgl/raw/main/examples/simpletexture.png "Wasp example showing simple texturing")


## How to build the examples:
While not strictly necessary, it may help to have some prior knowledge on WASM with Go. There are plenty
of great tutorials out there which help you understand how Go and WASM work together.

On the terminal of your choice, move into the project directory and run
```bash
GOOS=js GOARCH=wasm go build -o examples/hellocube.wasm examples/hellocube.go 
```
This will compile the cube example into a wasm file in `examples/hellocube.wasm` where it will be picked up
by the javascript loader in `examples/hellocube.html`.

Make sure that `examples` also contains `wasm_exec.js` which is shiped along the standard Go installation, such that you may just need to copy it over `cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./examples`.

To run the example, you need any webserver to serve the static content in the `examples` directory (HTML files and WASM files.). Any http server will do, but you may want to run [goexec][https://github.com/shurcooL/goexec]  
```bash
 goexec 'http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`.`)))'
```
in the `examples` directory. Side note: `goexec`is not on my `PATH` environment variable, such that I need to call it will the
full path
```bash
~/go/bin/goexec 'http.ListenAndServe(`:8080`, http.FileServer(http.Dir(`.`)))'
```
Then, open a browser and request `http://localhost:8080/hellocube.html` which will load the page with the embedded WASM
code, init the WebGL2 rendering context and show a rotating cube. 

You may want to check out the browsers built-on console to see any logging output.