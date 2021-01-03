# wasp: WebGL2 utility lib in Go on WASM
A small, stingy and light-weight (= little functionality) webgl2 layer in Go for WASM.

How to build an example:
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