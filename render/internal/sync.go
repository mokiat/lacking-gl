package internal

import (
	"github.com/go-gl/gl/v4.6-core/gl"
	"github.com/mokiat/lacking/render"
)

func NewFence() *Fence {
	return &Fence{
		id: gl.FenceSync(gl.SYNC_GPU_COMMANDS_COMPLETE, 0),
	}
}

type Fence struct {
	render.FenceObject
	id uintptr
}

func (f *Fence) Status() render.FenceStatus {
	switch gl.ClientWaitSync(f.id, gl.SYNC_FLUSH_COMMANDS_BIT, 0) {
	case gl.ALREADY_SIGNALED, gl.CONDITION_SATISFIED:
		return render.FenceStatusSuccess
	case gl.TIMEOUT_EXPIRED:
		return render.FenceStatusNotReady
	default:
		return render.FenceStatusDeviceLost
	}
}

func (f *Fence) Delete() {
	gl.DeleteSync(f.id)
	f.id = 0
}
