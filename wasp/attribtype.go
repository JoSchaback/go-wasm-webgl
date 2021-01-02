package wasp

type AttribType struct {
	name string
	size int
}

// POSITION defines the common position vertex attrib components.
var POSITION = AttribType{name: "POSITION", size: 3}
var RGB = AttribType{name: "RGB", size: 3}
var RGBA = AttribType{name: "RGBA", size: 4}
var UV = AttribType{name: "UV", size: 2}
var NORMAL = AttribType{name: "NORMAL", size: 3}

func (c AttribType) Name() string { return c.name }
func (c AttribType) Size() int    { return c.size }
