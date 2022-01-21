package main

import (
	"github.com/dedavidsamya/locationhistory/app"
	//"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {

	// use mux as request multiplexer
	router := mux.NewRouter()

	// define routes
	router.HandleFunc("/location/{order_id}", app.AddLocation).Methods("PUT")
	router.HandleFunc("/location/{order_id}", app.GetLocation).Methods("GET")
	router.HandleFunc("/location/{order_id}", app.DeleteLocations).Methods("DELETE")
	// starting server
	error := http.ListenAndServe("localhost:8000", router)
	if error != nil {
		log.Fatal(error)
	}

}
