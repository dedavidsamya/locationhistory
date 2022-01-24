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

// orders is a map that connects a key (order_id) to a slice of Location structs.
//so each order id is going to have a set of different {latitude, longitude} objects
var orders = map[string][]Location{}

func AddLocation(response http.ResponseWriter, request *http.Request) {
	// location here is going to store the Location passed by the client in the body of the request and unmarshalled by the UnmarshalLocation function.
	location, err := UnmarshalLocation(request)
	fmt.Println("location (after unmarshal):", location, err)
	if err != nil {
		badInput := NewError("bad request", errors.New("the format of the input is invalid"))
		JSON(response, http.StatusBadRequest, &badInput)
		return
	}
	//here, vars is going to contain every piece of data the client sent in the URL of the request
	vars := mux.Vars(request)
	orderId := vars["order_id"]
	fmt.Println("orderId: ", orderId)
	fmt.Println("orders before appending", orders)
	//here a specific struct of Location (the one retrieved from the body of the request and unmarshalled) is appended to the orders map, "inside" the specific orderId key retrieved from the URL.
	//normally the syntax for appending to a map is map["key"] = append(map["key"], object_to_be_appended)
	//but I inverted the order of the parameters of append to actually prepend instead of appending
	//I want the object Location added last to always go to the beginning of the order slice, because I need to retrieve the objects by descending order (last added first)
	orders[orderId] = append([]Location{*location}, orders[orderId]...)
	fmt.Println("orders after appending", orders)
	response.WriteHeader(http.StatusOK)
}

//this function receives the request and returns the data the client passed in the body, in the form of a Location struct, and an error
func UnmarshalLocation(request *http.Request) (*Location, error) {
	var location Location
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

func GetLocation(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	inputOrderId := vars["order_id"]
	fmt.Println(" inputOrderId: ", inputOrderId)

	locations := orders[inputOrderId]
	fmt.Println(" locations (for the order_id) before max:", locations)
	locationsLength := len(locations)
	fmt.Println("locationsLength", locationsLength)

	//this checks if location is contained in orders
	if v, found := orders[inputOrderId]; found {
		fmt.Println(" []Location{location} found in the orders slice:", v)
	}

	max, err := strconv.Atoi(request.URL.Query().Get("max"))
	fmt.Println("max:", max)
	if err != nil {
		badInput := NewError("bad_request", errors.New("the format of the input is invalid"))
		JSON(response, http.StatusBadRequest, &badInput)
		return
	}
	switch {
	case max > locationsLength:
		fmt.Println("max > locationsLength, locations:", locations)
		JSON(response, http.StatusOK, locations)

	case max > 0 && max < locationsLength:
		fmt.Println("max < locationsLength, locations[:max]: ", locations[:max])
		JSON(response, http.StatusOK, locations[:max])

		//CASE MAX > 0 IS NOT WORKING.
	case max < 0:
		badInput := NewError("bad_request", errors.New("this order does not exist"))
		JSON(response, http.StatusBadRequest, &badInput)
		return
	}
}

func DeleteLocations(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id := vars["order_id"]
	fmt.Println("id being deleted: ", id)
	delete(orders, id)
	JSON(response, http.StatusOK, "Order deleted")
	fmt.Println("orders map after deletion: ", orders)
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
