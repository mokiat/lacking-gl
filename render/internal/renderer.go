package internal

import (
	"fmt"

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
	isDefaultFramebuffer := r.framebuffer.id == 0

	gl.BindFramebuffer(gl.FRAMEBUFFER, r.framebuffer.id)
	gl.Viewport(
		int32(info.Viewport.X),
		int32(info.Viewport.Y),
		int32(info.Viewport.Width),
		int32(info.Viewport.Height),
	)

	for i, attachment := range info.Colors {
		if r.framebuffer.activeDrawBuffers[i] && (attachment.LoadOp == render.LoadOperationClear) {
			gl.ClearNamedFramebufferfv(r.framebuffer.id, gl.COLOR, int32(i), &attachment.ClearValue[0])
		}
	}

	clearDepth := info.DepthLoadOp == render.LoadOperationClear
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
			stencilValue := int32(info.StencilClearValue)
			gl.ClearNamedFramebufferiv(r.framebuffer.id, gl.STENCIL, 0, &stencilValue)
		}
	}

	r.invalidateAttachments = r.invalidateAttachments[:0]

	invalidateDepth := info.DepthStoreOp == render.StoreOperationDontCare
	invalidateStencil := info.StencilStoreOp == render.StoreOperationDontCare

	for i, attachment := range info.Colors {
		if r.framebuffer.activeDrawBuffers[i] && (attachment.StoreOp == render.StoreOperationDontCare) {
			if isDefaultFramebuffer {
				if i == 0 {
					r.invalidateAttachments = append(r.invalidateAttachments, gl.COLOR)
				}
			} else {
				r.invalidateAttachments = append(r.invalidateAttachments, gl.COLOR_ATTACHMENT0+uint32(i))
			}
		}
	}

	if invalidateDepth && invalidateStencil && !isDefaultFramebuffer {
		r.invalidateAttachments = append(r.invalidateAttachments, gl.DEPTH_STENCIL_ATTACHMENT)
	} else {
		if invalidateDepth {
			if isDefaultFramebuffer {
				r.invalidateAttachments = append(r.invalidateAttachments, gl.DEPTH)
			} else {
				r.invalidateAttachments = append(r.invalidateAttachments, gl.DEPTH_ATTACHMENT)
			}
		}
		if invalidateStencil {
			if isDefaultFramebuffer {
				r.invalidateAttachments = append(r.invalidateAttachments, gl.STENCIL)
			} else {
				r.invalidateAttachments = append(r.invalidateAttachments, gl.STENCIL_ATTACHMENT)
			}
		}
	}
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
	gl.ColorMask(true, true, true, true)
	gl.DepthMask(true)

	r.framebuffer = DefaultFramebuffer
}

func (r *Renderer) BindPipeline(pipeline render.Pipeline) {
	intPipeline := pipeline.(*Pipeline)
	r.executeCommandBindPipeline(CommandBindPipeline{
		ProgramID:        intPipeline.ProgramID,
		Topology:         intPipeline.Topology,
		CullTest:         intPipeline.CullTest,
		FrontFace:        intPipeline.FrontFace,
		LineWidth:        intPipeline.LineWidth,
		DepthTest:        intPipeline.DepthTest,
		DepthWrite:       intPipeline.DepthWrite,
		DepthComparison:  intPipeline.DepthComparison,
		StencilTest:      intPipeline.StencilTest,
		StencilOpFront:   intPipeline.StencilOpFront,
		StencilOpBack:    intPipeline.StencilOpBack,
		StencilFuncFront: intPipeline.StencilFuncFront,
		StencilFuncBack:  intPipeline.StencilFuncBack,
		StencilMaskFront: intPipeline.StencilMaskFront,
		StencilMaskBack:  intPipeline.StencilMaskBack,
		ColorWrite:       intPipeline.ColorWrite,
		BlendEnabled:     intPipeline.BlendEnabled,
		BlendColor:       intPipeline.BlendColor,
		BlendEquation:    intPipeline.BlendEquation,
		BlendFunc:        intPipeline.BlendFunc,
		VertexArray:      intPipeline.VertexArray,
	})
}

func (r *Renderer) Uniform1f(location render.UniformLocation, value float32) {
	r.executeCommandUniform1f(CommandUniform1f{
		Location: location.(int32),
		Value:    value,
	})
}

func (r *Renderer) Uniform1i(location render.UniformLocation, value int) {
	r.executeCommandUniform1i(CommandUniform1i{
		Location: location.(int32),
		Value:    int32(value),
	})
}

func (r *Renderer) Uniform3f(location render.UniformLocation, values [3]float32) {
	r.executeCommandUniform3f(CommandUniform3f{
		Location: location.(int32),
		Values:   values,
	})
}

func (r *Renderer) Uniform4f(location render.UniformLocation, values [4]float32) {
	r.executeCommandUniform4f(CommandUniform4f{
		Location: location.(int32),
		Values:   values,
	})
}

func (r *Renderer) UniformMatrix4f(location render.UniformLocation, values [16]float32) {
	r.executeCommandUniformMatrix4f(CommandUniformMatrix4f{
		Location: location.(int32),
		Values:   values,
	})
}

