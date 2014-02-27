package main

import (
	djlogic "github.com/campadrenalin/go-deje/logic"
	"reflect"
	"testing"
	//djmodel "github.com/campadrenalin/go-deje/model"
)

type GraphTest struct {
	Doc        djlogic.Document
	Result     GraphNode
	FailureMsg string
}

func (test *GraphTest) Run(t *testing.T) {
	root_node := NewRootNode()
	root_node.Populate(test.Doc)

	if !reflect.DeepEqual(test.Result, root_node) {
		t.Log(test.Result)
		t.Log(root_node)
		t.Fatal(test.FailureMsg)
	}
}

func TestGraph_Root(t *testing.T) {
	test := &GraphTest{
		djlogic.NewDocument(),
		GraphNode{
			"root",
			"root",
			map[string]interface{}{
				"about": graphRootExplanation,
			},
			make([]*GraphNode, 0),
		},
		"Empty document should just result in root node alone",
	}
	test.Run(t)
}
