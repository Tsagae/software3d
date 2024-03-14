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

		triangles := ClipTriangleAgainsPlanes(&t, r.parameters.viewFrustumSides)

		//fmt.Println("tris to draw: ", triangles)
		for _, triangle := range triangles {

			for i := 0; i < 3; i++ {
				p0 := projectPointOnViewPlane(&triangle[i].Position)
				p1 := projectPointOnViewPlane(&triangle[(i+1)%3].Position)
				scalePointOnScreen(&p0.X, &p0.Y, r.parameters.hw, r.parameters.hh, r.parameters.aspectRatio)
				scalePointOnScreen(&p1.X, &p1.Y, r.parameters.hw, r.parameters.hh, r.parameters.aspectRatio)
				drawLine(&p0, &p1, &r.imageBuffer)
			}
		}
		/*
			for i := 0; i < 3; i++ {
				p0, p1, notCulled := clipLineAgainstFrustum(t[i].Position, t[(i+1)%3].Position, &r.parameters.viewFrustumSides)

				if notCulled {
					p0 = projectPointOnViewPlane(&p0)
					p1 = projectPointOnViewPlane(&p1)
					scalePointOnScreen(&p0.X, &p0.Y, r.parameters.hw, r.parameters.hh, r.parameters.aspectRatio)
					scalePointOnScreen(&p1.X, &p1.Y, r.parameters.hw, r.parameters.hh, r.parameters.aspectRatio)
					drawLine(&p0, &p1, &r.imageBuffer, &r.zBuffer)
				}
			}*/
	}
}

// returns false if the line is completely outside the frustum
func clipLineAgainstFrustum(p0, p1 basics.Vector3, frustumSides *[]basics.Plane) (basics.Vector3, basics.Vector3, bool) {
	isKept := false
	for _, plane := range *frustumSides {
		p0, p1, isKept = ClipSegment(&p0, &p1, &plane)
		if !isKept {
			return p0, p1, false
		}
	}
	return p0, p1, true
}
