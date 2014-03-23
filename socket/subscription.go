package socket

import (
	deje "github.com/campadrenalin/go-deje"
	djbroad "github.com/campadrenalin/go-deje/broadcast"
	djlogic "github.com/campadrenalin/go-deje/logic"
	djmodel "github.com/campadrenalin/go-deje/model"
	djserv "github.com/campadrenalin/go-deje/services"
	djstate "github.com/campadrenalin/go-deje/state"
	"github.com/googollee/go-socket.io"
	"log"
)

// A Socket.IO connection may subscribe to a Document, and
// therefore listen to certain events that pertain to that
// Document.
//
// The Subscription object "lives" via a goroutine that
// delivers updates to Socket.IO connection. This is terminated
// upon disconnect.
type Subscription struct {
	Document   djlogic.Document
	IRC        djserv.IRCChannel
	Primitives *djbroad.Subscription
	stopper    chan struct{}
}

func NewSubscription(c *deje.DEJEController, url string) (*Subscription, error) {
	location := djmodel.IRCLocation{}
	err := location.ParseFrom(url)
	if err != nil {
		return nil, err
	}
	doc := c.GetDocument(location)
	return &Subscription{
		doc,
		c.Networker.GetChannel(location),
		doc.State.Subscribe(),
		make(chan struct{}),
	}, nil
}

// Send the current docstate as a full-replacement SET primitive.
func (s *Subscription) SendState(ns *socketio.NameSpace) {
	primitive := &djstate.SetPrimitive{
		Path:  []interface{}{},
		Value: s.Document.State.Export(),
	}
	ns.Emit("primitive", wrap_primitive(primitive))
}

func (s *Subscription) Run(ns *socketio.NameSpace) {
	defer log.Printf("Stopping Subscription for %v", ns.Id())
	for {
		select {
		case line := <-s.IRC.Incoming:
			ns.Emit("output", line)
		case prim := <-s.Primitives.Out():
			primitive := prim.(djstate.Primitive)
			ns.Emit("primitive", wrap_primitive(primitive))
		case <-s.stopper:
			return
		}
	}
}

func (s *Subscription) Stop() {
	close(s.stopper)
}

type PrimitiveWrapper struct {
	Type      string
	Arguments djstate.Primitive `json:"args"`
}

func wrap_primitive(p djstate.Primitive) PrimitiveWrapper {
	wrapper := PrimitiveWrapper{
		Arguments: p,
	}
	switch p.(type) {
	case *djstate.SetPrimitive:
		wrapper.Type = "SET"
	case *djstate.DeletePrimitive:
		wrapper.Type = "DELETE"
	default:
		wrapper.Type = "unknown type"
	}
	return wrapper
}
