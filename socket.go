package main

// Socket.IO connection stuff
import (
	"encoding/json"
	"errors"
	"github.com/campadrenalin/go-deje"
	djlogic "github.com/campadrenalin/go-deje/logic"
	"github.com/campadrenalin/go-deje/model"
	djstate "github.com/campadrenalin/go-deje/state"
	"github.com/googollee/go-socket.io"
	"log"
	"net/http"
)

type strchan chan string

func sio_irc_loop(c *deje.DEJEController, ns *socketio.NameSpace) {
	defer log.Printf("Ending sio_irc_loop for %v", ns.Id())

	// Set up stopper and cgetter chans
	stopper := ns.Session.Values["stopper"].(chan interface{})
	cgetter := ns.Session.Values["cgetter"].(chan strchan)

	// Create IRC channel listener (nil)
	var channel strchan

	// Listen for any
	for {
		select {
		case <-stopper:
			return
		case line := <-channel:
			ns.Emit("output", line)
		case channel = <-cgetter:
			continue // Do nothing
		}
	}
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

func get_document(c *deje.DEJEController, ns *socketio.NameSpace) (*djlogic.Document, error) {
	loc, ok := ns.Session.Values["location"]
	if !ok {
		return nil, errors.New("Not subscribed yet, cannot publish events")
	}
	location := loc.(model.IRCLocation)
	doc := c.GetDocument(location)
	return &doc, nil
}

func run_sio(controller *deje.DEJEController) {
	sock_config := &socketio.Config{}
	sock_config.HeartbeatTimeout = 2
	sock_config.ClosingTimeout = 4

	sio := socketio.NewSocketIOServer(sock_config)

	sio.On("connect", func(ns *socketio.NameSpace) {
		ns.Session.Values["stopper"] = make(chan interface{})
		ns.Session.Values["cgetter"] = make(chan strchan)
		go sio_irc_loop(controller, ns)

		id, endpoint := ns.Id(), ns.Endpoint()
		log.Printf("connected: %v in channel %v", id, endpoint)
	})
	sio.On("disconnect", func(ns *socketio.NameSpace) {
		stopper := ns.Session.Values["stopper"].(chan interface{})
		stopper <- nil

		id, endpoint := ns.Id(), ns.Endpoint()
		log.Printf("disconnected: %v in channel %v", id, endpoint)
	})
	sio.On("subscribe", func(ns *socketio.NameSpace, url string) {
		location := model.IRCLocation{}
		parse_err := location.ParseFrom(url)
		if parse_err != nil {
			ns.Emit("error", parse_err.Error())
			return
		}
		ns.Session.Values["location"] = location

		channel := controller.Networker.GetChannel(location)
		cgetter := ns.Session.Values["cgetter"].(chan strchan)
		cgetter <- channel.Incoming
		channel.Incoming <- "Subscribed to " + url

		doc, err := get_document(controller, ns)
		if err != nil {
			ns.Emit("error", err.Error())
			return
		}
		primitive := &djstate.SetPrimitive{
			Path:  []interface{}{},
			Value: doc.State.Export(),
		}
		ns.Emit("primitive", wrap_primitive(primitive))
		//log.Printf("Subscribed!")
	})
	sio.On("event", func(ns *socketio.NameSpace, evstr string) {
		log.Printf("Incoming event: %s", evstr)
		doc, err := get_document(controller, ns)
		if err != nil {
			ns.Emit("error", err.Error())
			return
		}

		var event model.Event
		err = json.Unmarshal([]byte(evstr), &event)
		if err != nil {
			ns.Emit("error", err.Error())
			ns.Emit("event_error", err.Error())
			log.Printf("Invalid event: %v", err)
			return
		}

		doc.Events.Register(event)
		ns.Emit("event_registered", event.Hash())
	})
	sio.On("stats_request", func(ns *socketio.NameSpace) {
		doc, err := get_document(controller, ns)
		if err != nil {
			ns.Emit("error", err.Error())
			return
		}
		data := map[string]interface{}{
			"num_events":  doc.Events.Length(),
			"num_quorums": doc.Quorums.Length(),
			"num_ts":      doc.Timestamps.Length(),
		}
		ns.Emit("stats", data)
	})
	sio.On("goto_request", func(ns *socketio.NameSpace, hash string) {
		doc, err := get_document(controller, ns)
		if err != nil {
			ns.Emit("error", err.Error())
			return
		}
		raw, ok := doc.Events.GetByKey(hash)
		if !ok {
			ns.Emit("error", "Hash "+hash+" not found")
		}
		event := djlogic.Event{raw.(model.Event), doc}
		err = event.Goto()
		if err != nil {
			ns.Emit("error", err.Error())
			return
		}
		primitive := &djstate.SetPrimitive{
			Path:  []interface{}{},
			Value: doc.State.Export(),
		}
		ns.Emit("primitive", wrap_primitive(primitive))
	})

	http.ListenAndServe(":3001", sio)
}
