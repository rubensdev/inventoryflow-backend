package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"github.com/rubensdev/inventoryflow-backend/internal/healthcheck"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	addr, found := os.LookupEnv("LISTEN_ADDR")
	if !found {
		addr = ":8080"
	}

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	healthCheckHandler := healthcheck.NewHealthCheckHandler(logger, healthcheck.SystemInfo{
		Env:     os.Getenv("ENV"),
		Version: os.Getenv("VERSION"),
	})

	router := httprouter.New()
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", healthCheckHandler.StatusHandler)

	log.Printf("Server listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, router))
}
