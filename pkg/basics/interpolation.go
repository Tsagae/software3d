package basics

// Allows extrapolation
func LerpVector3(a *Vector3, b *Vector3, t Scalar) Vector3 {
	//a + t(b - a) allows extrapolation
	tba := b.Sub(a).Mul(t)
	return a.Add(&tba)
}

func NLerpVector3(a *Vector3, b *Vector3, t Scalar) Vector3 {
	v := LerpVector3(a, b, t)
	v.ThisNormalize()
	return v
}

func mixQuaternion(a Quaternion, b Quaternion, t Scalar) Quaternion {
	/*
		Scalar d = sign(dot(a.im, b.im) + a.re * b.re); // shortest path!
		Quaternion res = a * (1 - t) + b * t * d; // linear interpolation...
		return normalize(res);  // ...then re-normalized (NLERP)
	*/
	d := Sign(a.Im.Dot(&b.Im) + a.Re*b.Re)
	a.ThisMulScalar(1 - t)
	b.ThisMulScalar(t * d)

	a.ThisAdd(&b)
	a.ThisNormalize()
	return a
}
