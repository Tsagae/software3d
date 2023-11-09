package graphics

import (
	"GoSDL/pkg/basics"
	"math"
)

type ZBuffer struct {
	buffer []basics.Scalar
	width  int
	height int
}

func NewZBuffer(width int, height int) ZBuffer {
	return ZBuffer{make([]basics.Scalar, width*height), width, height}
}

func (z *ZBuffer) Set(x int, y int, val basics.Scalar) {
	z.buffer[y*z.width+x] = val
}

func (z *ZBuffer) Get(x int, y int) basics.Scalar {
	return z.buffer[y*z.width+x]
}

func (z *ZBuffer) GetWidth() int {
	return z.width
}

func (z *ZBuffer) GetHeight() int {
	return z.height
}

func (z *ZBuffer) Clear() {
	//maybe reallocate the whole slice?
	inf := basics.Scalar(math.Inf(+1))
	lenZbuf := len(z.buffer)
	zBuf := z.buffer

	for i := 0; i < lenZbuf; i++ {
		zBuf[i] = inf
	}
}
