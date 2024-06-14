package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Hello World!\n")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr, found := os.LookupEnv("LISTEN_ADDR")
	if !found {
		addr = ":8080"
	}

	router := httprouter.New()
	router.GET("/", Index)

	log.Printf("Server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
