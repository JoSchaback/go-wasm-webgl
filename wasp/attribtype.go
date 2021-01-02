package wasp

type AttribType struct {
	name string
	size int
}

// POSITION defines the common position vertex attrib components.
var POSITION = AttribType{name: "POSITION", size: 3}

// RGB stands for Red, Green, Blue and represents color components in a vertex
var RGB = AttribType{name: "RGB", size: 3}

// RGBA stands for Red, Green, Blue, Alpha and represents color components in a vertex
var RGBA = AttribType{name: "RGBA", size: 4}

// UV represent texture mapping coordinates in a vertex
var UV = AttribType{name: "UV", size: 2}

// NORMAL represent normal coordinates for a vertex
var NORMAL = AttribType{name: "NORMAL", size: 3}

// Name exposes the private name field
func (c AttribType) Name() string { return c.name }

// Size exposes the private size field
func (c AttribType) Size() int { return c.size }
