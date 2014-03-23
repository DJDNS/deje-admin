package socket

import (
	deje "github.com/campadrenalin/go-deje"
	djlogic "github.com/campadrenalin/go-deje/logic"
	djmodel "github.com/campadrenalin/go-deje/model"
	djserv "github.com/campadrenalin/go-deje/services"
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
	Document djlogic.Document
	IRC      djserv.IRCChannel
	stopper  chan struct{}
}

func NewSubscription(c *deje.DEJEController, url string) (*Subscription, error) {
	location := djmodel.IRCLocation{}
	err := location.ParseFrom(url)
	if err != nil {
		return nil, err
	}
	return &Subscription{
		c.GetDocument(location),
		c.Networker.GetChannel(location),
		make(chan struct{}),
	}, nil
}

func (s *Subscription) Run(ns *socketio.NameSpace) {
	defer log.Printf("Stopping Subscription for %v", ns.Id())
	for {
		select {
		case line := <-s.IRC.Incoming:
			ns.Emit("output", line)
		case <-s.stopper:
			return
		}
	}
}

func (s *Subscription) Stop() {
	close(s.stopper)
}
