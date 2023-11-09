package main

/*
import (
	"GoSDL/pkg/basics"
	"fmt"
	"math"
)

type Plane struct {
	p *basics.Vector3 //point
	n *basics.Vector3 //vector
}

type Ray struct {
	start *basics.Vector3 //point
	dir   *basics.Vector3 //vector
}

// TESTS
func unitTestSumIsCommutative() bool {
	v := basics.NewVector3(1, 2, 3)
	w := basics.NewVector3(4, 5, 6)
	return v.Add(w).Equals(w.Add(v))
}

func unitTestCross() bool {
	v := basics.NewVector3(1, 2, 3)
	w := basics.NewVector3(4, 5, 6)
	r1 := v.Cross(w).Equals(w.Cross(v).Mul(-1))
	// Resulting vector is orthogonal to both
	r2 := v.Cross(w).Dot(v).IsZero()
	r3 := v.Cross(w).Dot(w).IsZero()
	return r1 && r2 && r3
}

// Reflect ray v off surface with normal n
func rayReflection(v *basics.Vector3, n *basics.Vector3) *basics.Vector3 {
	//return v - (2*dot(n, v))*n
	return v.Sub(n.Mul(n.Dot(v) * 2))
}

func intersection(plane *Plane, ray *Ray) *basics.Vector3 {
	var k basics.Scalar = plane.p.Sub(ray.start).Dot(plane.n) / ray.dir.Dot(plane.n)
	// todo: avoid division by 0
	return ray.start.Add(ray.dir.Mul(k))
}

func test() {
	   //v[0] = 2;
	   //Scalar foo = v[0]+5;
	   //v = v + w;
	   //unitTestSumIsCommutative();
	   //unitTestCross();
	   //std::cout << rayReflection(v , w).x << std::endl;
	v := basics.NewVector3(0, 0, 0)
	w := basics.NewVector3(2, 3, 5)
	v.X = 2
	v = v.Add(w)
	fmt.Println(unitTestSumIsCommutative())
	fmt.Println(unitTestCross())
	fmt.Println(rayReflection(v, w).X)
}

func quatTest() {
	v := basics.NewVector3(1, 0, 0)
	q1 := basics.NewRawQuaternion(basics.Scalar(math.Cos(math.Pi/4)), *basics.NewVector3(0, 0, 1).Mul(basics.Scalar(math.Sin(math.Pi / 4))))
	q1.ThisNormalize()
	q2 := basics.NewRawQuaternion(0.96639055, *basics.NewVector3(0, 0, 0.25707844))
	q2.ThisNormalize()
	q3 := basics.NewRawQuaternion(0.39401576, *basics.NewVector3(0, 0, -0.9191037))
	q3.ThisNormalize()

	vcopy := *v
	q1.Rotate(&vcopy)
	fmt.Println(vcopy)
	vcopy = *v
	q2.Rotate(&vcopy)
	fmt.Println(vcopy)
	vcopy = *v
	q3.Rotate(&vcopy)
	fmt.Println(vcopy)

}

func angleTest() {
	h := basics.Scalar(2)
	a := basics.Scalar(16) / basics.Scalar(9)
	hFovDeg := basics.Scalar(120)
	hFovRad := hFovDeg * math.Pi / 180
	beta := (math.Pi/2 - hFovRad/2)
	fmt.Println("beta: ", beta)
	h2s := (h / 2) * (h / 2)
	z := math.Sqrt(float64((h2s*a*a)/(beta*beta) - (h2s * (a*a + 1))))
	p1 := basics.NewVector3(h*a/2, h/2, basics.Scalar(z))
	gamma := math.Acos(float64(p1.Z * (1 / p1.Length())))
	vFovRad := gamma * 2
	fmt.Println("vfovRad: ", vFovRad)
	fmt.Println("hFovDeg: ", hFovDeg)
	fmt.Println("vFovDeg: ", vFovRad*180/math.Pi)
}

*/
