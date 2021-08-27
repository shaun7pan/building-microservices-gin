package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// request to /hello will be handled by this function
	http.HandleFunc("/hello", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("hello world")
	})

	// others will be handled by this funciton
	http.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		log.Println("Running Hello handler")

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Println("Error reading body", err)
			http.Error(rw, "Unable to read request body", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(rw, "Hello %s", b)

	})

	log.Println("Starting server")
	err := http.ListenAndServe(":9090", nil)
	log.Fatal(err)
}
