package entities

import (
	"github.com/stretchr/testify/assert"
	"github.com/tsagae/software3d/pkg/basics"
	"testing"
)

func TestSceneGraph_String(t *testing.T) {

	cubeObj := NewEmptyObject("cubeObj")

	sceneGraph := NewSceneGraph()
	err := sceneGraph.AddChild("world", NewSceneGraphNode(cubeObj, "cube"), basics.NewZeroTransform())
	assert.Nil(t, err, "err should be nil")
	//nodes := []*entities.SceneGraphNode{cubeNode}

	tr := basics.NewTransform(1, basics.NewIdentityQuaternion(), basics.NewVector3(1, 0, -1))

	err = sceneGraph.AddChild("cube", NewSceneGraphNode(cubeObj, "secondCube"), tr)
	if err != nil {
		assert.Nil(t, err, "err should be nil")
	}

	assert.Equal(t, "world: worldObj\n\tcube: cubeObj\n\t\tsecondCube: cubeObj\n", sceneGraph.String())
}
