package main

import (
	"github.com/martini-contrib/encoder"
	"net/http"
)

type GraphNode struct {
	Label    string                 `json:"label"`
	Type     string                 `json:"type"`
	Details  map[string]interface{} `json:"details"`
	Children []*GraphNode           `json:"children"`
}

func NewRootNode() GraphNode {
	return GraphNode{
		Label:    "root",
		Type:     "root",
		Details:  make(map[string]interface{}),
		Children: make([]*GraphNode, 0),
	}
}

func do_events_json(enc encoder.Encoder) (int, []byte) {
	root_node := NewRootNode()
	return http.StatusOK, encoder.Must(enc.Encode(root_node))
}
