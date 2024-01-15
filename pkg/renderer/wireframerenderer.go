package renderer

import (
	"github.com/tsagae/software3d/pkg/basics"
	"github.com/tsagae/software3d/pkg/entities"
	"github.com/tsagae/software3d/pkg/graphics"
)

func (r *RasterRender) RenderSceneGraphWireFrame(sceneGraph *entities.SceneGraph) *graphics.ImageBuffer {
	inverseCameraT := sceneGraph.GetNode("camera").WorldTransform()
	inverseCameraT.ThisInvert()
	itemsToRender, _ := getAllItemsToRender(sceneGraph, &inverseCameraT)

	for _, item := range itemsToRender {
		r.renderSingleItemWireFrame(item)
	}
	r.zBuffer.Clear()
	return &r.imageBuffer
}

func (r *RasterRender) renderSingleItemWireFrame(item renderItem) {
	mesh := item.modelObject.Mesh()
	iterator := mesh.Iterator()
	var t graphics.Triangle
	for iterator.HasNext() {
		// Translate triangle in view space
		t = iterator.Next()
		t.ThisApplyTransformation(&item.completeTransform)

		if t[0].Position.Z <= 0 {
			continue
		}
		if t[1].Position.Z <= 0 {
			continue
		}
		if t[2].Position.Z <= 0 {
			continue
		}
		projectTriangle(&t)

		// Back face culling
		forward := basics.Forward()
		triangleNormal := t.GetSurfaceNormal()
		if forward.Dot(&triangleNormal) >= 0 {
			//	continue
		}

		// Correct scaling for the aspect ratio
		scaleTriangleOnScreen(&t, r.parameters.hw, r.parameters.hh, r.parameters.aspectRatio)

		drawLine(&t[0].Position, &t[1].Position, &r.imageBuffer, &r.zBuffer)
		drawLine(&t[1].Position, &t[2].Position, &r.imageBuffer, &r.zBuffer)
		drawLine(&t[2].Position, &t[0].Position, &r.imageBuffer, &r.zBuffer)
	}
}
