package main

import (
	"log"
	"net/http"
)

func main() {
	log.Println("ListenAndServe 8888")
	http.ListenAndServe(":8888", http.HandlerFunc(hello))
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Println(*r)
}
