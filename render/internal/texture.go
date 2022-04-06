package internal

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/mokiat/lacking/render"
)

func NewColorTexture2D(info render.ColorTexture2DInfo) *Texture {
	var id uint32
	gl.CreateTextures(gl.TEXTURE_2D, 1, &id)
	gl.TextureParameteri(id, gl.TEXTURE_WRAP_S, glWrap(info.Wrapping))
	gl.TextureParameteri(id, gl.TEXTURE_WRAP_T, glWrap(info.Wrapping))
	gl.TextureParameteri(id, gl.TEXTURE_MIN_FILTER, glFilter(info.Filtering, info.Mipmapping))
	gl.TextureParameteri(id, gl.TEXTURE_MAG_FILTER, glFilter(info.Filtering, false)) // no mipmaps when magnification
	if info.Filtering == render.FilterModeAnisotropic {
		var maxAnisotropy float32
		gl.GetFloatv(gl.MAX_TEXTURE_MAX_ANISOTROPY, &maxAnisotropy)
		gl.TextureParameterf(id, gl.TEXTURE_MAX_ANISOTROPY, maxAnisotropy)
	}

	levels := glMipmapLevels(info.Width, info.Height, info.Mipmapping)
	internalFormat := glInternalFormat(info.Format, info.GammaCorrection)
	gl.TextureStorage2D(id, levels, internalFormat, int32(info.Width), int32(info.Height))

	if info.Data != nil {
		dataFormat := glDataFormat(info.Format)
		componentType := glDataComponentType(info.Format)
		gl.TextureSubImage2D(id, 0, 0, 0, int32(info.Width), int32(info.Height), dataFormat, componentType, gl.Ptr(info.Data))

		if info.Mipmapping {
			gl.GenerateTextureMipmap(id)
		}
	}

	return &Texture{
		id: id,
	}
}

func NewDepthTexture2D(info render.DepthTexture2DInfo) *Texture {
	var id uint32
	gl.CreateTextures(gl.TEXTURE_2D, 1, &id)
	gl.TextureParameteri(id, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TextureParameteri(id, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TextureParameteri(id, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TextureParameteri(id, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TextureStorage2D(id, 1, gl.DEPTH_COMPONENT32, int32(info.Width), int32(info.Height))
	return &Texture{
		id: id,
	}
}

func NewStencilTexture2D(info render.StencilTexture2DInfo) *Texture {
	var id uint32
	gl.CreateTextures(gl.TEXTURE_2D, 1, &id)
	gl.TextureParameteri(id, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TextureParameteri(id, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TextureParameteri(id, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TextureParameteri(id, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TextureStorage2D(id, 1, gl.STENCIL_INDEX8, int32(info.Width), int32(info.Height))
	return &Texture{
		id: id,
	}
}

func NewDepthStencilTexture2D(info render.DepthStencilTexture2DInfo) *Texture {
	var id uint32
	gl.CreateTextures(gl.TEXTURE_2D, 1, &id)
	gl.TextureParameteri(id, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TextureParameteri(id, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TextureParameteri(id, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TextureParameteri(id, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TextureStorage2D(id, 1, gl.DEPTH24_STENCIL8, int32(info.Width), int32(info.Height))
	return &Texture{
		id: id,
	}
}

type Texture struct {
	id uint32
}

func (t *Texture) Release() {
	gl.DeleteTextures(1, &t.id)
	t.id = 0
}

func glWrap(wrap render.WrapMode) int32 {
	switch wrap {
	case render.WrapModeClamp:
		return gl.CLAMP_TO_EDGE
	case render.WrapModeRepeat:
		return gl.REPEAT
	case render.WrapModeMirroredRepeat:
		return gl.MIRRORED_REPEAT
	default:
		return gl.CLAMP_TO_EDGE
	}
}

func glFilter(filter render.FilterMode, mipmaps bool) int32 {
	switch filter {
	case render.FilterModeNearest:
		if mipmaps {
			return gl.NEAREST_MIPMAP_NEAREST
		}
		return gl.NEAREST
	case render.FilterModeLinear, render.FilterModeAnisotropic:
		if mipmaps {
			return gl.LINEAR_MIPMAP_LINEAR
		}
		return gl.LINEAR
	default:
		return gl.NEAREST
	}
}

func glMipmapLevels(width, height int, mipmapping bool) int32 {
	if !mipmapping {
		return 1
	}
	count := int32(1)
	for width > 1 || height > 1 {
		width /= 2
		height /= 2
		count++
	}
	return count
}

func glInternalFormat(format render.DataFormat, gammaCorrection bool) uint32 {
	switch format {
	case render.DataFormatRGBA8:
		if gammaCorrection {
			return gl.SRGB8_ALPHA8
		}
		return gl.RGBA8
	case render.DataFormatRGBA32F:
		return gl.RGBA32F
	default:
		return gl.RGBA8
	}
}

func glDataFormat(format render.DataFormat) uint32 {
	switch format {
	default:
		return gl.RGBA
	}
}

func glDataComponentType(format render.DataFormat) uint32 {
	switch format {
	case render.DataFormatRGBA8:
		return gl.UNSIGNED_BYTE
	case render.DataFormatRGBA32F:
		return gl.FLOAT
	default:
		return gl.UNSIGNED_BYTE
	}
}
