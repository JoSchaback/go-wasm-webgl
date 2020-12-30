package wasp

type AttribType struct {
	name string
	size uint32
}

// POSITION defines the common position vertex attrib components.
var POSITION = AttribType{name: "POSITION", size: uint32(3)}
var RGB = AttribType{name: "RGB", size: uint32(3)}
var RGBA = AttribType{name: "RGBA", size: uint32(4)}
var UV = AttribType{name: "UV", size: uint32(2)}
var NORMAL = AttribType{name: "NORMAL", size: uint32(3)}

func (c AttribType) Name() string { return c.name }
func (c AttribType) Size() uint32 { return c.size }
