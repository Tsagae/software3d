package main

import (
	"fmt"
	"github.com/tsagae/software3d/pkg/basics"
	"github.com/tsagae/software3d/pkg/entities"
	"github.com/tsagae/software3d/pkg/graphics"
	"github.com/tsagae/software3d/pkg/renderer"
	"image"
	"image/color"
	"runtime"
	"time"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var winWidth, winHeight int = 800, 600

// tranformations
var forwardT basics.Transform = basics.NewTransform(1, basics.NewIdentityQuaternion(), basics.Forward().Mul(0.1))
var backwardT basics.Transform = basics.NewTransform(1, basics.NewIdentityQuaternion(), basics.Backward().Mul(0.1))
var rightT basics.Transform = basics.NewTransform(1, basics.NewIdentityQuaternion(), basics.Right().Mul(0.1))
var leftT basics.Transform = basics.NewTransform(1, basics.NewIdentityQuaternion(), basics.Left().Mul(0.1))

var cameraYaw basics.Scalar = 0
var cameraPitch basics.Scalar = 0

func init() {
	// GLFW: This is needed to arrange that main() runs on main thread.
	// See documentation for functions that are only allowed to be called from the main thread.
	runtime.LockOSThread()
}

func main() {
	/*
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil)) //go tool pprof -png http://localhost:6060/debug/pprof/heap > out.png
			http.HandleFunc("/debug/pprof/", pprof.Index)
			http.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
			http.HandleFunc("/debug/pprof/profile", pprof.Profile)
			http.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
			http.HandleFunc("/debug/pprof/trace", pprof.Trace)
		}()*/
	run()
}

func oGLUpdateFrame(window *glfw.Window, texture uint32, w int, h int, img *image.RGBA) {
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA8, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(img.Pix))

	gl.BlitFramebuffer(0, 0, int32(w), int32(h), 0, 0, int32(w), int32(h), gl.COLOR_BUFFER_BIT, gl.LINEAR)

	window.SwapBuffers()
	glfw.PollEvents()
}

func run() int {

	sceneGraph := setup()
	var objRenderer = renderer.NewRasterRenderer(sceneGraph.GetNode("camera"), 1, winWidth, winHeight)

	//renderRGBAxis(1, objRenderer, false)
	//renderer.DrawString(16, 16, fmt.Sprintf("FPS %.0f", displayFps), basics.Color{R: 0, G: 255, B: 0, A: 255})
	//renderer.DrawString(16, 25, fmt.Sprintf("ms %v", displayDelta), basics.Color{R: 0, G: 255, B: 0, A: 255})
	//objRenderer.RenderSceneGraph(&sceneGraph)

	err := glfw.Init()
	if err != nil {
		panic(err)
	}
	defer glfw.Terminate()

	window, err := glfw.CreateWindow(winWidth, winHeight, "My Window", nil, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()

	err = gl.Init()
	if err != nil {
		panic(err)
	}

	var texture uint32
	{
		gl.GenTextures(1, &texture)

		gl.BindTexture(gl.TEXTURE_2D, texture)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)

		gl.BindImageTexture(0, texture, 0, false, 0, gl.WRITE_ONLY, gl.RGBA8)
	}

	var framebuffer uint32
	{
		gl.GenFramebuffers(1, &framebuffer)
		gl.BindFramebuffer(gl.FRAMEBUFFER, framebuffer)
		gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, texture, 0)

		gl.BindFramebuffer(gl.READ_FRAMEBUFFER, framebuffer)
		gl.BindFramebuffer(gl.DRAW_FRAMEBUFFER, 0)
	}

	var imageBuffer *graphics.ImageBuffer
	var frames int = 1
	var startTime time.Time = time.Now()
	var elapsed time.Duration
	var elapsedSum time.Duration
	camera := sceneGraph.GetNode("camera")
	if camera == nil {
		panic("camera not found in scene graph")
	}

	//for !window.ShouldClose() {
	for {
		elapsed = time.Since(startTime)
		elapsedSum += elapsed

		if frames%20 == 0 && elapsed.Milliseconds() != 0 {
			fmt.Printf("avg ms: %v \n", elapsedSum.Milliseconds()/int64(frames))
			elapsedSum = 0
			frames = 0
			//fmt.Println("FPS: ", 1000/elapsed.Milliseconds(), "ms: ", elapsed.Milliseconds())
		}
		startTime = time.Now()
		frames++

		imageBuffer = objRenderer.RenderSceneGraph(&sceneGraph)
		//sg := entities.NewSceneGraph()
		//imageBuffer = objRenderer.RenderSceneGraph(&sg)
		mainLoop(sceneGraph)

		var w, h = window.GetSize()

		var img = image.NewRGBA(image.Rect(0, 0, w, h))

		// -------------------------
		// MODIFY OR LOAD IMAGE HERE
		img = imageBuffer.GetImage()
		/*
			// RESIZING
			// Set the expected size that you want:
			//dst := image.NewRGBA(image.Rect(0, 0, w, h))

			// Resize:
			im := resize.Resize(uint(w), uint(h), image.Image(img), resize.NearestNeighbor)
			if tmp, ok := im.(*image.RGBA); ok {
				img = tmp
			}
			//draw.NearestNeighbor.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)
		*/
		// -------------------------
		oGLUpdateFrame(window, texture, w, h, img)
		inputHandler(window, camera)

		imageBuffer.Clear()
	}
	return 0
}

