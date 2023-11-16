package entities

import (
	"errors"
	"github.com/tsagae/software3d/pkg/basics"
	"strings"
)

// The root has name "world"
type SceneGraph struct {
	root  *SceneGraphNode
	nodes map[string]*SceneGraphNode
}

type SceneGraphNode struct {
	nodeName          string
	parentNode        *SceneGraphNode
	childNodes        []*SceneGraphNode //children are not ordered and the order may change at runtime
	toParentTransform basics.Transform
	GameObject        GameObject
	//caching world transform and a bool to refresh it ?
}

func NewSceneGraph() SceneGraph {
	worldNode := newWorldNode()
	nodes := make(map[string]*SceneGraphNode)
	sceneGraph := SceneGraph{
		root:  worldNode,
		nodes: nodes,
	}
	nodes[worldNode.nodeName] = worldNode
	return sceneGraph
}

// Adds a child provided the name of the parent
func (sceneGraph *SceneGraph) AddChild(parentName string, childNode *SceneGraphNode, toParentTransform basics.Transform) (*SceneGraphNode, error) {
	_, ok := sceneGraph.nodes[childNode.nodeName]
	if ok {
		//node already exists
		return nil, errors.New("a node with the same name already exists")
	}
	parentNode, ok := sceneGraph.nodes[parentName]
	if ok {
		//parent node found
		childNode.toParentTransform = toParentTransform
		childNode.parentNode = parentNode
		parentNode.childNodes = append(parentNode.childNodes, childNode)
		sceneGraph.nodes[childNode.nodeName] = childNode
	} else {
		return nil, errors.New("the parent node does not exist")
	}
	return childNode, nil
}

// Removes a node and all of its children
func (sceneGraph *SceneGraph) RemoveChild(nodeName string) {
	nodeToDelete, ok := sceneGraph.nodes[nodeName]
	if ok {
		delete(sceneGraph.nodes, nodeName)
		children := nodeToDelete.parentNode.childNodes
		for i := 0; i < len(children); i++ {
			if children[i] == nodeToDelete {
				//swap with the last element and remove from slice
				children[i], children[len(children)-1] = children[len(children)-1], children[i]
				nodeToDelete.parentNode.childNodes = children[:len(children)-1]
				nodeToDelete.parentNode = nil
				return
			}
		}
	}
}

// Unordered list of node names
func (scenegGraph *SceneGraph) ListNodes() []string {
	keys := make([]string, 0, len(scenegGraph.nodes))
	for k := range scenegGraph.nodes {
		keys = append(keys, k)
	}
	return keys
}

func (sceneGraph *SceneGraph) ToString() string {
	return sceneGraph.root.sceneGraphToString()
}

// Returns a scene graph node from its name, nil if its not found
func (sceneGraph *SceneGraph) GetNode(nodeName string) *SceneGraphNode {
	node, ok := sceneGraph.nodes[nodeName]
	if ok {
		return node
	}
	return nil
}

// Returns root node of the scene graph
func (sceneGraph *SceneGraph) GetRoot() *SceneGraphNode {
	return sceneGraph.root
}

func newWorldNode() *SceneGraphNode {
	return NewSceneGraphNode(NewEmptyObject("worldObj"), "world")
}

func NewSceneGraphNode(gameObject GameObject, nodeName string) *SceneGraphNode {
	return &SceneGraphNode{
		nodeName:          nodeName,
		parentNode:        nil,
		childNodes:        make([]*SceneGraphNode, 0),
		toParentTransform: basics.NewZeroTransform(),
		GameObject:        gameObject,
	}
}

// Returns the name of the node
func (node *SceneGraphNode) GetName() string {
	return node.nodeName
}

// Returns the parent node
func (node *SceneGraphNode) GetParent() *SceneGraphNode {
	return node.parentNode
}

// Returns a slice containing pointers to all the child nodes
func (node *SceneGraphNode) GetChildren() []*SceneGraphNode {
	src := node.childNodes
	dst := make([]*SceneGraphNode, len(src))
	copy(dst, src)
	return dst
}

func (node *SceneGraphNode) CumulateWorldTransform(t *basics.Transform) {
	worldT := node.GetWorldTransform()
	worldT.ThisCumulate(t)
	parentWorldT := node.parentNode.GetWorldTransform()
	parentWorldT.ThisInvert()
	node.toParentTransform = worldT.Cumulate(&parentWorldT)
}

func (node *SceneGraphNode) CumulateBeforeLocalTranform(t *basics.Transform) {
	node.toParentTransform = t.Cumulate(&node.toParentTransform)
}

func (node *SceneGraphNode) CumulateLocalTransform(t *basics.Transform) {
	node.toParentTransform.ThisCumulate(t)
}

func (node *SceneGraphNode) GetWorldTransform() basics.Transform {
	tempNode := node
	worldT := basics.NewZeroTransform()
	for tempNode != nil {
		worldT.ThisCumulate(&tempNode.toParentTransform)
		tempNode = tempNode.parentNode
	}
	return worldT
}

func (node *SceneGraphNode) SetViewRotation(yaw basics.Scalar, pitch basics.Scalar) {
	//fmt.Println("yaw: ", yaw, " pitch: ", pitch)
	//fmt.Println("worldToLocal: ", o.LocalToWorldTransform)
	newRotation := basics.NewQuaternionFromEulerAngles(yaw, pitch, 0)
	//fmt.Println("newRotation", newRotation)
	node.toParentTransform.Rotation = newRotation
}

func (node *SceneGraphNode) GetOrientation() basics.Matrix3 {
	locToWRot := node.GetWorldTransform().Rotation
	matrix3 := basics.NewCanonicalMatrix3()
	matrix3[0] = locToWRot.Rotated(matrix3[0])
	matrix3[1] = locToWRot.Rotated(matrix3[1])
	matrix3[2] = locToWRot.Rotated(matrix3[2])
	return matrix3
}

func (node *SceneGraphNode) rescursiveToString(depth int) string {
	outStr := ""
	outStr += strings.Repeat("\t", depth)
	if node.GameObject.GetName() == "" {
		outStr += "unnamedGameObj"
	} else {
		outStr += node.GameObject.GetName()
	}
	outStr += "\n"
	for _, child := range node.childNodes {
		outStr += child.rescursiveToString(depth + 1)
	}
	return outStr
}

func (node *SceneGraphNode) sceneGraphToString() string {
	return node.rescursiveToString(0)
}

/*
	func (o *InstanceObject) LookAtPoint(point *basics.Vector3) {
		direction := point.Sub(o.LocalToWorldTransform.Position)
		o.LookAtDirection(&direction)
	}

// look at direction

	func (o *InstanceObject) LookAtDirection(direction *basics.Vector3) {
		direction.ThisNormalize()
		orientation := o.GetOrientation()
		lookAtRotation := orientation.LookAt(direction)
		o.LocalRotate(&o.WorldToLocalTransform.Rotation)
		o.LocalRotate(&lookAtRotation)
	}
*/
