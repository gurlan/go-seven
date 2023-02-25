package main

import (
	"gee/gee"
	"log"
	"net/http"
)

func main() {

	log.Fatal(http.ListenAndServe(":9999", &gee.Engine{}))
}
