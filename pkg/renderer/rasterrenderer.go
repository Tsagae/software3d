package renderer

import (
	"github.com/tsagae/software3d/pkg/basics"
	"github.com/tsagae/software3d/pkg/entities"
	"github.com/tsagae/software3d/pkg/graphics"
)

type RasterRender struct {
	parameters  RendererParameters
	zBuffer     graphics.ZBuffer
	imageBuffer graphics.ImageBuffer
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
			aspectRatio:            basics.Scalar(winWidth) / basics.Scalar(winHeight),
			hw:                     basics.Scalar(winWidth) / 2,
			hh:                     basics.Scalar(winHeight) / 2,
			inverseCameraTransform: inverseCameraT,
			viewFrustumSides:       getViewFrustumSides(basics.Scalar(winWidth) / basics.Scalar(winHeight)),
		},
		zBuffer:     graphics.NewZBuffer(winWidth, winHeight),
		imageBuffer: graphics.NewImageBuffer(winWidth, winHeight),
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

		triangles := ClipTriangleAgainsPlanes(&t, r.parameters.viewFrustumSides)

		for _, t := range triangles {

			lightTriangle(&t, &item, lights)

			projectTriangle(&t)

			// Back face culling
			triangleNormal := t.GetSurfaceNormal()
			forward := basics.Forward()
			if forward.Dot(triangleNormal) >= 0 {
				continue
			}

			// Correct scaling for the aspect ratio
			scaleTriangleOnScreen(&t, r.parameters.hw, r.parameters.hh, r.parameters.aspectRatio)

			rasterTriangle(t, r.parameters.winWidth, r.parameters.winHeight, &r.imageBuffer, &r.zBuffer)

		}
	}
}

func rasterTriangle(t graphics.Triangle, winWidth int, winHeight int, imageBuffer *graphics.ImageBuffer, zBuffer *graphics.ZBuffer) {
	// Bounding box
	maxX, minX, maxY, minY := getMaxMin(t[0].Position, t[1].Position, t[2].Position)
	minX = basics.Clamp(0, basics.Scalar(winWidth), basics.Floor(minX))
	minY = basics.Clamp(0, basics.Scalar(winHeight), basics.Floor(minY))

	maxX = basics.Clamp(0, basics.Scalar(winWidth), basics.Ceil(maxX))
	maxY = basics.Clamp(0, basics.Scalar(winHeight), basics.Ceil(maxY))

	//fmt.Printf("minX: %v minY: %v maxX: %v maxY: %v\n", minX, minY, maxX, maxY)

	// Test for each pixel in the bounding box from top left to bottom right
	for y := int(minY); y < int(maxY); y++ {
		for x := int(minX); x < int(maxX); x++ {
			target2D := basics.NewVector3(basics.Scalar(x), basics.Scalar(y), 0)
			// find weights for interpolation
			w0, w1, w2 := t.FindWeightsPosition(&target2D)
			if w0 < 0 || w1 < 0 || w2 < 0 {
				continue // point lands outside the triangle
			}
			point := t.InterpolatePosition(w0, w1, w2)
			// depth test
			if point.Z < 1 || zBuffer.Get(x, y) < point.Z { // if the point is behind the camera or the depth buffer has already something closer
				continue
			}
			zBuffer.Set(x, y, point.Z)

			// set color
			colorVector := t.InterpolateColor(w0, w1, w2)

			// Scaling to uint8 range
			colorVector = colorVector.Mul(255.0 / 65535.0) // was: colorVector.ThisMul(1 / 65535.0); colorVector.ThisMul(255.0)
			imageBuffer.Set(x, y, colorVector.ToColor())
		}
	}
}
