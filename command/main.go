package main

import (
	"log"

	"github.com/bettersun/mockservice"
)

func main() {
	log.Println("mockserice is running")

	mockservice.Load()
	mockservice.MockServiceCommand()

	log.Println("mockserice is stopped")
}
