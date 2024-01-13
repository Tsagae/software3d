package basics

type Plane struct {
	point     Vector3
	normal    Vector3
	coplanarA Vector3
	coplanarB Vector3
}

func NewPlaneFromPointNormal(point *Vector3, normal *Vector3) Plane {
	pNormal := normal.Normalized()
	temp := pNormal
	temp.X += 4129
	temp.Y += 4133
	temp.Z += 4139
	temp.ThisNormalize()

	coplanarA := pNormal.Cross(&temp).Normalized()
	coplanarB := pNormal.Cross(&coplanarA).Normalized()
	return Plane{
		point:     *point,
		normal:    pNormal,
		coplanarA: coplanarA,
		coplanarB: coplanarB,
	}
}

// TestPoint tests a point against the plane. Returns:
// 0 if the point is behind the plane,
// 1 if it's in front,
// 2 if it's on the plane
func (p *Plane) TestPoint(point *Vector3) uint8 {
	dist := point.Sub(&p.point)
	dist.ThisNormalize()
	if dist.Dot(&p.normal).IsZero() {
		return 2
	}
	if dist.Dot(&p.normal) > 0 {
		return 1
	}
	return 0
}

func (p *Plane) CoplanarVectors() (Vector3, Vector3) {
	return p.coplanarA, p.coplanarB
}

func (p *Plane) Point() Vector3 {
	return p.point
}

func (p *Plane) Normal() Vector3 {
	return p.normal
}