func inputHandler(window *glfw.Window, camera *entities.SceneGraphNode) {
	cameraDir := camera.GetOrientation()
	cameraDir[2].Y = 0
	cameraDir[2].ThisNormalize()
	movement := basics.NewZeroTransform()
	var tempMov basics.Vector3
	// Movement
	if window.GetKey(glfw.KeyW) == glfw.Press {
		tempMov = cameraDir[2].Mul(0.1)
		movement.Translation.ThisAdd(tempMov)
	}
	if window.GetKey(glfw.KeyS) == glfw.Press {
		tempMov = cameraDir[2].Mul(-0.1)
		movement.Translation.ThisAdd(tempMov)
	}
	if window.GetKey(glfw.KeyD) == glfw.Press {
		tempMov = cameraDir[0].Mul(0.1)
		movement.Translation.ThisAdd(tempMov)
	}
	if window.GetKey(glfw.KeyA) == glfw.Press {
		tempMov = cameraDir[0].Mul(-0.1)
		movement.Translation.ThisAdd(tempMov)
	}
	if window.GetKey(glfw.KeyQ) == glfw.Press {
		movement.Translation.ThisAdd(basics.Up().Mul(0.1))
	}
	if window.GetKey(glfw.KeyE) == glfw.Press {
		movement.Translation.ThisAdd(basics.Down().Mul(0.1))
	}
	// View Rotation
	if window.GetKey(glfw.KeyUp) == glfw.Press {
		cameraPitch -= 1
	}
	if window.GetKey(glfw.KeyDown) == glfw.Press {
		cameraPitch += 1
	}
	if window.GetKey(glfw.KeyRight) == glfw.Press {
		cameraYaw += 1
	}
	if window.GetKey(glfw.KeyLeft) == glfw.Press {
		cameraYaw -= 1
	}
	cameraPitch = basics.Clamp(-89, 89, cameraPitch)
	camera.SetViewRotation(cameraYaw, cameraPitch)
	camera.CumulateWorldTransform(&movement)
}

func mainLoop(sceneGraph entities.SceneGraph) {
	yRotationTransformation := basics.NewTransform(1, basics.NewQuaternionFromAngleAndAxis(0.3, basics.Up()), basics.NewVector3(0, 0, 0))
	xRot := basics.NewTransform(1, basics.NewQuaternionFromAngleAndAxis(1, basics.Right()), basics.ZeroVector())
	torusNode := sceneGraph.GetNode("torus")
	torusNode.CumulateBeforeLocalTranform(&yRotationTransformation)
	torusNode.CumulateBeforeLocalTranform(&xRot)
	/*
		//movement := basics.NewTransform(1, basics.NewIdentityQuaternion(), basics.NewVector3(-0.002, 0, +0.001))

		cube2Node := sceneGraph.GetNode("cube2")


		cube2Node.CumulateLocalTransform(&yRotationTransformation)
	*/
	//camera := sceneGraph.GetNode("camera")
	//inverseRot := basics.NewTransform(1, basics.NewQuaternionFromAngleAndAxis(1, basics.Up()), basics.ZeroVector())
	//camera.CumulateBeforeLocalTranform(&inverseRot)
}