func (r *Renderer) TextureUnit(index int, texture render.Texture) {
	r.executeCommandTextureUnit(CommandTextureUnit{
		Index:     uint32(index),
		TextureID: texture.(*Texture).id,
	})
}

func (r *Renderer) Draw(vertexOffset, vertexCount, instanceCount int) {
	r.executeCommandDraw(CommandDraw{
		VertexOffset:  int32(vertexOffset),
		VertexCount:   int32(vertexCount),
		InstanceCount: int32(instanceCount),
	})
}

func (r *Renderer) DrawIndexed(indexOffset, indexCount, instanceCount int) {
	r.executeCommandDrawIndexed(CommandDrawIndexed{
		IndexOffset:   int32(indexOffset),
		IndexCount:    int32(indexCount),
		InstanceCount: int32(instanceCount),
	})
}

func (r *Renderer) CopyContentToTexture(info render.CopyContentToTextureInfo) {
	intTexture := info.Texture.(*Texture)
	gl.CopyTextureSubImage2D(
		intTexture.id,
		int32(info.TextureLevel),
		int32(info.TextureX),
		int32(info.TextureY),
		int32(info.FramebufferX),
		int32(info.FramebufferY),
		int32(info.Width),
		int32(info.Height),
	)
	if info.GenerateMipmaps {
		gl.GenerateTextureMipmap(intTexture.id)
	}
}

func (r *Renderer) SubmitQueue(queue *CommandQueue) {
	for MoreCommands(queue) {
		header := PopCommand[CommandHeader](queue)
		switch header.Kind {
		case CommandKindBindPipeline:
			command := PopCommand[CommandBindPipeline](queue)
			r.executeCommandBindPipeline(command)
		case CommandKindTopology:
			command := PopCommand[CommandTopology](queue)
			r.executeCommandTopology(command)
		case CommandKindCullTest:
			command := PopCommand[CommandCullTest](queue)
			r.executeCommandCullTest(command)
		case CommandKindFrontFace:
			command := PopCommand[CommandFrontFace](queue)
			r.executeCommandFrontFace(command)
		case CommandKindLineWidth:
			command := PopCommand[CommandLineWidth](queue)
			r.executeCommandLineWidth(command)
		case CommandKindDepthTest:
			command := PopCommand[CommandDepthTest](queue)
			r.executeCommandDepthTest(command)
		case CommandKindDepthWrite:
			command := PopCommand[CommandDepthWrite](queue)
			r.executeCommandDepthWrite(command)
		case CommandKindDepthComparison:
			command := PopCommand[CommandDepthComparison](queue)
			r.executeCommandDepthComparison(command)
		case CommandKindUniform1f:
			command := PopCommand[CommandUniform1f](queue)
			r.executeCommandUniform1f(command)
		case CommandKindUniform1i:
			command := PopCommand[CommandUniform1i](queue)
			r.executeCommandUniform1i(command)
		case CommandKindUniform3f:
			command := PopCommand[CommandUniform3f](queue)
			r.executeCommandUniform3f(command)
		case CommandKindUniform4f:
			command := PopCommand[CommandUniform4f](queue)
			r.executeCommandUniform4f(command)
		case CommandKindUniformMatrix4f:
			command := PopCommand[CommandUniformMatrix4f](queue)
			r.executeCommandUniformMatrix4f(command)
		case CommandKindTextureUnit:
			command := PopCommand[CommandTextureUnit](queue)
			r.executeCommandTextureUnit(command)
		case CommandKindDraw:
			command := PopCommand[CommandDraw](queue)
			r.executeCommandDraw(command)
		case CommandKindDrawIndexed:
			command := PopCommand[CommandDrawIndexed](queue)
			r.executeCommandDrawIndexed(command)
		default:
			panic(fmt.Errorf("unknown command kind: %v", header.Kind))
		}
	}
	queue.Reset()
}

func (r *Renderer) executeCommandBindPipeline(command CommandBindPipeline) {
	gl.UseProgram(command.ProgramID)
	r.executeCommandTopology(command.Topology)
	r.executeCommandCullTest(command.CullTest)
	r.executeCommandFrontFace(command.FrontFace)
	r.executeCommandLineWidth(command.LineWidth)
	r.executeCommandDepthTest(command.DepthTest)
	r.executeCommandDepthWrite(command.DepthWrite)
	r.executeCommandDepthComparison(command.DepthComparison)
	r.executeCommandStencilTest(command.StencilTest)
	// TODO: Optimize if equal except for face
	r.executeCommandStencilFunc(command.StencilFuncFront)
	r.executeCommandStencilFunc(command.StencilFuncBack)
	// TODO: Optimize if equal except for face
	r.executeCommandStencilOperation(command.StencilOpFront)
	r.executeCommandStencilOperation(command.StencilOpBack)
	// TODO: Optimize if equal except for face
	r.executeCommandStencilMask(command.StencilMaskFront)
	r.executeCommandStencilMask(command.StencilMaskBack)
	r.executeCommandColorWrite(command.ColorWrite)
	if command.BlendEnabled {
		gl.Enable(gl.BLEND)
	} else {
		gl.Disable(gl.BLEND)
	}
	r.executeCommandBlendEquation(command.BlendEquation)
	r.executeCommandBlendFunc(command.BlendFunc)
	r.executeCommandBlendColor(command.BlendColor)
	r.executeCommandBindVertexArray(command.VertexArray)
}

