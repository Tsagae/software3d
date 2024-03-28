package renderer

import (
	"github.com/tsagae/software3d/pkg/basics"
	"github.com/tsagae/software3d/pkg/entities"
	"github.com/tsagae/software3d/pkg/graphics"
	"image/color"
)

type RasterRenderer struct {
	parameters                  Parameters
	zBuffer                     graphics.ZBuffer
	imageBuffer                 graphics.ImageBuffer
	lastTriCount                uint
	lastTriDiscardedCount       uint
	lastFragmentsInsideBBCount  uint
	lastFragmentsOutsideBBCount uint
}

func NewRasterRenderer(camera *entities.SceneGraphNode, planeZ basics.Scalar, winWidth int, winHeight int) *RasterRenderer {
	inverseCameraT := camera.WorldTransform()
	inverseCameraT.ThisInvert()
	return &RasterRenderer{
		parameters: Parameters{
			camera:                 camera,
			planeZ:                 planeZ,
			winWidth:               winWidth,
			winHeight:              winHeight,
			aspectRatio:            basics.Scalar(winWidth) / basics.Scalar(winHeight),
			hw:                     basics.Scalar(winWidth) / 2,
			hh:                     basics.Scalar(winHeight) / 2,
			inverseCameraTransform: inverseCameraT,
			viewFrustumSides:       getViewFrustumSides(basics.Scalar(winWidth) / basics.Scalar(winHeight)),
			renderMode:             RendermodeNormal,
		},
		zBuffer:     graphics.NewZBuffer(winWidth, winHeight),
		imageBuffer: graphics.NewImageBuffer(winWidth, winHeight),
	}
}

func (r *RasterRenderer) SetRenderMode(renderMode uint8) {
	r.parameters.renderMode = renderMode
}

func (r *RasterRenderer) RenderSceneGraph(sceneGraph *entities.SceneGraph) *graphics.ImageBuffer {
	r.lastTriDiscardedCount = 0
	r.lastFragmentsOutsideBBCount = 0
	r.lastFragmentsInsideBBCount = 0
	r.lastTriCount = 0
	r.zBuffer.Clear()

	inverseCameraT := sceneGraph.GetNode("camera").WorldTransform()
	inverseCameraT.ThisInvert()
	itemsToRender, lightsToRender := getAllItemsToRender(sceneGraph, &inverseCameraT)

	for _, item := range itemsToRender {
		switch r.parameters.renderMode {
		case RendermodeNormal:
			r.renderSingleItem(item, lightsToRender)
		case RendermodeWireframe:
			r.renderSingleItemWireFrame(item)
		default:
			panic("invalid Rendermode")
		}
	}
	return &r.imageBuffer
}

func (r *RasterRenderer) renderSingleItem(item renderItem, lights []renderLight) {
	mesh := item.modelObject.Mesh()
	iterator := mesh.Iterator()

	var nextFunc func() graphics.Triangle
	if item.modelObject.IgnoreMeshNormals() {
		nextFunc = func() graphics.Triangle {
			return iterator.NextWithFaceNormals()
		}
	} else {
		nextFunc = func() graphics.Triangle {
			return iterator.Next()
		}
	}

	for iterator.HasNext() {
		// Translate triangle in view space
		var t graphics.Triangle
		t = nextFunc()
		t.ThisApplyTransformation(&item.completeTransform)

		triangles := ClipTriangleAgainstPlanes(&t, r.parameters.viewFrustumSides)

		for _, t := range triangles {
			for _, vertex := range t {
				if vertex.Position.X.IsNaN() || vertex.Position.Y.IsNaN() || vertex.Position.Z.IsNaN() {
					panic("NaN found in vertex position") //assertion
				}
			}

			lightTriangle(&t, &item, lights)

			projectTriangle(&t)

			// Back face culling
			triangleNormal := t.GetSurfaceNormal()
			forward := basics.Forward()
			if forward.Dot(triangleNormal) > 0 {
				r.lastTriDiscardedCount++
				continue
			}

			// Correct scaling for the aspect ratio
			scaleTriangleOnScreen(&t, r.parameters.hw, r.parameters.hh, r.parameters.aspectRatio)

			r.rasterTriangle(t)
			r.lastTriCount++
		}
	}
}

func (r *RasterRenderer) renderSingleItemWireFrame(item renderItem) {
	mesh := item.modelObject.Mesh()
	iterator := mesh.Iterator()
	var t graphics.Triangle
	for iterator.HasNext() {

		// Translate triangle in view space
		t = iterator.Next()
		t.ThisApplyTransformation(&item.completeTransform)

		triangles := ClipTriangleAgainstPlanes(&t, r.parameters.viewFrustumSides)

		for _, t := range triangles {

			projectTriangle(&t)

			// Back face culling
			triangleNormal := t.GetSurfaceNormal()
			forward := basics.Forward()
			if forward.Dot(triangleNormal) > 0 {
				r.lastTriDiscardedCount++
				continue
			}

			// Correct scaling for the aspect ratio
			scaleTriangleOnScreen(&t, r.parameters.hw, r.parameters.hh, r.parameters.aspectRatio)

			for i := 0; i < 3; i++ {
				drawLine(&t[i].Position, &t[(i+1)%3].Position, color.RGBA{R: 255, G: 255, B: 255}, &r.imageBuffer)
			}

			r.lastTriCount++
		}
	}
}

func (r *RasterRenderer) rasterTriangle(t graphics.Triangle) {
	// Bounding box
	maxX, minX, maxY, minY := getMaxMin(t[0].Position, t[1].Position, t[2].Position)
	minX = basics.Clamp(0, basics.Scalar(r.parameters.winWidth), basics.Floor(minX))
	minY = basics.Clamp(0, basics.Scalar(r.parameters.winHeight), basics.Floor(minY))

	maxX = basics.Clamp(0, basics.Scalar(r.parameters.winWidth), basics.Ceil(maxX))
	maxY = basics.Clamp(0, basics.Scalar(r.parameters.winHeight), basics.Ceil(maxY))

	// Test for each pixel in the bounding box from top left to bottom right
	for y := int(minY); y < int(maxY); y++ {
		for x := int(minX); x < int(maxX); x++ {
			target2D := basics.NewVector3(basics.Scalar(x), basics.Scalar(y), 0)
			// find weights for interpolation
			w0, w1, w2 := basics.FindWeights2D(&t[0].Position, &t[1].Position, &t[2].Position, &target2D)
			if w0 < 0 || w1 < 0 || w2 < 0 {
				r.lastFragmentsOutsideBBCount++
				continue // point lands outside the triangle
			}
			r.lastFragmentsInsideBBCount++
			point := t.InterpolateVertexProps(w0, w1, w2)

			// depth test
			if r.zBuffer.Get(x, y) < point.Position.Z { // if the depth buffer has already something closer
				continue
			}

			r.zBuffer.Set(x, y, point.Position.Z)

			// Scaling to uint8 range
			point.Color = point.Color.Mul(255.0 / 65535.0) // was: colorVector.ThisMul(1 / 65535.0); colorVector.ThisMul(255.0)
			r.imageBuffer.Set(x, y, point.Color.ToColor())
		}
	}
}

func (r *RasterRenderer) LastFragmentsInsideBBCount() uint {
	return r.lastFragmentsInsideBBCount
}

func (r *RasterRenderer) LastFragmentsOutsideBBCount() uint {
	return r.lastFragmentsOutsideBBCount
}

func (r *RasterRenderer) LastTriCount() uint {
	return r.lastTriCount
}

func (r *RasterRenderer) LastTriDiscardedCount() uint {
	return r.lastTriDiscardedCount
}
