package graphics

import (
	"fmt"
	"image/color"
	"testing"
	time2 "time"
)

func BenchmarkImageBuffer_Clear(b *testing.B) {
	imageBuffer := NewImageBuffer(800, 600)
	fmt.Println("---------------Benchmark start---------------")
	b.ResetTimer()
	time := time2.Now()
	imageBuffer.Clear()
	fmt.Printf("clearing image buffer: %v\n", time2.Now().Sub(time))
	fmt.Println("----------------Benchmark end----------------")
}

func BenchmarkImageBuffer_Get(b *testing.B) {
	imageBuffer := NewImageBuffer(800, 600)
	fmt.Println("---------------Benchmark start---------------")
	b.ResetTimer()
	time := time2.Now()
	for i := 0; i < b.N; i++ {
		imageBuffer.Get(0, 0)
	}
	fmt.Printf("accessing %v pixels: %v\n", b.N, time2.Now().Sub(time))
	fmt.Println("----------------Benchmark end----------------")
}

func BenchmarkImageBuffer_Set(b *testing.B) {
	imageBuffer := NewImageBuffer(800, 600)
	fmt.Println("---------------Benchmark start---------------")
	b.ResetTimer()
	time := time2.Now()
	for i := 0; i < b.N; i++ {
		imageBuffer.Set(0, 0, color.RGBA{
			R: 255,
			G: 255,
			B: 255,
			A: 255,
		})
	}
	fmt.Printf("setting %v pixels: %v\n", b.N, time2.Now().Sub(time))
	fmt.Println("----------------Benchmark end----------------")
}

func BenchmarkImageBufferFillImage(b *testing.B) {
	imageBuffer := NewImageBuffer(800, 600)
	fmt.Println("---------------Benchmark start---------------")
	b.ResetTimer()
	time := time2.Now()
	for y := 0; y <= 600; y++ {
		for x := 0; x <= 800; x++ {
			imageBuffer.Set(x, y, color.RGBA{
				R: 255,
				G: 255,
				B: 255,
				A: 255,
			})
		}
	}
	fmt.Printf("filling image buffer: %v\n", time2.Now().Sub(time))
	fmt.Println("----------------Benchmark end----------------")
}
