package internal

import (
	"runtime"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/mokiat/lacking/render"
)

func NewRenderer() *Renderer {
	return &Renderer{
		framebuffer: DefaultFramebuffer,
	}
}

type Renderer struct {
	framebuffer           *Framebuffer
	invalidateAttachments []uint32
	primitive             uint32
	indexType             uint32
}

func (r *Renderer) BeginRenderPass(info render.RenderPassInfo) {
	gl.Enable(gl.FRAMEBUFFER_SRGB)
	gl.Enable(gl.CLIP_DISTANCE0)
	gl.Enable(gl.CLIP_DISTANCE1)
	gl.Enable(gl.CLIP_DISTANCE2)
	gl.Enable(gl.CLIP_DISTANCE3)

	r.framebuffer = info.Framebuffer.(*Framebuffer)

	gl.BindFramebuffer(gl.FRAMEBUFFER, r.framebuffer.id)
	gl.Viewport(
		int32(info.Viewport.X),
		int32(info.Viewport.Y),
		int32(info.Viewport.Width),
		int32(info.Viewport.Height),
	)

	// TODO
	// var rgba = [4]float32{
	// 	0.0,
	// 	0.3,
	// 	0.6,
	// 	1.0,
	// }
	// gl.ClearNamedFramebufferfv(r.framebuffer.id, gl.COLOR, 0, &rgba[0])

	clearDepth := info.StencilLoadOp == render.LoadOperationClear
	clearStencil := info.StencilLoadOp == render.LoadOperationClear

	if clearDepth && clearStencil {
		depthValue := info.DepthClearValue
		stencilValue := int32(info.StencilClearValue)
		gl.ClearNamedFramebufferfi(r.framebuffer.id, gl.DEPTH_STENCIL, 0, depthValue, stencilValue)
	} else {
		if clearDepth {
			depthValue := info.DepthClearValue
			gl.ClearNamedFramebufferfv(r.framebuffer.id, gl.DEPTH, 0, &depthValue)
		}
		if clearStencil {
			stencilValue := uint32(info.StencilClearValue)
			gl.ClearNamedFramebufferuiv(r.framebuffer.id, gl.STENCIL, 0, &stencilValue)
		}
	}

	invalidateDepth := info.DepthStoreOp == render.StoreOperationDontCare
	invalidateStencil := info.StencilStoreOp == render.StoreOperationDontCare

	r.invalidateAttachments = r.invalidateAttachments[:0]

	if invalidateDepth && invalidateStencil {
		r.invalidateAttachments = append(r.invalidateAttachments, gl.DEPTH_STENCIL_ATTACHMENT)
	} else {
		if invalidateDepth {
			r.invalidateAttachments = append(r.invalidateAttachments, gl.DEPTH_ATTACHMENT)
		}
		if invalidateStencil {
			r.invalidateAttachments = append(r.invalidateAttachments, gl.STENCIL_ATTACHMENT)
		}
	}
}

func (r *Renderer) BindPipeline(pipeline render.Pipeline) {
	if pipeline, ok := pipeline.(*Pipeline); ok {
		r.primitive = gl.TRIANGLE_FAN // FIXME

		switch pipeline.Culling {
		case render.CullModeNone:
			gl.Disable(gl.CULL_FACE)
		case render.CullModeBack:
			gl.Enable(gl.CULL_FACE)
			gl.CullFace(gl.BACK)
		case render.CullModeFront:
			gl.Enable(gl.CULL_FACE)
			gl.CullFace(gl.FRONT)
		case render.CullModeFrontAndBack:
			gl.Enable(gl.CULL_FACE)
			gl.CullFace(gl.FRONT_AND_BACK)
		}

		switch pipeline.FrontFace {
		case render.FaceOrientationCCW:
			gl.FrontFace(gl.CCW)
		case render.FaceOrientationCW:
			gl.FrontFace(gl.CW)
		}

		gl.LineWidth(pipeline.LineWidth)

		if pipeline.DepthTest {
			gl.Enable(gl.DEPTH_TEST)
			// gl.DepthFunc(xfunc uint32) // TODO
		} else {
			gl.Disable(gl.DEPTH_TEST)
		}
		if pipeline.DepthWrite {
			gl.DepthMask(true)
		} else {
			gl.DepthMask(false)
		}

		if pipeline.StencilTest {
			gl.Enable(gl.STENCIL_TEST)
		} else {
			gl.Disable(gl.STENCIL_TEST)
		}

		gl.Enable(gl.BLEND)
		gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

		gl.ColorMask(pipeline.ColorWrite[0], pipeline.ColorWrite[1], pipeline.ColorWrite[2], pipeline.ColorWrite[3])

		if program, ok := pipeline.Program.(*Program); ok {
			gl.UseProgram(program.id)
		}

		if vertexArray, ok := pipeline.VertexArray.(*VertexArray); ok {
			gl.BindVertexArray(vertexArray.id)
			r.indexType = vertexArray.indexFormat
		}
	}
}

func (r *Renderer) Uniform4f(location render.UniformLocation, values [4]float32) {
	gl.Uniform4f(location.(int32), values[0], values[1], values[2], values[3])
}

func (r *Renderer) Uniform1i(location render.UniformLocation, value int) {
	gl.Uniform1i(location.(int32), int32(value))
}

func (r *Renderer) UniformMatrix4f(location render.UniformLocation, values [16]float32) {
	slice := values[:]
	gl.UniformMatrix4fv(location.(int32), 1, false, &slice[0])
	runtime.KeepAlive(slice)
}

func (r *Renderer) TextureUnit(index int, texture render.Texture) {
	if texture, ok := texture.(*Texture); ok {
		gl.BindTextureUnit(uint32(index), texture.id)
	}
}

func (r *Renderer) Draw(vertexOffset, vertexCount, instanceCount int) {
	gl.DrawArraysInstanced(r.primitive, int32(vertexOffset), int32(vertexCount), int32(instanceCount))
}

func (r *Renderer) DrawIndexed(indexOffset, indexCount, instanceCount int) {
	gl.DrawElementsInstanced(r.primitive, int32(indexCount), r.indexType, gl.PtrOffset(indexOffset), int32(instanceCount))
}

func (r *Renderer) EndRenderPass() {
	if len(r.invalidateAttachments) > 0 {
		// TODO: When the viewport is just part of the framebuffer
		// we should use glInvalidateNamedFramebufferSubData
		gl.InvalidateNamedFramebufferData(r.framebuffer.id, 1, &r.invalidateAttachments[0])
	}

	// FIXME
	gl.Disable(gl.CLIP_DISTANCE0)
	gl.Disable(gl.CLIP_DISTANCE1)
	gl.Disable(gl.CLIP_DISTANCE2)
	gl.Disable(gl.CLIP_DISTANCE3)
	gl.Disable(gl.BLEND)
	gl.DepthMask(true)

	r.framebuffer = DefaultFramebuffer
}
