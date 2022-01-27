package main

import (
	http2 "github.com/dedavidsamya/locationhistory/app/http"

	//"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	// use mux as request multiplexer
	router := mux.NewRouter()

	// define routes
	router.HandleFunc("/location/{order_id}", http2.AddLocation).Methods("PUT")
	router.HandleFunc("/location/{order_id}", http2.GetLocation).Methods("GET")
	router.HandleFunc("/location/{order_id}", http2.DeleteLocations).Methods("DELETE")
	// starting server
	error := http.ListenAndServe("localhost:8000", router)
	if error != nil {
		log.Fatal(error)
	}

}
