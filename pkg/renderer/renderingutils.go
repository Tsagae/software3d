package renderer

import (
	"github.com/tsagae/software3d/pkg/basics"
	"github.com/tsagae/software3d/pkg/entities"
	"github.com/tsagae/software3d/pkg/graphics"
	"image/color"
)

// RendererParameters Near clip plane is always assumed to be at (0,0,1) looking at (0,0,1)
type RendererParameters struct {
	camera                 *entities.SceneGraphNode
	planeZ                 basics.Scalar
	winWidth               int
	winHeight              int
	aspectRatio            basics.Scalar
	hw                     basics.Scalar
	hh                     basics.Scalar
	inverseCameraTransform basics.Transform
	viewFrustumSides       [5]basics.Plane
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

func scaleTriangleOnScreen(triangle *graphics.Triangle, hw, hh, aspectRatio basics.Scalar) {
	for i := 0; i < 3; i++ {
		scalePointOnScreen(&triangle[i].Position.X, &triangle[i].Position.Y, hw, hh, aspectRatio)
	}
}

// modifies x and y
func scalePointOnScreen(x *basics.Scalar, y *basics.Scalar, hw basics.Scalar, hh basics.Scalar, aspectRatio basics.Scalar) {
	*x += 1 * aspectRatio
	*y += 1
	*x *= hw / aspectRatio
	*y *= hh
}

func projectPointOnViewPlane(p *basics.Vector3) basics.Vector3 {
	//camera is assumed to be in (0,0,0) in its local space
	//center of the view plane is assumed to be at (x: 0, y: 0, z: 1)
	d := *p
	d.X /= d.Z
	d.Y /= d.Z
	return d
}

func findWeights(v1, v2, v3, target *basics.Vector3) (basics.Scalar, basics.Scalar, basics.Scalar) {
	// most of this can be cached when finding weights inside the same triangle TODO
	den := (v2.Y-v3.Y)*(v1.X-v3.X) + (v3.X-v2.X)*(v1.Y-v3.Y)
	t1 := target.X - v3.X
	t2 := target.Y - v3.Y

	w1 := ((v2.Y-v3.Y)*t1 + (v3.X-v2.X)*t2) / den
	w2 := ((v3.Y-v1.Y)*t1 + (v1.X-v3.X)*t2) / den
	w3 := 1 - w1 - w2
	return w1, w2, w3
}

func interpolate3Vertices(v1, v2, v3 *basics.Vector3, w1, w2, w3 basics.Scalar) basics.Vector3 {
	return v1.Mul(w1).Add(v2.Mul(w2)).Add(v3.Mul(w3))
}

// Returns maxX, minX, maxY, minY
func getMaxMin(p0, p1, p2 basics.Vector3) (basics.Scalar, basics.Scalar, basics.Scalar, basics.Scalar) {
	maxX := max(p0.X, p1.X, p2.X)
	minX := min(p0.X, p1.X, p2.X)
	maxY := max(p0.Y, p1.Y, p2.Y)
	minY := min(p0.Y, p1.Y, p2.Y)
	return maxX, minX, maxY, minY
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

func lightTriangle(t *graphics.Triangle, item *renderItem, lights []renderLight) {
	ambientLightColor := basics.Vector3FromColor(color.RGBA{R: 30, G: 30, B: 30, A: 255})
	forward := basics.Forward()
	TriangleNormalsPhong(t, &forward, &ambientLightColor, item.modelObject.SpecularExponent(), lights, color.RGBA64{R: 1, G: 1, B: 1, A: 255}, item.modelObject.IgnoreSpecular())
}

func projectTriangle(t *graphics.Triangle) {
	// Translate triangle in clip space:
	// top left: (-1, +1) | bottom right: (+1, -1) | center: (0, 0)
	for i := 0; i < 3; i++ {
		t[i].Position = projectPointOnViewPlane(&t[i].Position)
	}
}

// Renders a line in clip space
func drawLine(v0, v1 *basics.Vector3, iBuf *graphics.ImageBuffer, zBuf *graphics.ZBuffer) {
	y0 := v0.Y
	y1 := v1.Y
	x0 := v0.X
	x1 := v1.X

	dx := x1 - x0
	dy := y1 - y0

	if basics.Abs(dx) > basics.Abs(dy) {
		if v0.X > v1.X {
			v0, v1 = v1, v0
		}
		if dx == 0 {
			return
		}
		genericDrawLine(v0.X, v1.X, v0.Y, dy/dx, basics.Scalar(iBuf.Width()), basics.Scalar(iBuf.Height()), func(a int, b int, c color.RGBA) {
			iBuf.Set(a, b, c)
		})
	} else {
		if v0.Y > v1.Y {
			v0, v1 = v1, v0
		}
		if dy == 0 {
			return
		}
		genericDrawLine(v0.Y, v1.Y, v0.X, dx/dy, basics.Scalar(iBuf.Height()), basics.Scalar(iBuf.Width()), func(a int, b int, c color.RGBA) {
			iBuf.Set(b, a, c)
		})
	}
}

// a0, a1: start and end points on the same axis
// b0: start on the other axis
// m: slope
// aMaxCanvas, bMaxCanvas first value outside of canvas for the two axis
// requires a0 <= a1
func genericDrawLine(a0, a1, b0, m, aMaxCanvas, bMaxCanvas basics.Scalar, setImage func(int, int, color.RGBA)) {
	b := b0
	for a := a0; a <= a1; a++ {
		if a < 0 || b < 0 || a >= aMaxCanvas || b >= bMaxCanvas {
			continue
		}
		setImage(int(a), int(b), color.RGBA{R: 255, G: 255, B: 255})
		b += m
	}
}

func getViewFrustumSides(aspectRatio basics.Scalar) [5]basics.Plane {
	{
		bottomLeft := basics.Vector3{-aspectRatio, -1, 1}
		bottomRight := basics.Vector3{+aspectRatio, -1, 1}
		topLeft := basics.Vector3{-aspectRatio, +1, 1}
		topRight := basics.Vector3{+aspectRatio, +1, 1}
		forward := basics.Forward()
		return [5]basics.Plane{
			basics.NewPlaneFromPoints(&bottomLeft, &basics.Vector3{}, &bottomRight), //bottom
			basics.NewPlaneFromPoints(&topRight, &basics.Vector3{}, &topLeft),       //top
			basics.NewPlaneFromPoints(&bottomRight, &basics.Vector3{}, &topRight),   //right
			basics.NewPlaneFromPoints(&topLeft, &basics.Vector3{}, &bottomLeft),     //left
			basics.NewPlaneFromPointNormal(&basics.Vector3{}, &forward),             //near
		}
	}
}
