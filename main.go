package main

// Set up DEJEController and run client servers
import (
	"github.com/campadrenalin/go-deje"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	controller := deje.NewDEJEController()

	// Listen for interruptions
	interrupter := make(chan os.Signal)
	signal.Notify(interrupter, syscall.SIGINT, syscall.SIGTERM)

	// Start services
	go run_http(controller)
	go run_sio(controller)

	// Wait for interrupt
	<-interrupter
}
