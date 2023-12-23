package graphics

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tsagae/software3d/pkg/basics"
	"math"
	"testing"
	time2 "time"
)

func TestZBuffer(t *testing.T) {
	zBuf := NewZBuffer(10, 10)
	value := basics.Scalar(15.23)
	zBuf.Set(5, 5, value)
	assert.Equal(t, value, zBuf.Get(5, 5), "Value is not set")

	zBuf.Set(0, 0, 10)
	zBuf.Clear()
	inf := basics.Scalar(math.Inf(+1))
	assert.Equal(t, inf, zBuf.Get(0, 0), "Clear does not clear the buffer")
	assert.Equal(t, inf, zBuf.Get(5, 9), "Clear does not clear the buffer")
}

func BenchmarkZBuffer_Clear(b *testing.B) {
	zBuf := NewZBuffer(800, 600)
	fmt.Println("---------------Benchmark start---------------")
	b.ResetTimer()
	time := time2.Now()
	zBuf.Clear()
	fmt.Printf("clearing zbuffer: %v\n", time2.Now().Sub(time))
	fmt.Println("----------------Benchmark end----------------")
}

func BenchmarkZBuffer_Get(b *testing.B) {
	zBuf := NewZBuffer(800, 600)
	fmt.Println("---------------Benchmark start---------------")
	b.ResetTimer()
	time := time2.Now()
	for i := 0; i < b.N; i++ {
		zBuf.Get(0, 0)
	}
	fmt.Printf("accessing %v pixels: %v\n", b.N, time2.Now().Sub(time))
	fmt.Println("----------------Benchmark end----------------")
}

func BenchmarkZBuffer_Set(b *testing.B) {
	zBuf := NewZBuffer(800, 600)
	valueToSet := basics.Scalar(10)
	fmt.Println("---------------Benchmark start---------------")
	b.ResetTimer()
	time := time2.Now()
	for i := 0; i < b.N; i++ {
		zBuf.Set(0, 0, valueToSet)
	}
	fmt.Printf("setting %v pixels: %v\n", b.N, time2.Now().Sub(time))
	fmt.Println("----------------Benchmark end----------------")
}

func BenchmarkZBuffer_FillBuffer(b *testing.B) {
	zBuf := NewZBuffer(800, 600)
	valueToSet := basics.Scalar(10)
	fmt.Println("---------------Benchmark start---------------")
	b.ResetTimer()
	time := time2.Now()
	for y := 0; y < 600; y++ {
		for x := 0; x < 800; x++ {
			zBuf.Set(x, y, valueToSet)
		}
	}
	fmt.Printf("filling zbuffer: %v\n", time2.Now().Sub(time))
	fmt.Println("----------------Benchmark end----------------")
}
