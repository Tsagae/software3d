package basics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFindWeights3(t *testing.T) {
	a := Vector3{-1, -1, 0}
	b := Vector3{1, -1, 0}
	c := Vector3{0, Sqrt(2), 0}
	target := Vector3{0, Sqrt(2), 0}
	w0, w1, w2 := FindWeights3(&a, &b, &c, &target)
	assert.True(t, w0.Equals(0))
	assert.True(t, w1.Equals(0))
	assert.True(t, w2.Equals(1))
}
