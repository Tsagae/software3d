package basics

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTestPoint(t *testing.T) {
	planeOrigin := NewVector3(0, 0, 0)
	planeNormal := NewVector3(0, 0, 3)
	plane := NewPlaneFromPointNormal(&planeOrigin, &planeNormal)

	pointInFront := NewVector3(5, 6, 3)
	assert.Equal(t, byte(1), plane.TestPoint(&pointInFront))

	pointOnBack := NewVector3(5, 6, -3)
	assert.Equal(t, byte(0), plane.TestPoint(&pointOnBack))

	pointOnPlane := NewVector3(5, 6, 0)
	assert.Equal(t, byte(2), plane.TestPoint(&pointOnPlane))

	pNormal := plane.Normal
	assert.True(t, pNormal.Length().Equals(1))
}

func TestGetCoplanarVectors(t *testing.T) {
	planeOrigin := NewVector3(0, 0, 0)
	planeNormal := NewVector3(0, 0, 1)
	plane := NewPlaneFromPointNormal(&planeOrigin, &planeNormal)

	v1, v2 := plane.CoplanarVectors()
	pNormal := plane.Normal
	assert.Truef(t, v1.Dot(pNormal).IsZero(), "Dot between %v and %v has to be zero for them to be coplanar", v1, pNormal)
	assert.True(t, v1.Length().Equals(1))

	assert.Truef(t, v2.Dot(pNormal).IsZero(), "Dot between %v and %v has to be zero for them to be coplanar", v2, pNormal)
	assert.True(t, v2.Length().Equals(1))

}

func TestNewPlaneFromPoints(t *testing.T) {
	a := NewVector3(-1, 0, 1)
	b := Vector3{}
	c := NewVector3(1, 0, 1)
	plane := NewPlaneFromPoints(&a, &b, &c)
	assert.True(t, plane.Normal.Equals(Vector3{0, 1, 0}))
}
