package entities

import (
	"github.com/tsagae/software3d/pkg/basics"
	"github.com/tsagae/software3d/pkg/graphics"
	"image/color"
)

type GameObject interface {
	Name() string
}

type EmptyObject struct {
	name string
}

type ModelObject struct {
	name              string
	mesh              graphics.Mesh
	ignoreMeshNormals bool
	specularExponent  basics.Scalar
	ignoreSpecular    bool
}

type CameraObject struct {
	name string
}

type FalloffFunction func(lightDistance basics.Scalar) basics.Scalar

type LightObject struct {
	name    string
	color   color.Color
	falloff FalloffFunction
}

func NewEmptyObject(name string) *EmptyObject {
	return &EmptyObject{
		name: name,
	}
}

func (e *EmptyObject) Name() string {
	return e.name
}

func NewModelObject(name string, mesh graphics.Mesh, ignoreMeshNormals bool, specularExponent basics.Scalar, ignoreSpecular bool) *ModelObject {
	return &ModelObject{
		mesh:              mesh,
		name:              name,
		ignoreMeshNormals: ignoreMeshNormals,
		specularExponent:  specularExponent,
		ignoreSpecular:    ignoreSpecular,
	}
}

func (m *ModelObject) Name() string {
	return m.name
}

func (m *ModelObject) Mesh() graphics.Mesh {
	return m.mesh
}

func (m *ModelObject) IgnoreMeshNormals() bool {
	return m.ignoreMeshNormals
}

func (m *ModelObject) SpecularExponent() basics.Scalar {
	return m.specularExponent
}

func (m *ModelObject) IgnoreSpecular() bool {
	return m.ignoreSpecular
}

func NewCameraObject(name string) *CameraObject {
	return &CameraObject{
		name: name,
	}
}

func (c *CameraObject) Name() string {
	return c.name
}

func NewLightObject(name string, lightColor color.Color, lightFallOff FalloffFunction) *LightObject {
	return &LightObject{
		name:    name,
		color:   lightColor,
		falloff: lightFallOff,
	}
}

func (l *LightObject) Name() string {
	return l.name
}

func (l *LightObject) Color() color.Color {
	return l.color
}

func (l *LightObject) FallOff() FalloffFunction {
	return l.falloff
}
