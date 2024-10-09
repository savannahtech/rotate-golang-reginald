package main

import (
	"fmt"
	"net/http"
)

func (a *serverApplication) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Service is healthy")
}
