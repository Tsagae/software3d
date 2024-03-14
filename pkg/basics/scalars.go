package basics

import "math"

type Scalar float64

const epsilon Scalar = 1e-5

// IsZero Equality and zero
func (a Scalar) IsZero() bool {
	return Scalar(math.Abs(float64(a))) < epsilon
}

func (a Scalar) Equals(b Scalar) bool {
	return (a - b).IsZero()
}

func (a Scalar) IsNaN() bool {
	return math.IsNaN(float64(a))
}