/*
	func renderRGBAxis(scale basics.Scalar, renderer *app.TriRenderer, unitLine bool) {
		renderer.RenderLine(basics.ZeroVector(), basics.Right().Mul(scale), basics.Red(255))
		renderer.RenderLine(basics.ZeroVector(), basics.Up().Mul(scale), basics.Green(255))
		renderer.RenderLine(basics.ZeroVector(), basics.Forward().Mul(scale), basics.Blue(255))
		if unitLine {
			renderer.RenderLine(basics.ZeroVector(), basics.NewVector3(1, 1, 1).Mul(scale), basics.NewColor(255, 255, 0, 255))
		}

}
*/
func loadMeshes() map[string]graphics.Mesh {
	meshes := make(map[string]graphics.Mesh)

	cubeMesh1, err := graphics.ReadMeshFromFile("meshes/cube.obj", basics.Vector3FromColor(color.RGBA{180, 25, 25, 255}))
	if err != nil {
		panic(err)
	}
	meshes["cube"] = cubeMesh1

	sphereMesh, err := graphics.ReadMeshFromFile("meshes/sphere.obj", basics.Vector3FromColor(color.RGBA{0, 0, 180, 255}))
	if err != nil {
		panic(err)
	}
	meshes["sphere"] = sphereMesh

	planeMesh, err := graphics.ReadMeshFromFile("meshes/plane.obj", basics.Vector3FromColor(color.RGBA{70, 50, 30, 255}))
	if err != nil {
		panic(err)
	}
	meshes["plane"] = planeMesh

	torusMesh, err := graphics.ReadMeshFromFile("meshes/torus.obj", basics.Vector3FromColor(color.RGBA{0, 180, 0, 255}))
	if err != nil {
		panic(err)
	}
	meshes["torus"] = torusMesh

	return meshes
}

func setup() entities.SceneGraph {
	var specularExp basics.Scalar = 20

	sceneGraph := entities.NewSceneGraph()

	meshes := loadMeshes()
	cameraObj := entities.NewCameraObject(
		"mainCamera",
	)
	rotateCameraT := basics.NewTransform(1, basics.NewQuaternionFromAngleAndAxis(-20, basics.Up()), basics.ZeroVector())
	sceneGraph.AddChild("world", entities.NewSceneGraphNode(cameraObj, "camera"), basics.NewTransform(1, basics.NewIdentityQuaternion(), basics.NewVector3(1.5, 1, -3)))
	cameraNode := sceneGraph.GetNode("camera")
	cameraNode.CumulateBeforeLocalTranform(&rotateCameraT)

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
	xRot := basics.NewTransform(1, basics.NewQuaternionFromAngleAndAxis(20, basics.Up()), basics.ZeroVector())
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

/*
func testSdlGeo(renderer *sdl.Renderer) {
	vertices := make([]sdl.Vertex, 0)
	color1 := sdl.Color{R: 255, G: 0, B: 0, A: 255}
	color2 := sdl.Color{R: 0, G: 255, B: 0, A: 255}
	color3 := sdl.Color{R: 0, G: 0, B: 255, A: 255}

	vertices = append(vertices, sdl.Vertex{
		Position: sdl.FPoint{X: 100.5, Y: 100.5},
		Color:    color1,
		TexCoord: sdl.FPoint{X: 0, Y: 0},
	})
	vertices = append(vertices, sdl.Vertex{
		Position: sdl.FPoint{X: 200.5, Y: 100.5},
		Color:    color2,
		TexCoord: sdl.FPoint{X: 0, Y: 0},
	})
	vertices = append(vertices, sdl.Vertex{
		Position: sdl.FPoint{X: 100.5, Y: 200.5},
		Color:    color3,
		TexCoord: sdl.FPoint{X: 0, Y: 0},
	})
	renderer.RenderGeometry(nil, vertices, []int32{0, 1, 2})
}
*/
