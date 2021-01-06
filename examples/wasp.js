/*  This JavaScript file holds helper functions that improve speed (or reduce garbage collection waste) when
 *  working with Wasp.
 */

// global constant for uploading matrices to shaders. This is often necessary several times per frame. The current
// go/wasm implementation creates new arrays which then need to be GC'ed. Keeping a global constant array
// for re-use is more efficient. Note that this also assumes that the array is accessed on a single-thread fashion.
const uniformMatrix4fvFloat32 = new Float32Array(16)

// function to be called from WASM to upload 4x4 matrices to shaders
function uniformMatrix4fv(gl, location, a0, a1, a2, a3, a4, a5, a6, a7, a8, a9, a10, a11, a12, a13, a14, a15) {
  uniformMatrix4fvFloat32[0] = a0
  uniformMatrix4fvFloat32[1] = a1
  uniformMatrix4fvFloat32[2] = a2
  uniformMatrix4fvFloat32[3] = a3
  uniformMatrix4fvFloat32[4] = a4
  uniformMatrix4fvFloat32[5] = a5
  uniformMatrix4fvFloat32[6] = a6
  uniformMatrix4fvFloat32[7] = a7
  uniformMatrix4fvFloat32[8] = a8
  uniformMatrix4fvFloat32[9] = a9
  uniformMatrix4fvFloat32[10] = a10
  uniformMatrix4fvFloat32[11] = a11
  uniformMatrix4fvFloat32[12] = a12
  uniformMatrix4fvFloat32[13] = a13
  uniformMatrix4fvFloat32[14] = a14
  uniformMatrix4fvFloat32[15] = a15
  gl.uniformMatrix4fv(location, false, uniformMatrix4fvFloat32)
}

const uniformMatrix3fvFloat32 = new Float32Array(9)
function uniformMatrix3fv(gl, location, a0, a1, a2, a3, a4, a5, a6, a7, a8) {
  uniformMatrix3fvFloat32[0] = a0
  uniformMatrix3fvFloat32[1] = a1
  uniformMatrix3fvFloat32[2] = a2
  uniformMatrix3fvFloat32[3] = a3
  uniformMatrix3fvFloat32[4] = a4
  uniformMatrix3fvFloat32[5] = a5
  uniformMatrix3fvFloat32[6] = a6
  uniformMatrix3fvFloat32[7] = a7
  uniformMatrix3fvFloat32[8] = a8
  gl.uniformMatrix3fv(location, false, uniformMatrix3fvFloat32)
}