package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/env"
	"github.com/shaun7pan/building-microservices-gin/product-images/files"
	"github.com/shaun7pan/building-microservices-gin/product-images/handlers"
)

var bindAddress = env.String("BIND_ADDRESS", false, ":9090", "Bind address for the server")
var logLevel = env.String("LOG_LEVEL", false, "debug", "Log output level for the server [debug, info, trace]")
var basePath = env.String("BASE_PATH", false, "./imagestore", "Base path to save images")

func main() {
	env.Parse()

	l := hclog.New(
		&hclog.LoggerOptions{
			Name:  "product-iamges",
			Level: hclog.LevelFromString(*logLevel),
		},
	)

	//create a logger for the server from default logger
	sl := l.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true})

	// create the storage class, use local storage
	// TODO: limit max filesize  to 5MB

	stor, err := files.NewLocal(*basePath, 5*1024*1024*1024)
	if err != nil {
		l.Error("Unable to create storage", "error", err)
		os.Exit(1)
	}

	// create the handlers
	fh := handlers.NewFiles(stor, l)

	// create router
	r := gin.Default()
	r.Use(limits.RequestSizeLimiter(stor.MaxFileSize()))
	r.Use(fh.GzipMiddleware)

	r.POST("/images/:id/:filename", fh.SaveFile)
	r.StaticFS("/images", http.Dir(*basePath))

	// create a new server
	s := http.Server{
		Addr:         *bindAddress,      // configure the bind address
		Handler:      r,                 // set the default handler
		ErrorLog:     sl,                // the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Info("Starting server", "bind_address", *bindAddress)

		err := s.ListenAndServe()
		if err != nil {
			l.Error("Unable to start server", "error", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	l.Info("Shutting down server with", "signal", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
