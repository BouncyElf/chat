package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BouncyElf/chat/handlers"

	"github.com/aofei/air"
)

func main() {
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		if err := air.Serve(); err != nil {
			air.ERROR(err.Error())
		}
	}()

	<-shutdownChan
	handlers.CloseSocket()
	air.INFO("shutting down the server")
	if air.DebugMode {
		air.Shutdown(time.Duration(1) * time.Second)
	} else {
		air.Shutdown(time.Duration(3) * time.Minute)
	}
	air.INFO("server gracefully stopped")
}
