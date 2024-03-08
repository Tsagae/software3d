package basics

// LerpVector3 Allows extrapolation
func LerpVector3(a *Vector3, b *Vector3, t Scalar) Vector3 {
	//a + t(b - a) allows extrapolation
	return a.Add(b.Sub(*a).Mul(t))
}

func NLerpVector3(a *Vector3, b *Vector3, t Scalar) Vector3 {
	return LerpVector3(a, b, t).Normalized()
}

func mixQuaternion(a Quaternion, b Quaternion, t Scalar) Quaternion {
	/*
		Scalar d = sign(dot(a.im, b.im) + a.re * b.re); // shortest path!
		Quaternion res = a * (1 - t) + b * t * d; // linear interpolation...
		return normalize(res);  // ...then re-normalized (NLERP)
	*/
	d := Sign(a.Im.Dot(b.Im) + a.Re*b.Re)
	a.ThisMulScalar(1 - t)
	b.ThisMulScalar(t * d)

	a.ThisAdd(&b)
	a.ThisNormalize()
	return a
}

func Interpolate3(v1, v2, v3 *Vector3, w1, w2, w3 Scalar) Vector3 {
	return v1.Mul(w1).Add(v2.Mul(w2)).Add(v3.Mul(w3))
}

func FindWeights3(v1, v2, v3, target *Vector3) (Scalar, Scalar, Scalar) {
	// most of this can be cached when finding weights inside the same triangle TODO
	// calculate vectors from point f to vertices p1, p2 and p3:
	f1 := v1.Sub(*target)
	f2 := v2.Sub(*target)
	f3 := v3.Sub(*target)
	// calculate the areas and factors (order of parameters doesn't matter):
	var a = Vector3.Cross(v1.Sub(*v2), v1.Sub(*v3)).Length() // main triangle area a
	var a1 = Vector3.Cross(f2, f3).Length() / a              // p1's triangle area / a
	var a2 = Vector3.Cross(f3, f1).Length() / a              // p2's triangle area / a
	var a3 = Vector3.Cross(f1, f2).Length() / a              // p3's triangle area / a
	return a1, a2, a3
}

func FindWeights2D(v1, v2, v3, target *Vector3) (Scalar, Scalar, Scalar) {
	// most of this can be cached when finding weights inside the same triangle TODO
	den := (v2.Y-v3.Y)*(v1.X-v3.X) + (v3.X-v2.X)*(v1.Y-v3.Y)
	t1 := target.X - v3.X
	t2 := target.Y - v3.Y

	w1 := ((v2.Y-v3.Y)*t1 + (v3.X-v2.X)*t2) / den
	w2 := ((v3.Y-v1.Y)*t1 + (v1.X-v3.X)*t2) / den
	w3 := 1 - w1 - w2
	return w1, w2, w3
}
