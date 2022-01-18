package server

import (
	"fmt"
	"log"
	"net/http"
)

type AppHandler func(http.ResponseWriter, *http.Request) *appError

func (fn AppHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		log.Println(err.Error())
		http.Error(w, err.responseText, err.status)
	}
}

type appError struct {
	err          error  // wrapped error
	responseText string // response text to user
	status       int    // HTTP status code
}

func (a *appError) Error() string {
	return fmt.Sprintf("%d (%s): %v", a.status, http.StatusText(a.status), a.err)
}
