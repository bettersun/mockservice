package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

/// HTTP请求方法
const httpMethodGet = "GET"
const httpMethodPost = "POST"
const httpMethodPut = "PUT"
const httpMethodDelete = "DELETE"
const httpMethodHead = "HEAD"

func main() {
	server := http.Server{
		Addr: "127.0.0.1:8012",
	}

	http.HandleFunc("/bettersun", home)
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/goodbye", goodbye)
	http.HandleFunc("/bettersun/hello", helloBS)
	http.HandleFunc("/bettersun/goodbye", goodbyeBS)

	log.Println("ListenAndServe 127.0.0.1:8012")
	server.ListenAndServe()
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome, bettersun")
}

func hello(w http.ResponseWriter, r *http.Request) {

	// 读取请求的Body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	log.Println(fmt.Sprintf("Request Body: %v", string(body)))

	if r.Method == httpMethodGet {
		fmt.Fprintf(w, "[GET] Hello, world.")
		return
	}
	if r.Method == httpMethodPost {
		fmt.Fprintf(w, "[Post] Hello, world.")
		return
	}
	if r.Method == httpMethodPut {
		fmt.Fprintf(w, "[Put] Hello, world.")
		return
	}
	if r.Method == httpMethodDelete {
		fmt.Fprintf(w, "[Delete] Hello, world.")
		return
	}

	fmt.Fprintf(w, "unsupport http method.[hello]")
}

func goodbye(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Goodbye, world.")
}

func helloBS(w http.ResponseWriter, r *http.Request) {

	// 读取请求的Body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	log.Println(fmt.Sprintf("Request Body: %v", string(body)))

	if r.Method == httpMethodGet {
		fmt.Fprintf(w, "[GET] Hello, bettersun.")
		return
	}
	if r.Method == httpMethodPost {
		fmt.Fprintf(w, "[Post] Hello, bettersun.")
		return
	}
	if r.Method == httpMethodPut {
		fmt.Fprintf(w, "[Put] Hello, bettersun.")
		return
	}
	if r.Method == httpMethodDelete {
		fmt.Fprintf(w, "[Delete] Hello, bettersun.")
		return
	}

	fmt.Fprintf(w, "unsupport http method.[bettersun]")
}

func goodbyeBS(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Goodbye, bettersun.")
}
