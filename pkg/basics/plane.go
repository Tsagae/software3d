package basics

type Plane struct {
	Point  Vector3
	Normal Vector3
}

func NewPlaneFromPointNormal(point *Vector3, normal *Vector3) Plane {
	pNormal := normal.Normalized()
	return Plane{
		Point:  *point,
		Normal: pNormal,
	}
}

// NewPlaneFromPoints return a plane with normal (a-b)x(c-b) normalized
func NewPlaneFromPoints(a *Vector3, b *Vector3, c *Vector3) Plane {
	v := a.Sub(b)
	w := c.Sub(b)
	return Plane{
		Point:  *a,
		Normal: v.Cross(&w).Normalized(),
	}
}

// TestPoint tests a Point against the plane. Returns:
// 0 if the Point is behind the plane,
// 1 if it's in front,
// 2 if it's on the plane
func (p *Plane) TestPoint(point *Vector3) uint8 {
	dist := point.Sub(&p.Point)
	dist.ThisNormalize()
	if dist.Dot(&p.Normal).IsZero() {
		return 2
	}
	if dist.Dot(&p.Normal) > 0 {
		return 1
	}
	return 0
}

func (p *Plane) CoplanarVectors() (Vector3, Vector3) {
	pNormal := p.Normal
	temp := pNormal
	temp.X += 4129
	temp.Y += 4133
	temp.Z += 4139
	temp.ThisNormalize()

	coplanarA := pNormal.Cross(&temp).Normalized()
	coplanarB := pNormal.Cross(&coplanarA).Normalized()
	return coplanarA, coplanarB
}
