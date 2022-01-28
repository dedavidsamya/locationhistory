package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dedavidsamya/locationhistory/app/db"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

const Limit int = 100

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// orders is a map that connects a key (order_id) to a slice of Location structs.
//so each order id is going to have a set of different {latitude, longitude} objects

func AddLocation(response http.ResponseWriter, request *http.Request) {
	//1. retrieve the order_id passed by the client on the URL.
	vars := mux.Vars(request)
	orderId := vars["order_id"]
	fmt.Println("orderId: ", orderId)

	//2. create an order that will be stored in the db
	//this func receives an id as argument, checks if it is valid (empty or not, already exists or not) and stores a new Order struct in db with the id passed.
	order, err := db.CreateOrder(orderId)
	if err != nil {
		errors.New("AddLocation: error within db.CreateOrder()")
	}
	fmt.Println("order created: ", order)
	fmt.Println("orders map: ", db.Orders)

	//3. read and unmarshall the body of the request, and store it into the location variable.
	location, err := UnmarshalLocation(request)
	fmt.Println("location (after unmarshal):", location, err)
	if err != nil {
		badInput := NewError("bad request", errors.New("the format of the input is invalid"))
		JSON(response, http.StatusBadRequest, &badInput)
		return
	}
	//4. append the location passed by the client to the order previously created
	err = db.AddLocation(orderId, *location)
	if err != nil {
		errors.New("AddLocation: error within db.AddLocation()")
	}
	fmt.Println("Orders map after addLocation: ", db.Orders)
	//5. send status response to client
	response.WriteHeader(http.StatusOK)
}

//this function receives the request and returns the data the client passed in the body, in the form of a Location struct, and an error
func UnmarshalLocation(request *http.Request) (*db.Location, error) {
	var location db.Location
	//here body refers to the data the client passed in the body. ReadAll is going to read and "transfer" this data to the variable body.
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		log.Printf("Error when reading body of the request: %v", err)
		return nil, err
	}
	// the unmarshal function transfers whatever is in the body variable to the location variable (which is now a pointer, so it is a specific place in memory)
	err = json.Unmarshal(body, &location)
	if err != nil {
		log.Printf("Error when unmarsahlling body of the request: %v", err)
		return nil, err
	}

	return &location, err
}

//
//func GetLocation(response http.ResponseWriter, request *http.Request) {
//	vars := mux.Vars(request)
//	inputOrderId := vars["order_id"]
//	fmt.Println(" inputOrderId: ", inputOrderId)
//
//	locations := orders[inputOrderId]
//	fmt.Println(" locations (for the order_id) before max:", locations)
//	locationsLength := len(locations)
//	fmt.Println("locationsLength", locationsLength)
//
//	//this checks if location is contained in orders
//	if v, found := orders[inputOrderId]; found {
//		fmt.Println(" []Location{location} found in the orders slice:", v)
//	}
//
//	// if max exists
//	queries := request.URL.Query()
//	maxQuery := queries.Get("max")
//
//	max, err := strconv.Atoi(maxQuery)
//	if maxQuery == "" {
//		max = Limit
//	}
//	fmt.Println("max:", max)
//	if err != nil {
//		badInput := NewError("bad_request", errors.New("the format of the input is invalid"))
//		JSON(response, http.StatusBadRequest, &badInput)
//		return
//	}
//
//	switch {
//	//add a case with a limit of locations returned
//	case max > locationsLength:
//		fmt.Println("max > locationsLength, locations:", locations)
//		JSON(response, http.StatusOK, locations[:Limit])
//		break
//
//	case max >= 1 && max < locationsLength:
//		fmt.Println("max < locationsLength, locations[:max]: ", locations[:max])
//		JSON(response, http.StatusOK, locations[:max])
//		break
//
//	case max < 1:
//		badInput := NewError("bad_request", errors.New("this order does not exist"))
//		JSON(response, http.StatusBadRequest, &badInput)
//		break
//
//	default:
//		badInput := NewError("bad_request", errors.New("this order does not exist"))
//		JSON(response, http.StatusBadRequest, &badInput)
//		break
//	}
//}
//
//func DeleteLocations(response http.ResponseWriter, request *http.Request) {
//	vars := mux.Vars(request)
//	id := vars["order_id"]
//	fmt.Println("id being deleted: ", id)
//	delete(orders, id)
//	JSON(response, http.StatusOK, "Order deleted")
//	fmt.Println("orders map after deletion: ", orders)
//}

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
