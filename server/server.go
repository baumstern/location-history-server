package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

// Location stores coordinate
type Location struct {
	Lat string `json:"lat"`
	Lng string `json:"lng"`
}

// Server stores location history by order id
type Server struct {
	orders map[string][]Location // key is order id
}

// NewServer creates Server instance
func NewServer() *Server {
	return &Server{orders: make(map[string][]Location)}
}

// HandleLocation branches operation by HTTP method (PUT, GET, DELETE)
func (s *Server) HandleLocation(w http.ResponseWriter, r *http.Request) *appError {
	path := r.URL.Path
	orderId := path[len(`/location/`):]
	if len(orderId) == 0 {
		return &appError{errors.New("order id is missing"), "order id is missing", http.StatusBadRequest}
	}

	switch r.Method {
	case http.MethodPut:
		// get location
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return &appError{errors.New("failed to read request body"), "failed to read request body", http.StatusInternalServerError}
		}
		if len(body) == 0 {
			return &appError{errors.New("empty body"), "body is missing", http.StatusBadRequest}
		}
		location := Location{}

		// TODO: validate request body format
		err = json.Unmarshal(body, &location)
		if err != nil {
			return &appError{errors.New("failed to unmarshal request body"), "could not process provided location value", http.StatusInternalServerError}
		}

		s.putLocation(orderId, location)

	case http.MethodGet:
		// get max parameter if exist
		param := r.URL.Query().Get("max")

		// -1 means there is no limit on the number of locations
		max := -1
		if param != "" {
			converted, err := strconv.Atoi(param)
			if err != nil {
				return &appError{errors.New("failed to convert string interger to integer variable"), "failed to process GET parameter", http.StatusInternalServerError}
			}
			max = converted
		}
		result, appErr := s.getLocation(orderId, max)
		if appErr != nil {
			return appErr
		}

		// TODO: create custom response type
		m, err := json.Marshal(struct {
			Order_id string     `json:"order_id"`
			History  []Location `json:"history"`
		}{
			Order_id: orderId,
			History:  result,
		})
		if err != nil {
			return &appError{err, "json marshal failed", http.StatusInternalServerError}
		}

		io.WriteString(w, string(m))

	case http.MethodDelete:
		s.deleteLocation(orderId)
	}

	return nil
}

func (s *Server) putLocation(orderId string, location Location) {
	if _, ok := s.orders[orderId]; !ok {
		s.orders[orderId] = []Location{}
	}

	s.orders[orderId] = append(s.orders[orderId], location)
	fmt.Println("orderID: ", orderId)
	fmt.Println(s.orders[orderId])
}

// If max < 0, there is no limit on the number of locations
func (s *Server) getLocation(orderId string, max int) ([]Location, *appError) {
	if _, ok := s.orders[orderId]; !ok {
		return nil, &appError{errors.New("order doesn't exist"), "order doesn't exist", http.StatusNotFound}
	}

	result := []Location{}

	switch {
	case max == 0:
		return []Location{}, nil

	case max > 0:
		if max > len(s.orders[orderId]) {
			max = len(s.orders[orderId])
		}

		tail := len(s.orders[orderId]) - 1
		for i := 0; i < max; i++ {
			result = append(result, s.orders[orderId][tail-i])
		}

	case max < 0:
		result = s.orders[orderId]
	}
	return result, nil
}

func (s *Server) deleteLocation(orderId string) (ok bool, _ *appError) {
	if _, ok := s.orders[orderId]; !ok {
		return false, &appError{errors.New("order doesn't exist"), "order doesn't exist", http.StatusNotFound}
	}

	delete(s.orders, orderId)
	return true, nil
}
