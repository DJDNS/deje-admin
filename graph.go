package main

import (
	djlogic "github.com/campadrenalin/go-deje/logic"
	djmodel "github.com/campadrenalin/go-deje/model"
	"github.com/martini-contrib/encoder"
	"net/http"
)

type GraphNode struct {
	Label    string                 `json:"label"`
	Type     string                 `json:"type"`
	Details  map[string]interface{} `json:"details"`
	Children []*GraphNode           `json:"children"`
}

var graphRootExplanation = `Root node.

Any object with no parent, or an invalid/unknown parent,
is treated like a child of root for graphing purposes.`

func NewRootNode() GraphNode {
	return GraphNode{
		Label:    "root",
		Type:     "root",
		Details:  map[string]interface{}{"about": graphRootExplanation},
		Children: make([]*GraphNode, 0),
	}
}

func NewEventNode(ev djmodel.Event) GraphNode {
	return GraphNode{
		Label:    ev.HandlerName,
		Type:     "event",
		Details:  make(map[string]interface{}),
		Children: make([]*GraphNode, 0),
	}
}

func (gn GraphNode) GetRootEvents(doc djlogic.Document) []djmodel.Event {
	ev_map := make(map[string]djmodel.Event)
	result := make([]djmodel.Event, 0)
	doc.Events.SerializeTo(ev_map)
	for _, event := range ev_map {
		result = append(result, event)
	}
	return result
}

func (gn GraphNode) Populate(doc djlogic.Document) {
	root_events := gn.GetRootEvents(doc)
	for _, event := range root_events {
		ev_node := NewEventNode(event)
		gn.Children = append(gn.Children, &ev_node)
	}
}

func do_events_json(doc djlogic.Document, enc encoder.Encoder) (int, []byte) {
	root_node := NewRootNode()
	root_node.Populate(doc)
	return http.StatusOK, encoder.Must(enc.Encode(root_node))
}
