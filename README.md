# wasp: WebGL2 utility lib in Go on WASM
A small, stingy and light-weight (= little functionality) webgl2 layer in Go for WASM.

How to build an example:
On the terminal of your choice, move into the project directory and run
```bash
GOOS=js GOARCH=wasm go build -o examples/hellocube.wasm examples/hellocube.go 
```
This will compile the cube example into a wasm file in `examples/hellocube.wasm` where it will be picked up
by the javascript loader in `examples/hellocube.html`.

To run the example, you need any webserver to serve the static content in the `examples` directory (HTML files and WASM files.). Any http server will do, but you may want to run  
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