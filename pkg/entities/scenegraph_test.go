package entities

import (
	"GoSDL/pkg/basics"
	"testing"
)

func TestTranformations(t *testing.T) {

	/*
		cubeMesh, err := graphics.ReadMeshFromFile("cube.obj", *basics.White(155))
		if err != nil {
			panic(err)
		}
	*/

	cubeObj := NewEmptyObject("cubeObj")

	sceneGraph := NewSceneGraph()
	sceneGraph.AddChild("world", NewSceneGraphNode(cubeObj, "cube"), basics.NewZeroTransform())
	//nodes := []*entities.SceneGraphNode{cubeNode}
	cubeNode := sceneGraph.GetNode("cube")

	yRotationTransformation := basics.NewTransform(1, basics.NewQuaternionFromAngleAndAxis(1, basics.Up()), basics.NewVector3(0, 0, 0))

	tr := basics.NewTransform(1, basics.NewIdentityQuaternion(), basics.NewVector3(1, 0, -1))

	sceneGraph.AddChild("cubeNode", NewSceneGraphNode(cubeObj, "secondCube"), tr)

	for i := 0; i < 100; i++ {
		cubeNode.CumulateLocalTransform(&yRotationTransformation)
	}

	if false {
		t.Errorf("error")
	}
}

func TestLocToWorldTransformTree(t *testing.T) {

}
