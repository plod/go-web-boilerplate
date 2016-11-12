package main

import (
	"net/http"
	"fmt"
)

func routing(){
	r.HandleFunc("/", demoHandler)
}

func demoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!")
}
