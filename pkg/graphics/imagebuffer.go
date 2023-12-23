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

func (iBuf *ImageBuffer) Get(x int, y int) color.Color { // image.Image interface
	return iBuf.innerImage.At(x, y) // TODO change the image's underlying slice instead
}

func (iBuf *ImageBuffer) Set(x int, y int, c color.Color) { // draw.Image interface
	iBuf.innerImage.Set(x, y, c) // TODO change the image's underlying slice instead
}

func (iBuf *ImageBuffer) GetWidth() int {
	return iBuf.width
}

func (iBuf *ImageBuffer) GetHeight() int {
	return iBuf.height
}

func (iBuf *ImageBuffer) Clear() {
	pix := iBuf.innerImage.Pix
	for i, _ := range pix {
		pix[i] = 0
	}
}

func (iBuf *ImageBuffer) GetImage() *image.RGBA {
	return iBuf.innerImage
}
