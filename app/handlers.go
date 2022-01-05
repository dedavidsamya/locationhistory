package app

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Order struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
}

var orders = map[string][]Order{}

func AddLocation(response http.ResponseWriter, request *http.Request) {
	order := UnmarshalOrder(request)
	//fmt.Println("order:", order)

	vars := mux.Vars(request)
	orderId := vars["order_id"]

	//fmt.Println("orders before", orders)
	orders[orderId] = append([]Order{order}, orders[orderId]...)
	//fmt.Println("orders after", orders)

	response.WriteHeader(http.StatusOK)

}

func UnmarshalOrder(request *http.Request) Order {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Fatal(err)
	}

	var order Order
	err = json.Unmarshal(body, &order)
	if err != nil {
		log.Fatal(err)
	}

	return order
}

//error handling
//example: what if the value of max is not 0, but negative, or a number that is out of the scope of the slice

func GetLocation(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	orderId := vars["order_id"]
	max, err := strconv.Atoi(request.URL.Query().Get("max"))
	if err != nil {
		http.NotFound(response, request)
	}
	if max != 0 {
		fmt.Println(orders[orderId][:max+1])
	}
}

func DeleteLocations(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["order_id"]
	delete(orders, id)
	respondWithJSON(response, http.StatusOK, nil)
	fmt.Println(orders)

}

func respondWithJSON(response http.ResponseWriter, statusCode int, data interface{}) {
	result, _ := json.Marshal(data)
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)
	response.Write(result)
}