package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/nicholasjackson/env"
	"github.com/shaun7pan/building-microservices-gin/product-api/data"
	"github.com/shaun7pan/building-microservices-gin/product-api/handlers"
)

var bindAddress = env.String("BINDADDRESS", false, ":9090", "Bind address for the server")

func main() {

	env.Parse()
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	// Create a gin router with default middleware
	r := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("sku", data.Sku)
	}

	// new handlers
	ph := handlers.NewProducts(l)

	// Use middleware
	r.Use(ph.CustomMiddleware)
	r.Use(ph.BuildCustomMiddleware())
	// handle cors
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:3000"},
	}))

	r.GET("/products", ph.ListAll)
	r.GET("/products/:id", ph.ListSingle)

	r.PUT("/products/:id", ph.Update)

	r.POST("/products", ph.Create)
	r.DELETE("/products", ph.Delete)

	//serve swagger docs
	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	// r.StaticFile("/swagger.yaml", "./swagger.yaml")
	r.GET("/docs", gin.WrapH(sh))
	r.GET("/swagger.yaml", gin.WrapH(http.FileServer(http.Dir("./"))))

	//create a new server
	s := http.Server{
		Addr:    *bindAddress,
		Handler: r,
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

	//gracefully shutdown the server, waiting max 30 seconds for current
	//operations to complete

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	// catching ctx.Done(). timeout of 30 secconds.

	select {
	case <-ctx.Done():
		log.Println("timeout of 30 seconds.")
	}
	log.Println("Server exiting")

}
