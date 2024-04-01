package renderer

import (
	"fmt"
	"github.com/tsagae/software3d/pkg/basics"
	"github.com/tsagae/software3d/pkg/graphics"
	"testing"
	time2 "time"
)

func BenchmarkRasterRenderer_rasterTriangle(b *testing.B) {
	vertices := [3]basics.Vector3{
		{
			X: 0,
			Y: 0,
			Z: 0,
		}, {
			X: 400,
			Y: 600,
			Z: 0,
		}, {
			X: 800,
			Y: 0,
			Z: 0,
		},
	}
	colors := [3]basics.Vector3{
		{
			X: 65535,
		}, {
			Y: 65535,
		}, {
			Z: 65535,
		},
	}
	triangle := graphics.NewTriangle(vertices, colors)
	renderer := NewRasterRenderer(nil, 1, 800, 600)

	fmt.Println("---------------Benchmark start---------------")
	b.ResetTimer()
	time := time2.Now()

	for i := 0; i < b.N; i++ {
		renderer.rasterTriangle(triangle)
	}

	fmt.Printf("rasterizing %v triangles: %v\n", b.N, time2.Now().Sub(time))
	fmt.Println("----------------Benchmark end----------------")
}
