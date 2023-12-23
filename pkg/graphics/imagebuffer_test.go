package graphics

import (
	"fmt"
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
