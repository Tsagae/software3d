package renderer

import (
	"github.com/tsagae/software3d/pkg/basics"
	"github.com/tsagae/software3d/pkg/entities"
	"github.com/tsagae/software3d/pkg/graphics"
)

type RendererParameters struct {
	camera                 *entities.SceneGraphNode
	planeZ                 basics.Scalar
	winWidth               int
	winHeight              int
	planeNormal            basics.Vector3
	planeZero              basics.Vector3
	aspectRatio            basics.Scalar
	hw                     basics.Scalar
	hh                     basics.Scalar
	inverseCameraTransform basics.Transform
}

type renderItem struct {
	modelObject       *entities.ModelObject
	completeTransform basics.Transform
	//distanceFromCamera basics.Scalar //probably unnecessary, could use the z of cameraViewTransform
}

type renderLight struct {
	light    *entities.LightObject
	position basics.Vector3 //position in camera space
}

func scaleTriangleOnScreen(triangle *graphics.Triangle, hw basics.Scalar, hh basics.Scalar, aspectRatio basics.Scalar) {
	for i := 0; i < 3; i++ {
		scalePointOnScreen(&triangle[i].Position.X, &triangle[i].Position.Y, hw, hh, aspectRatio)
	}
}

func getAllItemsToRender(sceneGraph *entities.SceneGraph, inverseCameraTransform *basics.Transform) ([]renderItem, []renderLight) {
	node := sceneGraph.GetRoot()
	queue := node.Children()
	nodesToRender := make([]renderItem, 0, len(queue))
	lightsToRender := make([]renderLight, 0)

	for len(queue) != 0 {
		node := queue[0]
		queue = queue[1:]
		objectWorldT := node.WorldTransform()
		objectCameraT := objectWorldT.Cumulate(inverseCameraTransform)
		//objectCameraT.ThisCumulate(&objRotT)
		// TODO optimize repeated transforms, non renderable entities could be removed here

		switch v := node.GameObject.(type) {
		case *entities.ModelObject:
			nodesToRender = append(nodesToRender, renderItem{
				modelObject:       v,
				completeTransform: objectCameraT,
			})
		case *entities.LightObject:
			lightsToRender = append(lightsToRender, renderLight{
				v,
				objectCameraT.Translation,
			})
		}

		queue = append(queue, node.Children()...)
	}
	return nodesToRender, lightsToRender
}

// modifies x and y
func scalePointOnScreen(x *basics.Scalar, y *basics.Scalar, hw basics.Scalar, hh basics.Scalar, aspectRatio basics.Scalar) {
	*x += 1 * aspectRatio
	*y += 1
	*x *= hw / aspectRatio
	*y *= hh
}

func projectPointOnViewPlane(p *basics.Vector3, planeZero *basics.Vector3, planeNormal *basics.Vector3) basics.Vector3 {
	//camera is assumed to be in (0,0,0) in its local space
	depth := p.Z
	d := p.Normalized()
	dDotN := d.Dot(planeNormal)
	k := planeZero.Dot(planeNormal) / dDotN // o.planeZero.Dot(&o.planeNormal) can be pre-computed
	d.ThisMul(k)                            //d is now p projected on the plane
	d.Z = depth
	//d.Y *= -1
	//d.X *= -1
	return d
}
