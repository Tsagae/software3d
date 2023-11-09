package basics

import (
	"math"
)

const toDeg = 180 / math.Pi
const toRad = math.Pi / 180

func Sign(n Scalar) Scalar {
	if n < 0 {
		return -1
	}
	return 1
}

func Sqrt(n Scalar) Scalar {
	return Scalar(math.Sqrt(float64(n)))
}

func Abs(n Scalar) Scalar {
	return Scalar(math.Abs(float64(n)))
}

func DegToRad(n Scalar) Scalar {
	return n * toRad
}

func RadToDeg(n Scalar) Scalar {
	return n * toDeg
}

func Atan2(y Scalar, x Scalar) Scalar {
	return Scalar(math.Atan2(float64(y), float64(x)))
}

func Sin(n Scalar) Scalar {
	return Scalar(math.Sin(float64(n)))
}

func Cos(n Scalar) Scalar {
	return Scalar(math.Cos(float64(n)))
}

func Asin(n Scalar) Scalar {
	return Scalar(math.Asin(float64(n)))
}

func Acos(n Scalar) Scalar {
	return Scalar(math.Acos(float64(n)))
}

func Clamp(min Scalar, max Scalar, n Scalar) Scalar {
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}

// if n < max return n, return max otherwise
func ClampMax(max Scalar, n Scalar) Scalar {
	if n > max {
		return max
	}
	return n
}

// if n > min return n, return min otherwise
func ClampMin(min Scalar, n Scalar) Scalar {
	if n < min {
		return min
	}
	return n
}

func Pow(x Scalar, y Scalar) Scalar {
	return Scalar(math.Pow(float64(x), float64(y)))
}

func Floor(n Scalar) Scalar {
	return Scalar(math.Floor(float64(n)))
}

func Round(n Scalar) Scalar {
	return Scalar(math.Round(float64(n)))
}

func Ceil(n Scalar) Scalar {
	return Scalar(math.Ceil(float64(n)))
}
