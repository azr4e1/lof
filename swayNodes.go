package lof

import (
	"encoding/json"
)

type ContainerType string
type FilterCriteria func(*BaseNode) bool

const (
	SimpleContainer     ContainerType = "con"
	FloatingContainer   ContainerType = "floating_con"
	OutputContainer     ContainerType = "output"
	RootContainer       ContainerType = "root"
	WorkspaceContainer  ContainerType = "workspace"
	ScratchpadContainer ContainerType = "__i3_scratch"
)

type BaseNode struct {
	Id        int           `json:"id"`
	Name      string        `json:"name"`
	AppId     string        `json:"app_id"`
	Marks     []string      `json:"marks"`
	Type      ContainerType `json:"type"`
	Workspace string        `json:"workspace"`
	Focused   bool          `json:"focused"`
}

type Node struct {
	BaseNode
	Nodes         []*Node `json:"nodes"`
	FloatingNodes []*Node `json:"floating_nodes"`
}

type FlattenedNodes []*BaseNode

func (bn *BaseNode) IsContainer() bool {
	if bn == nil {
		return false
	}

	if bn.Type == SimpleContainer || bn.Type == FloatingContainer {
		return true
	}

	return false
}

func (bn *BaseNode) IsTrueWindow() bool {
	if bn == nil {
		return false
	}

	if bn.IsContainer() && bn.Name != "" {
		return true
	}
	return false
}

// TODO: to test
func (n *Node) Flatten() FlattenedNodes {
	var flatNodes = []*BaseNode{}
	if n == nil {
		return flatNodes
	}

	// append current
	flatNodes = append(flatNodes, &BaseNode{
		Id:        n.Id,
		Name:      n.Name,
		AppId:     n.AppId,
		Marks:     n.Marks,
		Type:      n.Type,
		Workspace: n.Workspace,
		Focused:   n.Focused,
	})

	workspace := n.Workspace
	if n.Type == WorkspaceContainer {
		workspace = n.Name
	}
	for _, node := range append(n.Nodes, n.FloatingNodes...) {
		node.Workspace = workspace
		newNodes := node.Flatten()

		flatNodes = append(flatNodes, newNodes...)
	}

	return flatNodes
}

// TODO: to test
func (fn FlattenedNodes) Filter(criteria FilterCriteria) FlattenedNodes {
	filterNodes := []*BaseNode{}
	for _, el := range fn {
		if criteria(el) {
			filterNodes = append(filterNodes, el)
		}
	}

	return filterNodes
}

// TODO: to test
func (fn FlattenedNodes) ToJSON() ([]byte, error) {
	data, err := json.Marshal(fn)

	if err != nil {
		return nil, err
	}

	return data, nil
}
