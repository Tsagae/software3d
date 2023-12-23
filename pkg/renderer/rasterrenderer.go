package renderer

import (
	"github.com/tsagae/software3d/pkg/basics"
	"github.com/tsagae/software3d/pkg/entities"
	"github.com/tsagae/software3d/pkg/graphics"
	"image/color"
)

type RasterRender struct {
	parameters  RendererParameters
	zBuffer     graphics.ZBuffer
	imageBuffer graphics.ImageBuffer
}

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

func NewRasterRenderer(camera *entities.SceneGraphNode, planeZ basics.Scalar, winWidth int, winHeight int) *RasterRender {
	inverseCameraT := camera.WorldTransform()
	inverseCameraT.ThisInvert()
	return &RasterRender{
		parameters: RendererParameters{
			camera:                 camera,
			planeZ:                 planeZ,
			winWidth:               winWidth,
			winHeight:              winHeight,
			planeNormal:            basics.NewVector3(0, 0, -1),
			planeZero:              basics.NewVector3(0, 0, planeZ),
			aspectRatio:            basics.Scalar(winWidth) / basics.Scalar(winHeight),
			hw:                     basics.Scalar(winWidth) / 2,
			hh:                     basics.Scalar(winHeight) / 2,
			inverseCameraTransform: inverseCameraT,
		},
		zBuffer:     graphics.NewZBuffer(int(winWidth), int(winHeight)),
		imageBuffer: graphics.NewImageBuffer(int(winWidth), int(winHeight)),
	}
}

func (r *RasterRender) RenderSceneGraph(sceneGraph *entities.SceneGraph) *graphics.ImageBuffer {
	inverseCameraT := sceneGraph.GetNode("camera").WorldTransform()
	inverseCameraT.ThisInvert()
	itemsToRender, lightsToRender := getAllItemsToRender(sceneGraph, &inverseCameraT)

	for _, item := range itemsToRender {
		r.renderSingleItem(item, lightsToRender)
	}
	r.zBuffer.Clear()
	return &r.imageBuffer
}

func (r *RasterRender) renderSingleItem(item renderItem, lights []renderLight) {
	mesh := item.modelObject.Mesh()
	ignoreMeshNormals := item.modelObject.IgnoreMeshNormals()
	iterator := mesh.Iterator()
	for iterator.HasNext() {
		// Translate triangle in view space
		var t graphics.Triangle
		if ignoreMeshNormals {
			t = iterator.NextWithFaceNormals()
		} else {
			t = iterator.Next()
		}
		t.ThisApplyTransformation(&item.completeTransform)
		// Triangles too close to the camera are discarded
		if getClosestZ(&t) < 0.3 {
			continue
		}

		lightTriangle(&t, &item, lights)

		projectTriangle(&t, &r.parameters.planeZero, &r.parameters.planeNormal)

		// Back face culling
		triangleNormal := t.GetSurfaceNormal()
		if r.parameters.planeNormal.Dot(&triangleNormal) <= 0 {
			continue
		}

		// Correct scaling for the aspect ratio
		scaleTriangleOnScreen(&t, r.parameters.hw, r.parameters.hh, r.parameters.aspectRatio)

		rasterTriangle(t, r.parameters.winWidth, r.parameters.winHeight, &r.imageBuffer, &r.zBuffer)
	}
}

func rasterTriangle(t graphics.Triangle, winWidth int, winHeight int, imageBuffer *graphics.ImageBuffer, zBuffer *graphics.ZBuffer) {
	// Bounding box
	maxX, minX, maxY, minY := getMaxMin(t[0].Position, t[1].Position, t[2].Position)

	minX = basics.Clamp(0, basics.Scalar(winWidth), basics.Floor(minX))
	minY = basics.Clamp(0, basics.Scalar(winHeight), basics.Floor(minY))

	maxX = basics.Clamp(0, basics.Scalar(winWidth), basics.Ceil(maxX))
	maxY = basics.Clamp(0, basics.Scalar(winHeight), basics.Ceil(maxY))

	// Test for each pixel in the bounding box from top left to bottom right
	for y := int(minY); y < int(maxY); y++ {
		for x := int(minX); x < int(maxX); x++ {
			target2D := basics.NewVector3(basics.Scalar(x), basics.Scalar(y), 0)
			// find weights for interpolation
			w0, w1, w2 := findWeights(&t[0].Position, &t[1].Position, &t[2].Position, &target2D)
			if w0 < 0 || w1 < 0 || w2 < 0 {
				continue // point lands outside the triangle
			}
			point := interpolate3Vertices(&t[0].Position, &t[1].Position, &t[2].Position, w0, w1, w2)
			// depth test
			if point.Z < 0 || zBuffer.Get(x, y) < point.Z { // if the point is behind the camera or the depth buffer has already something closer
				continue
			}
			zBuffer.Set(x, y, point.Z)

			// set color
			color0 := &t[0].Color
			color1 := &t[1].Color
			color2 := &t[2].Color
			colorVector := interpolate3Vertices(color0, color1, color2, w0, w1, w2)
			// Scaling to uint8 range
			colorVector.ThisMul(255.0 / 65535.0) // was: colorVector.ThisMul(1 / 65535.0); colorVector.ThisMul(255.0)
			imageBuffer.Set(x, y, colorVector.ToColor())
		}
	}
}

