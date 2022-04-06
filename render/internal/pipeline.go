package internal

import (
	"fmt"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/mokiat/lacking/render"
)

func NewPipeline(info render.PipelineInfo) *Pipeline {
	intProgram := info.Program.(*Program)
	intVertexArray := info.VertexArray.(*VertexArray)

	pipeline := &Pipeline{
		ProgramID: intProgram.id,
		VertexArray: CommandBindVertexArray{
			VertexArrayID: intVertexArray.id,
			IndexFormat:   intVertexArray.indexFormat,
		},
	}

	switch info.Topology {
	case render.TopologyPoints:
		pipeline.Topology.Topology = gl.POINTS
	case render.TopologyLineStrip:
		pipeline.Topology.Topology = gl.LINE_STRIP
	case render.TopologyLineLoop:
		pipeline.Topology.Topology = gl.LINE_LOOP
	case render.TopologyLines:
		pipeline.Topology.Topology = gl.LINES
	case render.TopologyTriangleStrip:
		pipeline.Topology.Topology = gl.TRIANGLE_STRIP
	case render.TopologyTriangleFan:
		pipeline.Topology.Topology = gl.TRIANGLE_FAN
	case render.TopologyTriangles:
		pipeline.Topology.Topology = gl.TRIANGLES
	}

	switch info.Culling {
	case render.CullModeNone:
		pipeline.CullTest.Enabled = false
	case render.CullModeBack:
		pipeline.CullTest.Enabled = true
		pipeline.CullTest.Face = gl.BACK
	case render.CullModeFront:
		pipeline.CullTest.Enabled = true
		pipeline.CullTest.Face = gl.FRONT
	case render.CullModeFrontAndBack:
		pipeline.CullTest.Enabled = true
		pipeline.CullTest.Face = gl.FRONT_AND_BACK
	}

	switch info.FrontFace {
	case render.FaceOrientationCCW:
		pipeline.FrontFace.Orientation = gl.CCW
	case render.FaceOrientationCW:
		pipeline.FrontFace.Orientation = gl.CW
	}

	pipeline.LineWidth.Width = info.LineWidth

	pipeline.DepthTest.Enabled = info.DepthTest
	pipeline.DepthWrite.Enabled = info.DepthWrite
	pipeline.DepthComparison.Mode = glEnumFromComparison(info.DepthComparison)

	pipeline.StencilTest.Enabled = info.StencilTest

	pipeline.StencilOpFront.Face = gl.FRONT
	pipeline.StencilOpFront.StencilFail = glEnumFromStencilOp(info.StencilFront.StencilFailOp)
	pipeline.StencilOpFront.DepthFail = glEnumFromStencilOp(info.StencilFront.DepthFailOp)
	pipeline.StencilOpFront.Pass = glEnumFromStencilOp(info.StencilFront.PassOp)

	pipeline.StencilOpBack.Face = gl.BACK
	pipeline.StencilOpBack.StencilFail = glEnumFromStencilOp(info.StencilBack.StencilFailOp)
	pipeline.StencilOpBack.DepthFail = glEnumFromStencilOp(info.StencilBack.DepthFailOp)
	pipeline.StencilOpBack.Pass = glEnumFromStencilOp(info.StencilBack.PassOp)

	pipeline.StencilFuncFront.Face = gl.FRONT
	pipeline.StencilFuncFront.Func = glEnumFromComparison(info.StencilFront.Comparison)
	pipeline.StencilFuncFront.Ref = info.StencilFront.Reference
	pipeline.StencilFuncFront.Mask = info.StencilFront.ComparisonMask

	pipeline.StencilFuncBack.Face = gl.BACK
	pipeline.StencilFuncBack.Func = glEnumFromComparison(info.StencilBack.Comparison)
	pipeline.StencilFuncBack.Ref = info.StencilBack.Reference
	pipeline.StencilFuncBack.Mask = info.StencilBack.ComparisonMask

	pipeline.StencilMaskFront.Face = gl.FRONT
	pipeline.StencilMaskFront.Mask = info.StencilFront.WriteMask

	pipeline.StencilMaskBack.Face = gl.BACK
	pipeline.StencilMaskBack.Mask = info.StencilBack.WriteMask

	pipeline.ColorWrite.Mask = info.ColorWrite

	pipeline.BlendEnabled = info.BlendEnabled
	pipeline.BlendColor.Color = [4]float32{ // TODO: Add ToArray method on sprec
		info.BlendColor.X,
		info.BlendColor.Y,
		info.BlendColor.Z,
		info.BlendColor.W,
	}

	return pipeline
}

func glEnumFromComparison(comparison render.Comparison) uint32 {
	switch comparison {
	case render.ComparisonNever:
		return gl.NEVER
	case render.ComparisonLess:
		return gl.LESS
	case render.ComparisonEqual:
		return gl.EQUAL
	case render.ComparisonLessOrEqual:
		return gl.LEQUAL
	case render.ComparisonGreater:
		return gl.GREATER
	case render.ComparisonNotEqual:
		return gl.NOTEQUAL
	case render.ComparisonGreaterOrEqual:
		return gl.GEQUAL
	case render.ComparisonAlways:
		return gl.ALWAYS
	default:
		panic(fmt.Errorf("unknown comparison: %d", comparison))
	}
}

func glEnumFromStencilOp(op render.StencilOperation) uint32 {
	switch op {
	case render.StencilOperationKeep:
		return gl.KEEP
	case render.StencilOperationZero:
		return gl.ZERO
	case render.StencilOperationReplace:
		return gl.REPLACE
	case render.StencilOperationIncrease:
		return gl.INCR
	case render.StencilOperationIncreaseWrap:
		return gl.INCR_WRAP
	case render.StencilOperationDecrease:
		return gl.DECR
	case render.StencilOperationDecreaseWrap:
		return gl.DECR_WRAP
	case render.StencilOperationInvert:
		return gl.INVERT
	default:
		panic(fmt.Errorf("unknown op: %d", op))
	}
}

type Pipeline struct {
	ProgramID        uint32
	Topology         CommandTopology
	CullTest         CommandCullTest
	FrontFace        CommandFrontFace
	LineWidth        CommandLineWidth
	DepthTest        CommandDepthTest
	DepthWrite       CommandDepthWrite
	DepthComparison  CommandDepthComparison
	StencilTest      CommandStencilTest
	StencilOpFront   CommandStencilOperation
	StencilOpBack    CommandStencilOperation
	StencilFuncFront CommandStencilFunc
	StencilFuncBack  CommandStencilFunc
	StencilMaskFront CommandStencilMask
	StencilMaskBack  CommandStencilMask
	ColorWrite       CommandColorWrite
	BlendEnabled     bool
	BlendColor       CommandBlendColor
	VertexArray      CommandBindVertexArray
}

func (p *Pipeline) Release() {
}
