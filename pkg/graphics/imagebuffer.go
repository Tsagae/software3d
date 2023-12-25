package graphics

import (
	"image/color"
)

type RGB struct {
	R uint8
	G uint8
	B uint8
}

type ImageBuffer struct {
	innerImage []RGB
	width      int
	height     int
}

func NewImageBuffer(width int, height int) ImageBuffer {
	return ImageBuffer{
		make([]RGB, width*height),
		width,
		height,
	}
}

// Get gets the color of the buffer at x, y with 255 alpha. There is no check for out of bounds values for efficiency reasons
func (iBuf *ImageBuffer) Get(x int, y int) color.RGBA {
	rgbColor := iBuf.innerImage[iBuf.width*y+x]
	return color.RGBA{
		R: rgbColor.R,
		G: rgbColor.G,
		B: rgbColor.B,
		A: 255,
	}
}

// Set sets the color of the buffer at x, y with the value c. There is no check for out of bounds values for efficiency reasons
func (iBuf *ImageBuffer) Set(x int, y int, c color.RGBA) {
	//iBuf.innerImage[iBuf.width*y+x] = *(*RGB)(unsafe.Pointer(&c)) // it looks like this is not more efficient
	iBuf.innerImage[iBuf.width*y+x] = RGB{
		R: c.R,
		G: c.G,
		B: c.B,
	}
}

func (iBuf *ImageBuffer) Width() int {
	return iBuf.width
}

func (iBuf *ImageBuffer) Height() int {
	return iBuf.height
}

func (iBuf *ImageBuffer) Clear() {
	pix := iBuf.innerImage
	for i := range pix {
		pix[i] = RGB{}
	}
}

func (iBuf *ImageBuffer) GetImage() []RGB {
	return iBuf.innerImage
}
