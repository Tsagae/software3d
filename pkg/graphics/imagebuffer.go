package graphics

import (
	"image"
	"image/color"
	"unsafe"
)

type ImageBuffer struct {
	innerImage image.RGBA
	width      int
	height     int
}

func NewImageBuffer(width int, height int) ImageBuffer {
	return ImageBuffer{image.RGBA{
		Pix:    make([]uint8, width*height*4),
		Stride: width,
		Rect:   image.Rect(0, 0, width, height),
	}, width, height}
}

func (iBuf *ImageBuffer) Get(x int, y int) color.RGBA {
	return *(*color.RGBA)(unsafe.Pointer(&iBuf.innerImage.Pix[(iBuf.width*y+x)*4]))
}

// Set sets the color of the buffer at x, y with the value c
func (iBuf *ImageBuffer) Set(x int, y int, c color.RGBA) {
	*(*color.RGBA)(unsafe.Pointer(&iBuf.innerImage.Pix[(iBuf.width*y+x)*4])) = c
}

func (iBuf *ImageBuffer) GetWidth() int {
	return iBuf.width
}

func (iBuf *ImageBuffer) GetHeight() int {
	return iBuf.height
}

func (iBuf *ImageBuffer) Clear() {
	pix := iBuf.innerImage.Pix
	for i := range pix {
		pix[i] = 0
	}
}

func (iBuf *ImageBuffer) GetImage() image.RGBA {
	return iBuf.innerImage
}