func (r *Renderer) executeCommandTopology(command CommandTopology) {
	r.primitive = command.Topology
}

func (r *Renderer) executeCommandCullTest(command CommandCullTest) {
	if command.Enabled {
		gl.Enable(gl.CULL_FACE)
		gl.CullFace(command.Face)
	} else {
		gl.Disable(gl.CULL_FACE)
	}
}

func (r *Renderer) executeCommandFrontFace(command CommandFrontFace) {
	gl.FrontFace(command.Orientation)
}

func (r *Renderer) executeCommandLineWidth(command CommandLineWidth) {
	if command.Width > 0.0 {
		gl.LineWidth(command.Width)
	}
}

func (r *Renderer) executeCommandDepthTest(command CommandDepthTest) {
	if command.Enabled {
		gl.Enable(gl.DEPTH_TEST)
	} else {
		gl.Disable(gl.DEPTH_TEST)
	}
}

func (r *Renderer) executeCommandDepthWrite(command CommandDepthWrite) {
	gl.DepthMask(command.Enabled)
}

func (r *Renderer) executeCommandDepthComparison(command CommandDepthComparison) {
	gl.DepthFunc(command.Mode)
}

func (r *Renderer) executeCommandStencilTest(command CommandStencilTest) {
	if command.Enabled {
		gl.Enable(gl.STENCIL_TEST)
	} else {
		gl.Disable(gl.STENCIL_TEST)
	}
}

func (r *Renderer) executeCommandStencilOperation(command CommandStencilOperation) {
	gl.StencilOpSeparate(
		command.Face,
		command.StencilFail,
		command.DepthFail,
		command.Pass,
	)
}

func (r *Renderer) executeCommandStencilFunc(command CommandStencilFunc) {
	gl.StencilFuncSeparate(
		command.Face,
		command.Func,
		int32(command.Ref),
		command.Mask,
	)
}

func (r *Renderer) executeCommandStencilMask(command CommandStencilMask) {
	gl.StencilMaskSeparate(
		command.Face,
		command.Mask,
	)
}

func (r *Renderer) executeCommandColorWrite(command CommandColorWrite) {
	gl.ColorMask(command.Mask[0], command.Mask[1], command.Mask[2], command.Mask[3])
}

func (r *Renderer) executeCommandBlendColor(command CommandBlendColor) {
	gl.BlendColor(
		command.Color[0],
		command.Color[1],
		command.Color[2],
		command.Color[3],
	)
}

func (r *Renderer) executeCommandBlendEquation(command CommandBlendEquation) {
	gl.BlendEquationSeparate(
		command.ModeRGB,
		command.ModeAlpha,
	)
}

func (r *Renderer) executeCommandBlendFunc(command CommandBlendFunc) {
	gl.BlendFuncSeparate(
		command.SourceFactorRGB,
		command.DestinationFactorRGB,
		command.SourceFactorAlpha,
		command.DestinationFactorAlpha,
	)
}

func (r *Renderer) executeCommandBindVertexArray(command CommandBindVertexArray) {
	gl.BindVertexArray(command.VertexArrayID)
	r.indexType = command.IndexFormat
}

func (r *Renderer) executeCommandUniform1f(command CommandUniform1f) {
	gl.Uniform1f(
		command.Location,
		command.Value,
	)
}

func (r *Renderer) executeCommandUniform1i(command CommandUniform1i) {
	gl.Uniform1i(
		command.Location,
		command.Value,
	)
}

func (r *Renderer) executeCommandUniform3f(command CommandUniform3f) {
	gl.Uniform3f(
		command.Location,
		command.Values[0],
		command.Values[1],
		command.Values[2],
	)
}

func (r *Renderer) executeCommandUniform4f(command CommandUniform4f) {
	gl.Uniform4f(
		command.Location,
		command.Values[0],
		command.Values[1],
		command.Values[2],
		command.Values[3],
	)
}

func (r *Renderer) executeCommandUniformMatrix4f(command CommandUniformMatrix4f) {
	slice := command.Values[:]
	gl.UniformMatrix4fv(
		command.Location,
		1,
		false,
		&slice[0],
	)
}

func (r *Renderer) executeCommandTextureUnit(command CommandTextureUnit) {
	gl.BindTextureUnit(
		command.Index,
		command.TextureID,
	)
}

func (r *Renderer) executeCommandDraw(command CommandDraw) {
	gl.DrawArraysInstanced(
		r.primitive,
		command.VertexOffset,
		command.VertexCount,
		command.InstanceCount,
	)
}

func (r *Renderer) executeCommandDrawIndexed(command CommandDrawIndexed) {
	gl.DrawElementsInstanced(
		r.primitive,
		command.IndexCount,
		r.indexType,
		gl.PtrOffset(int(command.IndexOffset)),
		command.InstanceCount,
	)
}
