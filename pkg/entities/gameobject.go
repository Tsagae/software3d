package entities

import (
	"github.com/tsagae/software3d/pkg/basics"
	"github.com/tsagae/software3d/pkg/graphics"
	"image/color"
)

type GameObject interface {
	GetName() string
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

func (e *EmptyObject) GetName() string {
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

func (m *ModelObject) GetName() string {
	return m.name
}

func (m *ModelObject) GetMesh() graphics.Mesh {
	return m.mesh
}

func (m *ModelObject) GetIgnoreMeshNormals() bool {
	return m.ignoreMeshNormals
}

func (m *ModelObject) GetSpecularExponent() basics.Scalar {
	return m.specularExponent
}

func (m *ModelObject) GetIgnoreSpecular() bool {
	return m.ignoreSpecular
}

func NewCameraObject(name string) *CameraObject {
	return &CameraObject{
		name: name,
	}
}

func (c *CameraObject) GetName() string {
	return c.name
}

func NewLightObject(name string, lightColor color.Color, lightFallOff FalloffFunction) *LightObject {
	return &LightObject{
		name:    name,
		color:   lightColor,
		falloff: lightFallOff,
	}
}

func (l *LightObject) GetName() string {
	return l.name
}

func (l *LightObject) GetColor() color.Color {
	return l.color
}

func (l *LightObject) GetFallOff() FalloffFunction {
	return l.falloff
}
