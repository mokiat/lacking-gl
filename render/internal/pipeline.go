package internal

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/mokiat/lacking/render"
)

func NewPipeline(info render.PipelineInfo) *Pipeline {
	intProgram := info.Program.(*Program)
	intVertexArray := info.VertexArray.(*VertexArray)

	pipeline := &Pipeline{
		ProgramID:     intProgram.id,
		VertexArrayID: intVertexArray.id,
		IndexFormat:   intVertexArray.indexFormat,
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
	switch info.DepthComparison {
	case render.ComparisonNever:
		pipeline.DepthComparison.Mode = gl.NEVER
	case render.ComparisonLess:
		pipeline.DepthComparison.Mode = gl.LESS
	case render.ComparisonEqual:
		pipeline.DepthComparison.Mode = gl.EQUAL
	case render.ComparisonLessOrEqual:
		pipeline.DepthComparison.Mode = gl.LEQUAL
	case render.ComparisonGreater:
		pipeline.DepthComparison.Mode = gl.GREATER
	case render.ComparisonNotEqual:
		pipeline.DepthComparison.Mode = gl.NOTEQUAL
	case render.ComparisonGreaterOrEqual:
		pipeline.DepthComparison.Mode = gl.GEQUAL
	case render.ComparisonAlways:
		pipeline.DepthComparison.Mode = gl.ALWAYS
	}

	return pipeline
}

type Pipeline struct {
	ProgramID       uint32
	Topology        CommandTopology
	CullTest        CommandCullTest
	FrontFace       CommandFrontFace
	LineWidth       CommandLineWidth
	DepthTest       CommandDepthTest
	DepthWrite      CommandDepthWrite
	DepthComparison CommandDepthComparison
	VertexArrayID   uint32
	IndexFormat     uint32
}

func (p *Pipeline) Release() {
}
