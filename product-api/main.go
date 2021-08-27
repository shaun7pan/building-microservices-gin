package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/nicholasjackson/env"
	"github.com/shaun7pan/building-microservices-gin/product-api/handlers"
)

var bindAddress = env.String("BINDADDRESS", false, ":9090", "Bind address for the server")

func main() {

	env.Parse()
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	// new handlers
	h := handlers.NewHello(l)
	g := handlers.NewGoodbye(l)

	// create new serve mux
	sm := http.NewServeMux()

	// register handlers
	sm.Handle("/", h)
	sm.Handle("/goodbye", g)

	//create a new server
	s := http.Server{
		Addr:         *bindAddress,
		Handler:      sm,
		ErrorLog:     l,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	//start server
	go func() {
		l.Printf("Starting server on port: %s", *bindAddress)
		err := s.ListenAndServe()

		if err != nil {
			log.Fatal(err)
		}
	}()

	//trap sigterm and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block untill a sinal is received

	sig := <-c
	l.Println("Got signal:", sig)

}
