package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var addr = ":9090"

func main() {
	l := log.New(os.Stdout, "products-api", log.LstdFlags)

	// create the handlers
	ph := NewProductHandler(l)

	sm := http.NewServeMux()
	sm.Handle("/", ph)

	s := http.Server{
		Addr:         addr,              // configure the bind address
		Handler:      sm,                // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	go func() {
		l.Println("Starting server on port :9090")
		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting the server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// block until a signal is received.
	sig := <-c
	log.Println("Got signal: ", sig)
	// gracefully shutdown the server, waiting max 30 seconds for current operations to finish
	ctx, _ := context.WithTimeout(context.Background(), 20*time.Second)
	s.Shutdown(ctx)
}
