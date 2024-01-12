package basics

import "github.com/stretchr/testify/assert"
import "testing"

func TestTestPoint(t *testing.T) {
	planeOrigin := NewVector3(0, 0, 0)
	planeNormal := NewVector3(0, 0, 1)
	plane := NewPlaneFromPointNormal(&planeOrigin, &planeNormal)

	pointInFront := NewVector3(5, 6, 3)
	assert.Equal(t, byte(1), plane.TestPoint(&pointInFront))

	pointOnBack := NewVector3(5, 6, -3)
	assert.Equal(t, byte(0), plane.TestPoint(&pointOnBack))

	pointOnPlane := NewVector3(5, 6, 0)
	assert.Equal(t, byte(2), plane.TestPoint(&pointOnPlane))
}
