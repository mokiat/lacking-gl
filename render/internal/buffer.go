package internal

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/mokiat/lacking/render"
)

func NewVertexBuffer(info render.BufferInfo) *Buffer {
	return newBuffer(info)
}

func NewIndexBuffer(info render.BufferInfo) *Buffer {
	return newBuffer(info)
}

func newBuffer(info render.BufferInfo) *Buffer {
	var id uint32
	gl.CreateBuffers(1, &id)

	flags := glBufferFlags(info.Dynamic)
	if info.Data != nil {
		gl.NamedBufferStorage(id, len(info.Data), gl.Ptr(&info.Data[0]), flags)
	} else {
		gl.NamedBufferStorage(id, info.Size, nil, flags)
	}
	return &Buffer{
		id: id,
	}
}

type Buffer struct {
	id uint32
}

func (b *Buffer) Update(info render.BufferUpdateInfo) {
	gl.NamedBufferSubData(b.id, info.Offset, len(info.Data), gl.Ptr(info.Data))
}

func (b *Buffer) Release() {
	gl.DeleteBuffers(1, &b.id)
	b.id = 0
}

func glBufferFlags(dynamic bool) uint32 {
	var flags uint32
	if dynamic {
		flags |= gl.DYNAMIC_STORAGE_BIT
	}
	return flags
}
