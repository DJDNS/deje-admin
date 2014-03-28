package main

// Set up DEJEController and run client servers
import (
	"github.com/campadrenalin/deje-admin/socket"
	"github.com/campadrenalin/go-deje"
	djserv "github.com/campadrenalin/go-deje/services"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	controller := deje.NewDEJEController()
	controller.Networker = djserv.RealIRCService{}

	// Listen for interruptions
	interrupter := make(chan os.Signal)
	signal.Notify(interrupter, syscall.SIGINT, syscall.SIGTERM)

	// Start services
	go run_http(controller)
	go socket.Run(controller)

	// Wait for interrupt
	<-interrupter
}
