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

func NewRasterRenderer(camera *entities.SceneGraphNode, planeZ basics.Scalar, winWidth int, winHeight int) *RasterRender {
	inverseCameraT := camera.GetWorldTransform()
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
		imageBuffer: graphics.NewImageBuffer(int(winWidth), int(winHeight), color.Black),
	}
}

func (r *RasterRender) RenderSceneGraph(sceneGraph *entities.SceneGraph) *graphics.ImageBuffer {
	inverseCameraT := sceneGraph.GetNode("camera").GetWorldTransform()
	inverseCameraT.ThisInvert()
	itemsToRender, lightsToRender := getAllItemsToRender(sceneGraph, &inverseCameraT)
	_ = lightsToRender
	for _, item := range itemsToRender {
		mesh := item.modelObject.GetMesh()
		triangles := mesh.GetTriangles(item.modelObject.GetIgnoreMeshNormals())
		ignoreSpecular := item.modelObject.GetIgnoreSpecular()
		for _, t := range triangles {
			// Translate triangle in view space
			t.ThisApplyTransformation(&item.completeTransform)

			// Tiangles too close to the camera are discarded
			if getClosestZ(&t) < 0.3 {
				continue
			}

			ambientLightColor := basics.Vector3FromColor(color.RGBA{30, 30, 30, 255})
			forward := basics.Forward()
			TriangleNormalsPhong(&t, &forward, &ambientLightColor, item.modelObject.GetSpecularExponent(), lightsToRender, color.RGBA64{1, 1, 1, 255}, ignoreSpecular)

			// Translate triangle in clip space:
			// top left: (-1, +1) | bottom right: (+1, -1) | center: (0, 0)
			//oneIn := false
			for i := 0; i < 3; i++ {
				t[i].Position = projectPointOnViewPlane(&t[i].Position, &r.parameters.planeZero, &r.parameters.planeNormal)
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

			// Back face culling
			triangleNormal := t.GetSurfaceNormal()
			if r.parameters.planeNormal.Dot(&triangleNormal) <= 0 {
				continue
			}

			// Correct scaling for the aspect ratio
			scaleTriangleOnScreen(&t, r.parameters.hw, r.parameters.hh, r.parameters.aspectRatio)

			// Bounding box
			max, min := getMaxMin(t[0].Position, t[1].Position, t[2].Position)
			min.X = basics.Clamp(0, basics.Scalar(r.parameters.winWidth), basics.Floor(min.X))
			min.Y = basics.Clamp(0, basics.Scalar(r.parameters.winHeight), basics.Floor(min.Y))

			max.X = basics.Clamp(0, basics.Scalar(r.parameters.winWidth), basics.Ceil(max.X))
			max.Y = basics.Clamp(0, basics.Scalar(r.parameters.winHeight), basics.Ceil(max.Y))

			// Test for each pixel in the bounding box from top left to bottom right
			for y := int(min.Y); y < int(max.Y); y++ {
				for x := int(min.X); x < int(max.X); x++ {
					target2D := basics.NewVector3(basics.Scalar(x), basics.Scalar(y), 0)
					// find weights for interpolation
					w0, w1, w2 := findWeights(&t[0].Position, &t[1].Position, &t[2].Position, &target2D)
					if w0 < 0 || w1 < 0 || w2 < 0 {
						continue // point lands outside of the triangle
					}
					point := interpolate3Vertices(&t[0].Position, &t[1].Position, &t[2].Position, w0, w1, w2)
					// depth test
					if point.Z < 0 || r.zBuffer.Get(x, y) < point.Z { // if the point is behind the camera or the depth buffer has already something closer
						continue
					}
					r.zBuffer.Set(x, y, point.Z)

					// set color
					color0 := &t[0].Color
					color1 := &t[1].Color
					color2 := &t[2].Color
					colorVector := interpolate3Vertices(color0, color1, color2, w0, w1, w2)
					// Scaling to uint8 range
					colorVector.ThisMul(1 / 65535.0)
					colorVector.ThisMul(255.0)
					r.imageBuffer.Set(x, y, colorVector.ToColor())
				}
			}
		}
	}
	r.zBuffer.Clear()
	return &r.imageBuffer
}

func (r *RasterRender) RenderLine(p0 basics.Vector3, p1 basics.Vector3, color color.Color) {

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

// Retuns the top-left and the bottom right vertex
func getMaxMin(vectors ...basics.Vector3) (basics.Vector3, basics.Vector3) {
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
