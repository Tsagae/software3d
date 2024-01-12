package basics

type Plane struct {
	Point  Vector3
	Normal Vector3
}

func NewPlaneFromPointNormal(point *Vector3, normal *Vector3) Plane {
	return Plane{
		Point:  *point,
		Normal: normal.Normalized(),
	}
}

// TestPoint tests a point against the plane. Returns:
// 0 if the point is behind the plane,
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
