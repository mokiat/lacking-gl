package internal

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/mokiat/lacking/log"
	"github.com/mokiat/lacking/render"
)

func NewFramebuffer(info render.FramebufferInfo) *Framebuffer {
	var id uint32
	gl.CreateFramebuffers(1, &id)

	var drawBuffers []uint32
	for i, attachment := range info.ColorAttachments {
		if colorAttachment, ok := attachment.(*Texture); ok {
			attachmentID := gl.COLOR_ATTACHMENT0 + uint32(i)
			gl.NamedFramebufferTexture(id, attachmentID, colorAttachment.id, 0)
			drawBuffers = append(drawBuffers, attachmentID)
		}
	}

	if depthStencilAttachment, ok := info.DepthStencilAttachment.(*Texture); ok {
		gl.NamedFramebufferTexture(id, gl.DEPTH_STENCIL_ATTACHMENT, depthStencilAttachment.id, 0)
	} else {
		if depthAttachment, ok := info.DepthAttachment.(*Texture); ok {
			gl.NamedFramebufferTexture(id, gl.DEPTH_ATTACHMENT, depthAttachment.id, 0)
		}
		if stencilAttachment, ok := info.StencilAttachment.(*Texture); ok {
			gl.NamedFramebufferTexture(id, gl.STENCIL_ATTACHMENT, stencilAttachment.id, 0)
		}
	}

	gl.NamedFramebufferDrawBuffers(id, int32(len(drawBuffers)), &drawBuffers[0])

	status := gl.CheckNamedFramebufferStatus(id, gl.FRAMEBUFFER)
	if status != gl.FRAMEBUFFER_COMPLETE {
		log.Error("Framebuffer is incomplete")
	}

	return &Framebuffer{
		id: id,
	}
}

var DefaultFramebuffer = &Framebuffer{
	id: 0,
}

type Framebuffer struct {
	id uint32
}

func (f *Framebuffer) Release() {
	gl.DeleteFramebuffers(1, &f.id)
	f.id = 0
}
