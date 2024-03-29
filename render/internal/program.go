package internal

import (
	"errors"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/mokiat/lacking/log"
	"github.com/mokiat/lacking/render"
)

func NewProgram(info render.ProgramInfo) *Program {
	program := &Program{
		id: gl.CreateProgram(),
	}
	if vertexShader, ok := info.VertexShader.(*Shader); ok {
		gl.AttachShader(program.id, vertexShader.id)
		defer gl.DetachShader(program.id, vertexShader.id)
	}
	if fragmentShader, ok := info.FragmentShader.(*Shader); ok {
		gl.AttachShader(program.id, fragmentShader.id)
		defer gl.DetachShader(program.id, fragmentShader.id)
	}
	if err := program.link(); err != nil {
		log.Error("Program link error: %v", err)
	}
	// NOTE: Texture bindings are to be done in GLSL through
	// `layout(binding = 2) uniform ...`.
	// NOTE: Buffer bindings are to be done in GLSL through
	// `layout(binding = 2, std140) uniform ...`
	return program
}

type Program struct {
	render.ProgramObject
	id uint32
}

func (p *Program) UniformLocation(name string) render.UniformLocation {
	nullTerminatedName := name + "\x00"
	result := gl.GetUniformLocation(p.id, gl.Str(nullTerminatedName))
	runtime.KeepAlive(nullTerminatedName)
	return result
}

func (p *Program) Release() {
	gl.DeleteProgram(p.id)
	p.id = 0
}

func (p *Program) link() error {
	gl.LinkProgram(p.id)
	if !p.isLinkSuccessful() {
		return errors.New(p.getInfoLog())
	}
	return nil
}

func (p *Program) isLinkSuccessful() bool {
	var status int32
	gl.GetProgramiv(p.id, gl.LINK_STATUS, &status)
	return status != gl.FALSE
}

func (p *Program) getInfoLog() string {
	var logLength int32
	gl.GetProgramiv(p.id, gl.INFO_LOG_LENGTH, &logLength)

	log := strings.Repeat("\x00", int(logLength+1))
	gl.GetProgramInfoLog(p.id, logLength, nil, gl.Str(log))
	runtime.KeepAlive(log)
	return log
}
