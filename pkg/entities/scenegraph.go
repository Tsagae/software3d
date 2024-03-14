package entities

import (
	"errors"
	"fmt"
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

func NewSceneGraph() *SceneGraph {
	worldNode := newWorldNode()
	nodes := make(map[string]*SceneGraphNode)
	sceneGraph := SceneGraph{
		root:  worldNode,
		nodes: nodes,
	}
	nodes[worldNode.nodeName] = worldNode
	return &sceneGraph
}

// AddChild Adds a child provided the name of the parent. Returns an error if a node with the same name already exists or the parent node does not exist
func (sceneGraph *SceneGraph) AddChild(parentName string, childNode *SceneGraphNode, toParentTransform basics.Transform) error {
	if _, ok := sceneGraph.nodes[childNode.nodeName]; ok {
		//node already exists
		return errors.New("a node with the same name already exists")
	}

	if parentNode, ok := sceneGraph.nodes[parentName]; ok {
		//parent node found
		childNode.toParentTransform = toParentTransform
		childNode.parentNode = parentNode
		parentNode.childNodes = append(parentNode.childNodes, childNode)
		sceneGraph.nodes[childNode.nodeName] = childNode
	} else {
		return errors.New("the parent node does not exist")
	}
	return nil
}

// RemoveChild Removes a node and all of its children
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

// ListNodes Unordered list of node names
func (sceneGraph *SceneGraph) ListNodes() []string {
	keys := make([]string, 0, len(sceneGraph.nodes))
	for k := range sceneGraph.nodes {
		keys = append(keys, k)
	}
	return keys
}

func (sceneGraph *SceneGraph) String() string {
	return sceneGraph.root.sceneGraphToString()
}

// GetNode Returns a scene graph node given its name, nil if it's not found
func (sceneGraph *SceneGraph) GetNode(nodeName string) *SceneGraphNode {
	node, ok := sceneGraph.nodes[nodeName]
	if ok {
		return node
	}
	return nil
}

// GetRoot Returns the root node of the scene graph
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

// Name GetName Returns the name of the node
func (node *SceneGraphNode) Name() string {
	return node.nodeName
}

// Parent Returns the parent node
func (node *SceneGraphNode) Parent() *SceneGraphNode {
	return node.parentNode
}

// Children Returns a slice containing pointers to all the child nodes
func (node *SceneGraphNode) Children() []*SceneGraphNode {
	src := node.childNodes
	dst := make([]*SceneGraphNode, len(src))
	copy(dst, src)
	return dst
}

func (node *SceneGraphNode) CumulateWorldTransform(t *basics.Transform) {
	worldT := node.WorldTransform()
	worldT.ThisCumulate(t)
	parentWorldT := node.parentNode.WorldTransform()
	parentWorldT.ThisInvert()
	node.toParentTransform = worldT.Cumulate(&parentWorldT)
}

func (node *SceneGraphNode) CumulateBeforeLocalTranform(t *basics.Transform) {
	node.toParentTransform = t.Cumulate(&node.toParentTransform)
}

func (node *SceneGraphNode) CumulateLocalTransform(t *basics.Transform) {
	node.toParentTransform.ThisCumulate(t)
}

func (node *SceneGraphNode) WorldTransform() basics.Transform {
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

func (node *SceneGraphNode) Orientation() basics.Matrix3 {
	locToWRot := node.WorldTransform().Rotation
	matrix3 := basics.NewCanonicalMatrix3()
	matrix3[0] = locToWRot.Rotated(matrix3[0])
	matrix3[1] = locToWRot.Rotated(matrix3[1])
	matrix3[2] = locToWRot.Rotated(matrix3[2])
	return matrix3
}

func (node *SceneGraphNode) rescursiveToString(depth int) string {
	var sb strings.Builder
	sb.WriteString(strings.Repeat("\t", depth))
	if node.GameObject.Name() == "" {
		sb.WriteString("unnamedGameObj")
	} else {
		sb.WriteString(fmt.Sprintf("%v: %v", node.nodeName, node.GameObject.Name()))
	}
	sb.WriteString("\n")
	for _, child := range node.childNodes {
		sb.WriteString(child.rescursiveToString(depth + 1))
	}
	return sb.String()
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
		orientation := o.Orientation()
		lookAtRotation := orientation.LookAt(direction)
		o.LocalRotate(&o.WorldToLocalTransform.Rotation)
		o.LocalRotate(&lookAtRotation)
	}
*/
