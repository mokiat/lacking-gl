package internal

import (
	"fmt"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/mokiat/lacking/render"
)

func NewVertexArray(info render.VertexArrayInfo) *VertexArray {
	var id uint32
	gl.CreateVertexArrays(1, &id)

	for index, binding := range info.Bindings {
		if buffer, ok := binding.VertexBuffer.(*Buffer); ok {
			gl.VertexArrayVertexBuffer(id, uint32(index), buffer.id, 0, int32(binding.Stride))
		}
	}

	for _, attribute := range info.Attributes {
		gl.EnableVertexArrayAttrib(id, uint32(attribute.Location))
		count, compType, normalized := glAttribParams(attribute.Format)
		gl.VertexArrayAttribFormat(id, uint32(attribute.Location), count, compType, normalized, uint32(attribute.Offset))
		gl.VertexArrayAttribBinding(id, uint32(attribute.Location), uint32(attribute.Binding))
	}

	if indexBuffer, ok := info.IndexBuffer.(*Buffer); ok {
		gl.VertexArrayElementBuffer(id, indexBuffer.id)
	}

	return &VertexArray{
		id:          id,
		indexFormat: gl.UNSIGNED_SHORT, // FIXME
	}
}

type VertexArray struct {
	id          uint32
	indexFormat uint32 // TODO
}

func (a *VertexArray) Release() {
	gl.DeleteVertexArrays(1, &a.id)
	a.id = 0
}

func glAttribParams(format render.VertexAttributeFormat) (int32, uint32, bool) {
	switch format {
	case render.VertexAttributeFormatR32F:
		return 1, gl.FLOAT, false
	case render.VertexAttributeFormatRG32F:
		return 2, gl.FLOAT, false
	case render.VertexAttributeFormatRGB32F:
		return 3, gl.FLOAT, false
	case render.VertexAttributeFormatRGBA32F:
		return 4, gl.FLOAT, false

	case render.VertexAttributeFormatR16F:
		return 1, gl.HALF_FLOAT, false
	case render.VertexAttributeFormatRG16F:
		return 2, gl.HALF_FLOAT, false
	case render.VertexAttributeFormatRGB16F:
		return 3, gl.HALF_FLOAT, false
	case render.VertexAttributeFormatRGBA16F:
		return 4, gl.HALF_FLOAT, false

	case render.VertexAttributeFormatR16S:
		return 1, gl.SHORT, false
	case render.VertexAttributeFormatRG16S:
		return 2, gl.SHORT, false
	case render.VertexAttributeFormatRGB16S:
		return 3, gl.SHORT, false
	case render.VertexAttributeFormatRGBA16S:
		return 4, gl.SHORT, false

	case render.VertexAttributeFormatR16SN:
		return 1, gl.SHORT, true
	case render.VertexAttributeFormatRG16SN:
		return 2, gl.SHORT, true
	case render.VertexAttributeFormatRGB16SN:
		return 3, gl.SHORT, true
	case render.VertexAttributeFormatRGBA16SN:
		return 4, gl.SHORT, true

	case render.VertexAttributeFormatR16U:
		return 1, gl.UNSIGNED_SHORT, false
	case render.VertexAttributeFormatRG16U:
		return 2, gl.UNSIGNED_SHORT, false
	case render.VertexAttributeFormatRGB16U:
		return 3, gl.UNSIGNED_SHORT, false
	case render.VertexAttributeFormatRGBA16U:
		return 4, gl.UNSIGNED_SHORT, false

	case render.VertexAttributeFormatR16UN:
		return 1, gl.UNSIGNED_SHORT, true
	case render.VertexAttributeFormatRG16UN:
		return 2, gl.UNSIGNED_SHORT, true
	case render.VertexAttributeFormatRGB16UN:
		return 3, gl.UNSIGNED_SHORT, true
	case render.VertexAttributeFormatRGBA16UN:
		return 4, gl.UNSIGNED_SHORT, true

	case render.VertexAttributeFormatR8S:
		return 1, gl.BYTE, false
	case render.VertexAttributeFormatRG8S:
		return 2, gl.BYTE, false
	case render.VertexAttributeFormatRGB8S:
		return 3, gl.BYTE, false
	case render.VertexAttributeFormatRGBA8S:
		return 4, gl.BYTE, false

	case render.VertexAttributeFormatR8SN:
		return 1, gl.BYTE, true
	case render.VertexAttributeFormatRG8SN:
		return 2, gl.BYTE, true
	case render.VertexAttributeFormatRGB8SN:
		return 3, gl.BYTE, true
	case render.VertexAttributeFormatRGBA8SN:
		return 4, gl.BYTE, true

	case render.VertexAttributeFormatR8U:
		return 1, gl.UNSIGNED_BYTE, false
	case render.VertexAttributeFormatRG8U:
		return 2, gl.UNSIGNED_BYTE, false
	case render.VertexAttributeFormatRGB8U:
		return 3, gl.UNSIGNED_BYTE, false
	case render.VertexAttributeFormatRGBA8U:
		return 4, gl.UNSIGNED_BYTE, false

	case render.VertexAttributeFormatR8UN:
		return 1, gl.UNSIGNED_BYTE, true
	case render.VertexAttributeFormatRG8UN:
		return 2, gl.UNSIGNED_BYTE, true
	case render.VertexAttributeFormatRGB8UN:
		return 3, gl.UNSIGNED_BYTE, true
	case render.VertexAttributeFormatRGBA8UN:
		return 4, gl.UNSIGNED_BYTE, true

	default:
		panic(fmt.Errorf("unknown attribute format: %d", format))
	}
}
