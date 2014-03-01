package main

import (
	djlogic "github.com/campadrenalin/go-deje/logic"
	djmodel "github.com/campadrenalin/go-deje/model"
	"reflect"
	"testing"
)

type GraphTest interface {
	GetExpected() GraphNode
	GetResult() GraphNode
	GetMessage() string
}

type GraphTestNode struct {
	Node       GraphNode
	Expected   GraphNode
	FailureMsg string
}

func (gtn *GraphTestNode) GetResult() GraphNode   { return gtn.Node }
func (gtn *GraphTestNode) GetExpected() GraphNode { return gtn.Expected }
func (gtn *GraphTestNode) GetMessage() string     { return gtn.FailureMsg }

type GraphTestDocument struct {
	Doc        djlogic.Document
	Expected   GraphNode
	FailureMsg string
}

func (gtd *GraphTestDocument) GetExpected() GraphNode { return gtd.Expected }
func (gtd *GraphTestDocument) GetMessage() string     { return gtd.FailureMsg }
func (gtd *GraphTestDocument) GetResult() GraphNode {
	root_node := NewRootNode()
	root_node.PopulateRoot(gtd.Doc)
	return root_node
}

func RunGraphTest(test GraphTest, t *testing.T) {
	expected, result := test.GetExpected(), test.GetResult()
	if !reflect.DeepEqual(expected, result) {
		t.Log(expected)
		t.Log(result)
		t.Fatal(test.GetMessage())
	}
}

func TestGraph_NewRootNode(t *testing.T) {
	test := &GraphTestNode{
		NewRootNode(),
		GraphNode{
			"root",
			"root",
			map[string]interface{}{
				"about": graphRootExplanation,
			},
			make([]GraphNode, 0),
		},
		"Expected different value for root node",
	}
	RunGraphTest(test, t)
}

func TestGraph_Root(t *testing.T) {
	test := &GraphTestDocument{
		djlogic.NewDocument(),
		NewRootNode(),
		"Empty document should just result in root node alone",
	}
	RunGraphTest(test, t)
}

func TestGraph_GetRootEvents(t *testing.T) {
	doc := djlogic.NewDocument()
	event := djmodel.NewEvent("example")
	event.Arguments["hello"] = "world"
	doc.Events.Register(event)

	root := NewRootNode()
	root_events := root.GetRootEvents(doc)
	expected := []djmodel.Event{event}
	if !reflect.DeepEqual(root_events, expected) {
		t.Fatalf("Expected %#v, got %#v", expected, root_events)
	}
}

func TestGraph_OneEvent(t *testing.T) {
	doc := djlogic.NewDocument()
	event := djmodel.NewEvent("example")
	event.Arguments["hello"] = "world"
	doc.Events.Register(event)

	root := NewRootNode()
	root.Children = append(root.Children, NewEventNode(event))
	test := &GraphTestDocument{
		doc,
		root,
		"Graph should have one event as child of root",
	}
	RunGraphTest(test, t)
}

func TestGraph_ChainAndFork(t *testing.T) {
	doc := djlogic.NewDocument()

	first := djmodel.NewEvent("first")
	second := djmodel.NewEvent("second")
	third := djmodel.NewEvent("third")
	fork := djmodel.NewEvent("fork")

	second.ParentHash = first.Hash()
	third.ParentHash = second.Hash()
	fork.ParentHash = first.Hash()

	for _, ev := range []djmodel.Event{first, second, third, fork} {
		doc.Events.Register(ev)
	}

	root := NewRootNode()
	node1 := NewEventNode(first)
	node2 := NewEventNode(second)
	node3 := NewEventNode(third)
	nodeF := NewEventNode(fork)

	node2.Children = []GraphNode{node3}
	node1.Children = []GraphNode{node2, nodeF}
	root.Children = []GraphNode{node1}
	test := &GraphTestDocument{
		doc,
		root,
		"Graph should have structure matching event tree",
	}
	RunGraphTest(test, t)
}
