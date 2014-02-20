package main

// Socket.IO connection stuff
import (
	"github.com/campadrenalin/go-deje"
	"github.com/googollee/go-socket.io"
	"log"
	"net/http"
	"time"
)

func sio_irc_loop(c *deje.DEJEController, ns *socketio.NameSpace) {
	// Create stopper signal
	stopper := make(chan interface{})
	ns.Session.Values["stopper"] = stopper

	// Create channel listener (nil)
	var channel chan string
	ns.Session.Values["channel"] = channel

	// Listen for any
	for {
		select {
		case <-stopper:
			break
		case line := <-channel:
			ns.Emit("output", line)
		case <-time.After(1 * time.Second):
			channel = ns.Session.Values["channel"].(chan string)
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
		id, endpoint := ns.Id(), ns.Endpoint()
		log.Printf("connected: %v in channel %v", id, endpoint)
		go sio_irc_loop(controller, ns)
	})
	irc.On("disconnect", func(ns *socketio.NameSpace) {
		id, endpoint := ns.Id(), ns.Endpoint()
		log.Printf("disconnected: %v in channel %v", id, endpoint)
	})
	irc.On("subscribe", func(ns *socketio.NameSpace, url string) {
		channel := make(chan string)
		ns.Session.Values["channel"] = channel
		channel <- "Subscribed to " + url
	})

	http.ListenAndServe(":3001", sio)
}
