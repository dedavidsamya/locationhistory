package db

import (
	"errors"
)

type Location struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lng"`
}

type Order struct {
	ID        string
	Locations []Location
}

var orders = map[string]Order{}

func CreateOrder(id string) (*Order, error) {
	if id == "" {
		return nil, errors.New("CreateOrder: empty ID")
	}

	order, _ := GetOrder(id)
	if order != nil {
		return order, errors.New("order_already_exist")
	}

	o := Order{ID: id}
	err := insertOrder(&o)
	if err != nil {
		return nil, err
	}

	return &o, nil
}

func UpdateOrder(order *Order) error {
	if order.ID == "" {
		return errors.New(" UpdateOrder: id is empty")
	}
	o, _ := GetOrder(order.ID)
	if o == nil {
		return errors.New("UpdateOrder: o is nil")
	}

	err := insertOrder(order)
	if err != nil {
		return err
	}

	return nil
}

func GetOrder(id string) (*Order, error) {
	if id == "" {
		return nil, errors.New("GetOrder: id is empty")
	}
	o := orders[id]

	if o.ID == "" {
		return nil, errors.New("GetOrder: order does not exist")
	}
	return &o, nil
}

func AddLocation(id string, location Location) error {
	if id == "" {
		return errors.New("InsertLocation: id is empty")
	}

	order, err := GetOrder(id)
	if err != nil {
		return err
	}

	order.Locations = append(order.Locations, location)
	err = UpdateOrder(order)
	if err != nil {
		return errors.New("InsertLocation: UpdateOrder returned an error")
	}
	return nil
}

func insertOrder(order *Order) error {
	if order == nil {
		return errors.New("insertOrder: order is nil")
	}

	orders[order.ID] = *order
	return nil
}
