package main

// Socket.IO connection stuff
import (
	"github.com/campadrenalin/go-deje"
	"github.com/googollee/go-socket.io"
	"log"
	"net/http"
)

func sio_onConnect(ns *socketio.NameSpace) {
	log.Printf("connected: %v in channel %v", ns.Id(), ns.Endpoint())
	ns.Emit("output", "Hello from server")
}

func sio_onDisconnect(ns *socketio.NameSpace) {
	log.Printf("disconnected: %v in channel %v", ns.Id(), ns.Endpoint())
}

func run_sio(controller *deje.DEJEController) {
	sock_config := &socketio.Config{}
	sock_config.HeartbeatTimeout = 2
	sock_config.ClosingTimeout = 4

	sio := socketio.NewSocketIOServer(sock_config)

	sio.Of("/irc").On("connect", sio_onConnect)
	sio.Of("/irc").On("disconnect", sio_onDisconnect)

	http.ListenAndServe(":3001", sio)
}
