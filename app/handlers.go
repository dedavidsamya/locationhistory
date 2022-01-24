package app

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Location struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
}

var orders = map[string][]Location{}

func AddLocation(response http.ResponseWriter, request *http.Request) {
	// location here is going to store the Location passed by the client in the body of the request.
	location, err := UnmarshalOrder(request)
	fmt.Println("l. 28: location:", location, err)
	if err != nil {
		badInput := NewError("bad request", errors.New("the format of the input is invalid"))
		JSON(response, http.StatusBadRequest, &badInput)
		return
	}
	vars := mux.Vars(request)
	//this is commented because it should not exist. instead
	//of checking if empty string, I should trim strings received as input
	//if vars["order_id"] == " " {
	//	emptyString := NewError("bad request", errors.New("location number can't be empty"))
	//	JSON(response, http.StatusBadRequest, &emptyString)
	//	fmt.Println("l. 45: empty string as order_id")
	//	return
	//}
	orderId := vars["order_id"]
	fmt.Println("l. 37: orders before", orders)
	orders[orderId] = append([]Location{*location}, orders[orderId]...)
	fmt.Println("l. 39: orders after", orders)
	response.WriteHeader(http.StatusOK)
}

func UnmarshalOrder(request *http.Request) (*Location, error) {
	//this function receives the request and returns the data the client passed in the body, in the form of a Location struct, and an error
	var location Location
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Printf(" l. 53: Error when reading body of the request: %v", err)
		return nil, err
	}

	err = json.Unmarshal(body, &location)
	if err != nil {
		log.Printf("l. 59: Error when unmarsahlling body of the request: %v", err)
		return nil, err
	}

	return &location, err
}

func GetLocation(response http.ResponseWriter, request *http.Request) {
	//this function is completely wrong. The client queries a specific number of last locations for the location id.
	// The slice orders is a map containing structs of Order linked to a string (order_id).
	// when I retrieve orders, I am retrieving the whole map with different Orders (each one containing different locations)
	//So instead of retrieving orders, I should be retrieving each Order.
	//this is why I was trying to find the KEY, so I would be able to retrieve each Order by its number (the string order_id)

	//the way this function is working now, it is returning the whole orders map and not only the specific element of that map
	//that is a specific order_id (key of the map)

	vars := mux.Vars(request)
	InputOrderId := vars["order_id"]
	order := orders[InputOrderId]
	fmt.Println(" order before max:", order)

	//this checks if order is contained in orders
	if v, found := orders[InputOrderId]; found {
		fmt.Println(" found:", v)
		max, err := strconv.Atoi(request.URL.Query().Get("max"))
		fmt.Println("max:", max)
		fmt.Println("length of orders: ", len(orders))

		if err != nil {
			badInput := NewError("bad_request", errors.New("the format of the input is invalid"))
			JSON(response, http.StatusBadRequest, &badInput)
			return
		}
		switch {
		case max > 0 && max > len(orders):
			fmt.Println("order max > len", order)
			JSON(response, http.StatusOK, order)

		case max > 0 && max <= len(orders):
			fmt.Println("order max < len", order[:max])
			JSON(response, http.StatusOK, order[:max])

		default:
			badInput := NewError("bad_request", errors.New("this order does not exist"))
			JSON(response, http.StatusBadRequest, &badInput)
			return
		}
	}
}

func DeleteLocations(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["order_id"]
	fmt.Println("id: ", id)
	//here, id is exactly what I want it to be, is the id that client passes on "order_id"
	//this if is commented because it is the same issue at the put endpoint
	//if id == "" {
	//	// this is not working. When I pass an empty space on the query field,
	//	//I still receive the output of l.105, as if my order had been deleted
	//	//but there was no order id, so nothing to delete
	//	badInput := NewError("bad_request", errors.New("the order number can't be empty"))
	//	JSON(response, http.StatusBadRequest, badInput)
	//	fmt.Println("l. 108:", vars["order_id"])
	//} else {
	delete(orders, id)
	JSON(response, http.StatusOK, "Order deleted")
	fmt.Println("l. 112: orders map: ", orders)
}

func JSON(w http.ResponseWriter, code int, obj interface{}) {
	if obj != nil {
		bb, err := json.Marshal(obj)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(bb)
		return
	}
	w.WriteHeader(code)
}

func NewError(code string, err error) Error {
	e := Error{
		Code:    code,
		Message: err.Error(),
	}
	return e
}
