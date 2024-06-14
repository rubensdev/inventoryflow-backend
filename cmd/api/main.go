package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Hello World!\n")
}

func main() {
	router := httprouter.New()
	router.GET("/", Index)

	log.Println("Server listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
