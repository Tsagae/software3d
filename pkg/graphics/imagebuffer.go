package graphics

import (
	"image"
	"image/color"
)

type ImageBuffer struct {
	innerImage   *image.RGBA
	width        int
	height       int
	defaultColor color.Color
}

func NewImageBuffer(width int, height int, defaultColor color.Color) ImageBuffer {
	return ImageBuffer{image.NewRGBA(image.Rect(0, 0, width, height)), width, height, defaultColor}
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
	for y := 0; y < iBuf.height; y++ {
		for x := 0; x < iBuf.width; x++ {
			iBuf.Set(x, y, color.Black)
		}
	}
}

func (iBuf *ImageBuffer) GetImage() *image.RGBA {
	return iBuf.innerImage
}