func projectTriangle(t *graphics.Triangle, planeZero *basics.Vector3, planeNormal *basics.Vector3) {
	// Translate triangle in clip space:
	// top left: (-1, +1) | bottom right: (+1, -1) | center: (0, 0)
	//oneIn := false
	for i := 0; i < 3; i++ {
		t[i].Position = projectPointOnViewPlane(&t[i].Position, planeZero, planeNormal)
		/*
			if (t[i].Position.X > -r.parameters.aspectRatio*2 && t[i].Position.X < r.parameters.aspectRatio*2) && (t[i].Position.Y > -2 && t[i].Position.Y < 2) { // if a vertex is inside the view frustum | edge case with big triangles close to the screen
				oneIn = true
			}
		*/
	}

	/*
		if !oneIn {
			continue
		}
	*/
}

func lightTriangle(t *graphics.Triangle, item *renderItem, lights []renderLight) {
	ambientLightColor := basics.Vector3FromColor(color.RGBA{30, 30, 30, 255})
	forward := basics.Forward()
	TriangleNormalsPhong(t, &forward, &ambientLightColor, item.modelObject.SpecularExponent(), lights, color.RGBA64{1, 1, 1, 255}, item.modelObject.IgnoreSpecular())

}

func (r *RasterRender) RenderLine(p0 basics.Vector3, p1 basics.Vector3, color color.Color) {

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

func findWeights(v1 *basics.Vector3, v2 *basics.Vector3, v3 *basics.Vector3, target *basics.Vector3) (basics.Scalar, basics.Scalar, basics.Scalar) {
	// most of this can be cached when finding weights inside the same triangle TODO
	den := (v2.Y-v3.Y)*(v1.X-v3.X) + (v3.X-v2.X)*(v1.Y-v3.Y)
	t1 := (target.X - v3.X)
	t2 := (target.Y - v3.Y)

	w1 := ((v2.Y-v3.Y)*t1 + (v3.X-v2.X)*t2) / den
	w2 := ((v3.Y-v1.Y)*t1 + (v1.X-v3.X)*t2) / den
	w3 := 1 - w1 - w2
	return w1, w2, w3
}

func interpolate3Vertices(v1 *basics.Vector3, v2 *basics.Vector3, v3 *basics.Vector3, w1 basics.Scalar, w2 basics.Scalar, w3 basics.Scalar) basics.Vector3 {
	point := v1.Mul(w1)
	temp := v2.Mul(w2)
	point.ThisAdd(temp)
	temp = v3.Mul(w3)
	point.ThisAdd(temp)
	return point
}

// Returns maxX, minX, maxY, minY
func getMaxMin(p0, p1, p2 basics.Vector3) (basics.Scalar, basics.Scalar, basics.Scalar, basics.Scalar) {
	maxX := max(p0.X, p1.X, p2.X)
	minX := min(p0.X, p1.X, p2.X)
	maxY := max(p0.Y, p1.Y, p2.Y)
	minY := min(p0.Y, p1.Y, p2.Y)
	return maxX, minX, maxY, minY
}

// Retuns the top-left and the bottom right vertex
func oldMaxMin(vectors ...basics.Vector3) (basics.Vector3, basics.Vector3) {
	maxX := vectors[0].X
	maxY := vectors[0].Y
	minX := maxX
	minY := maxY
	for i := 1; i < 3; i++ {
		if maxX < vectors[i].X {
			maxX = vectors[i].X
		} else if minX > vectors[i].X {
			minX = vectors[i].X
		}
		if maxY < vectors[i].Y {
			maxY = vectors[i].Y
		} else if minY > vectors[i].Y {
			minY = vectors[i].Y
		}
	}
	return basics.NewVector3(maxX, maxY, 0), basics.NewVector3(minX, minY, 0)
}

func getClosestZ(triangle *graphics.Triangle) basics.Scalar {
	minZ := triangle[0].Position.Z
	if triangle[1].Position.Z < minZ {
		minZ = triangle[1].Position.Z
	}
	if triangle[2].Position.Z < minZ {
		minZ = triangle[2].Position.Z
	}
	return minZ
}
