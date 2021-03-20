package main

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/romiras/bv/filters"
)

type StreamHandler struct {
	Filters []filters.IFilter
}

func NewStreamHandler() *StreamHandler {

	return &StreamHandler{
		Filters: make([]filters.IFilter, 0),
	}
}

func (sh *StreamHandler) SetFilters(filters []string) {
	if len(filters) > 0 {
		for _, filterName := range filters {
			switch filterName {
			}
		}
	}
}

func (sh *StreamHandler) applyFilters(reader io.Reader) (io.Reader, error) {
	// TODO Implement chain filtering here
	return reader, nil
}

func (sh *StreamHandler) Handle(w http.ResponseWriter, r *http.Request) {
	defer func() {
		stop <- struct{}{}
		close(stop)
	}()

	ctx := r.Context()
	endChan := make(chan struct{})

	filteredReader, err := sh.applyFilters(os.Stdin)
	if err != nil {
		return
	}
	go sh.writeOutput(w, filteredReader, endChan)

	select {
	case <-endChan:
	case <-ctx.Done():
		err := ctx.Err()
		log.Printf("Client disconnected: %s\n", err)
	}
}

func (sh *StreamHandler) writeOutput(w http.ResponseWriter, input io.Reader, endChan chan struct{}) {
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
