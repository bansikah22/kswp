package analyzer

import (
	"github.com/bansikah22/kswp/pkg/models"
)

type Node struct {
	Resource models.Resource
	Children []*Node
}

func BuildDependencyGraph(resources []models.Resource) *Node {
	root := &Node{
		Resource: models.Resource{
			Name: "root",
			Kind: "Root",
		},
	}
	// This is a placeholder for the dependency analysis logic.
	// For now, we will just add all resources as children of the root.
	for _, res := range resources {
		root.Children = append(root.Children, &Node{Resource: res})
	}
	return root
}
