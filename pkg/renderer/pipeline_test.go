package renderer

import (
	"fmt"
	"github.com/tsagae/software3d/pkg/basics"
	"github.com/tsagae/software3d/pkg/entities"
	"github.com/tsagae/software3d/pkg/graphics"
	"image/color"
	"os"
	"testing"
	time2 "time"
)

func BenchmarkSampleScene(b *testing.B) {
	sceneGraph := SampleScene()
	fmt.Println("---------------Benchmark start---------------")
	fmt.Println("SampleScene: ", sceneGraph.String())
	var objRenderer = NewRasterRenderer(sceneGraph.GetNode("camera"), 1, 800, 600)
	b.ResetTimer()
	time := time2.Now()
	objRenderer.RenderSceneGraph(sceneGraph)
	fmt.Printf("Scene graph: %v\n", time2.Now().Sub(time))

	time = time2.Now()
	objRenderer.imageBuffer.Clear()
	fmt.Printf("clearing image buffer: %v\n", time2.Now().Sub(time))
	fmt.Println("----------------Benchmark end----------------")
}

func SampleScene() *entities.SceneGraph {
	var specularExp basics.Scalar = 600

	sceneGraph := entities.NewSceneGraph()

	meshes := loadMeshes()
	cameraObj := entities.NewCameraObject(
		"mainCamera",
	)
	rotateCameraT := basics.NewTransform(1, basics.NewQuaternionFromAngleAndAxis(-20, basics.Up()), basics.Vector3{})
	sceneGraph.AddChild("world", entities.NewSceneGraphNode(cameraObj, "camera"), basics.NewTransform(1, basics.NewIdentityQuaternion(), basics.NewVector3(1.5, 1, -3)))
	cameraNode := sceneGraph.GetNode("camera")
	cameraNode.CumulateBeforeLocalTranform(&rotateCameraT)

	quadObj := entities.NewModelObject("quad", meshes["quad"], true, specularExp, true)
	sceneGraph.AddChild("world", entities.NewSceneGraphNode(quadObj, "quad"), basics.NewTransform(5, basics.NewQuaternionFromEulerAngles(0, 90, 0), basics.NewVector3(0, 0, 0)))

	planeObj := entities.NewModelObject("planeObj", meshes["plane"], true, specularExp, true)

	cubeObj := entities.NewModelObject("cubeObj", meshes["cube"], true, specularExp, false)

	torusObj := entities.NewModelObject("torusObj", meshes["torus"], false, specularExp, false)
	sceneGraph.AddChild("world", entities.NewSceneGraphNode(torusObj, "torus"), basics.NewTransform(1, basics.NewIdentityQuaternion(), basics.NewVector3(3, 1, 3)))

	sphereObj := entities.NewModelObject("sphereObj", meshes["sphere"], false, specularExp, false)
	sceneGraph.AddChild("world", entities.NewSceneGraphNode(cubeObj, "cube"), basics.NewTransform(1, basics.NewIdentityQuaternion(), basics.NewVector3(0, 0, 0)))
	sceneGraph.AddChild("world", entities.NewSceneGraphNode(sphereObj, "sphere"), basics.NewTransform(0.6, basics.NewIdentityQuaternion(), basics.NewVector3(-1, 1, 1)))
	//sceneGraph.AddChild("world", entities.NewSceneGraphNode(planeObj, "plane2"), basics.NewTransform(2, basics.NewQuaternionFromAngleAndAxis(45, basics.Up()), basics.NewVector3(1, -3, 5)))

	//sceneGraph.AddChild("cube", entities.NewSceneGraphNode(cubeObj, "cube2"), basics.NewTransform(0.5, basics.NewIdentityQuaternion(), basics.NewVector3(1, 1, 1)))
	sceneGraph.AddChild("world", entities.NewSceneGraphNode(planeObj, "plane"), basics.NewTransform(5, basics.NewIdentityQuaternion(), basics.NewVector3(0, -2, 0)))
	sceneGraph.AddChild("world", entities.NewSceneGraphNode(planeObj, "plane2"), basics.NewTransform(5, basics.NewQuaternionFromAngleAndAxis(90, basics.Forward()), basics.NewVector3(10, -2, 10)))
	yRotationTransformation := basics.NewTransform(1, basics.NewQuaternionFromAngleAndAxis(20, basics.Right()), basics.NewVector3(0, 0, 0))
	torusNode := sceneGraph.GetNode("torus")
	torusNode.CumulateBeforeLocalTranform(&yRotationTransformation)
	xRot := basics.NewTransform(1, basics.NewQuaternionFromAngleAndAxis(20, basics.Up()), basics.Vector3{})
	torusNode.CumulateBeforeLocalTranform(&xRot)

	// Lighting
	simpleFallOff := func(lightDistance basics.Scalar) basics.Scalar {
		return basics.Clamp(0, 1, 1-(lightDistance/basics.Scalar(50)))
	}
	_ = simpleFallOff

	noFallOff := func(lightDistance basics.Scalar) basics.Scalar {
		return 1
	}
	_ = noFallOff

	easeOutQuint := func(lightDistance basics.Scalar) basics.Scalar {
		lightDistance = basics.ClampMin(1, lightDistance)
		//val := 1 - (lightDistance / 100)
		val := basics.Pow(lightDistance, 100)
		//fmt.Println(val)
		return basics.Clamp(0, 1, val)
	}
	_ = easeOutQuint

	sceneGraph.AddChild("world", entities.NewSceneGraphNode(entities.NewLightObject("light1", color.RGBA{150, 150, 150, 255}, simpleFallOff), "ligh1"), basics.NewTransform(1, basics.NewIdentityQuaternion(), basics.NewVector3(0, 5, 0)))
	sceneGraph.AddChild("world", entities.NewSceneGraphNode(entities.NewLightObject("light2", color.RGBA{80, 150, 20, 255}, simpleFallOff), "ligh2"), basics.NewTransform(1, basics.NewIdentityQuaternion(), basics.NewVector3(2, 2, 2)))

	return sceneGraph
}

func loadMeshes() map[string]graphics.Mesh {

	meshes := make(map[string]graphics.Mesh)

	meshes["cube"] = readMeshFromFile("../../meshes/cube.obj", color.RGBA{R: 180, G: 25, B: 25, A: 255})

	meshes["sphere"] = readMeshFromFile("../../meshes/sphere.obj", color.RGBA{B: 180, A: 255})

	meshes["plane"] = readMeshFromFile("../../meshes/lowpolyplane.obj", color.RGBA{R: 70, G: 50, B: 30, A: 255})

	meshes["torus"] = readMeshFromFile("../../meshes/torus.obj", color.RGBA{G: 180, A: 255})

	meshes["quad"] = readMeshFromFile("../../meshes/quad.obj", color.RGBA{R: 200, G: 200, B: 30, A: 255})

	return meshes
}

func readMeshFromFile(fileName string, meshColor color.RGBA) graphics.Mesh {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	mesh, err := graphics.NewMeshFromReader(f, basics.Vector3FromColor(meshColor))
	if err != nil {
		panic(err)
	}
	err = f.Close()
	if err != nil {
		panic(err)
	}
	return mesh
}
