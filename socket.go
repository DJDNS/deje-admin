package main

// Socket.IO connection stuff
import (
	"github.com/campadrenalin/go-deje"
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

func run_sio(controller *deje.DEJEController) {
	sock_config := &socketio.Config{}
	sock_config.HeartbeatTimeout = 2
	sock_config.ClosingTimeout = 4

	sio := socketio.NewSocketIOServer(sock_config)

	irc := sio.Of("/irc")
	irc.On("connect", func(ns *socketio.NameSpace) {
		ns.Session.Values["stopper"] = make(chan interface{})
		ns.Session.Values["cgetter"] = make(chan strchan)
		go sio_irc_loop(controller, ns)

		id, endpoint := ns.Id(), ns.Endpoint()
		log.Printf("connected: %v in channel %v", id, endpoint)
	})
	irc.On("disconnect", func(ns *socketio.NameSpace) {
		stopper := ns.Session.Values["stopper"].(chan interface{})
		stopper <- nil

		id, endpoint := ns.Id(), ns.Endpoint()
		log.Printf("disconnected: %v in channel %v", id, endpoint)
	})
	irc.On("subscribe", func(ns *socketio.NameSpace, url string) {
		channel := make(chan string)
		cgetter := ns.Session.Values["cgetter"].(chan strchan)
		cgetter <- channel
		channel <- "Subscribed to " + url
	})

	http.ListenAndServe(":3001", sio)
}
