package internal

import (
	"unsafe"

	"github.com/mokiat/lacking/render"
)

func NewCommandQueue() *CommandQueue {
	return &CommandQueue{
		data: make([]byte, 1024*1024), // TODO: Start smaller and allow for growth
	}
}

type CommandQueue struct {
	data        []byte
	writeOffset uintptr
	readOffset  uintptr
}

func (q *CommandQueue) Reset() {
	q.readOffset = 0
	q.writeOffset = 0
}

func (q *CommandQueue) BindPipeline(pipeline render.Pipeline) {
	PushCommand(q, CommandHeader{
		Kind: CommandKindBindPipeline,
	})
	intPipeline := pipeline.(*Pipeline)
	PushCommand(q, CommandBindPipeline{
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
}

func (q *CommandQueue) Uniform4f(location render.UniformLocation, values [4]float32) {
	PushCommand(q, CommandHeader{
		Kind: CommandKindUniform4f,
	})
	PushCommand(q, CommandUniform4f{
		Location: location.(int32),
		Values:   values,
	})
}

func (q *CommandQueue) Uniform1i(location render.UniformLocation, value int) {
	PushCommand(q, CommandHeader{
		Kind: CommandKindUniform1i,
	})
	PushCommand(q, CommandUniform1i{
		Location: location.(int32),
		Value:    int32(value),
	})
}

func (q *CommandQueue) UniformMatrix4f(location render.UniformLocation, values [16]float32) {
	PushCommand(q, CommandHeader{
		Kind: CommandKindUniformMatrix4f,
	})
	PushCommand(q, CommandUniformMatrix4f{
		Location: location.(int32),
		Values:   values,
	})
}

func (q *CommandQueue) TextureUnit(index int, texture render.Texture) {
	PushCommand(q, CommandHeader{
		Kind: CommandKindTextureUnit,
	})
	PushCommand(q, CommandTextureUnit{
		Index:     uint32(index),
		TextureID: texture.(*Texture).id,
	})
}

func (q *CommandQueue) Draw(vertexOffset, vertexCount, instanceCount int) {
	PushCommand(q, CommandHeader{
		Kind: CommandKindDraw,
	})
	PushCommand(q, CommandDraw{
		VertexOffset:  int32(vertexOffset),
		VertexCount:   int32(vertexCount),
		InstanceCount: int32(instanceCount),
	})
}

func (q *CommandQueue) DrawIndexed(indexOffset, indexCount, instanceCount int) {
	PushCommand(q, CommandHeader{
		Kind: CommandKindDrawIndexed,
	})
	PushCommand(q, CommandDraw{
		VertexOffset:  int32(indexOffset),
		VertexCount:   int32(indexCount),
		InstanceCount: int32(instanceCount),
	})
}

func (q *CommandQueue) Release() {
}

func MoreCommands(queue *CommandQueue) bool {
	return queue.writeOffset > queue.readOffset
}

func PushCommand[T any](queue *CommandQueue, command T) {
	target := (*T)(unsafe.Add(unsafe.Pointer(&queue.data[0]), queue.writeOffset))
	*target = command
	queue.writeOffset += unsafe.Sizeof(command)
}

func PopCommand[T any](queue *CommandQueue) T {
	target := (*T)(unsafe.Add(unsafe.Pointer(&queue.data[0]), queue.readOffset))
	command := *target
	queue.readOffset += unsafe.Sizeof(command)
	return command
}

type CommandKind uint8

const (
	CommandKindBindPipeline CommandKind = iota
	CommandKindTopology
	CommandKindCullTest
	CommandKindFrontFace
	CommandKindLineWidth
	CommandKindDepthTest
	CommandKindDepthWrite
	CommandKindDepthComparison
	CommandKindUniform4f
	CommandKindUniform1i
	CommandKindUniformMatrix4f
	CommandKindTextureUnit
	CommandKindDraw
	CommandKindDrawIndexed
)

type CommandHeader struct {
	Kind CommandKind
}

type CommandBindPipeline struct {
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

	// StencilTest                 bool
	// StencilFront                StencilOperationState
	// StencilBack                 StencilOperationState
	// ColorWrite                  [4]bool
	// BlendEnabled                bool
	// BlendColor                  sprec.Vec4
	// BlendSourceColorFactor      BlendFactor
	// BlendDestinationColorFactor BlendFactor
	// BlendSourceAlphaFactor      BlendFactor
	// BlendDestinationAlphaFactor BlendFactor
	// BlendOpColor                BlendOperation
	// BlendOpAlpha                BlendOperation
}

type CommandTopology struct {
	Topology uint32
}

type CommandCullTest struct {
	Enabled bool
	Face    uint32
}

type CommandFrontFace struct {
	Orientation uint32
}

type CommandLineWidth struct {
	Width float32
}

type CommandDepthTest struct {
	Enabled bool
}

type CommandDepthWrite struct {
	Enabled bool
}

type CommandDepthComparison struct {
	Mode uint32
}

type CommandUniform4f struct {
	Location int32
	Values   [4]float32
}

type CommandUniform1i struct {
	Location int32
	Value    int32
}

type CommandUniformMatrix4f struct {
	Location int32
	Values   [16]float32
}

type CommandTextureUnit struct {
	Index     uint32
	TextureID uint32
}

type CommandDraw struct {
	VertexOffset  int32
	VertexCount   int32
	InstanceCount int32
}

type CommandDrawIndexed struct {
	IndexOffset   int32
	IndexCount    int32
	InstanceCount int32
}
