package internal

import (
	"fmt"
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

func (r *Renderer) BindPipeline(pipeline render.Pipeline) {
	intPipeline := pipeline.(*Pipeline)

	r.executeCommandBindPipeline(CommandBindPipeline{
		ProgramID:       intPipeline.ProgramID,
		Topology:        intPipeline.Topology,
		CullTest:        intPipeline.CullTest,
		FrontFace:       intPipeline.FrontFace,
		LineWidth:       intPipeline.LineWidth,
		DepthTest:       intPipeline.DepthTest,
		DepthWrite:      intPipeline.DepthWrite,
		DepthComparison: intPipeline.DepthComparison,
		VertexArrayID:   intPipeline.VertexArrayID,
		IndexFormat:     intPipeline.IndexFormat,
	})

	// 	if pipeline.StencilTest {
	// 		gl.Enable(gl.STENCIL_TEST)
	// 	} else {
	// 		gl.Disable(gl.STENCIL_TEST)
	// 	}

	// 	gl.Enable(gl.BLEND)
	// 	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// 	gl.ColorMask(pipeline.ColorWrite[0], pipeline.ColorWrite[1], pipeline.ColorWrite[2], pipeline.ColorWrite[3])
}

func (r *Renderer) Uniform4f(location render.UniformLocation, values [4]float32) {
	r.executeCommandUniform4f(CommandUniform4f{
		Location: location.(int32),
		Values:   values,
	})
}

func (r *Renderer) Uniform1i(location render.UniformLocation, value int) {
	r.executeCommandUniform1i(CommandUniform1i{
		Location: location.(int32),
		Value:    int32(value),
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
		case CommandKindUniform4f:
			command := PopCommand[CommandUniform4f](queue)
			r.executeCommandUniform4f(command)
		case CommandKindUniform1i:
			command := PopCommand[CommandUniform1i](queue)
			r.executeCommandUniform1i(command)
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
	r.executeCommandTopology(command.Topology)
	r.executeCommandCullTest(command.CullTest)
	r.executeCommandFrontFace(command.FrontFace)
	r.executeCommandLineWidth(command.LineWidth)
	r.executeCommandDepthTest(command.DepthTest)
	r.executeCommandDepthWrite(command.DepthWrite)
	r.executeCommandDepthComparison(command.DepthComparison)

	// if pipeline.StencilTest {
	// 	gl.Enable(gl.STENCIL_TEST)
	// } else {
	// 	gl.Disable(gl.STENCIL_TEST)
	// }

	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// gl.ColorMask(pipeline.ColorWrite[0], pipeline.ColorWrite[1], pipeline.ColorWrite[2], pipeline.ColorWrite[3])

	gl.UseProgram(command.ProgramID)
	gl.BindVertexArray(command.VertexArrayID)
	r.indexType = command.IndexFormat
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
	gl.LineWidth(command.Width)
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

func (r *Renderer) executeCommandUniform4f(command CommandUniform4f) {
	gl.Uniform4f(
		command.Location,
		command.Values[0],
		command.Values[1],
		command.Values[2],
		command.Values[3],
	)
}

func (r *Renderer) executeCommandUniform1i(command CommandUniform1i) {
	gl.Uniform1i(
		command.Location,
		command.Value,
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
	runtime.KeepAlive(slice)
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
