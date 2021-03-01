package main

import (
	"bufio"
	"context"
	"io"
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
	router := http.NewServeMux()
	router.HandleFunc("/", stream)
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

func stream(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	endChan := make(chan struct{})

	go writeOutput(w, os.Stdin, endChan)

	select {
	case <-endChan:
	case <-ctx.Done():
		err := ctx.Err()
		log.Printf("Client disconnected: %s\n", err)
	}

	stop <- struct{}{}
	close(stop)
}

func writeOutput(w http.ResponseWriter, input io.ReadCloser, endChan chan struct{}) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming not supported", http.StatusInternalServerError)
		return
	}

	hdr := w.Header()
	hdr.Set("Cache-Control", "no-cache")
	hdr.Set("Connection", "keep-alive")
	hdr.Set("X-Content-Type-Options", "nosniff")

	reader := bufio.NewReader(input)

	for {
		in, _, err := reader.ReadLine()
		if err != nil && err == io.EOF {
			break
		}

		in = append(in, []byte("\n")...)
		if _, err := w.Write(in); err != nil {
			log.Fatalln(err)
		}

		flusher.Flush()
	}

	endChan <- struct{}{}
	close(endChan)
}
