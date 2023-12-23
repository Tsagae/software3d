package graphics

import (
	"image"
	"image/color"
)

type ImageBuffer struct {
	innerImage *image.RGBA
	width      int
	height     int
}

func NewImageBuffer(width int, height int) ImageBuffer {
	return ImageBuffer{image.NewRGBA(image.Rect(0, 0, width, height)), width, height}
}

func (iBuf *ImageBuffer) Get(x int, y int) color.RGBA {
	pixelAddr := iBuf.width*y + x
	pixel := iBuf.innerImage.Pix[pixelAddr : pixelAddr+4 : pixelAddr+4]
	return color.RGBA{
		R: pixel[0],
		G: pixel[1],
		B: pixel[2],
		A: pixel[3],
	}
}

// Set sets the color of the buffer at x, y with the value c
func (iBuf *ImageBuffer) Set(x int, y int, c color.RGBA) {
	pixelAddr := iBuf.width*y + x
	pixel := iBuf.innerImage.Pix[pixelAddr : pixelAddr+4 : pixelAddr+4]
	pixel[0] = c.R
	pixel[1] = c.G
	pixel[2] = c.B
	pixel[3] = c.A
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

func (iBuf *ImageBuffer) GetImage() *image.RGBA {
	return iBuf.innerImage
}
