package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"time"
)

var stop chan struct{}
var port = "8080"

func main() {
	reg := NewRegistry()

	sh := NewStreamHandler()
	sh.SetFilters(reg.filters)

	router := http.NewServeMux()
	router.HandleFunc("/", sh.Handle)
	http.Handle("/", router)

	ctx, cancel := context.WithCancel(context.Background())
	server := &http.Server{
		Addr:        ":" + port,
		Handler:     router,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}

	// Setting up capturing shutdown
	stop = make(chan struct{})

	// Run server
	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalln(err)
		}
	}()

	go func() {
		cmd := exec.Command("xdg-open", "http://localhost:"+port) // Linux desktop
		if err := cmd.Run(); err != nil {
			log.Fatalln(err)
			// return
		}
	}()

	// Wait for ListenAndServe goroutine to close.

	// Waiting for stop
	<-stop

	gracefullCtx, cancelShutdown := context.WithTimeout(context.Background(), time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(gracefullCtx); err != nil {
		log.Fatalln(err)
	}

	// manually cancel context if not using httpServer.RegisterOnShutdown(cancel)
	cancel()

	defer os.Exit(0)
}
