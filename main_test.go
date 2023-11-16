package main

import (
	"github.com/tsagae/software3d/pkg/graphics"
	"testing"
)

var zBuf graphics.ZBuffer = graphics.NewZBuffer(1024, 1024)

func BenchmarkZBuffer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		zBuf.Clear()
	}
}

func TestRun(t *testing.T) {

	if 0 != 0 {

		t.Error("run retuned error")
	}
}
